# find-game-we-can-play-API

這是一個基於 Go 語言開發的遊戲查詢 API，允許用戶根據玩家人數、遊戲類型和平台篩選可用的遊戲。

## 功能特點

- 根據玩家數量篩選遊戲
- 支援遊戲類型和平台篩選
- 只返回支援線上模式的遊戲
- 資料來自 `games.json`

## 安裝與運行

### 1. 克隆專案

```sh
git clone https://github.com/your-repo/game-api.git
cd game-api
```

### 2. 準備 `games.json`

請確保與程式同一目錄下有 `games.json`，格式如下：

```json
[
  {
    "name": "Game Title",
    "min_players": 1,
    "max_players": 4,
    "platforms": ["PC", "PS5"],
    "genre": "RPG",
    "online": true
  }
]
```

### 3. 運行 API 伺服器

```sh
go run main.go
```

API 伺服器將運行於 `http://localhost:8080`

## API 使用說明

### 取得符合條件的遊戲

**請求方式：** `GET /games`

**查詢參數：**

| 參數         | 類型     | 說明       | 必填 |
| ---------- | ------ | -------- | -- |
| `players`  | int    | 玩家數量     | ✅  |
| `genre`    | string | 遊戲類型（選填） | ❌  |
| `platform` | string | 遊戲平台（選填） | ❌  |

**請求範例：**

```sh
curl "http://localhost:8080/games?players=2&genre=RPG&platform=PC"
```

**回應範例：**

```json
[
  {
    "name": "Game Title",
    "min_players": 1,
    "max_players": 4,
    "platforms": ["PC", "PS5"],
    "genre": "RPG",
    "online": true
  }
]
```

## 授權

本專案採用 GPL-3.0 license，詳見 `LICENSE` 文件。

## 貢獻

歡迎提交 Issue 或 Pull Request 來改進此專案！

