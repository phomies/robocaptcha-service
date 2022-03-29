package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
)

func sqsSendNotification(messageBody string) {
	sqsClient := sqs.New(sqsGetSession())
	queueUrl := os.Getenv("AWS_SQS_NOTIFICATION_URL")
	fmt.Println(sqsClient.Client.Endpoint)

	fmt.Printf("Sending notification")

	_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String(messageBody),
	})

	if err != nil {
		panic(err)
	}
}

func sqsGetSession() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-southeast-1"),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ID"),
				os.Getenv("AWS_SECRET"),
				"",
			),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}
