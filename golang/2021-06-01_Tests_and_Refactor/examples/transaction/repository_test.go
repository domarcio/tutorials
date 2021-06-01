package transaction

import (
	"errors"
	"golang/2021-06-01_Tests_and_Refactor/examples/transaction/mock"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/shopspring/decimal"
)

func TestGetTransaction(t *testing.T) {
	cases := []struct {
		name                    string
		transactionID           string
		customerID              string
		expectedTransaction     *Transanction
		expectedError           error
		mockDynamoDBGetItemMock func() (*dynamodb.GetItemOutput, error)
	}{
		{
			name:                "an aws error",
			transactionID:       "foo",
			customerID:          "bar",
			expectedTransaction: nil,
			expectedError:       errors.New("aws error"),
			mockDynamoDBGetItemMock: func() (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("aws error")
			},
		},
		{
			name:                "transaction not found",
			transactionID:       "foo",
			customerID:          "bar",
			expectedTransaction: nil,
			expectedError:       ErrTransactionNotFound,
			mockDynamoDBGetItemMock: func() (*dynamodb.GetItemOutput, error) {
				out := &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{},
				}
				return out, nil
			},
		},
		{
			name:                "successful",
			transactionID:       "foo",
			customerID:          "bar",
			expectedTransaction: &Transanction{ID: "foo", CustomerID: "bar", Amount: Currency{decimal.NewFromFloat(10.5)}},
			expectedError:       nil,
			mockDynamoDBGetItemMock: func() (*dynamodb.GetItemOutput, error) {
				out := &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						"ID":         {S: aws.String("foo")},
						"CustomerID": {S: aws.String("bar")},
						"Amount":     {N: aws.String("10.5")},
					},
				}
				return out, nil
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := &mock.DynamoDBMock{
				GetItemMock: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
					return c.mockDynamoDBGetItemMock()
				},
			}
			repo := NewRepository(mock)
			transaction, err := repo.GetTransaction(c.transactionID, c.customerID)

			if c.expectedError != nil && c.expectedError.Error() != err.Error() {
				t.Errorf("expected error %s, got %s", c.expectedError.Error(), err.Error())
			} else if c.expectedTransaction != nil && !reflect.DeepEqual(c.expectedTransaction, transaction) {
				t.Errorf("expected transaction %+v, got %+v", c.expectedTransaction, transaction)
			}
		})
	}
}
