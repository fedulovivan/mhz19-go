migrate-reset: migrate-down migrate-up

migrate-down:
	DB_REV=05 make migrate-down-single
	DB_REV=04 make migrate-down-single
	DB_REV=03 make migrate-down-single
	DB_REV=02 make migrate-down-single
	DB_REV=01 make migrate-down-single
	DB_REV=00 make migrate-down-single

migrate-up:
	DB_REV=00 make migrate-up-single
	DB_REV=01 make migrate-up-single
	DB_REV=02 make migrate-up-single
	DB_REV=03 make migrate-up-single
	DB_REV=04 make migrate-up-single
	DB_REV=05 make migrate-up-single

migrate-up-single:
	sqlite3 ./sqlite/database.bin < ./sql/$(DB_REV)-up.sql

migrate-down-single:
	sqlite3 ./sqlite/database.bin < ./sql/$(DB_REV)-down.sql

docker-migrate-up-single:
	docker exec -i $(NAME)-1 sqlite3 ./sqlite/database.bin < ./sql/$(DB_REV)-up.sql

docker-migrate-down-single:
	docker exec -i $(NAME)-1 sqlite3 ./sqlite/database.bin < ./sql/$(DB_REV)-down.sql

migrate-dump:
	sqlite3 ./sqlite/database.bin .dump > ./sql/$(DATE)-dump.sql