package model

import (
	"github.com/baderkha/flavenue/internal/pkg/lib/position"
)

type PropertyType = string

const (
	SingleFamilyPropertyType PropertyType = "SINGLE_FAMILY"
	TownHomesPropertyType    PropertyType = "TOWNHOUSE"
	MultFamPropertyType      PropertyType = "MULTI_FAMILY"
	CondosPropertyType       PropertyType = "CONDO"
	LandPropertyType         PropertyType = "LOT"
)

type Address struct {
	Street1    string `json:"street_1" db:"street_1" gorm:"type:varchar(255)"`
	Street2    string `json:"street_2" db:"street_2" gorm:"type:varchar(255)"`
	Street3    string `json:"street_3" db:"street_3" gorm:"type:varchar(255)"`
	City       string `json:"city" db:"city" gorm:"type:varchar(100)"`
	State      string `json:"state" db:"state" gorm:"type:varchar(2)"`
	Province   string `json:"province" db:"province" gorm:"type:varchar(255)"`
	PostalCode string `json:"postal_code" db:"postal_code" gorm:"type:varchar(12);index:z_code"`
	Country    string `json:"country" db:"country" gorm:"type:varchar(10)"`
}

type PropertySpecs struct {
	SquareFootage *int         `json:"square_footage" db:"square_footage" gorm:"type:int(11)"`
	Bedrooms      int          `json:"bedrooms" db:"bedrooms" gorm:"type:int(11)"`
	Bathrooms     int          `json:"bathrooms" db:"bathrooms" gorm:"type:int(11)"`
	Type          PropertyType `json:"type" db:"type" gorm:"type:varchar(30)"`
	YearBuilt     int          `json:"year_built" db:"year_built" gorm:"type:int(11)"`
	ParkingSpots  int          `json:"parking_spots" db:"parking_spots" gorm:"type:int(11)"`
	HasGarage     *bool        `json:"has_garage" db:"has_garage"`
}

type CoordinateLocation struct {
	position.Coordinates
}

type Listing struct {
	BaseOwned
	Distance    float64 `json:"distance_km" db:"distance"`
	Label       string  `json:"label" db:"label" gorm:"type:varchar(255)"`
	Description string  `json:"description" db:"description"`
	Address
	PropertySpecs
	position.Coordinates
}

func (l Listing) TableName() string {
	return "listings"
}
