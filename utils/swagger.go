package utils

import (
	"fmt"
	"strings"
)

type swaggerImpl struct {
}

type ISwagger interface {
	// PrintSummaryInfo prints Swagger summary information parameters
	//
	// Used to generate Swagger comments
	PrintSummaryInfo(param SwaggerSummaryInfoParam) string
	// PrintControllerInfo prints Swagger controller information
	//
	// Used to generate Swagger comments
	PrintControllerInfo(param SwaggerControllerInfoParam) string
}

var swi ISwagger = &swaggerImpl{}

// Swagger instantiates Swagger utilities
func Swagger() ISwagger {
	return swi
}

var (
	summaryInfoTpl = `
// ----------- api doc definition -----------------

// @title          %s
// @version        %s
// @description   %s
// @termsOfService %s

// @contact.name  %s
// @contact.url   %s
// @contact.email %s

// @license.name MIT
// @license.url  %s

// @Host     %s
// @BasePath %s
%s
// ----------- api doc definition -----------------`
	needAuthorizeTpl = `// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                Access Token protects our entity endpoints`

	controllerInfoTpl = `
// %s %s
// @Tags        %s
// @Summary     %s
// @Description %s
%s
// @Accept      application/json
// @Produce     json
%s
// @Param %s %s true "parameters"
// @Success     200   {object} %s
// @Router      %s [%s]`

	controllerInfoSecurityTpl      = `// @Security    ApiKeyAuth`
	controllerInfoAuthorizationTpl = `// @Param Authorization header string true "Authorization token"`
	controllerInfoQueryTpl         = `parameter query`
	controllerInfoPayloadTpl       = `payload body`
)

// SwaggerSummaryInfoParam summary information parameters
type SwaggerSummaryInfoParam struct {
	Title            string //Title
	Version          string //Version
	Description      string //Description
	TermOfServiceUrl string //Terms of service URL
	ContactName      string //Contact
	ContactUrl       string //Contact website
	ContactEmail     string //Contact email
	LicenseUrl       string //License URL
	Host             string //API host address
	BasePath         string //Common API base path
	NeedAuthorize    bool   //Whether token authorization is required
}

// PrintSummaryInfo prints Swagger summary information parameters
//
// Used to generate Swagger comments
func (swaggerImpl) PrintSummaryInfo(param SwaggerSummaryInfoParam) string {
	var needAuthorizeString string
	if param.NeedAuthorize {
		needAuthorizeString = needAuthorizeTpl
	}

	return fmt.Sprintf(summaryInfoTpl,
		param.Title, param.Version, param.Description, param.TermOfServiceUrl, param.ContactName, param.ContactUrl,
		param.ContactEmail, param.LicenseUrl, param.Host, param.BasePath, needAuthorizeString)
}

// SwaggerControllerInfoParam controller information parameters
type SwaggerControllerInfoParam struct {
	FunctionName       string //Method name
	FunctionDesc       string //Method description
	Tag                string //Category tags
	Summary            string //Summary title
	Description        string //API description
	Method             string //Request method
	RequestParamString string //Request parameter string
	ResponseBodyString string //Response body string
	ApiPath            string //API path
	NeedAuthorize      bool   //Whether authorization is required
}

// PrintControllerInfo prints Swagger controller information
//
// Used to generate Swagger comments
func (swaggerImpl) PrintControllerInfo(param SwaggerControllerInfoParam) string {
	var (
		requestString       string
		securityString      string
		needAuthorizeString string
	)
	if param.NeedAuthorize {
		securityString = controllerInfoSecurityTpl
		needAuthorizeString = controllerInfoAuthorizationTpl
	}
	if len(param.Method) == 0 {
		param.Method = "get"
	}
	param.Method = strings.ToLower(param.Method)
	if param.Method != "get" {
		requestString = controllerInfoPayloadTpl
	} else {
		requestString = controllerInfoQueryTpl
	}

	return fmt.Sprintf(controllerInfoTpl,
		param.FunctionName, param.FunctionDesc, param.Tag, param.Summary, param.Description, securityString,
		needAuthorizeString, requestString, param.RequestParamString, param.ResponseBodyString, param.ApiPath, param.Method)
}
