# go-tools-s3

Workload for general purposes

It upload a file from localmachine to S3. 

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
