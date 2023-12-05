# QuickCertS

## 語言

<p>
    <a href="./README.md"><img alt="English" src="https://img.shields.io/badge/English-000000?style=for-the-badge"></img></a>
    <a href="./README-zhHant.md"><img alt="繁體中文" src="https://img.shields.io/badge/繁體中文-000000?style=for-the-badge"></img></a>
    <a href="./README-zhHans.md"><img alt="简体中文" src="https://img.shields.io/badge/简体中文-000000?style=for-the-badge"></img></a>
</p>

## 介紹

此專案（QuickCertS, QCS）提供開發者快速建立一個憑證伺服器，可為您開發的應用核發憑證（金鑰＆簽名）。於您的應用中嵌入公鑰，並同時使用簽名，驗證金鑰是否是由您架設的憑證伺服器核發，達到＂產品啟用＂的功能。

同時也提供有效期的臨時許可，若您的應用需要＂試用期＂或是產品基於週期授權而非永久時，可使用相關 API 達成。

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

## 建置

- ### Docker

> 使用 docker 以及 docker compose 快速啟動伺服器

確保您已在您的作業系統上安裝 `docker` 以及 `docker compose` 於專案根目錄下輸入：

```sh
docker compose up --build -d
```

即可完成架設，若未更改配置設定，預設啟動埠號 `:33333`

- ### 執行檔

> 使用 Release 中提供的執行檔

- 建立一個 PostgreSQL 的資料庫，並將相關配置設置到 `path_to_qcs/configs/database.toml`。

- 至 Release 根據您的作業系統選擇要下載的壓縮檔，於專案根目錄執行 `./init/Init(.exe)`。

- 於專案根目錄執行 `server(.exe)`。

即可完成架設，若未更改配置設定，預設啟動埠號 `:33333`

- ### 原始碼

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

## SDK

> SDK & 範例

您可於 `path_to_qcs/sdk` 查看 SDK 以及使用範例，目前支援 Python, TypeScript, Golang。

- #### Python

於 `path_to_qcs/sdk/python` 開啟終端，輸入：

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

- #### TypeScript

於 `path_to_qcs/sdk/typescript` 開啟終端，輸入：

```sh
npm i

npm run start # SDK Usage
npm run verify # Verify RSA signature.
```

- #### Golang

於 `path_to_qcs/sdk/go` 開啟終端，輸入：

```sh
cd ./example

go run usage.go # SDK Usage.
# If you want to run the verification test case, you can call VerifyExample().
```
