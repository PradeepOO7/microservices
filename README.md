## Simple Go Webserver Implematation using Microservices 

## Prequisties
    #Docker
    Golang version 1.18 or latest

## The microservices i have build includes the following functionality:

    A Front End service, that just displays web pages;

    An Authentication service, with a Postgres database;

    A Logging service, with a MongoDB database;

    A Listener service, which receives messages from RabbitMQ and acts upon them;

    A Broker service, which is an optional single point of entry into the microservice cluster;

    A Mail service, which takes a JSON payload, converts into a formatted email, and send it out.


  
## Steps to run the Microservice:-
   ### Step 1 
    Start your docker engine first.
  
  ### Step 2
  #### change your current directory to different microservices directory and run go mod tidy e.g.
    cd broker-service 
    go mod tidy
    
  
  ### Step 3 
  #### change your current directory to project directory and run:-
    make up_build //starts docker-compose
    make down //removes services from docker
    
 ### Step 4
 #### change current directory to front-end and run:-
    make start
 
 #### Now the front end server is started open http:localhost:8080/ in your browser and test the microservices.
    
  
  
