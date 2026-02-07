up:
	docker-compose up --build

down:
	docker-compose down

restart:
	docker-compose down
	docker-compose up --build

logs:
	docker-compose logs -f

build:
	docker-compose build

clean:
	rm -f data/users.csv