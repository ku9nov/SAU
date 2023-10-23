package info

import (
	db "SAU/mongod"
	"SAU/server/utils"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindLatestVersion(c *gin.Context, repository db.AppRepository, db *mongo.Database) {
	_, err := utils.ValidateParams(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, ctxErr := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer ctxErr()

	// Request on repository
	checkResult, err := repository.CheckLatestVersion(c.Query("app_name"), c.Query("version"), c.Query("channel"), c.Query("platform"), c.Query("arch"), ctx)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !checkResult.Found {
		if len(checkResult.Artifacts) == 0 {
			c.JSON(http.StatusOK, gin.H{"update_available": false, "error": "Not found"})
		} else {
			errorMsg := "Not found"
			if err != nil {
				errorMsg = err.Error()
			}
			c.JSON(http.StatusOK, gin.H{"update_available": false, "error": errorMsg})
		}
		return
	}

	response := gin.H{"update_available": true}
	for _, artifact := range checkResult.Artifacts {
		if artifact.Package != "" && artifact.Link != "" {
			key := "update_url_" + strings.TrimPrefix(artifact.Package, ".")
			response[key] = artifact.Link
		}
	}

	c.JSON(http.StatusOK, response)
}
