package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY_ID")
	useSSL, _ := strconv.ParseBool(os.Getenv("MINIO_SECURE"))

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")
		bucket := paths[1]
		fileName := paths[len(paths)-1]
		fileFullPath := strings.Join(paths[2:], "/")

		object, err := minioClient.GetObject(context.Background(), bucket, fileFullPath, minio.GetObjectOptions{})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer object.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		_, err = io.Copy(w, object)
		if err != nil {
			log.Println(err)
			w.Header().Del("Content-Disposition")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})

	log.Println("server started at :8080")
	panic(http.ListenAndServe(":8080", nil))
}
