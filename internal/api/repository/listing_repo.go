package repository

import (
	"errors"
	"fmt"

	"github.com/baderkha/flavenue/internal/api/model"
	"github.com/baderkha/flavenue/internal/pkg/lib/position"
	"github.com/baderkha/library/pkg/store/repository"
	"github.com/jftuga/geodist"
	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

const (
	MaxGeoHashPrecision = 8
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
type BoxedAreaQuery = geohash.Box

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

func (m *MYSQListing) Create(mdl *model.Listing) error {
	var (
		geoHashes = make([]*model.GeoHashListing, 0, MaxGeoHashPrecision)
	)
	tx := repository.GormTransaction{
		DB: m.DB,
	}
	rpoListings := m.WithTransaction(&tx)
	rpoGeo := m.geoHashRepo.WithTransaction(&tx)

	hash := geohash.EncodeWithPrecision(mdl.Latitude, mdl.Longtitude, MaxGeoHashPrecision)
	for i := 0; i < len(hash); i++ {
		var gmdl model.GeoHashListing
		gmdl.New()
		gmdl.ListingID = mdl.ID
		gmdl.Precision = i + 1
		gmdl.GeoHash = hash[:i]
		geoHashes = append(geoHashes, &gmdl)
	}

	err := rpoGeo.BulkCreate(geoHashes)
	if err != nil {
		tx.RollBack()
		return err
	}
	err = rpoListings.Create(mdl)
	if err != nil {
		tx.RollBack()
		return err
	}
	return tx.Commit()
}

func (m *MYSQListing) Update(mdl *model.Listing) error {

	var (
		geoHashes = make([]*model.GeoHashListing, 0, MaxGeoHashPrecision)
	)

	tx := m.DB.Begin()

	hash := geohash.EncodeWithPrecision(mdl.Latitude, mdl.Longtitude, MaxGeoHashPrecision)
	for i := 0; i < len(hash); i++ {
		var gmdl model.GeoHashListing
		gmdl.New()
		gmdl.ListingID = mdl.ID
		gmdl.Precision = i + 1
		gmdl.GeoHash = hash[:i]
		geoHashes = append(geoHashes, &gmdl)
	}
	err := tx.Unscoped().Where("listing_id=?", mdl.ID).Delete(&model.GeoHashListing{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.CreateInBatches(geoHashes, MaxGeoHashPrecision).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Updates(mdl).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error

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

	_, distKM := geodist.HaversineDistance(geodist.Coord{
		Lat: p.MaxLat,
		Lon: p.MaxLng,
	}, geodist.Coord{
		Lat: p.MinLat,
		Lon: p.MinLng,
	})
	distKM = distKM / 2
	var (
		precision    = position.DetermineGeoHashPrecision(float32(distKM))
		lt, lng      = p.Center()
		midHash      = geohash.EncodeWithPrecision(lt, lng, precision)
		allNeighbors = geohash.Neighbors(midHash)
		lstTName     = model.Listing{}.TableName()
		geoTName     = model.GeoHashListing{}.TableName()
		res          = make([]*model.Listing, 0, 1000)
	)

	allNeighbors = append(allNeighbors, midHash)

	err := s.DB.Raw(
		fmt.Sprintf(`
	SELECT 
		%s.*
	FROM %s 
	INNER JOIN %s
		ON %s.id=%s.listing_id
	WHERE 
			%s.geo_hash in(?)
		AND
			%s.precision=?
	`,
			lstTName,
			lstTName,
			geoTName,
			lstTName,
			geoTName,
			geoTName,
			geoTName,
		),
		allNeighbors,
		precision,
	).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
