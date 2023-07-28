![Go Language](https://img.shields.io/badge/Go-1.20-blue.svg)
![Docker](https://img.shields.io/badge/Docker-24.0.2-blue.svg)
![Keycloak](https://img.shields.io/badge/Keycloak-21.1.2-blue.svg)

# OAuth Client with Keycloak Integration

<p align="center">
  <img src="./src/img/logo.png" width="30%" alt="demo header">
</p>

<!-- CONTACTS -->
<div>
    <p align="center">
        <br />
        <a href="#running-the-application">About Project</a>
        ·
        <a href=mailto:"konstantin.priluchnyi@gmail.com?subject=report%20Bug">Report Bug</a>
        ·
        <a href=mailto:"konstantin.priluchnyi@gmail.com?subject=Request%20Feature">Request Feature</a>
    </p>
    <p align="center">
        <a href=mailto:"konstantin.priluchnyi@gmail.com">
            <img src="https://img.shields.io/badge/Gmail-D14836?style=for-the-badge&logo=gmail&logoColor=white">
        </a>
        <a href="https://t.me/konstantin_evo">
            <img src="https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white">
        </a>
        <a href="https://www.linkedin.com/in/konstantin-evo/">
            <img
                src="https://img.shields.io/badge/linkedin-%230077B5.svg?style=for-the-badge&logo=linkedin&logoColor=white">
        </a>
    </p>
</div>

This repository contains an OAuth client built in Go that integrates with Keycloak for authentication and authorization.

## Running the Application

To run the application, make sure you have Docker and Docker Compose installed on your system. Then, navigate to the
root directory of the project and run the following command:

```bash
docker-compose up
```

This will build the Docker images and start the containers for both the Keycloak authentication server and the Go OAuth
client. Once the containers are up and running, you can access the following resources:

- Go OAuth client by opening your web browser and visiting `http://localhost:8080`.
- Keycloak Auth server by opening your web browser and visiting `http://localhost:8081`.

The client is set up to use the test user credentials mentioned below for authentication and authorization. You can
modify the client to integrate with other users and roles as needed.

## Keycloak credentials

To access the Keycloak admin console, use the following credentials:

- Username: `admin`
- Password: `admin`

For testing authentication and authorization, you can use the following test user:

- Username: `myUser`
- Password: `1111`

## Project Structure

The project directory structure is as follows:

```
├── README.MD
├── docker-compose.yml
└── src
    ├── client
    │   ├── Dockerfile
    │   ├── go.mod
    │   ├── go.sum
    │   ├── handlers.go
    │   ├── main.go
    │   └── template
    │       └── index.html
    └── keycloak
        ├── Dockerfile
        └── data
            └── realm-export.json
```

The `docker-compose.yml` file defines the services for the Keycloak authentication server and the Go OAuth client.

The Go OAuth client code is located in the `src/client` directory. It uses the Gorilla Mux library for handling routes
and HTTP requests.

The Keycloak configuration files are in the `src/keycloak` directory, including the `Dockerfile` and a sample realm
configuration file (`realm-export.json`).
