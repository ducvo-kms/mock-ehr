package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"kms-connect.com/mock-ehr/pkg/domain/fhir/r4"
	"kms-connect.com/mock-ehr/pkg/domain/qa"
)

var CONFIGURATION = make(map[string]*qa.Config)

func CernerController(engine *gin.Engine) {
	auth := engine.Group("/cerner/authentication")

	auth.POST("/tenants/:tenant/protocols/oauth2/profiles/smart-v1/token", authentication)

	cerner := engine.Group("/cerner/r4/:tenant")
	// Get
	cerner.GET("/:resource/:id", getResourceById)
	cerner.POST("/:resource/:id/configuration", configurationGetResourceById)
	cerner.DELETE("/:resource/:id/configuration", deleteconfigurationGetResourceById)

	// Search
	cerner.GET("/:resource", searchResource)
	cerner.POST("/:resource/search/configuration", configurationSearchResource)
	cerner.DELETE("/:resource/search/configuration", deleteconfigurationSearchResource)

	//Create
	cerner.POST("/:resource", createResource)
	cerner.POST("/:resource/create/configuration", configurationCreateResource)
	cerner.DELETE("/:resource/create/configuration", deleteconfigurationCreateResource)

	//Get all configuration
	cerner.GET("/configuration", func(context *gin.Context) {
		context.JSON(http.StatusOK, CONFIGURATION)
	})
}

func authentication(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"access_token": "FAKE_ACCESS_TOKEN",
		"scope":        "FAKE_SCOPE",
		"token_type":   "Bearer",
		"expires_in":   570,
	})
}

func searchResource(context *gin.Context) {
	resource := context.Param("resource")

	data := r4.Bundle{
		ResourceType: "Bundle",
		Type:         "searchset",
		Entry:        []r4.Entry{{Resource: r4.Resource{Id: "1", ResourceType: resource}}},
	}

	key := keySearchConfiguration(context)
	config := CONFIGURATION[key]

	if config == nil {
		context.JSON(http.StatusOK, data)
	} else {
		time.Sleep(time.Duration(config.WaitIn * 1_000_000))
		context.JSON(config.StatusCode, data)
	}
}

func createResource(context *gin.Context) {
	resource := context.Param("resource")
	data := r4.Resource{Id: "1", ResourceType: resource}

	key := keyCreateConfiguration(context)
	config := CONFIGURATION[key]

	if config == nil {
		context.Header("Location", resource+"/1")
		context.JSON(http.StatusCreated, data)
	} else {
		time.Sleep(time.Duration(config.WaitIn * 1_000_000))
		context.Header("Location", resource+"/1")
		context.JSON(config.StatusCode, data)
	}
}

func getResourceById(context *gin.Context) {
	id := context.Param("id")
	resource := context.Param("resource")
	data := r4.Resource{Id: id, ResourceType: resource}
	key := keyGetConfiguration(context)
	config := CONFIGURATION[key]

	if config == nil {
		context.JSON(http.StatusOK, data)
	} else {
		time.Sleep(time.Duration(config.WaitIn * 1_000_000))
		context.JSON(config.StatusCode, data)
	}

}

func configurationGetResourceById(context *gin.Context) {
	key := keyGetConfiguration(context)
	var body qa.Config

	if err := context.BindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "only allow configuration wait_in: int64 and status_code: int32"})
		return
	}

	CONFIGURATION[key] = &body
	context.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func configurationCreateResource(context *gin.Context) {
	key := keyCreateConfiguration(context)
	var body qa.Config

	if err := context.BindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "only allow configuration wait_in: int64 and status_code: int32"})
		return
	}

	CONFIGURATION[key] = &body
	context.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func configurationSearchResource(context *gin.Context) {
	key := keySearchConfiguration(context)
	var body qa.Config

	if err := context.BindJSON(&body); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "only allow configuration wait_in: int64 and status_code: int32"})
		return
	}

	CONFIGURATION[key] = &body
	context.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func deleteconfigurationGetResourceById(context *gin.Context) {
	delete(CONFIGURATION, keyGetConfiguration(context))

	context.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func deleteconfigurationCreateResource(context *gin.Context) {
	delete(CONFIGURATION, keyCreateConfiguration(context))

	context.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func deleteconfigurationSearchResource(context *gin.Context) {
	delete(CONFIGURATION, keySearchConfiguration(context))

	context.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func keyGetConfiguration(context *gin.Context) string {
	id := context.Param("id")
	resource := context.Param("resource")
	key := "get/" + resource + "/" + id
	return key
}

func keyCreateConfiguration(context *gin.Context) string {
	resource := context.Param("resource")
	key := "create/" + resource
	return key
}

func keySearchConfiguration(context *gin.Context) string {
	resource := context.Param("resource")
	key := "search/" + resource
	return key
}
