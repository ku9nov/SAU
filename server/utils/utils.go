package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	Database DatabaseSetting
	Server   ServerSettings
}
type DatabaseSetting struct {
	Url        string
	DbName     string
	Collection string
}

type ServerSettings struct {
	Port string
}

func IsValidInputAppName(input string) bool {
	// Only allow letters and numbers, no spaces or special characters
	validName := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return validName.MatchString(input)
}
func IsValidInputVersion(input string) bool {
	// Only allow letters and numbers, no spaces or special characters
	validVersion := regexp.MustCompile(`^[0-9.-]+$`)
	return validVersion.MatchString(input)
}

func createS3Client() *s3.Client {
	env := viper.GetViper()
	// Manually set the AWS credentials and region
	creds := credentials.NewStaticCredentialsProvider(env.GetString("S3_ACCESS_KEY"), env.GetString("S3_SECRET_KEY"), "")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(env.GetString("S3_REGION")))
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}

	// Create an S3 client using the AWS config
	return s3.NewFromConfig(cfg)
}

func UploadToS3(appName, version string, file *multipart.FileHeader, c *gin.Context, env *viper.Viper) string {
	// // Create an S3 client using another func
	s3Client := createS3Client()

	var extension string
	// Extract base filename and extension
	baseFileName := file.Filename
	dotIndex := strings.Index(baseFileName, ".")
	if dotIndex > -1 {
		extension = baseFileName[dotIndex:]
	}
	// Generate new file name
	newFileName := fmt.Sprintf("%s-%s%s", appName, version, extension)

	link := fmt.Sprintf("%s/%s/%s", env.GetString("S3_ENDPOINT"), appName, newFileName)

	// Open the file for reading
	fileReader, err := file.Open()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file for reading"})
	}
	s3Key := appName + "/" + newFileName

	// Upload file to S3
	_, err = s3Client.PutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket: aws.String(env.GetString("S3_BUCKET_NAME")),
		Key:    aws.String(s3Key),
		Body:   fileReader,
	})
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file to S3"})

	}
	return link
}

func DeleteFromS3(objectKey string, c *gin.Context, env *viper.Viper) {

	s3Client := createS3Client()

	// Delete object from bucket
	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(env.GetString("S3_BUCKET_NAME")),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file from S3"})
	}

	fmt.Printf("Object '%s' deleted from bucket '%s'\n", objectKey, env.GetString("S3_BUCKET_NAME"))
}
