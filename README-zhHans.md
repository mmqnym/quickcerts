# QuickCerts

<p>
    <img alt="English" src="https://img.shields.io/badge/English-000000?style=for-the-badge">
        <a href="./README.md"></a>
    </img>
    <img alt="繁體中文" src="https://img.shields.io/badge/繁體中文-000000?style=for-the-badge">
        <a href="./README-zhHant.md"></a>
    </img>
    <img alt="简体中文" src="https://img.shields.io/badge/简体中文-000000?style=for-the-badge">
        <a href="./README-zhHans.md"></a>
    </img>
</p>

## 介绍

本项目（QuickCertS，QCS）旨在帮助开发者快速建立一个证书服务器，用于为您的应用程序颁发证书（密钥和签名）。您可以嵌入公钥到您的应用中，并同时使用签名来验证密钥是否由您架设的证书服务器颁发，实现产品激活功能。

此外，还提供了临时许可证的支持，如果您的应用需要试用期或者产品需要基于许可周期而不是永久授权，您可以使用相关 API 来实现这一功能。

## 配置

- 您可以在 `path_to_qcs/configs/allowlist.toml` 文件中配置管理员的用户名和令牌，用于管理员 API 的身份验证。

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

- `path_to_qcs/init.sql` 中可以设置数据库的时区，建议使用与本地或云端相同的时区，以避免混淆。

## 建置

#### Docker

> 使用 docker 以及 docker compose 快速啟動伺服器

確保您已在您的作業系統上安裝 `docker` 以及 `docker compose` 於專案根目錄下輸入：

```sh
docker compose up --build -d
```

即可完成架設，若未更改配置設定，預設啟動埠號 `:33333`

#### 執行檔

> 使用 Release 中提供的執行檔

- 建立一個 PostgreSQL 的資料庫，並將相關配置設置到 `path_to_qcs/configs/database.toml`。

- 至 Release 根據您的作業系統選擇要下載的壓縮檔，於專案根目錄執行 `./init/Init(.exe)`。

- 於專案根目錄執行 `server(.exe)`。

即可完成架設，若未更改配置設定，預設啟動埠號 `:33333`

#### 原始碼

> 由原始碼編譯後執行或直接執行

請使用 Golang `>= 1.21.1` 編譯後執行或直接執行：

```sh
go run ./init/Init.go
go run ./server.go
```

即可完成架設，若未更改配置設定，預設啟動埠號 `:33333`

## API

> 啟用伺服器後可至以下網址查閱：

預設：http://localhost:33333/swagger/index.html

若有使用 TLS 或不同的埠號請自行切換網址。

## 技術

架構：

- 伺服器框架：Gin Web Framework
- 資料庫：PostgreSQL

公私鑰儲存規範：PKCS8

簽名：

| SHA2    | SHA3     |
| ------- | -------- |
| SHA-256 | SHA3-256 |
| SHA-384 | SHA3-384 |
| SHA-512 | SHA3-512 |

使用 RSA-PSS 填充自動計算的長度
