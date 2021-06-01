package main

import (
	"golang/2021-06-01_Tests_and_Refactor/examples/transaction"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	endpoint = "http://localhost:8000"
	region   = "us-east-1"
)

func main() {
	config := &aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Endpoint:                      aws.String(endpoint),
		Region:                        aws.String(region),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		log.Fatalln(err)
	}

	db := dynamodb.New(sess)
	repo := transaction.NewRepository(db)

	t, err := repo.GetTransaction("f4302dff-5e50-4d07-b201-8e675a583c2a", "d07b14ae-55dc-4f68-917f-cd9857994b16")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("ID: %s, Customer ID: %s, Amount %s", t.ID, t.CustomerID, t.Amount.StringFixed(2))
}
