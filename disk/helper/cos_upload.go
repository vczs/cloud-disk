package helper

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"cloud-disk/disk/define"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func CosUploadFile(r *http.Request) (string, error) {
	url, _ := url.Parse(define.CosUrl)
	baseUrl := &cos.BaseURL{BucketURL: url}
	client := cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	name := define.CosFolderPath + "/" + GetUid() + path.Ext(fileHeader.Filename)

	_, err = client.Object.Put(context.Background(), name, file, nil)
	if err != nil {
		return "", err
	}

	return define.CosUrl + "/" + name, nil
}

// ///////////////////////////////////////////////分片上传////////////////////////////////////////////////////////////////////////
func CosInitPart(ext string) (string, string, error) {
	url, _ := url.Parse(define.CosUrl)
	baseUrl := &cos.BaseURL{BucketURL: url}
	client := cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})
	key := define.CosFolderPath + "/" + GetUid() + ext
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		return "", "", nil
	}
	return key, v.UploadID, nil
}

func CosUploadPart(r *http.Request, key, uploadId string, number int) (string, error) {
	url, _ := url.Parse(define.CosUrl)
	baseUrl := &cos.BaseURL{BucketURL: url}
	client := cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
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
	url, _ := url.Parse(define.CosUrl)
	baseUrl := &cos.BaseURL{BucketURL: url}
	client := cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
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
