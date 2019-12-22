# Minimart API

[![Golang](https://golang.org/lib/godoc/images/go-logo-blue.svg)](https://golang.org/)

Minimart API is a simple microservice capable of handling REST API and authorize users using an open-source authentication module.

### Tech

Minimart API uses a number of open source projects to work properly:

* [Golang](https://golang.org/) - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.
* [GORM](https://gorm.io/) - The fantastic ORM library for Golang
* [gorilla/mux](https://github.com/gorilla/mux) - Package mux implements a request router and dispatcher.
* [PostgreSQL](https://www.postgresql.org/) - The World's Most Advanced Open Source Relational Database
* [Docker](https://www.docker.com/) - Securely build, share and run modern applications anywhere
* [jwt-go](https://github.com/dgrijalva/jwt-go) - A go (or 'golang' for search engine friendliness) implementation of JSON Web Tokens

### Installation

Minimart API requires [Docker](https://www.docker.com/) and [docker-compose](https://docs.docker.com/compose/) to run.

Install Docker and docker-compose to start the server
 - [Docker Desktop on Windows](https://docs.docker.com/docker-for-windows/install/)
 - [Docker on Linux](https://docs.docker.com/install/linux/docker-ce/centos/)
 - [Docker Desktop on MacOS](https://docs.docker.com/docker-for-mac/install/)
 - [Install docker-compose](https://docs.docker.com/compose/install/)

```sh
$ cd minimart-api
$ docker-compose up
```

### Usage
    - POST "https://{HOST}:9988/api/authenticate"
        {
            "username": myuser,
            "password": mypass
        }
    - GET "https://{HOST}:9988/api/products"
    - GET "https://{HOST}:9988/api/carts"
    - POST "https://{HOST}:9988/api/carts"
        {
            "id": 23232,
            "quantity": 5
        }
    - PUT "https://{HOST}:9988/api/carts/{productId}"
        {
            "id": 23232,
            "quantity": 5
        }
    - DELETE "https://{HOST}:9988/api/carts/{productId}"

### Todos

 - Write MORE Tests
 - Integrate with open-source authentication module target: Keycloak
 - Validate credentials against DB

