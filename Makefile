PWD				 		 := $(shell pwd)
OS				 		 := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH			 		 := $(shell uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')
DB_PORT		 		 := 3308
GO_VERSION 		 := $(shell cat .go-version)
PYTHON_VERSION := $(shell cat .python-version)

init/tools: tbls/install sql-migrate/install

tbls/install:
	@if [ ! -x "$(command -v tbls)" ]; then\
		go install github.com/k1LoW/tbls@latest;\
	fi
	tbls version

sql-migrate/install:
	@if [ ! -x "$(command -v sql-migrate)" ]; then\
		go install github.com/rubenv/sql-migrate/...@latest;\
	fi
	sql-migrate --version

base/init: database/init docker/fetch/run docker/api/run

database/init: database/up sleep docker/mysql/migrate 
	- mysql -h 127.0.0.1 -P $(DB_PORT) -uroot -proot gimme_scholarship < seeds/seed.sql
	- mysql -h 127.0.0.1 -P $(DB_PORT) -uroot -proot gimme_scholarship_test < seeds/sees.sql

database/up:
	docker compose up -d

database/down:
	docker compose down 

sleep:
	until (mysqladmin ping -h 127.0.0.1 -P $(DB_PORT) -uroot -proot --silent) do echo 'waiting for mysql connection...' && sleep 2; done

docker/mysql/migrate:
	sql-migrate up -env="development"
	sql-migrate up -env="test"

docker/api/build:
	docker build -t gimme-scholarship-api -f './cmd/api/Dockerfile' . --build-arg GO_VERSION=$(GO_VERSION)

docker/migrate/build:
	docker build -t gimme-scholarship-migrate -f './cmd/migrate/Dockerfile' . --build-arg GO_VERSION=$(GO_VERSION)

docker/task/build:
	docker build -t gimme-scholarship-task -f './cmd/task/Dockerfile' . --build-arg GO_VERSION=$(GO_VERSION)

docker/fetch/build:
	docker build -t gimme-scholarship-fetch -f './cmd/fetch/Dockerfile' ./cmd/fetch --build-arg PYTHON_VERSION=$(PYTHON_VERSION)

docker/api/run: docker/api/build
	docker run -p 8080:8080 --net=gs-net -e DB_HOST=db -e DB_PORT=3306 gimme-scholarship-api

docker/fetch/run: docker/fetch/build
	docker run --net=gs-net -e DB_HOST=db -e DB_PORT=3306 gimme-scholarship-fetch

gen/tbls:
	tbls doc --rm-dist