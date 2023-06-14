package api

import (
	"github.com/baderkha/flavenue/internal/api/controller"
	"github.com/baderkha/flavenue/internal/api/routes"
)

func Run() {
	ctrlr := controller.NewRest()
	rtr := routes.Build(ctrlr)
	rtr.Run(":8082")
}
