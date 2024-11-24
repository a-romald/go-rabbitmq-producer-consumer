# go-rabbitmq-producer-consumer

## Project description

The project represents golang application uses rabbitmq and microservices to reverse words. Frontend application on localhost:8000 shows form to type any word. JS AJAX sends it to server on same port that pushes this word to queue in RabbitMQ. Then listener microservice consumes this word from the queue and sends with POST-request to worker microservice. Worker microservice pushes this data to channel in goroutine. Another goroutine in main func pulls it and sends to open websocket. Open page on localhost:8000 connected to websocket that receives reverse form of the word and shows it on opened page.

## Project setup
```
docker-compose up -d
```
#### Open project
```
localhost:8000
```

Open page and type any word to show result.
