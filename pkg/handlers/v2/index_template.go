package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/template"
)

func ListIndexTemplate(c *gin.Context) {
	pattern := c.Query("pattern")
	templates, err := core.ListTemplates(pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

func UpdateIndexTemplate(c *gin.Context) {
	name := c.Param("target")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template.name should be not empty"})
		return
	}

	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := template.Request(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = core.NewTemplate(name, template)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "template " + name + " created",
	})
}

func GetIndexTemplate(c *gin.Context) {
	name := c.Param("target")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template.name should be not empty"})
		return
	}
	template, exists, err := core.LoadTemplate(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "template " + name + " does not exists"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func DeleteIndexTemplate(c *gin.Context) {
	name := c.Param("target")
	err := core.DeleteTemplate(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
