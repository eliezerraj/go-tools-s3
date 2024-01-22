package main

import (
	"context"
	"flag"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"log"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	log.Println("Starting S3 download file")

	bucket_name		:= flag.String("bucket_name","","")
	file_name		:= flag.String("file_name","1","")

	flag.Parse()
	log.Printf("bucket_name: %s file_name: %s \n", *bucket_name, *file_name)

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Couldn't load default configuration %v", err)
		panic(err)
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(*bucket_name),
		Key:    aws.String(*file_name),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", *bucket_name, *file_name, err)
		panic(err)
	}
	defer result.Body.Close()

	file, err := os.Create(*file_name)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", *file_name, err)
		panic(err)
	}
	defer file.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", *file_name, err)
		panic(err)
	}

	_, err = file.Write(body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", *file_name, err)
		panic(err)
	}

	file_read, err := os.Open(*file_name)
	if err != nil {
		log.Println("Erro open file: ", err.Error())
		panic(err)
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file_read)	
	if err != nil {
		log.Println("Hash error %v /n", err)
		panic(err)
	}

	resultTag, err := s3Client.GetObjectTagging(context.TODO(), &s3.GetObjectTaggingInput{
		Bucket: aws.String(*bucket_name),
		Key:    aws.String(*file_name),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", *bucket_name, *file_name, err)
		panic(err)
	}

	log.Printf("tag from s3: %v  %v \n", string(*resultTag.TagSet[0].Key), string(*resultTag.TagSet[0].Value))
	log.Printf("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	log.Printf("sha256sum from s3  : %v \n",string(*resultTag.TagSet[0].Value))
	log.Printf("sha256sum from file: %v \n", hex.EncodeToString(hash.Sum(nil)))

	log.Println("Done...")
}