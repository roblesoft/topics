# Endpoints

POST /api/v1/users/register
- username
- password

POST /api/v1/users/login
- username
- password

GET /api/v1/channel/:username

POST /api/v1/channel/:username
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

- [*] Gorm installation
- [*] User model
- [*] Registration Endpoint
- [*] JWT handler
- [*] Login endpoint

- [*] gorilla installation
- [*] message model
- [*] message creation endpoint

- [] websocket index endpoint
- [] store pending messages

- [] public channels
