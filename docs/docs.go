// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "name": "Rizal Arfiyan",
            "url": "https://rizalrfiyan.com",
            "email": "rizal.arfiyan.23@gmail.com"
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
        "/": {
            "get": {
                "description": "Base Home",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "home"
                ],
                "summary": "Get Base Home based on parameter",
                "operationId": "get-base-home",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    }
                }
            }
        },
        "/auth/google": {
            "get": {
                "description": "Auth Google Redirection",
                "tags": [
                    "auth"
                ],
                "summary": "Get Auth Google Redirection based on parameter",
                "operationId": "get-auth-google",
                "responses": {
                    "307": {
                        "description": "Temporary Redirect"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/google/callback": {
            "get": {
                "description": "Auth Google Callback Redirection",
                "tags": [
                    "auth"
                ],
                "summary": "Get Auth Google Callback Redirection based on parameter",
                "operationId": "get-auth-google-callback",
                "responses": {
                    "307": {
                        "description": "Temporary Redirect"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/verification": {
            "post": {
                "description": "Auth Verification",
                "tags": [
                    "auth"
                ],
                "summary": "Post Auth Verification based on parameter",
                "operationId": "post-auth-verification",
                "parameters": [
                    {
                        "description": "Data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AuthVerification"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.BaseResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.AuthVerification"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Base Health",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "home"
                ],
                "summary": "Get Base Health based on parameter",
                "operationId": "get-base-health",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "constants.AuthVerificationStep": {
            "type": "integer",
            "enum": [
                1,
                2,
                3
            ],
            "x-enum-varnames": [
                "AuthVerificationRegister",
                "AuthVerificationOtp",
                "AuthVerificationDone"
            ]
        },
        "request.AuthVerification": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string",
                    "example": "2YbPyusF2G06BFQLamoKFXvGgPd"
                }
            }
        },
        "response.AuthVerification": {
            "type": "object",
            "properties": {
                "step": {
                    "$ref": "#/definitions/constants.AuthVerificationStep"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "response.BaseResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 999
                },
                "data": {},
                "message": {
                    "type": "string",
                    "example": "Message!"
                }
            }
        }
    },
    "securityDefinitions": {
        "AccessToken": {
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
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "BE Revend API",
	Description:      "This is a API documentation of BE Revend",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
