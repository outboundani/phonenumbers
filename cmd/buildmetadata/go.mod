module github.com/outboundani/phonenumbers/cmd/buildmetadata

go 1.20

replace github.com/outboundani/phonenumbers => ../..

require (
	github.com/outboundani/phonenumbers v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.31.0
)

require golang.org/x/text v0.12.0 // indirect
