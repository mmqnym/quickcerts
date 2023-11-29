# QuickCertS

## Language

<p>
    <a href="./README.md"><img alt="English" src="https://img.shields.io/badge/English-000000?style=for-the-badge"></img></a>
    <a href="./README-zhHant.md"><img alt="繁體中文" src="https://img.shields.io/badge/繁體中文-000000?style=for-the-badge"></img></a>
    <a href="./README-zhHans.md"><img alt="简体中文" src="https://img.shields.io/badge/简体中文-000000?style=for-the-badge"></img></a>
</p>

## Introduction

This project (QuickCertS, QCS) aims to help developers quickly establish a certificate server for issuing certificates (keys and signatures) for your applications. You can embed a public key into your application and use a signature to verify whether the key is issued by the certificate server you have set up, achieving "product activation" functionality.

Additionally, QCS provides support for temporary permission. If your application requires a "trial period" or periodic authorization rather than permanent authorization, you can use the relevant API to achieve this.

## Configuration

- You can configure the names and tokens for administrators in the `path_to_qcs/configs/allowlist.toml` file, which is used for administrator authentication in the admin API.

- You can change the default configuration to your desired settings in the `path_to_qcs/configs/database.toml` file. However, if you later start the server using `docker compose`, you will need to change the `docker compose` file accordingly.

  ```yml
  services:
  qcs-db: # <- container name corresponding to host name
    build:
      context: .
      dockerfile: Dockerfile.database
    networks:
      - qcs_subnet
    ports:
      - "33332:5432"
    environment:
      POSTGRES_USER: quickcerts
      POSTGRES_PASSWORD: password # <- It is recommended to change the database password
      POSTGRES_DB: quickcerts
  ```

- The `path_to_qcs/configs/server.toml` file contains all the relevant settings  
  for the server. It is recommended to configure it before running the server officially.

- In the `path_to_qcs/init.sql` file, you can set the time zone for the database.
  It is recommended to use the same time zone as your local or cloud environment to avoid confusion.

## Running

- ### Docker

> Quickly start the server using docker and docker compose

Ensure that you have installed docker and docker compose on your OS. Run the following command in the project's root directory:

```sh
docker compose up --build -d
```

The server will be built and started. If you have not changed the configuration settings, the server will run on the default port `:33333`.

- ### Executables

> Use the executable files provided in the Release

- Create a PostgreSQL database and set the relevant configuration in the `path_to_qcs/configs/database.toml` file.

- In the Release, choose the compressed file to download based on your OS, and run `./init/Init(.exe)` in the project's root directory.

- Run `server(.exe)` in the project's root directory.

The server will be built and started. If you have not changed the configuration settings, the server will run on the default port `:33333`.

- ### Source code

> Build and run from source code, or run directly

Please use Golang version `>= 1.21.1` to compile and run, or run directly:

```sh
go run ./init/Init.go
go run ./server.go
```

The server will be built and started. If you have not changed the configuration settings, the server will run on the default port `:33333`.

## API

> After starting the server, you can access the API documentation at the following URL:

Default：http://localhost:33333/swagger/index.html

If you are using TLS or a different port, please adjust the URL accordingly.

## Technology

Architecture:

- Server Framework: Gin Web Framework
- Database: PostgreSQL

Public and private key storage standard: PKCS8

Signature:

| SHA2    | SHA3     |
| ------- | -------- |
| SHA-256 | SHA3-256 |
| SHA-384 | SHA3-384 |
| SHA-512 | SHA3-512 |

Automatic length calculation with RSA-PSS padding
