package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var credsPath, bucketName, requiredFileName, stringToAppend *string

func init() {
	credsPath = flag.String("cred", ".aws/", "Creds Path")
	bucketName = flag.String("b", "somebucket-123321", "Bucket name")
	requiredFileName = flag.String("f", "index.txt", "Filename to replace")
	stringToAppend = flag.String("a",
		fmt.Sprintf("%s - appended string", time.Now().Format(time.Stamp)),
		"String to append")	
}

func main() {
	flag.Parse()

	client, err := getAwsClient(*credsPath)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	bucketObjects, err := getBucketObjectsList(client, bucketName)
	if err != nil {
		log.Fatalf("failed to receive bucket %s objects list, %v", *bucketName, err)
	}

	log.Println("first page results:")
	for _, object := range bucketObjects.Contents {
		if *object.Key != *requiredFileName {
			continue
		}

		log.Printf("found file, key=%s size=%d", *requiredFileName, object.Size)
		fileToBeChanged, err := getBucketObjectToFile(client, bucketName, object.Key)
		if err != nil {
			log.Fatalf("error saving to local file, %v", err)
		}
		defer deleteTmpFile(fileToBeChanged)

		_, err = fmt.Fprintf(fileToBeChanged, "\n%s", *stringToAppend)
		if err != nil {
			log.Fatalf("error writing to local file, %v", err)
		}

		if putFileIntoBucketObject(client, *requiredFileName) != nil {
			log.Fatalf("error replacing S3-file, %v", err)
		}
		log.Printf("successfully replaced %v", *requiredFileName)
	}
}
