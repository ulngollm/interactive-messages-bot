include .env

build:
	CGO_ENABLED=0 go build -o ${BOT_BINARY_NAME}

start:
	./${BOT_BINARY_NAME} > log &

remote-start:
	ssh ${SERVER} "cd ${DEPLOY_DIR} && ./${BOT_BINARY_NAME} > log &"

deploy:
	ssh ${SERVER} "mkdir -p ${DEPLOY_DIR}"
	scp ${BOT_BINARY_NAME} ${SERVER}:${DEPLOY_DIR}/${BOT_BINARY_NAME}
	scp .env ${SERVER}:${DEPLOY_DIR}/
	scp Makefile ${SERVER}:${DEPLOY_DIR}/