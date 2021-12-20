
build-grpc:
	go build -v -o ./bin/xds-client .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/xds-client-linux .

docker-build-grpc: build-grpc
	docker build . -t allenlsy/grpc-xds-client:latest

docker-push-grpc: docker-build-grpc
	docker push allenlsy/grpc-xds-client:latest