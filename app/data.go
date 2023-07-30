package app

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
)

const (
	// MEASUREYOURLIFE_TABLE_NAME string = "measureyourlife-test"
	MEASUREYOURLIFE_TABLE_NAME            string = "measureyourlife"
	MEASUREYOURLIFE_TABLE_KEY             string = "Key"
	MEASUREYOURLIFE_TABLE_SORT_KEY        string = "SortKey"
	MEASUREYOURLIFE_TABLE_METRICS_ATTR    string = "metrics"
	MEASUREYOURLIFE_TABLE_UPDATED_AT_ATTR string = "udpatedAt"
)

type dayStatsItem struct {
	SortKey string
	Metrics []metricValue
}

func updateDayStats(userId string, date string, dayStats dayStatsData) error {
	// get service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return logAndConvertError(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// define keys
	hashKey := fmt.Sprintf("DAYSTATS#%s", userId)
	sortKey := date

	// encode data

	metrics, err := attributevalue.MarshalList(dayStats.MetricValues)
	if err != nil {
		return logAndConvertError(err)
	}

	// query input
	input := &dynamodb.PutItemInput{
		TableName: aws.String(MEASUREYOURLIFE_TABLE_NAME),
		Item: map[string]types.AttributeValue{
			MEASUREYOURLIFE_TABLE_KEY:          &types.AttributeValueMemberS{Value: hashKey},
			MEASUREYOURLIFE_TABLE_SORT_KEY:     &types.AttributeValueMemberS{Value: sortKey},
			MEASUREYOURLIFE_TABLE_METRICS_ATTR: &types.AttributeValueMemberL{Value: metrics},
		},
		ReturnValues: types.ReturnValueNone,
	}

	// run query
	_, err = svc.PutItem(context.TODO(), input)
	if err != nil {
		return logAndConvertError(err)
	}

	// done
	return nil
}

func getDayStats(userId string, date string) (*dayStatsData, error) {
	// get service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, logAndConvertError(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// define keys
	hashKey := fmt.Sprintf("DAYSTATS#%s", userId)
	sortKey := date

	// query expression
	projection := expression.NamesList(
		expression.Name(MEASUREYOURLIFE_TABLE_SORT_KEY),
		expression.Name(MEASUREYOURLIFE_TABLE_METRICS_ATTR))
	expr, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		return nil, logAndConvertError(err)
	}

	// query input
	input := &dynamodb.GetItemInput{
		TableName: aws.String(MEASUREYOURLIFE_TABLE_NAME),
		Key: map[string]types.AttributeValue{
			MEASUREYOURLIFE_TABLE_KEY:      &types.AttributeValueMemberS{Value: hashKey},
			MEASUREYOURLIFE_TABLE_SORT_KEY: &types.AttributeValueMemberS{Value: sortKey},
		},
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	}

	// run query
	result, err := svc.GetItem(context.TODO(), input)
	if err != nil {
		return nil, logAndConvertError(err)
	}

	// re-pack the results
	if result.Item == nil {
		return nil, nil
	}
	item := dayStatsItem{}
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, logAndConvertError(err)
	}
	dayStats := dayStatsData{
		MetricValues: item.Metrics,
	}

	return &dayStats, nil
}

func logAndConvertError(err error) error {
	log.Printf("%v", err)
	return fmt.Errorf("service unavailable")
}
