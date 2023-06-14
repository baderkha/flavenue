package main

import (
	"encoding/json"
	"os"

	"github.com/baderkha/flavenue/internal/api/model"
	"github.com/baderkha/flavenue/internal/pkg/cfg"
	"github.com/baderkha/flavenue/internal/pkg/lib/position"
	"github.com/mmcloughlin/geohash"
)

func main() {
	var data []map[string]interface{}
	db := cfg.GetDB()
	b, err := os.ReadFile("sample_data/data.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &data)
	listings := make([]*model.Listing, 0, len(data))
	geoHashes := make([]*model.GeoHashListing, 0, len(data)*12)
	if err != nil {
		panic(err)
	}

	for _, item := range data {
		var mdl model.Listing
		mdl.New()
		zID, ok := item["zpid"].(string)
		if !ok {
			continue
		}
		mdl.ID = zID
		mdl.Label = item["address"].(string)

		mdl.Street1 = item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["streetAddress"].(string)
		mdl.City = item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["city"].(string)
		mdl.State = item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["state"].(string)
		mdl.Country = item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["country"].(string)
		mdl.PostalCode = item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["zipcode"].(string)
		bdrms, ok := item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["bedrooms"].(float64)
		if ok {
			mdl.Bedrooms = int(bdrms)
		}
		bthrm, ok := item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["bathrooms"].(float64)
		if ok {
			mdl.Bathrooms = int(bthrm)
		}
		area, ok := item["area"].(int)
		if ok {
			mdl.SquareFootage = &area
		}

		mdl.Type = item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["homeType"].(string)
		mdl.Coordinates = *position.NewCoordinates(item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["latitude"].(float64), item["hdpData"].(map[string]interface{})["homeInfo"].(map[string]interface{})["longitude"].(float64))

		listings = append(listings, &mdl)

		for i := 1; i <= 12; i++ {
			var gMdl model.GeoHashListing
			gMdl.New()
			gMdl.GeoHash = geohash.EncodeWithPrecision(mdl.Latitude, mdl.Longtitude, uint(i))
			gMdl.ListingID = mdl.ID
			gMdl.Precision = i
			geoHashes = append(geoHashes, &gMdl)
		}
	}

	err = db.CreateInBatches(&listings, 500).Error
	if err != nil {
		panic(err)
	}
	err = db.CreateInBatches(&geoHashes, 500).Error
	if err != nil {
		panic(err)
	}

}
