package transaction

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/shopspring/decimal"
)

type Currency struct {
	decimal.Decimal
}

func (c *Currency) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	d, err := decimal.NewFromString(*av.N)
	if err != nil {
		return err
	}

	*c = Currency{
		d,
	}

	return nil
}

type Transanction struct {
	ID         string
	CustomerID string
	Amount     Currency
}
