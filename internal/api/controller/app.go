package controller

import (
	"github.com/baderkha/flavenue/internal/api/repository"
	"github.com/baderkha/flavenue/internal/pkg/cfg"
)

type Rest struct {
	ListingRest
}

func NewRest() *Rest {
	var (
		db = cfg.GetDB()
	)

	return &Rest{
		ListingRest: ListingRest{
			repo: repository.NewMYSQListing(db),
		},
	}
}
