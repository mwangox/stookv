package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"stoo-kv/config"
	"stoo-kv/internal/crypto"
	"strings"
)

func HandleGeneralError(c *gin.Context, message string) {
	HandleError(c, StatusGeneralError, message)
}

func HandleError(c *gin.Context, statusCode int, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  statusCode,
		"message": message,
	})
}

func HandleSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusSuccess,
		"message": "Success",
		"data":    data,
	})
}
func CheckEncryption(value string, config *config.Config) (string, error) {
	encPrefix := config.EncryptPrefix
	if encPrefix == "" {
		encPrefix = "{ENC} "
	}
	encValue, status := strings.CutPrefix(value, encPrefix)
	if status {
		valueByte, err := crypto.Decrypt([]byte(encValue), config.EncryptKey)
		if err != nil {
			return "", err
		}
		value = string(valueByte)
	}
	return value, nil
}

func ParseValues(values map[string]string, config *config.Config) map[string]string {
	parsedValues := make(map[string]string)
	for k, v := range values {
		v, err := CheckEncryption(v, config)
		if err != nil {
			log.Printf("Failed to decrypt the value: %v", err)
			v = "****NOT VALID****"
		}
		parsedValues[k] = v
	}
	return parsedValues
}

const (
	StatusSuccess      = 0
	StatusGeneralError = -1
	StatusNotFound     = -2
)
