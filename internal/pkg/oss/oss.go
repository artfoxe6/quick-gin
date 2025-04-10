package oss

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/artfoxe6/quick-gin/internal/app/config"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	c *oss.Client
}

var client Client

func GetClient() Client {
	if client.c == nil {
		var err error
		client.c, err = oss.New(config.Oss.Endpoint, config.Oss.AccessKeyId, config.Oss.AccessKeySecret)
		if err != nil {
			log.Fatalln("oss init error: ", err.Error())
		}
	}
	return client
}

func (client Client) bucket() *oss.Bucket {
	bucket, err := client.c.Bucket(config.Oss.BucketName)
	if err != nil {
		log.Fatalln(errors.New("oss get bucket error"))
	}
	return bucket
}

func (client Client) Upload(filename string, reader io.Reader) string {
	if err := client.bucket().PutObject(filename, reader); err == nil {
		return config.Oss.CdnUrl + "/" + filename
	}
	return ""
}

func (client Client) UploadImgFromUrl(url string) string {
	if response, err := http.Get(url); err == nil {
		fileName := fmt.Sprintf("%s/%s/%d_%d%s", config.Oss.ImageDir, time.Now().Format("20060102"), time.Now().Unix(), response.ContentLength, ".jpg")
		defer response.Body.Close()
		client.bucket().PutObject(fileName, response.Body)
		return config.Oss.CdnUrl + fileName
	}
	return ""
}
