build:
	#生成proto
	protoc --go_out=paths=source_relative:. \
        --micro_out=paths=source_relative:. \
        proto/consignment/consignment.proto

	#根据当前目录的Dockerfile，创建一个名为consignment-service的容器
	docker build -t consignment-service .
run:
	docker run -p 50051:50051 \
	-e MICRO_SERVER_ADDRESS=:50051 \
	-e MICRO_REGISTRY=mdns consignment-service