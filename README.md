# go-tools-s3

Workload for general purposes

/upload   => Calc the sha256sum, set a tag and upload if to S3. 
/download => Download the file, get the sha256sum from tag and calc the the sha256sun from the file 

## How to use

Set the below variables

+ export AWS_REGION=
+ export AWS_ACCESS_KEY_ID=
+ export AWS_SECRET_ACCESS_KEY=
+ export AWS_SESSION_TOKEN=

or load the info from a yaml file

config.yaml

    AWS_REGION=
    AWS_ACCESS_KEY_ID=
    AWS_SECRET_ACCESS_KEY=

## How to use

        go run . --bucket_name 908671954593-eliezer-my-bucket-test --file_name file01.txt

## Check

    cat file01.txt | sha256sum
    946040acf82ba0547e32b167f0efc206eae99234a350ebc94d9829b88cc8a787  -

## Docto

https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#Client.GetObjectTagging