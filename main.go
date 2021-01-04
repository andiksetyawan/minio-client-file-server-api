package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main()  {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY_ID")
	useSSL,_ := strconv.ParseBool(os.Getenv("MINIO_SECURE"))

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/{bucket}/{id}", func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		object, err := minioClient.GetObject(ctx, vars["bucket"], vars["id"], minio.GetObjectOptions{})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		b, err := ioutil.ReadAll(object)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)
	})

	log.Println("server started at :8080")
	panic(http.ListenAndServe(":8080", r))
}

