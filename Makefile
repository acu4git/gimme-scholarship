PWD				:= $(shell pwd)
OS				:= $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH			:= $(shell uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')
DB_PORT		:= 3308

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

database/init: database/up sleep docker/mysql/migrate
	- mysql -h 127.0.0.1 -P $(DB_PORT) -uroot -proot gimme_scholarship < seeds/seed.sql

database/up:
	docker compose up -d

database/down:
	docker compose down 

sleep:
	until (mysqladmin ping -h 127.0.0.1 -P $(DB_PORT) -uroot -proot --silent) do echo 'waiting for mysql connection...' && sleep 2; done

docker/mysql/migrate:
	sql-migrate up -env="development"
