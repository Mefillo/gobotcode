package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var TableName = "BotData"

type Item struct {
	ID     string   `json:"id"`
	Films  []string `json:"films"`
	Status string   `json:"status"`
	Choice int      `json:"choice"`
}

func Save(item Item) (err error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new item: %s", err)
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return err
	}
	return nil
}

func Get(key string) (item Item, err error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// getting Item
	fmt.Printf("!!!START GETTING")
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(key)},
		},
	})
	fmt.Printf("!!!RESULT %+v", result)
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
		return item, err
	}

	if result.Item == nil {
		fmt.Print("Could not find, creating")
		item.ID = key
		err = Save(item)
		return
		// return item, merry.New(msg)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Printf("Failed to unmarshal Record, %v", err)
		return item, err
	}

	fmt.Printf("\n!!! GOT ITEM: %+v\n", item)

	return item, err
}
