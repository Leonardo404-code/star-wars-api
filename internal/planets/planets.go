package planets

import (
	"ame-challenge/internal/database"
	"ame-challenge/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/peterhellberg/swapi"
)

func GetPlanets(c *gin.Context) {
	client := swapi.DefaultClient

	var planets []models.Planet

	var formatPlanet models.Planet

	for i := 1; i <= 6; i++ {
		newPlanet, err := client.Planet(i)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})

			return
		}

		formatPlanet.Name = newPlanet.Name
		formatPlanet.Ground = newPlanet.Terrain
		formatPlanet.Weather = newPlanet.Climate
		formatPlanet.MoviesNumber = len(newPlanet.FilmURLs)

		planets = append(planets, formatPlanet)
	}

	// if err := database.DBConn.Find(&planets); err.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error,
	// 	})
	// }

	c.JSON(200, planets)
}

func GetPlanetById(c *gin.Context) {
	client := swapi.DefaultClient

	planetID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	if planetID < 1 || planetID > 6 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "planet id must be betwenn 0 and 7",
		})

		return
	}

	newPlanet, err := client.Planet(planetID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error in find planet: %v" + err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"name":          newPlanet.Name,
		"weather":       newPlanet.Climate,
		"terrain":       newPlanet.Terrain,
		"movies_number": len(newPlanet.FilmURLs),
	})
}

func GetPlanetByName(c *gin.Context) {
	c.String(200, "Hi")
}

func CreatePlanet(c *gin.Context) {
	var planets models.Planet

	if err := c.ShouldBind(&planets); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})

		return
	}

	if planets.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name is required",
		})

		return
	}

	if planets.Weather == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "weather is required",
		})

		return
	}

	if planets.Ground == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ground is required",
		})

		return
	}

	if planets.MoviesNumber < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "movies_number is required",
		})

		return
	}

	if err := database.DBConn.Create(&planets); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error,
		})
	}

	c.JSON(http.StatusOK, planets)
}

func DeletePlanet(c *gin.Context) {
	planetID := c.Param("id")

	var planet models.Planet

	planetFind := database.DBConn.Find(&planet, planetID).RowsAffected

	if planetFind == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "planet not found",
		})

		return
	}

	if err := database.DBConn.Delete(&planet, planetID); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error,
		})

		return
	}

	c.JSON(http.StatusOK, planet)
}
