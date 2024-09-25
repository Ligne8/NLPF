package controllers

import (
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CheckpointController struct {
	Db *gorm.DB
}

// GetAllCheckpoints Retrieve all checkpoints
//
//		@Summary      List checkpoints
//		@Description  get all checkpoints
//		@Tags         checkpoints
//		@Accept       json
//		@Produce      json
//		@Success      200  {array}   models.Checkpoint
//	 	@Failure   	  500 "Unable to retrieve checkpoints"
//		@Router       /checkpoints [get]
func (controller *CheckpointController) GetAllCheckpoints(c *gin.Context) {
	var checkpoints []models.Checkpoint
	if err := controller.Db.Find(&checkpoints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve checkpoints"})
		return
	}
	c.JSON(http.StatusOK, checkpoints)
}

// GetCitiesByCountry Retrieve cities by country
//
//	@Summary      List cities by country
//	@Description  get all cities for a given country
//	@Tags         checkpoints
//	@Accept       json
//	@Produce      json
//	@Param        country  path      string  true  "Country"
//	@Success      200  {array}   models.City
//	@Failure	  404	"Country not found"
//	@Router       /checkpoints/cities/{country} [get]
func (controller *CheckpointController) GetCitiesByCountry(c *gin.Context) {
	country := c.Param("country")
	switch models.Country(country) {
	case models.CountryFrance:
		c.JSON(http.StatusOK, []models.City{
			models.CityParis, models.CityMarseille, models.CityPerpignan, models.CityStrasbourg, models.CityLyon,
		})
	case models.CountryItaly:
		c.JSON(http.StatusOK, []models.City{
			models.CityRome, models.CityFlorence, models.CityMilan, models.CityComo, models.CityNaples,
		})
	case models.CountrySwitzerland:
		c.JSON(http.StatusOK, []models.City{
			models.CityGeneva, models.CityZurich, models.CityBern, models.CityLausanne, models.CityChatelSaintDenis,
		})
	case models.CountrySpain:
		c.JSON(http.StatusOK, []models.City{
			models.CityMadrid, models.CityBarcelona, models.CitySeville, models.CityLloretDelMar, models.CityMalaga,
		})
	case models.CountryPortugal:
		c.JSON(http.StatusOK, []models.City{
			models.CityLisbon, models.CityPorto, models.CityBraga, models.CityLeiria, models.CityEvora,
		})
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
	}
}

// GetCountryByCity Retrieve the country by city
//
//	@Summary      Get country by city
//	@Description  get the country for a given city
//	@Tags         checkpoints
//	@Accept       json
//	@Produce      json
//	@Param        city  path      string  true  "City"
//	@Success      200  {string}  string  "Country"
//	@Failure	  404	"City not found"
//	@Router       /checkpoints/country/{city} [get]
func (controller *CheckpointController) GetCountryByCity(c *gin.Context) {
	city := c.Param("city")
	switch models.City(city) {
	case models.CityParis, models.CityMarseille, models.CityPerpignan, models.CityStrasbourg, models.CityLyon:
		c.JSON(http.StatusOK, models.CountryFrance)
	case models.CityRome, models.CityFlorence, models.CityMilan, models.CityComo, models.CityNaples:
		c.JSON(http.StatusOK, models.CountryItaly)
	case models.CityGeneva, models.CityZurich, models.CityBern, models.CityLausanne, models.CityChatelSaintDenis:
		c.JSON(http.StatusOK, models.CountrySwitzerland)
	case models.CityMadrid, models.CityBarcelona, models.CitySeville, models.CityLloretDelMar, models.CityMalaga:
		c.JSON(http.StatusOK, models.CountrySpain)
	case models.CityLisbon, models.CityPorto, models.CityBraga, models.CityLeiria, models.CityEvora:
		c.JSON(http.StatusOK, models.CountryPortugal)
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
	}
}
