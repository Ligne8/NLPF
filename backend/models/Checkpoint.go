package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
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
	return nil
}

func CreateCheckpoints(db *gorm.DB) {
	// Liste des checkpoints à créer par pays
	checkpoints := []Checkpoint{
		// France
		{Id: uuid.New(), Name: CityParis, Country: CountryFrance},
		{Id: uuid.New(), Name: CityMarseille, Country: CountryFrance},
		{Id: uuid.New(), Name: CityPerpignan, Country: CountryFrance},
		{Id: uuid.New(), Name: CityStrasbourg, Country: CountryFrance},
		{Id: uuid.New(), Name: CityLyon, Country: CountryFrance},

		// Italie
		{Id: uuid.New(), Name: CityRome, Country: CountryItaly},
		{Id: uuid.New(), Name: CityFlorence, Country: CountryItaly},
		{Id: uuid.New(), Name: CityMilan, Country: CountryItaly},
		{Id: uuid.New(), Name: CityComo, Country: CountryItaly},
		{Id: uuid.New(), Name: CityNaples, Country: CountryItaly},

		// Suisse
		{Id: uuid.New(), Name: CityGeneva, Country: CountrySwitzerland},
		{Id: uuid.New(), Name: CityZurich, Country: CountrySwitzerland},
		{Id: uuid.New(), Name: CityBern, Country: CountrySwitzerland},
		{Id: uuid.New(), Name: CityLausanne, Country: CountrySwitzerland},
		{Id: uuid.New(), Name: CityChatelSaintDenis, Country: CountrySwitzerland},

		// Espagne
		{Id: uuid.New(), Name: CityMadrid, Country: CountrySpain},
		{Id: uuid.New(), Name: CityBarcelona, Country: CountrySpain},
		{Id: uuid.New(), Name: CitySeville, Country: CountrySpain},
		{Id: uuid.New(), Name: CityLloretDelMar, Country: CountrySpain},
		{Id: uuid.New(), Name: CityMalaga, Country: CountrySpain},

		// Portugal
		{Id: uuid.New(), Name: CityLisbon, Country: CountryPortugal},
		{Id: uuid.New(), Name: CityPorto, Country: CountryPortugal},
		{Id: uuid.New(), Name: CityBraga, Country: CountryPortugal},
		{Id: uuid.New(), Name: CityLeiria, Country: CountryPortugal},
		{Id: uuid.New(), Name: CityEvora, Country: CountryPortugal},
	}

	for _, checkpoint := range checkpoints {
		if err := db.FirstOrCreate(&checkpoint, Checkpoint{Name: checkpoint.Name, Country: checkpoint.Country}).Error; err != nil {
			log.Printf("Failed to create checkpoint %s: %v", checkpoint.Name, err)
		} else {
			log.Printf("Checkpoint created: %s, %s", checkpoint.Name, checkpoint.Country)
		}
	}
}
