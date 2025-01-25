include .env

build:
	CGO_ENABLED=0 go build -o ${BOT_BINARY_NAME}

start:
	./${BOT_BINARY_NAME} > log &

deploy:
	scp ${BOT_BINARY_NAME} ${SERVER}:${DEPLOY_DIR}
	scp .env ${SERVER}:${DEPLOY_DIR}
	scp Makefile ${SERVER}:${DEPLOY_DIR}