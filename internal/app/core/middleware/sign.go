package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/artfoxe6/quick-gin/internal/pkg/kit"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var excludes = []string{
	"/v1/file",
	"/ping",
}

func Sign(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/v1") {
			c.Next()
			return
		}
		if kit.InArray(c.Request.URL.Path, excludes) {
			c.Next()
			return
		}
		var (
			ignoreSignature  = c.GetHeader("IgnoreSignature")
			requestSignature = c.GetHeader("Signature")
			requestTimestamp = c.GetHeader("Timestamp")
			contentType      = c.GetHeader("Content-TypeId")
			signData         string
			params           = make(map[string]string)
		)
		if ignoreSignature != "" {
			c.Next()
			return
		}
		if requestSignature == "" {
			c.Next()
			return
		}

		if requestTimestamp == "" || requestSignature == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "invalid signature"})
			return
		}
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				switch contentType {
				case "application/x-www-form-urlencoded", "multipart/form-data":
					bodyString := string(bodyBytes)
					bodyParams := strings.Split(bodyString, "&")
					for _, param := range bodyParams {
						parts := strings.Split(param, "=")
						if len(parts) == 2 {
							params[parts[0]] = parts[1]
						}
					}
					signData = ParseSignData(params)
				default:
					signData = string(bodyBytes)
					re := regexp.MustCompile(`[\s\n\r]+`)
					signData = re.ReplaceAllString(signData, "")
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}
		if c.Request.Method == http.MethodGet {
			for key, values := range c.Request.URL.Query() {
				if len(values) > 0 {
					params[key] = values[0]
				}
			}
			signData = ParseSignData(params)
		}

		signData += fmt.Sprintf("%s%s%s", c.Request.URL.Path, key, requestTimestamp)
		hash := sha256.Sum256([]byte(signData))
		signature := hex.EncodeToString(hash[:])

		if signature != requestSignature {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "invalid signature"})
			return
		}
		requestTimestampInt, err := strconv.Atoi(requestTimestamp)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "signature expired"})
			return
		}
		if time.Now().Unix()-int64(requestTimestampInt) > 86400 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "signature expired"})
			return
		}
		c.Next()
	}
}

func ParseSignData(params map[string]string) string {
	signStrings := []string{}
	sortedKeys := make([]string, 0, len(params))
	for key := range params {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		signStrings = append(signStrings, fmt.Sprintf("%s=%s", key, params[key]))
	}
	return strings.Join(signStrings, "&")
}

func GenerateSignatureGet(path string, params map[string]any, timestamp int64, secretKey string) string {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var signStrings []string
	for _, key := range keys {
		value := fmt.Sprintf("%v", params[key])
		signStrings = append(signStrings, fmt.Sprintf("%s=%s", key, value))
	}
	signData := strings.Join(signStrings, "&")

	signData += fmt.Sprintf("%s%s%d", path, secretKey, timestamp)
	hash := sha256.Sum256([]byte(signData))
	signature := hex.EncodeToString(hash[:])
	fmt.Println(signature, timestamp)
	return signature
}

func GenerateSignaturePost(path string, signData string, timestamp int64, secretKey string) string {
	re := regexp.MustCompile(`[\s\n\r]+`)
	signData = re.ReplaceAllString(signData, "")
	signData += fmt.Sprintf("%s%s%d", path, secretKey, timestamp)
	hash := sha256.Sum256([]byte(signData))
	signature := hex.EncodeToString(hash[:])
	fmt.Println(signature, timestamp)
	return signature
}
