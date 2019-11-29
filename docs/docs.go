// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-11-29 00:24:40.7047956 -0300 -03 m=+0.037001501

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/planets": {
            "get": {
                "description": "Gets all added planets",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a planet into the database",
                "parameters": [
                    {
                        "description": "Add Planet Request",
                        "name": "message",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/views.AddPlanetRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    }
                }
            }
        },
        "/planets/forecast": {
            "get": {
                "description": "Obtains a previously generated forecast",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    }
                }
            }
        },
        "/solar_systems": {
            "get": {
                "description": "Gets all added solar systems",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a solar system into the database",
                "parameters": [
                    {
                        "description": "Add Solar System Request",
                        "name": "message",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/views.AddSolarSystemRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    }
                }
            }
        },
        "/solar_systems/": {
            "post": {
                "description": "Generates forecasts given added system and planets",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/views.BaseResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "views.AddPlanetRequest": {
            "type": "object",
            "properties": {
                "clockwise": {
                    "type": "boolean"
                },
                "initial_degrees": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "radio": {
                    "type": "number"
                },
                "solar_system_id": {
                    "type": "integer"
                },
                "speed_by_day": {
                    "type": "number"
                }
            }
        },
        "views.AddSolarSystemRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "views.BaseResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
