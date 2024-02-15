mongo:
	docker run --name mongo_auth -e MONGODB_DATABASE=admin -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=qwerty -p 27019:27017 -d mongo


.PHONY: mongo 