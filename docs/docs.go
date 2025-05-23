// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Muhammed Ibrahim",
            "email": "aminmuhammad18@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Logs a user in and returns a JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Authenticate user",
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.AuthUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Create a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "UserRegister info",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UserRegister"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    }
                }
            }
        },
        "/todos": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieve all todos for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todos"
                ],
                "summary": "Get all todos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Todo"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create a new todo item for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todos"
                ],
                "summary": "Create a new todo",
                "parameters": [
                    {
                        "description": "Todo details",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TodoCreate"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Todo"
                        }
                    }
                }
            }
        },
        "/todos/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get a todo item by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todos"
                ],
                "summary": "Get a single todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Todo ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Todo"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update a todo item by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todos"
                ],
                "summary": "Update a todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Todo ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated todo data",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TodoUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Todo"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete a todo item by ID",
                "tags": [
                    "todos"
                ],
                "summary": "Delete a todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Todo ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        }
    },
    "definitions": {
        "errors.APIError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "details": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "model.AuthUser": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                }
            }
        },
        "model.Todo": {
            "description": "Todo represents a task that a user wants to track",
            "type": "object",
            "required": [
                "title"
            ],
            "properties": {
                "completed": {
                    "type": "boolean",
                    "example": false
                },
                "createdAt": {
                    "type": "string",
                    "example": "2022-01-01T12:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "5f8d0614db5c5c7b3a18f201"
                },
                "title": {
                    "type": "string",
                    "example": "Buy groceries"
                },
                "updatedAt": {
                    "type": "string",
                    "example": "2022-01-01T12:00:00Z"
                },
                "userId": {
                    "type": "string",
                    "example": "5f8d0614db5c5c7b3a18f200"
                }
            }
        },
        "model.TodoCreate": {
            "description": "TodoCreate is used when creating a new todo item",
            "type": "object",
            "required": [
                "title"
            ],
            "properties": {
                "completed": {
                    "type": "boolean",
                    "example": false
                },
                "title": {
                    "type": "string",
                    "example": "Buy groceries"
                }
            }
        },
        "model.TodoUpdate": {
            "description": "TodoUpdate is used when updating an existing todo item",
            "type": "object",
            "properties": {
                "completed": {
                    "type": "boolean",
                    "example": true
                },
                "description": {
                    "type": "string",
                    "example": "Need to get milk and eggs"
                },
                "title": {
                    "type": "string",
                    "example": "Buy more groceries"
                },
                "updatedAt": {
                    "type": "string",
                    "example": "2022-01-02T12:00:00Z"
                }
            }
        },
        "model.User": {
            "type": "object",
            "required": [
                "email",
                "fullName",
                "password"
            ],
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 3
                },
                "id": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "model.UserRegister": {
            "type": "object",
            "required": [
                "email",
                "fullName",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 3
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
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
	Schemes:          []string{"http", "https"},
	Title:            "Todo API",
	Description:      "This is a simple todo app backend API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
