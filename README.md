
# Insider Automatic Message Sending System


This project implements an automatic message-sending system that retrieves unsent messages from a database and sends two messages every two minutes. The system ensures that messages are not resent and provides API endpoints to control the message-sending process and retrieve sent messages.

  

## Overview
  

The system processes unsent messages stored in a database and sends them automatically. Once a message is sent, it will not be resent. Additionally, the system provides an API to start/stop the automatic message sender and retrieve the list of sent messages.

  

## Features
  

- Automatically retrieves and sends two unsent messages every two minutes.

- Ensures that messages are sent only once.

- Provides API endpoints to retrieve sent messages.

- Caches sent messages returned from webhook.site

- API documentation is available via Swagger.

  

## Prerequisites

- Go: Version 1.23+

- Docker/Docker-compose has been used for development.

- MySQL: A relational database to store messages. (no installation required)

- Redis Used for caching sent message IDs and timestamps. (no installation required)

  
  

## Installation & Setup

  
```
git clone https://github.com/mehmetalisavas/message-sender.git
cd message-sender
```

  
Start application with:
  
`make docker-build`


  
Run tests with:

`make test-docker`

  

## API Endpoints

#### LIST SENT MESSAGES

 
Basic:

`curl -X GET "http://localhost:8080/messages" -H "Content-Type: application/json"`

  

With Limit & Offset:
`curl -X GET "http://localhost:8080/messages?limit=10&offset=20" -H "Content-Type: application/json"`

  

#### START / STOP MESSAGE SENDING


Start message sending (default)
  
`curl -X GET "http://localhost:8080/process_message?command=start" -H "Content-Type: application/json"`

  

Stop message sending
`curl -X GET "http://localhost:8080/process_message?command=stop" -H "Content-Type: application/json"`


#### View Swagger Docs 
`curl -X GET "http://localhost:8080/swagger/index.html"`


## Architecture & Design

The system is designed to be scalable and efficient, leveraging goroutines for concurrent message processing.




### High-Level Architecture

##### API SERVER

**Crobjob**
**(System Scheduler)**

| -- Start/Stop flag

| -- Fetch messages from DB

| -- Send messages asynchronously

| -- Consume messages asynchronously

  
**App Connections**

|---> Remote Call (Notification Service)

|---> Database (Stores messages)

|---> Redis (Caches sent message IDs)
  
  
  

## Key Design Decisions

- Concurrency: Uses goroutines to send messages asynchronously without blocking execution.

- Custom Scheduler: Instead of relying on external cron packages, a native Go timer schedules tasks every two minutes.
Scheduler is designed as an extendible, Plus producer & consumer parts are introduced for single responsibility purpose. 
**Producers** will fetch required data from DB and send it to Consumer via channels.
**Consumers** will consume from related channels and process messages.

- Data Integrity: Ensures messages are sent only once by updating the database after sending.

- Scalability: The system is designed to handle high throughput with minimal resource usage.


## Assumptions
- Webhook.site is considered as a `notification service` for tracking sent messages.
- Response will include different message Id for different messages.(Redis doesn't have to override in that case.)
- Character limits are enforced at the database level to prevent overly long messages.
- Newly added records will only be picked up in the next processing cycle, records will be picked up in order (according to created_at)
- No external cron jobs or scheduling libraries are used; instead, a native Go timer handles scheduling.
- If message fails after multiple retry, marked as 'failed', and it should be handled in a different scope


## Future Improvements
- Persistent Storage: Store logs and analytics data for long-term tracking.
- No external logging framework has been used, but future improvements may include structured logging libraries.
- Retry Mechanism for failed messages(It's ignored for this illustration)
- Tracing & Logging: Implement distributed tracing for better observability.
- Prometheus Integration: Add metrics and monitoring using Prometheus.
- Unit & Integration Tests: Extend test coverage for better reliability.
- CI/CD Pipeline: Implement automated testing, linting, deployment, and monitoring through CI/CD tools like GitHub Actions or Jenkins.

