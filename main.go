package main

import(
	"fmt"
	"flag"
	"os"
	"bytes"

	"github.com/spf13/viper"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

var aws_region = ""
var aws_access_id = ""
var aws_acess_token = ""

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

	bucket_name 	:= flag.String("bucket_name","","")
	file_name 		:= flag.String("file_name","1","")

	flag.Parse()
	fmt.Printf("bucket_name: %s file_name: %s  \n", *bucket_name, *file_name)

	aws_region 		:= getEnvVar("AWS_REGION")
	aws_access_id 	:= getEnvVar("AWS_ACCESS_ID")
	aws_access_secret := getEnvVar("AWS_ACCESS_SECRET")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(aws_region),
		Credentials: credentials.NewStaticCredentials( aws_access_id , aws_access_secret , ""),},
	)
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

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(*bucket_name),
        Key:                  aws.String(*file_name),
        ACL:                  aws.String("private"),
        Body:                 bytes.NewReader(buffer),
        ContentLength:        aws.Int64(size),
        ContentDisposition:   aws.String("attachment"),
        ServerSideEncryption: aws.String("AES256"),
    })
	
	if err != nil {
		fmt.Println("Erro upload file: ",err.Error())
		os.Exit(1)
	}

	fmt.Println("Done...")
}