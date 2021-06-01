package mock

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoDBMock struct {
	dynamodbiface.DynamoDBAPI

	GetItemMock func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

func (m *DynamoDBMock) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return m.GetItemMock(input)
}
