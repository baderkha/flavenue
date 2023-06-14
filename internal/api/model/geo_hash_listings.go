package model

type GeoHashListing struct {
	Base
	GeoHash   string `json:"geo_hash" db:"geo_hash" gorm:"type:varchar(12);index:geo_idx"`
	Precision int    `json:"precision" db:"precision" gorm:"type:int(11);index"`
	ListingID string `json:"listing_id" db:"listing_id" gorm:"type:varchar(50);index:geo_idx"`
}

func (l GeoHashListing) TableName() string {
	return "geo_hash_listings"
}

func (l GeoHashListing) GetAccountID() string {
	return ""
}
