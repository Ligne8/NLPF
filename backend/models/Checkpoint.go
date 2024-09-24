package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Country string

const (
	CountryFrance      Country = "France"
	CountryItaly       Country = "Italy"
	CountrySwitzerland Country = "Switzerland"
	CountrySpain       Country = "Spain"
	CountryPortugal    Country = "Portugal"
)

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

type Checkpoint struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name    City      `json:"name" gorm:"not null" binding:"required"`
	Country Country   `json:"country" gorm:"not null" binding:"required"`
}

func (checkpoint *Checkpoint) BeforeCreate(tx *gorm.DB) (err error) {
	checkpoint.Id = uuid.New()
	if !isCountryValid(checkpoint.Country) {
		return errors.New("invalid country")
	}
	if !isCityInCountry(checkpoint.Name, checkpoint.Country) {
		return errors.New("city does not belong to the specified country or is invalid")
	}
	return nil
}

func isCountryValid(country Country) bool {
	switch country {
	case CountryFrance, CountryItaly, CountrySwitzerland, CountrySpain, CountryPortugal:
		return true
	default:
		return false
	}
}

func isCityInCountry(city City, country Country) bool {
	switch country {
	case CountryFrance:
		return city == CityParis || city == CityMarseille || city == CityPerpignan || city == CityStrasbourg || city == CityLyon
	case CountryItaly:
		return city == CityRome || city == CityFlorence || city == CityMilan || city == CityComo || city == CityNaples
	case CountrySwitzerland:
		return city == CityGeneva || city == CityZurich || city == CityBern || city == CityLausanne || city == CityChatelSaintDenis
	case CountrySpain:
		return city == CityMadrid || city == CityBarcelona || city == CitySeville || city == CityLloretDelMar || city == CityMalaga
	case CountryPortugal:
		return city == CityLisbon || city == CityPorto || city == CityBraga || city == CityLeiria || city == CityEvora
	default:
		return false
	}
}
