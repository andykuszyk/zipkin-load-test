build:
	go build cmd/main.go

build_linux:
	env GOOS=linux GOARCH=amd64 go build -o sqs-load-test cmd/main.go


run-sqs:
	docker run -d -p 4100:4100 pafortin/goaws
	./main -test sqs -sqsEndpoint "http://localhost:4100"

run-zipkin:
	docker run -d -p 9411:9411 openzipkin/zipkin
	./main -test zipkin
