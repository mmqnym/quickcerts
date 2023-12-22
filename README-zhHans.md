# QuickCertS

<p align="center">
    <img alt="app version" src="https://img.shields.io/github/v/release/mmq88/quickcerts"></img>
    <a href="https://app.codecov.io/gh/mmq88/quickcerts"><img alt="Codecov" src="https://img.shields.io/codecov/c/github/mmq88/quickcerts?logo=codecov&logoColor=%23F01F7A&label=codecov"></a>
    <a href="https://github.com/mmq88/quickcerts/actions/workflows/codeql.yml"><img src="https://github.com/mmq88/quickcerts/workflows/CodeQL/badge.svg" alt="CodeQL status"></a>
    <a href="https://goreportcard.com/report/github.com/mmq88/quickcerts"><img src="https://goreportcard.com/badge/github.com/mmq88/quickcerts"></a>
    <img alt="license" src="https://img.shields.io/github/license/mmq88/quickcerts"></img>
</p>

<p align="center">
    <img alt="go version" src="https://img.shields.io/github/go-mod/go-version/mmq88/quickcerts"></img>
    <img alt="python version" src="https://img.shields.io/badge/Python-v3.9.13-blue"></img>
    <img alt="node version" src="https://img.shields.io/badge/Node-v18.16.0-blue"></img>
</p>

## 语言

<p>
    <a href="./README.md"><img alt="English" src="https://img.shields.io/badge/English-6498cc?style=for-the-badge"></img></a>
    <a href="./README-zhHant.md"><img alt="繁體中文" src="https://img.shields.io/badge/繁體中文-6498cc?style=for-the-badge"></img></a>
    <a href="./README-zhHans.md"><img alt="简体中文" src="https://img.shields.io/badge/简体中文-6498cc?style=for-the-badge"></img></a>
</p>

## 介绍

本项目（QuickCertS，QCS）旨在帮助开发者快速建立一个证书服务器，用于为您的应用程序颁发证书（密钥和签名）。您可以嵌入公钥到您的应用中，并同时使用签名来验证密钥是否由您架设的证书服务器颁发，实现产品激活功能。

此外，还提供了临时许可证的支持，如果您的应用需要试用期或者产品需要基于许可周期而不是永久授权，您可以使用相关 API 来实现这一功能。

## 技术

架构：

- 服务器框架：Gin Web Framework
- 緩存：Redis
- 数据库：PostgreSQL

公私钥存储规范：PKCS8

签名：

| SHA2    | SHA3     |
| ------- | -------- |
| SHA-256 | SHA3-256 |
| SHA-384 | SHA3-384 |
| SHA-512 | SHA3-512 |

使用 RSA-PSS 填充自动计算的长度

## 配置

- 您可以在 `path_to_qcs/configs/allowlist.toml` 文件中配置管理员的用户名和令牌，用于管理员 API 的身份验证。

- 您可以在 `path_to_qcs/configs/cache.toml` 中将默认的配置更改为您期望的配置。

- 您可以在 `path_to_qcs/configs/database.toml` 文件中更改默认配置为您所期望的配置。但如果您之后使用 Docker Compose 启动服务器，则需要相应更改 Docker Compose 的配置。

  ```yml
  services:
  qcs-db: # <- 容器名称对应主机名称
    build:
      context: .
      dockerfile: Dockerfile.database
    networks:
      - qcs_subnet
    ports:
      - "33332:5432"
    environment:
      POSTGRES_USER: quickcerts
      POSTGRES_PASSWORD: password # <- 建议更改数据库密码
      POSTGRES_DB: quickcerts
  ```

- `path_to_qcs/configs/server.toml` 包含了服务器的所有相关设置，建议在正式运行之前完成配置。
  ※请将以下的设置改为 `false`。

  ```toml
  LOG_TEST_MODE = false
  ```

- `path_to_qcs/init.sql` 中可以设置数据库的时区，建议使用与本地或云端相同的时区，以避免混淆。

- 如果您了解如何使用 Redis，可于 `path_to_qcs/redis.conf` 更动 Redis 的默认值。

## 构建

- ### Docker

> 使用 docker 和 docker compose 快速启动服务器

确保您已经在您的操作系统上安装了 `docker` 和 `docker compose`。在项目根目录中运行以下命令：

```sh
docker compose up --build -d
```

即可完成构建，若未更改配置设置，默认运行端口 `:33333`

### 可执行文件

> 使用 Release 提供的可执行文件

- 创建一个 PostgreSQL 数据库，并将相关配置设置到 `path_to_qcs/configs/database.toml`。

- 至 Release 根据您的操作系统选择要下载的压缩文件，然后在项目根目录中运行 `./init/Init(.exe)`。

- 在项目根目录中运行 `server(.exe)`。

即可完成构建，若未更改配置设置，默认运行端口 `:33333`

### 源代码

> 从源代码编译并运行，或直接运行

请使用 Golang 版本 `>= 1.21.1` 进行编译和运行，或者直接运行：

```sh
go run ./init/Init.go
go run ./server.go
```

即可完成构建，若未更改配置设置，默认运行端口 `:33333`

## API

> 启动服务器后访问以下网址以查看：

默认：http://localhost:33333/swagger/index.html

如果使用 TLS 或不同端口，请相应调整网址。

## SDK

> SDK & 示例

您可以在 `path_to_qcs/sdk` 查看 SDK 以及使用示例，目前支持 Python、TypeScript、Golang。

- ### Python

於 `path_to_qcs/sdk/python` 打开终端，输入：

```sh
# Here uses pyenv + virtualenv + pip,
# you can also use your preferred environment/package management tool.
virtualenv -p "path_to_python" venv
./venv/Script/activate
pip install -r "./requirements.txt"
cd ./example

python ./usage.py # SDK Usage
python ./verify.py # Verify RSA signature.
```

- ### TypeScript

於 `path_to_qcs/sdk/typescript` 打开终端，输入：

```sh
npm i

npm run start # SDK Usage
npm run verify # Verify RSA signature.
```

- ### Golang

於 `path_to_qcs/sdk/go` 打开终端，输入：

```sh
cd ./example

go run usage.go # SDK Usage.
# If you want to run the verification test case, you can call VerifyExample().
```
