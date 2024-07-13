package delete

import (
	"context"
	db "faynoSync/mongod"
	"faynoSync/server/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func DeleteApp(c *gin.Context, repository db.AppRepository) {
	env := viper.GetViper()
	ctx, ctxErr := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer ctxErr()

	// Convert string to ObjectID
	objID, err := primitive.ObjectIDFromHex(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//request on repository
	links, result, err := repository.DeleteApp(objID, ctx)
	if err != nil {
		logrus.Error(err)
	}

	for _, link := range links {
		subLink := strings.TrimPrefix(link, env.GetString("S3_ENDPOINT"))
		utils.DeleteFromS3(subLink, c, viper.GetViper())
	}
	c.JSON(http.StatusOK, gin.H{"deleteAppResult.DeletedCount": result})
}

func DeleteChannel(c *gin.Context, repository db.AppRepository) {
	deleteEntity(c, repository, "channel")
}

func DeleteArch(c *gin.Context, repository db.AppRepository) {
	deleteEntity(c, repository, "arch")
}

func DeletePlatform(c *gin.Context, repository db.AppRepository) {
	deleteEntity(c, repository, "platform")
}

func deleteEntity(c *gin.Context, repository db.AppRepository, itemType string) {
	ctx, ctxErr := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer ctxErr()

	// Convert string to ObjectID
	objID, err := primitive.ObjectIDFromHex(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var result interface{}
	// var err error
	switch itemType {
	case "channel":
		result, err = repository.DeleteChannel(objID, ctx)
	case "platform":
		result, err = repository.DeletePlatform(objID, ctx)
	case "arch":
		result, err = repository.DeleteArch(objID, ctx)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type"})
		return
	}
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete " + itemType})
		return
	}
	var tag language.Tag
	titleCase := cases.Title(tag)

	capitalizedItemType := titleCase.String(itemType)
	c.JSON(http.StatusOK, gin.H{"delete" + capitalizedItemType + "Result.DeletedCount": result})
}
