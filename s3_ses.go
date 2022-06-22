// s3 → lambda → ses

package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

const Region = "ap-northeast-1"

var (
	From = "from@example.com"
	To1  = "to1@example.com"
	To2  = "to2@example.com"

	Subject = "From SES"

	ObjectKey string
)

func handler(c context.Context, e events.S3Event) {
	for _, record := range e.Records {
		ObjectKey = record.S3.Object.Key
	}

	cfg, err := config.LoadDefaultConfig(c, config.WithRegion(Region))

	if err != nil {
		fmt.Print("ERROR config === ", err)
		return
	}

	client := sesv2.NewFromConfig(cfg)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: &From,
		Destination: &types.Destination{
			ToAddresses: []string{To1, To2},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: &ObjectKey,
					},
				},
				Subject: &types.Content{
					Data: &Subject,
				},
			},
		},
	}

	res, err := client.SendEmail(c, input)

	if err != nil {
		fmt.Print("ERROR send mail === ", err)
		return
	}

	print("SUCCESS send main message id === ", res.MessageId)
}

func main() {
	lambda.Start(handler)
}
