# Endpoints

POST /api/v1/users/register
- username
- password

POST /api/v1/users/login
- username
- password

GET /api/v1/users/:username/messages
- username
- body

# tables

    users
----------
id
username
password

messages
-----------
body
recipient_id
sender_id

# DEPENDENCIES

* Gorilla
* GORM
* Golang-JWT
* bcrypt


# TASK LIST

- [x] Gorm installation
- [x] User model
- [x] Registration Endpoint
- [x] JWT handler
- [x] Login endpoint

- [x] gorilla installation
- [x] message model
- [x] message creation endpoint

- [x] websocket index endpoint
- [x] message queue

- [ ] public channels
