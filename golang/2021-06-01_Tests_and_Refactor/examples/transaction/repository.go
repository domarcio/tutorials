package transaction

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrTransactionNotFound error = errors.New("transaction does not exists")
)

type Repository struct {
	db dynamodbiface.DynamoDBAPI
}

func NewRepository(db dynamodbiface.DynamoDBAPI) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetTransaction(id, customerID string) (*Transanction, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
			"CustomerID": {
				S: aws.String(customerID),
			},
		},
		TableName: aws.String("Transaction"),
	}

	result, err := r.db.GetItem(input)
	if err != nil {
		return nil, err
	}

	if len(result.Item) <= 0 {
		return nil, ErrTransactionNotFound
	}

	transaction := &Transanction{}
	err = dynamodbattribute.UnmarshalMap(result.Item, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
