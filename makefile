all: build_server build_client

build_server:
	cd ./server && go mod tidy && go build -o ../serverapp

build_client:
	cd ./client && go mod tidy && go build -o ../clientapp