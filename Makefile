build: 
	@echo "Building app..."
	env GOOS=linux CGO_ENABLED=0 go build -o echoRest ./api 
	@echo "Done"

up_build:	build
	@echo "Stopping docker images if running, building when required and starting docker images..."
	sudo docker-compose down
	sudo docker-compose up -d --build
	@echo "Done"