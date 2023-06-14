package routes

import (
	"github.com/baderkha/easy-gin/v1/easygin"
	"github.com/baderkha/flavenue/internal/api/controller"
	"github.com/gin-gonic/gin"
)

func Build(app *controller.Rest) *gin.Engine {
	rtr := gin.Default()
	v1 := rtr.Group("/api/v1")
	listings := v1.Group("/listings")
	{
		listings.GET("/search/_relative", easygin.To(app.GetListingsRelToLoc, easygin.BindQuery))
		listings.GET("/search/_map_boundary", easygin.To(app.GetListingBoundBox, easygin.BindQuery))
	}

	return rtr
}
