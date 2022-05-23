package planets

import (
	"ame-challenge/internal/database"
	"ame-challenge/pkg/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/peterhellberg/swapi"
	"github.com/spf13/viper"
)

// Get Planets by request much more fast
func GetPlanetsByURL(c *gin.Context) {
	url := viper.GetString("PLANETS_URL")

	var planetsObject models.Response

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

	// if err := database.DBConn.Find(&planetsDatabase); err.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err,
	// 	})

	// 	return
	// }

	c.JSON(200, planetsObject)
}

// Get Planets by Client
func GetPlanets(c *gin.Context) {
	client := swapi.DefaultClient

	var planets []models.Planet

	var planetsInDatabase models.Planet

	for i := 1; i <= 6; i++ {
		newPlanet, err := client.Planet(i)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})

			return
		}

		planetsInDatabase = models.Planet{
			Name:         newPlanet.Name,
			Terrain:      newPlanet.Terrain,
			Climate:      newPlanet.Climate,
			MoviesNumber: len(newPlanet.FilmURLs),
		}

		planets = append(planets, planetsInDatabase)

		// if err := database.DBConn.Create(&planetsInDatabase); err.Error != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err.Error,
		// 	})
		// }
	}

	if err := database.DBConn.Find(&planetsInDatabase); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error,
		})
	}

	planets = append(planets, planetsInDatabase)

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
