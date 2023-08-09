module github.com/nyaruka/phonenumbers/cmd/phoneserver

go 1.19

replace github.com/nyaruka/phonenumbers => ../../

require (
	github.com/aws/aws-lambda-go v1.13.1
	github.com/nyaruka/phonenumbers v1.1.7
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/text v0.12.0 // indirect
)
