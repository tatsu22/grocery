run-api:
	go run main.go

start-container:
	sudo docker run -p 8080:8080 --name grocery -d --network grocery-network testing123