package request

import (
	"regexp"

	"github.com/baderkha/easy-gin/v1/easygin"
	"github.com/baderkha/flavenue/internal/api/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ easygin.IRequest = &AddListing{}

const (
	MaxDistanceAllowedKM = 1000
)

var (
	LatValid = []validation.Rule{
		validation.Required,
		validation.Max(90.0),
		validation.Min(-90.0),
	}
	LngValid = []validation.Rule{
		validation.Required,
		validation.Min(-180.0),
		validation.Max(180.0),
	}
)

type BaseRequest struct {
}

func (b BaseRequest) ValidationErrorFormat(err error) any {
	return map[string]interface{}{
		"data":    nil,
		"message": err.Error(),
	}
}

type GetListingsRelativeToLocation struct {
	BaseRequest
	Latitude   float64 `json:"lat" form:"lat"`
	Longtitude float64 `json:"long" form:"long"`
	DistanceKM int     `json:"distance_km" form:"distance_km"`
}

func (g GetListingsRelativeToLocation) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.Latitude, LatValid...),
		validation.Field(&g.Longtitude, LngValid...),
		validation.Field(&g.DistanceKM, validation.Required, validation.Min(1), validation.Max(MaxDistanceAllowedKM)),
	)
}

type GetListingBoundBox struct {
	BaseRequest
	MinLat float64 `form:"min_lat"`
	MinLng float64 `form:"min_lng"`
	MaxLat float64 `form:"max_lat"`
	MaxLng float64 `form:"max_lng"`
}

func (g GetListingBoundBox) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.MinLat, LatValid...),
		validation.Field(&g.MinLng, LngValid...),
		validation.Field(&g.MaxLat, LatValid...),
		validation.Field(&g.MaxLng, LngValid...),
	)
}

type AddListing struct {
	model.Listing
	BaseRequest
}

func (a *AddListing) Validate() error {
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Street1, validation.Required, validation.Length(5, 50)),
		validation.Field(&a.Street2, validation.Required, validation.Length(5, 50)),
		validation.Field(&a.Street3, validation.Required, validation.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.City, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		validation.Field(&a.State, validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		validation.Field(&a.Country, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
		validation.Field(&a.PostalCode, validation.Required),
		validation.Field(&a.Description, validation.Required),
		validation.Field(&a.Label, validation.Required),
		validation.Field(&a.Latitude, LatValid...),
		validation.Field(&a.Longtitude, LngValid...),
	)
}

type PutListing = AddListing
