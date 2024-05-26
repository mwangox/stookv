package api

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"stoo-kv/config"
	"stoo-kv/internal/crypto"
	"stoo-kv/internal/store"
)

type Handler struct {
	storage store.Store
	config  *config.Config
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewHandler(storage store.Store, config *config.Config) *Handler {
	return &Handler{
		config:  config,
		storage: storage,
	}
}
func (h Handler) GetHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	profile := c.Param("profile")
	key := c.Param("key")
	value, err := h.storage.Get(fmt.Sprintf("%s::%s::%s", namespace, profile, key))
	if err != nil {
		log.Printf("Failed to read key from storage: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}
	if value == "" {
		HandleError(c, StatusNotFound, "Data not found from storage")
		return
	}

	value, err = CheckEncryption(value, h.config)
	if err != nil {
		log.Printf("Failed to decrypt the value: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}
	HandleSuccess(c, value)
}

func (h Handler) GetByNamespaceAndProfileHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	profile := c.Param("profile")
	values, err := h.storage.GetByNameSpaceAndProfile(namespace, profile)
	h.valuesProcessor(c, values, err)
}

//
//func (h Handler) GetAllHandler(c *gin.Context) {
//	values, err := h.storage.GetAll()
//	h.valuesProcessor(c, values, err)
//}

func (h Handler) SetHandler(c *gin.Context) {
	h.set(c, false)
}

func (h Handler) SetSecretHandler(c *gin.Context) {
	h.set(c, true)
}

func (h Handler) set(c *gin.Context, isSecret bool) {
	namespace := c.Param("namespace")
	profile := c.Param("profile")
	kv := &KV{}
	if err := c.ShouldBindJSON(kv); err != nil {
		log.Printf("Failed to decode data: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}

	value := kv.Value
	key := kv.Key
	if key == "" {
		log.Printf("Key cannot be empty")
		HandleGeneralError(c, "Key name cannot be empty")
		return
	}

	if isSecret {
		ciphertext, err := crypto.Encrypt([]byte(value), h.config.Application.EncryptKey)
		if err != nil {
			log.Printf("Failed to encrypt data: %v", err)
			HandleGeneralError(c, err.Error())
			return
		}

		encPrefix := h.config.Application.EncryptPrefix
		if encPrefix == "" {
			encPrefix = "{ENC} "
		}

		value = encPrefix + hex.EncodeToString(ciphertext)
	}

	if err := h.storage.Set(fmt.Sprintf("%s::%s::%s", namespace, profile, key), value); err != nil {
		log.Printf("Failed to store data into storage: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}
	HandleSuccess(c, "Key set successfully")
}

func (h Handler) DeleteHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	profile := c.Param("profile")
	key := c.Request.URL.Query().Get("key")
	if err := h.storage.Delete(fmt.Sprintf("%s::%s::%s", namespace, profile, key)); err != nil {
		log.Printf("Failed to remove data from storage: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}
	HandleSuccess(c, "Key removed successfully")
}

func (h Handler) EncryptHandler(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		log.Printf("Failed to parse data: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}

	ciphertext, err := crypto.Encrypt(data, h.config.Application.EncryptKey)
	if err != nil {
		log.Printf("Failed to encrypt data: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}
	c.String(http.StatusOK, hex.EncodeToString(ciphertext))
}

func (h Handler) DecryptHandler(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		log.Printf("Failed to parse data: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}

	plaintext, err := crypto.Decrypt(data, h.config.Application.EncryptKey)
	if err != nil {
		log.Printf("Failed to decrypt data: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}
	c.String(http.StatusOK, string(plaintext))
}

func (h Handler) valuesProcessor(c *gin.Context, values map[string]string, err error) {
	if err != nil {
		log.Printf("Failed to read keys from storage: %v", err)
		HandleGeneralError(c, err.Error())
		return
	}
	if len(values) == 0 {
		log.Printf("Keys not found from storage")
		HandleError(c, StatusNotFound, "Keys not found from storage")
		return
	}

	HandleSuccess(c, ParseValues(values, h.config))
}
