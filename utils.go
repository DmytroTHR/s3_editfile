package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getAwsClient(credsPath string) (*s3.Client, error) {
	var cfg aws.Config
	var err error
	if _, errStat := os.Stat(credsPath); errStat == nil {
		log.Printf("trying configs from %s directory\n", credsPath)
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithSharedCredentialsFiles([]string{credsPath + "credentials"}),
			config.WithSharedConfigFiles([]string{credsPath + "config"}))
	} else {
		log.Println("trying default configs")
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		return &s3.Client{}, err
	}

	return s3.NewFromConfig(cfg), nil
}

func getBucketObjectsList(client *s3.Client, bucketName *string) (*s3.ListObjectsV2Output, error) {
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: bucketName,
	})

	if err != nil {
		return &s3.ListObjectsV2Output{}, err
	}

	return output, nil
}

func getBucketObjectToFile(client *s3.Client, bucketName, objectKey *string) (*os.File, error) {
	theObject, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: bucketName,
		Key:    objectKey,
	})
	if err != nil {
		log.Printf("error getting object from file, %v", err)
		return &os.File{}, err
	}

	fileToBeChanged, err := os.Create(*objectKey)
	if err != nil {
		log.Printf("error creating local file, %v", err)
		return &os.File{}, err
	}

	_, err = io.Copy(fileToBeChanged, theObject.Body)
	if err != nil {
		log.Printf("error copying to local file, %v", err)
		return &os.File{}, err
	}

	return fileToBeChanged, nil
}

func putFileIntoBucketObject(client *s3.Client, fileName string) error {
	fileToBeChanged, err := os.Open(fileName)
	if err != nil {
		log.Printf("error reopening local file, %v", err)
		return err
	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: bucketName,
		Key:    requiredFileName,
		Body:   fileToBeChanged})
	if err != nil {
		log.Printf("error replacing S3-file, %v", err)
		return err
	}

	return nil
}

func deleteTmpFile(file *os.File) {
	file.Close()
	os.Remove(file.Name())
}
