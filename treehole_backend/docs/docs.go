// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/capacha/get": {
            "get": {
                "tags": [
                    "登录业务接口"
                ],
                "summary": "获取验证码",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/capacha/verify": {
            "get": {
                "tags": [
                    "登录业务接口"
                ],
                "summary": "验证验证码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "capachaId",
                        "name": "capachaId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "capachaVal",
                        "name": "capachaVal",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/forgetPassword/modifyPassword": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "登录业务接口"
                ],
                "summary": "密码修改",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.ModifyPasswordForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/forgetPassword/verifyEmailCode": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "登录业务接口"
                ],
                "summary": "忘记密码-验证验证码",
                "parameters": [
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.VerifyEmailCodeForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/getNoteInfo": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "公共接口"
                ],
                "summary": "获取帖子详细信息",
                "parameters": [
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.GetNoteInfoFrom"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/hello": {
            "get": {
                "tags": [
                    "公共接口"
                ],
                "summary": "首页",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "登录业务接口"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.LoginForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/note/createNote": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "帖子业务接口"
                ],
                "summary": "创建帖子",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.CreateNoteForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/note/deleteNote": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "帖子业务接口"
                ],
                "summary": "删除帖子",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.GetNoteInfoFrom"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/note/getNoteList": {
            "get": {
                "tags": [
                    "帖子业务接口"
                ],
                "summary": "获取发布帖子列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "size",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "登录业务接口"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.RegisterForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/searchNotes": {
            "get": {
                "tags": [
                    "公共接口"
                ],
                "summary": "搜索帖子",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "size",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "keyword",
                        "name": "keyword",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/searchNotesScoreOrTime": {
            "get": {
                "tags": [
                    "公共接口"
                ],
                "summary": "按照热度或时间获取帖子信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "size",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "type",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/sendEmailCode": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "登录业务接口"
                ],
                "summary": "发送邮件验证码",
                "parameters": [
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.SendCodeForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/uploadLocal": {
            "post": {
                "tags": [
                    "公共接口"
                ],
                "summary": "上传图片",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "文件",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "文件类型",
                        "name": "type",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/user/getUserInfo": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户业务接口"
                ],
                "summary": "获取用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/user/modifyUserInfo": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户业务接口"
                ],
                "summary": "修改用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.ModifyUserInfoForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        },
        "/user/userModifyPassword": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户业务接口"
                ],
                "summary": "更换密码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "发送参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.UserModifyPasswordForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.H"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "utils.CreateNoteForm": {
            "type": "object",
            "required": [
                "content",
                "title"
            ],
            "properties": {
                "content": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "urls": {
                    "type": "string"
                }
            }
        },
        "utils.GetNoteInfoFrom": {
            "type": "object",
            "required": [
                "note_identity"
            ],
            "properties": {
                "note_identity": {
                    "type": "string"
                }
            }
        },
        "utils.H": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "msg": {
                    "type": "string"
                }
            }
        },
        "utils.LoginForm": {
            "type": "object",
            "required": [
                "password",
                "repassword",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "repassword": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "utils.ModifyPasswordForm": {
            "type": "object",
            "required": [
                "password",
                "repassword",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "repassword": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "utils.ModifyUserInfoForm": {
            "type": "object",
            "required": [
                "age",
                "sex",
                "url",
                "username"
            ],
            "properties": {
                "age": {
                    "type": "string"
                },
                "sex": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "utils.RegisterForm": {
            "type": "object",
            "required": [
                "code",
                "email",
                "password",
                "repassword",
                "username"
            ],
            "properties": {
                "code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "repassword": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "utils.SendCodeForm": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "utils.UserModifyPasswordForm": {
            "type": "object",
            "required": [
                "nowpassword",
                "password",
                "repassword"
            ],
            "properties": {
                "nowpassword": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "repassword": {
                    "type": "string"
                }
            }
        },
        "utils.VerifyEmailCodeForm": {
            "type": "object",
            "required": [
                "code",
                "email"
            ],
            "properties": {
                "code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
