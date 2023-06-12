// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "University Library Basel, Informatik",
            "url": "https://ub.unibas.ch",
            "email": "it-ub@unibas.ch"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/{collection}/": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "creates new media item in database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "media"
                ],
                "summary": "new media entry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Collection Name",
                        "name": "collection",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Add new media",
                        "name": "NewMediaRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.NewMediaItemRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.HTTPResultMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.HTTPResultMessage"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/rest.HTTPResultMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.HTTPResultMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.NewMediaItemRequestBody": {
            "type": "object",
            "properties": {
                "signature": {
                    "description": "Collection string ` + "`" + `json:\"collection\" example:\"erara\" format:\"string\"` + "`" + `",
                    "type": "string",
                    "format": "string",
                    "example": "sig-4711"
                },
                "urn": {
                    "type": "string",
                    "format": "string",
                    "example": "vfs://digispace/data/test.zip/image.tif"
                }
            }
        },
        "rest.HTTPResultMessage": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Bearer Authentication with JWT",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Media Server Ingest API",
	Description:      "Ingesting Media files",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}