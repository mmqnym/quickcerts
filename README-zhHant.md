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

## 介紹

此專案（QuickCertS, QCS）提供開發者快速建立一個憑證伺服器，可為您開發的應用核發憑證（金鑰＆簽名）。於您的應用中嵌入公鑰，並同時使用簽名，驗證金鑰是否是由您架設的憑證伺服器核發，達到＂產品啟用＂的功能。

同時也提供有效期的臨時許可，若您的應用需要＂試用期＂或是產品基於週期授權而非永久時，可使用相關 API 達成。

## 配置

- 您可於 `path_to_qcs/configs/allowlist.toml` 中設置您要配置給管理員的名稱以及通行令牌，用於管理員用 API。

- 您可於 `path_to_qcs/configs/database.toml` 中將預設的配置更改為您期望的配置，但若於之後使用 `docker compose` 啟動伺服器，須要同樣更改以下 `docker compose` 的相關配置。

  ```yml
  services:
  qcs-db: # <- container name 對應於 host name
    build:
      context: .
      dockerfile: Dockerfile.database
    networks:
      - qcs_subnet
    ports:
      - "33332:5432"
    environment:
      POSTGRES_USER: quickcerts
      POSTGRES_PASSWORD: password # <- 建議更改資料庫密碼
      POSTGRES_DB: quickcerts
  ```

- `path_to_qcs/configs/server.toml` 包含全部伺服器的相關設定，建議正式執行前完成配置。

- `path_to_qcs/init.sql` 中可以替資料庫設定時區，建議使用與本地或雲端相同的時區，避免混亂。

## 构建

#### Docker

> 使用 docker 和 docker compose 快速启动服务器

确保您已经在您的操作系统上安装了 `docker` 和 `docker compose`。在项目根目录中运行以下命令：

```sh
docker compose up --build -d
```

即可完成构建，若未更改配置设置，默认运行端口 `:33333`

#### 可执行文件

> 使用 Release 提供的可执行文件

- 创建一个 PostgreSQL 数据库，并将相关配置设置到 `path_to_qcs/configs/database.toml`。

- 至 Release 根据您的操作系统选择要下载的压缩文件，然后在项目根目录中运行 `./init/Init(.exe)`。

- 在项目根目录中运行 `server(.exe)`。

即可完成构建，若未更改配置设置，默认运行端口 `:33333`

#### 源代码

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

## 技术

架构：

- 服务器框架：Gin Web Framework
- 数据库：PostgreSQL

公私钥存储规范：PKCS8

签名：

| SHA2    | SHA3     |
| ------- | -------- |
| SHA-256 | SHA3-256 |
| SHA-384 | SHA3-384 |
| SHA-512 | SHA3-512 |

使用 RSA-PSS 填充自动计算的长度
