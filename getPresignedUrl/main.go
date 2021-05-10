package main

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var svc *s3.S3

func init() {
	sess, _ := session.NewSession()

	svc = s3.New(sess)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(r *Request) (*Response, error) {
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Key: aws.String(r.QueryStringParameters["key"]),
		// use your own bucket name below instead of "100daysofcode-upload"
		Bucket: aws.String("100daysofcode-upload"),
	})

	presignedURL, err := req.Presign(15 * time.Minute)

	if err != nil {
		return nil, err
	}

	resp := &Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            presignedURL,
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
