default: build push

build:
	docker build -t fern4lvarez/piladb:latest .

push:
	docker push fern4lvarez/piladb:latest

pull:
	docker pull fern4lvarez/piladb

run:
	docker run -d --name piladb -p 1205:1205 fern4lvarez/piladb

bash:
	docker exec -it piladb bash

logs:
	docker logs -f piladb
