package main

import(
	"fmt"
	"flag"
	"os"
	"io"
	"bytes"
	"crypto/md5"
	"encoding/hex"

	"github.com/spf13/viper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

var aws_region = ""
var aws_access_id = ""
var aws_acess_token = ""

// If use env variables from a yaml file
func getEnvVar(key string) string {
	fmt.Printf("Loading enviroment variable %s .. \n", key)
	
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file %s \n", err)
		os.Exit(1)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		fmt.Printf("Invalid type \n")
		os.Exit(1)
	}
	return value
}

func main(){
	fmt.Println("Starting S3 loading tool")

	bucket_name		:= flag.String("bucket_name","","")
	file_name		:= flag.String("file_name","1","")

	flag.Parse()
	fmt.Printf("bucket_name: %s file_name: %s  \n", *bucket_name, *file_name)

	aws_region 			:= os.Getenv("AWS_REGION")

	// Use the yaml file as info source
	//aws_region 			:= os.Getenv("AWS_REGION")  //getEnvVar("AWS_REGION")
	//aws_access_id 		:= os.Getenv("AWS_ACCESS_KEY_ID") //getEnvVar("AWS_ACCESS_KEY_ID")
	//aws_access_secret 	:= os.Getenv("AWS_SECRET_ACCESS_KEY") //getEnvVar("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(aws_region),
		//Credentials: credentials.NewStaticCredentials( aws_access_id , aws_access_secret , ""), // Use the yaml file as info source
	},)
	if err != nil {
		fmt.Println("Erro Create aws Session: ",err.Error())
		os.Exit(1)
	}

	file, err := os.Open(*file_name)
	if err != nil {
		fmt.Println("Erro open file: ",err.Error())
		os.Exit(1)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		panic(err)
	}

	var s3_tag = "md5=" + hex.EncodeToString(hash.Sum(nil))

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{Bucket:               aws.String(*bucket_name),
													Key:                  aws.String(*file_name),
													ACL:                  aws.String("private"),
													Body:                 bytes.NewReader(buffer),
													ContentLength:        aws.Int64(size),
													ContentDisposition:   aws.String("attachment"),
													Tagging:			  aws.String(s3_tag),
													ServerSideEncryption: aws.String("AES256"),})
	if err != nil {
		fmt.Println("Erro upload file: ",err.Error())
		os.Exit(1)
	}

	fmt.Println("Done...")
}