package main

import (
	"github.com/baderkha/flavenue/internal/api/model"
	"github.com/baderkha/flavenue/internal/pkg/cfg"
)

func main() {
	var (
		db = cfg.GetDB()
	)

	db.AutoMigrate(&model.Listing{})
	db.AutoMigrate(&model.GeoHashListing{})
}
