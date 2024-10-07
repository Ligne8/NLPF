package models

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Country represents a country
// @Description Represents a country
type Country string

const (
	CountryFrance      Country = "France"
	CountryItaly       Country = "Italy"
	CountrySwitzerland Country = "Switzerland"
	CountrySpain       Country = "Spain"
	CountryPortugal    Country = "Portugal"
)

// City represents a city name
// @Description Represents a city name within a country
type City string

const (
	// Villes en France
	CityParis      City = "Paris"
	CityMarseille  City = "Marseille"
	CityPerpignan  City = "Perpignan"
	CityStrasbourg City = "Strasbourg"
	CityLyon       City = "Lyon"

	// Villes en Italie
	CityRome     City = "Rome"
	CityFlorence City = "Florence"
	CityMilan    City = "Milan"
	CityComo     City = "Como"
	CityNaples   City = "Naples"

	// Villes en Suisse
	CityGeneva           City = "Geneva"
	CityZurich           City = "Zurich"
	CityBern             City = "Bern"
	CityLausanne         City = "Lausanne"
	CityChatelSaintDenis City = "Chatel-Saint-Denis"

	// Villes en Espagne
	CityMadrid       City = "Madrid"
	CityBarcelona    City = "Barcelona"
	CitySeville      City = "Seville"
	CityLloretDelMar City = "Lloret del Mar"
	CityMalaga       City = "Malaga"

	// Villes au Portugal
	CityLisbon City = "Lisbon"
	CityPorto  City = "Porto"
	CityBraga  City = "Braga"
	CityLeiria City = "Leiria"
	CityEvora  City = "Evora"
)

// Checkpoint represents a geographic checkpoint
// @Description Represents a checkpoint with a city and a country
type Checkpoint struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      City      `json:"name" gorm:"not null" binding:"required"`
	Country   Country   `json:"country" gorm:"not null" binding:"required"`
	Longitude float64   `json:"longitude" gorm:"not null" binding:"required"`
	Latitude  float64   `json:"latitude" gorm:"not null" binding:"required"`
}

func (checkpoint *Checkpoint) BeforeCreate(tx *gorm.DB) (err error) {
	checkpoint.Id = uuid.New()
	return nil
}

func CreateCheckpoints(db *gorm.DB) {
	// Liste des checkpoints à créer par pays
	checkpoints := []Checkpoint{
		// France
		{Id: uuid.New(), Name: CityParis, Country: CountryFrance, Latitude: 48.8566, Longitude: 2.3522},
		{Id: uuid.New(), Name: CityMarseille, Country: CountryFrance, Latitude: 43.2965, Longitude: 5.3698},
		{Id: uuid.New(), Name: CityPerpignan, Country: CountryFrance, Latitude: 42.6887, Longitude: 2.8948},
		{Id: uuid.New(), Name: CityStrasbourg, Country: CountryFrance, Latitude: 48.5734, Longitude: 7.7521},
		{Id: uuid.New(), Name: CityLyon, Country: CountryFrance, Latitude: 45.7640, Longitude: 4.8357},

		// Italy
		{Id: uuid.New(), Name: CityRome, Country: CountryItaly, Latitude: 41.9028, Longitude: 12.4964},
		{Id: uuid.New(), Name: CityFlorence, Country: CountryItaly, Latitude: 43.7696, Longitude: 11.2558},
		{Id: uuid.New(), Name: CityMilan, Country: CountryItaly, Latitude: 45.4642, Longitude: 9.1900},
		{Id: uuid.New(), Name: CityComo, Country: CountryItaly, Latitude: 45.8081, Longitude: 9.0852},
		{Id: uuid.New(), Name: CityNaples, Country: CountryItaly, Latitude: 40.8518, Longitude: 14.2681},

		// Switzerland
		{Id: uuid.New(), Name: CityGeneva, Country: CountrySwitzerland, Latitude: 46.2044, Longitude: 6.1432},
		{Id: uuid.New(), Name: CityZurich, Country: CountrySwitzerland, Latitude: 47.3769, Longitude: 8.5417},
		{Id: uuid.New(), Name: CityBern, Country: CountrySwitzerland, Latitude: 46.9481, Longitude: 7.4474},
		{Id: uuid.New(), Name: CityLausanne, Country: CountrySwitzerland, Latitude: 46.5197, Longitude: 6.6323},
		{Id: uuid.New(), Name: CityChatelSaintDenis, Country: CountrySwitzerland, Latitude: 46.5270, Longitude: 6.8985},

		// Spain
		{Id: uuid.New(), Name: CityMadrid, Country: CountrySpain, Latitude: 40.4168, Longitude: -3.7038},
		{Id: uuid.New(), Name: CityBarcelona, Country: CountrySpain, Latitude: 41.3851, Longitude: 2.1734},
		{Id: uuid.New(), Name: CitySeville, Country: CountrySpain, Latitude: 37.3891, Longitude: -5.9845},
		{Id: uuid.New(), Name: CityLloretDelMar, Country: CountrySpain, Latitude: 41.6994, Longitude: 2.8455},
		{Id: uuid.New(), Name: CityMalaga, Country: CountrySpain, Latitude: 36.7213, Longitude: -4.4214},

		// Portugal
		{Id: uuid.New(), Name: CityLisbon, Country: CountryPortugal, Latitude: 38.7223, Longitude: -9.1393},
		{Id: uuid.New(), Name: CityPorto, Country: CountryPortugal, Latitude: 41.1579, Longitude: -8.6291},
		{Id: uuid.New(), Name: CityBraga, Country: CountryPortugal, Latitude: 41.5454, Longitude: -8.4265},
		{Id: uuid.New(), Name: CityLeiria, Country: CountryPortugal, Latitude: 39.7436, Longitude: -8.8071},
		{Id: uuid.New(), Name: CityEvora, Country: CountryPortugal, Latitude: 38.5710, Longitude: -7.9137},
	}

	for _, checkpoint := range checkpoints {
		var existing Checkpoint
		if err := db.Where("name = ? AND country = ?", checkpoint.Name, checkpoint.Country).First(&existing).Error; err == nil {
			log.Printf("Checkpoint already exists: %s, %s", checkpoint.Name, checkpoint.Country)
			continue
		}

		if err := db.Create(&checkpoint).Error; err != nil {
			log.Printf("Failed to create checkpoint %s: %v", checkpoint.Name, err)
		} else {
			log.Printf("Checkpoint created: %s, %s", checkpoint.Name, checkpoint.Country)
		}
	}
}
