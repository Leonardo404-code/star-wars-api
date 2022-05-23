package routes

import (
	"ame-challenge/internal/planets"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/url", planets.GetPlanetsByURL)

	r.GET("/", planets.GetPlanets)

	r.GET("/planet", planets.GetPlanetByName)

	r.GET("/planet/:id", planets.GetPlanetById)

	r.POST("/", planets.CreatePlanet)

	r.DELETE("/:id", planets.DeletePlanet)

	return r
}
