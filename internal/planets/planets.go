package planets

import (
	"ame-challenge/internal/database"
	"ame-challenge/pkg/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/peterhellberg/swapi"
)

// Get Planets by request much more fast
func GetPlanets(c *gin.Context) {
	url := "https://swapi.dev/api/planets/?format=json"

	var (
		planetsObject models.Response

		planetsDatabase []models.Planet
	)

	resp, err := http.Get(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	json.Unmarshal(respBody, &planetsObject)

	if err := database.DBConn.Find(&planetsDatabase); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}

	c.JSON(200, gin.H{
		"api":      planetsObject.Result,
		"database": planetsDatabase,
	})
}

func GetPlanetById(c *gin.Context) {
	client := swapi.DefaultClient

	planetID, err := strconv.Atoi(c.Param("id"))

	var planetDatabase models.Planet

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	newPlanet, err := client.Planet(planetID)

	if planetID < 1 || planetID > 6 {
		newPlanet = swapi.Planet{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error in find planet: %v" + err.Error(),
		})

		return
	}

	if err := database.DBConn.Where("id = ?", planetID).Find(&planetDatabase); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error,
		})

		return
	}

	c.JSON(200, gin.H{
		"database": planetDatabase,
		"api": models.Planet{
			Name:    newPlanet.Name,
			Climate: newPlanet.Climate,
			Terrain: newPlanet.Terrain,
			Films:   len(newPlanet.FilmURLs),
		},
	})
}

func GetPlanetByName(c *gin.Context) {
	name := c.Query("name")
	client := swapi.DefaultClient

	var (
		planetsDatabase models.Planet
		planet          swapi.Planet
	)

	if err := database.DBConn.Table("planets").Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).
		Find(&planetsDatabase); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}

	switch {
	case name == "Tatooine":
		planet, _ = client.Planet(1)
	case name == "Alderaan":
		planet, _ = client.Planet(2)
	case name == "Yavin IV" || name == "Yavin":
		planet, _ = client.Planet(3)
	case name == "Hoth":
		planet, _ = client.Planet(4)
	case name == "Dagobah":
		planet, _ = client.Planet(5)
	case name == "Bespin":
		planet, _ = client.Planet(6)
	default:
		planet = swapi.Planet{}
	}

	c.JSON(http.StatusOK, gin.H{
		"database": planetsDatabase,
		"api": models.Planet{
			Name:    planet.Name,
			Climate: planet.Climate,
			Terrain: planet.Terrain,
			Films:   len(planet.FilmURLs),
		},
	})
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

	if planets.Climate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "weather is required",
		})

		return
	}

	if planets.Terrain == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ground is required",
		})

		return
	}

	if planets.Films < 1 {
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

	if err := database.DBConn.Where("id = ?", planetID).Delete(&planet); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error,
		})

		return
	}

	c.JSON(http.StatusOK, planet)
}
