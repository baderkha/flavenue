package controller

import (
	"github.com/baderkha/easy-gin/v1/easygin"
	"github.com/baderkha/flavenue/internal/api/repository"
	"github.com/baderkha/flavenue/internal/api/request"

	"github.com/baderkha/flavenue/internal/pkg/lib/position"
	"github.com/mmcloughlin/geohash"
)

type ListingRest struct {
	repo repository.IListing
}

func (l *ListingRest) AddListing(r *request.AddListing) *easygin.Response {
	err := l.repo.Create(&r.Listing)
	if err != nil {
		return easygin.Err(err)
	}
	return easygin.Res(r.Listing)
}

func (l *ListingRest) UpdateListing(r *request.PutListing) *easygin.Response {
	err := l.repo.Update(&r.Listing)
	if err != nil {
		return easygin.Err(err)
	}
	return easygin.Res(r.Listing)
}

func (l *ListingRest) GetListingsRelToLoc(r request.GetListingsRelativeToLocation) *easygin.Response {
	res, err := l.repo.GetAllRelativeToPosition(&repository.RelativePositionQuery{
		Coordinates:      *position.NewCoordinates(r.Latitude, r.Longtitude),
		RadiusDistanceKM: r.DistanceKM,
	})
	if err != nil {
		return easygin.Err(err)
	}
	return easygin.Res(res)
}

func (l *ListingRest) GetListingBoundBox(r request.GetListingBoundBox) *easygin.Response {
	res, err := l.repo.GetAllBoxedPosition(&geohash.Box{
		MinLat: r.MinLat,
		MaxLat: r.MaxLat,
		MinLng: r.MinLng,
		MaxLng: r.MaxLng,
	})
	if err != nil {
		return easygin.Err(err)
	}
	return easygin.Res(res)
}
