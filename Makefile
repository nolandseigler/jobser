
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
ETC_WORDSERWEB = docker/etc/wordserweb
DEV_PRIV_KEY_PATH := $(ETC_WORDSERWEB)/keys/priv.rsa.pem
DEV_PUB_KEY_PATH := $(ETC_WORDSERWEB)/keys/pub.rsa.pem

.PHONY: up up-d down down-v

gen-keys:
	mkdir -p $(ETC_WORDSERWEB)/keys
	openssl genrsa -out $(DEV_PRIV_KEY_PATH) 4096
	openssl rsa -in $(DEV_PRIV_KEY_PATH) -outform PEM -pubout -out $(DEV_PUB_KEY_PATH)

keys:
	test -s $(DEV_PRIV_KEY_PATH) || $(MAKE) gen-keys

up:	keys
	docker-compose up --build --force-recreate  --remove-orphans

up-d: keys
	docker-compose up --build --force-recreate  --remove-orphans -d

down:
	docker-compose down

down-v:
	docker-compose down -v
