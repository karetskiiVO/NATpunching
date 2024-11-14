all:
	build_server
	build_client

build_server:
	cd ./server && go build -o ../serverapp

build_client:
	cd ./client && go build -o ../clientapp