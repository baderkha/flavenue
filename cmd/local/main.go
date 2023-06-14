package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/baderkha/flavenue/internal/api/model"
	"github.com/baderkha/flavenue/internal/api/repository"
	"github.com/baderkha/flavenue/internal/pkg/lib/position"
	"github.com/davecgh/go-spew/spew"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "main:password@tcp(127.0.0.1:6001)/main?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	db.Logger = newLogger
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.Listing{})
	db.AutoMigrate(&model.GeoHashListing{})

	rpo := repository.NewMYSQListing(db)
	t := time.Now()
	res, err := rpo.GetAllRelativeToPosition(&repository.RelativePositionQuery{
		Coordinates:      *position.NewCoordinates(40.797009, -74.110291),
		RadiusDistanceKM: 20,
	})
	duration := time.Since(t)
	fmt.Println(duration)
	if err != nil {
		panic(err)
	}
	spew.Dump(res[len(res)-10])
}
