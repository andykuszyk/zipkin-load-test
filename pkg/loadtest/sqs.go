package loadtest

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/hashicorp/go-uuid"
	"log"
)

type SqsTest struct {
	SQSRegion   *string
	SQSEndpoint *string

	SQSService *sqs.SQS
	QueueName  string
	QueueURL   string
}

func (r *SqsTest) Setup() {
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region:   r.SQSRegion,
			Endpoint: r.SQSEndpoint,
		}))
	r.SQSService = sqs.New(sess)

	if r.QueueName == "" {
		randomId, _ := uuid.GenerateUUID()
		queueName := "sqs-load-test-" + randomId

		log.Printf("Using random queue: %s", queueName)
		result, err := r.SQSService.CreateQueue(&sqs.CreateQueueInput{
			QueueName: &queueName,
		})

		if err != nil {
			log.Fatal(err)
		}

		r.QueueURL = *result.QueueUrl
	} else {

		log.Printf("Using configured queue: %s", r.QueueName)
		result, err := r.SQSService.GetQueueUrl(&sqs.GetQueueUrlInput{
			QueueName: &r.QueueName,
		})

		if err != nil {
			log.Fatal(err)
		}

		r.QueueURL = *result.QueueUrl
	}
}

func (r *SqsTest) Run() {
	_, err := r.SQSService.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("a test message"),
		QueueUrl:    &r.QueueURL,
	})

	if err != nil {
		log.Fatal("Error", err)
	}

	var receive *sqs.ReceiveMessageOutput
	receive, err = r.SQSService.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &r.QueueURL,
		MaxNumberOfMessages: aws.Int64(10),
		//VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds: aws.Int64(20),
	})

	if err != nil {
		log.Fatal("Error", err)
		return
	}

	if len(receive.Messages) == 0 {
		log.Fatal("Zero messages received")
	}

}

func (r *SqsTest) Teardown() {

	if r.QueueName == "" {
		_, err := r.SQSService.DeleteQueue(&sqs.DeleteQueueInput{
			QueueUrl: &r.QueueURL,
		})

		if err != nil {
			log.Fatal("Error", err)
			return
		}
	}
}
