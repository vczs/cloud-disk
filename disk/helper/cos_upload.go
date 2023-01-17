package helper

import (
	"bytes"
	"cloud-disk/disk/internal/config"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func CosUploadFile(r *http.Request) (string, error) {
	url, _ := url.Parse(config.CosUrl)
	baseUrl := &cos.BaseURL{BucketURL: url}
	client := cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.TencentSecretId,
			SecretKey: config.TencentSecretKey,
		},
	})

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	name := config.CosFolderPath + "/" + GetUid() + path.Ext(fileHeader.Filename)

	_, err = client.Object.Put(context.Background(), name, file, nil)
	if err != nil {
		return "", err
	}

	return config.CosUrl + "/" + name, nil
}

/////////////////////////////////////////////////分片上传////////////////////////////////////////////////////////////////////////
func CosInitPart(ext string) (string, string, error) {
	url, _ := url.Parse(config.CosUrl)
	baseUrl := &cos.BaseURL{BucketURL: url}
	client := cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.TencentSecretId,
			SecretKey: config.TencentSecretKey,
		},
	})
	key := config.CosFolderPath + "/" + GetUid() + ext
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		return "", "", nil
	}
	return key, v.UploadID, nil
}

func CosUploadPart(r *http.Request, key, uploadId string, number int) (string, error) {
	url, _ := url.Parse(config.CosUrl)
	baseUrl := &cos.BaseURL{BucketURL: url}
	client := cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.TencentSecretId,
			SecretKey: config.TencentSecretKey,
		},
	})

	f, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	resp, err := client.Object.UploadPart(
		context.Background(), key, uploadId, number, bytes.NewReader(buf.Bytes()), nil,
	)
	if err != nil {
		panic(err)
	}

	return strings.Trim(resp.Header.Get("ETag"), "\""), nil
}

func CosCompletePart(key, uploadId string, c []cos.Object) error {
	url, _ := url.Parse(config.CosUrl)
	baseUrl := &cos.BaseURL{BucketURL: url}
	client := cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.TencentSecretId,
			SecretKey: config.TencentSecretKey,
		},
	})

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, c...)

	_, _, err := client.Object.CompleteMultipartUpload(context.Background(), key, uploadId, opt)
	if err != nil {
		return err
	}

	return nil
}
