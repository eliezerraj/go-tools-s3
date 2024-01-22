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
	log.Println("Starting S3 upload file")

	bucket_name		:= flag.String("bucket_name","","")
	file_name		:= flag.String("file_name","1","")

	flag.Parse()
	log.Printf("bucket_name: %s file_name: %s \n", *bucket_name, *file_name)

	// Open File
	file, err := os.Open(*file_name)
	if err != nil {
		log.Println("Erro open file: ", err.Error())
		panic(err)
	}
	defer file.Close()

	// Calc sha256sum
	hash := sha256.New()
	_, err = io.Copy(hash, file)	
	if err != nil {
		log.Println("Hash error %v /n", err)
		panic(err)
	}

	var s3_tag = "sha256sum=" + hex.EncodeToString(hash.Sum(nil))
	log.Printf("%s \n",s3_tag)

	//Prepare to log in S3
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Couldn't load default configuration %v", err)
		panic(err)
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	// Set back the file cursor in file 1, the file cursor was changed during the sha256 calc
	file.Seek(0,io.SeekStart)
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(*bucket_name),
		Key:    aws.String(*file_name),
		Body:   file,
		Tagging: aws.String(s3_tag),
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",file, bucket_name, file_name, err)
		panic(err)
	}

	log.Println("Done...")
}