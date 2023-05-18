# Topics project
# Endpoints
## Users
### Register
Register a user with username and password.

**POST** `/api/v1/users/register`

*Request body*
```json
{
	"username":  "string",
	"password":  "string"
}
```
*Responses*
```json
{
	"message": "registration success"
}
```

### Login
Login user with username and password.

**POST** `/api/v1/users/login`

*Request body*
```json
{
	"username":  "string",
	"password":  "string"
}
```
*Responses*
```json
{
	"token": "string"
}
```

## Messages
### User messages
Send a message to specific user with the username param or your can read your message inbox if the username in the param is yours.

**Websocket** `/api/v1/users/:username/messages`

*Message body*
```json
{
	"content": "string"
}
```
*Responses*
```json
{
	"content": "string",
	"username": "string"
}
```

## Channels
### Channel
You can Subscribe to a channel with a specific name in the channel param, use the command field in the message body with this options for a specific command:
* 0 = Subscribe
* 1 = Unsubscribe
* 2 = Send message

**Websocket** `/api/v1/channels/:channel`

*Message body*
```json
{
	"command": 0,
	"content": "string"
}
```
*Responses*
```json
{
	"content": "string",
	"channel": "string"
}
```
# Tables

| users     |          |
|-----------|----------|
| id        | uint     |
| username  | string   |
| password  | string   |
| CreatedAt | datetime |
| UpdatedAt | datetime |

| Messages |        |
|----------|--------|
| content  | string |
| channel  | string |
| username | string |
| command  | int    |
| error    |        |

# Run project locally

Run project with docker compose
```bash
make local.run
```

Stop all containers
```bash
make local.run
```

Build project
```bash
make local.build
```

# DEPENDENCIES

* Gorilla
* GORM
* Golang-JWT
* Bcrypt
* Rabbitmq/amqp091-go
* Go-redis/redis/v7

# TASK LIST

## User resources
- [x] Gorm installation
- [x] User model
- [x] Registration Endpoint
- [x] JWT handler
- [x] Login endpoint

## Websockets
- [x] gorilla installation
- [x] message model
- [x] message creation endpoint

## User messages
- [x] websocket index endpoint
- [x] message queue

## Channels
- [x] Users channels
