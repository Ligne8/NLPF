// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/checkpoints": {
            "get": {
                "description": "get all checkpoints",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "checkpoints"
                ],
                "summary": "List checkpoints",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Checkpoint"
                            }
                        }
                    },
                    "500": {
                        "description": "Unable to retrieve checkpoints"
                    }
                }
            }
        },
        "/checkpoints/cities/{country}": {
            "get": {
                "description": "get all cities for a given country",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "checkpoints"
                ],
                "summary": "List cities by country",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Country",
                        "name": "country",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.City"
                            }
                        }
                    },
                    "404": {
                        "description": "Country not found"
                    }
                }
            }
        },
        "/checkpoints/country/{city}": {
            "get": {
                "description": "get the country for a given city",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "checkpoints"
                ],
                "summary": "Get country by city",
                "parameters": [
                    {
                        "type": "string",
                        "description": "City",
                        "name": "city",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Country",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "City not found"
                    }
                }
            }
        },
        "/user": {
            "post": {
                "description": "create a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "user",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Invalid request"
                    },
                    "500": {
                        "description": "Unable to retrieve user"
                    }
                }
            }
        },
        "/user/:id": {
            "get": {
                "description": "get a user by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Invalid userId"
                    },
                    "500": {
                        "description": "Unable to retrieve user"
                    }
                }
            },
            "delete": {
                "description": "deletes a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Invalid request"
                    },
                    "500": {
                        "description": "Unable to retrieve user or unable to delete user"
                    }
                }
            },
            "patch": {
                "description": "updates a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "user",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Invalid request"
                    },
                    "500": {
                        "description": "Unable to retrieve user"
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "get all users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "List users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    },
                    "500": {
                        "description": "Unable to retrieve users"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Checkpoint": {
            "description": "Represents a checkpoint with a city and a country",
            "type": "object",
            "required": [
                "country",
                "latitude",
                "longitude",
                "name"
            ],
            "properties": {
                "country": {
                    "$ref": "#/definitions/models.Country"
                },
                "id": {
                    "type": "string"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "name": {
                    "$ref": "#/definitions/models.City"
                }
            }
        },
        "models.City": {
            "description": "Represents a city name within a country",
            "type": "string",
            "enum": [
                "Paris",
                "Marseille",
                "Perpignan",
                "Strasbourg",
                "Lyon",
                "Rome",
                "Florence",
                "Milan",
                "Como",
                "Naples",
                "Geneva",
                "Zurich",
                "Bern",
                "Lausanne",
                "Chatel-Saint-Denis",
                "Madrid",
                "Barcelona",
                "Seville",
                "Lloret del Mar",
                "Malaga",
                "Lisbon",
                "Porto",
                "Braga",
                "Leiria",
                "Evora"
            ],
            "x-enum-varnames": [
                "CityParis",
                "CityMarseille",
                "CityPerpignan",
                "CityStrasbourg",
                "CityLyon",
                "CityRome",
                "CityFlorence",
                "CityMilan",
                "CityComo",
                "CityNaples",
                "CityGeneva",
                "CityZurich",
                "CityBern",
                "CityLausanne",
                "CityChatelSaintDenis",
                "CityMadrid",
                "CityBarcelona",
                "CitySeville",
                "CityLloretDelMar",
                "CityMalaga",
                "CityLisbon",
                "CityPorto",
                "CityBraga",
                "CityLeiria",
                "CityEvora"
            ]
        },
        "models.Country": {
            "description": "Represents a country",
            "type": "string",
            "enum": [
                "France",
                "Italy",
                "Switzerland",
                "Spain",
                "Portugal"
            ],
            "x-enum-varnames": [
                "CountryFrance",
                "CountryItaly",
                "CountrySwitzerland",
                "CountrySpain",
                "CountryPortugal"
            ]
        },
        "models.Role": {
            "type": "string",
            "enum": [
                "trader",
                "traffic_manager",
                "client",
                "admin"
            ],
            "x-enum-varnames": [
                "RoleTrader",
                "RoleTrafficManager",
                "RoleClient",
                "RoleAdmin"
            ]
        },
        "models.User": {
            "type": "object",
            "required": [
                "firstname",
                "lastname"
            ],
            "properties": {
                "firstname": {
                    "type": "string"
                },
                "hashed_password": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "lastname": {
                    "type": "string"
                },
                "role": {
                    "$ref": "#/definitions/models.Role"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
