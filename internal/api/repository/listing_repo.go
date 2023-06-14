package repository

import (
	"errors"
	"fmt"

	"github.com/baderkha/flavenue/internal/api/model"
	"github.com/baderkha/flavenue/internal/pkg/lib/position"
	"github.com/baderkha/library/pkg/store/repository"
	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

var (
	ErrorCouldNotAddGeoKey    = errors.New("could not add a geospatial point to the db, this is unexpected")
	ErrorCouldNotUpdateGeoKey = errors.New("could not update a geospatial point to the db, this is unexpected")
)

type RelativePositionQuery struct {
	position.Coordinates
	RadiusDistanceKM int
}

// BoxedAreaQuery : you give it a box and it will query for items within that area
type BoxedAreaQuery struct {
	NPos float64
	SPos float64
	EPos float64
	WPos float64
}

type IListing interface {
	repository.ICrud[model.Listing]
	GetAllRelativeToPosition(p *RelativePositionQuery) ([]*model.Listing, error)
	GetAllBoxedPosition(p *BoxedAreaQuery) ([]*model.Listing, error)
}

var _ IListing = &MYSQListing{}

func NewMYSQListing(db *gorm.DB) *MYSQListing {
	return &MYSQListing{
		geoHashRepo: repository.CrudGorm[model.GeoHashListing]{
			DB: db,
		},
		CrudGorm: repository.CrudGorm[model.Listing]{
			DB: db,
		},
	}
}

type MYSQListing struct {
	repository.CrudGorm[model.Listing]
	geoHashRepo repository.CrudGorm[model.GeoHashListing]
}

func (s *MYSQListing) GetAllRelativeToPosition(p *RelativePositionQuery) ([]*model.Listing, error) {
	var (
		precision = position.DetermineGeoHashPrecision(float32(p.RadiusDistanceKM))
		midHash   = geohash.EncodeWithPrecision(p.Latitude, p.Longtitude, precision)
		lstTName  = model.Listing{}.TableName()
		geoTName  = model.GeoHashListing{}.TableName()
		res       = make([]*model.Listing, 0, 1000)
	)

	allNeighbors := geohash.Neighbors(midHash)
	allNeighbors = append(allNeighbors, midHash)

	err := s.DB.Raw(
		fmt.Sprintf(`
	SELECT 
		%s.*,
		(ST_Distance_Sphere(
				Point(?,?), 
				Point(%s.longtitude, %s.latitude)
				)/1000)
				as distance
	FROM %s 
	INNER JOIN %s
		ON %s.id=%s.listing_id
	WHERE 
			%s.geo_hash in(?)
		AND
			%s.precision=?
	HAVING
		distance <= ?
	ORDER BY
		distance asc
	`,
			lstTName,
			lstTName,
			lstTName,
			lstTName,
			geoTName,
			lstTName,
			geoTName,
			geoTName,
			geoTName,
		),
		p.Longtitude,
		p.Latitude,
		allNeighbors,
		precision,
		p.RadiusDistanceKM,
	).Find(&res).Error

	return res, err
}

func (s *MYSQListing) GetAllBoxedPosition(p *BoxedAreaQuery) ([]*model.Listing, error) {
	return nil, nil
}
