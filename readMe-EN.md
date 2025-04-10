# find-game-we-can-play-API

This is a Go-based game search API that allows users to filter available games based on the number of players, game genre, and platform.

## Features

- Filter games by the number of players
- Support for genre and platform filtering
- Returns only online-supported games
- Data is sourced from `games.json`

## Installation & Running

### 1. Clone the repository

```sh
git clone https://github.com/your-repo/game-api.git
cd game-api
```

### 2. Prepare `games.json`

Ensure `games.json` is in the same directory as the program, formatted as follows:

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

### 3. Run the API server

```sh
go run main.go
```

The API server will be available at `http://localhost:8080`

## API Usage

### Retrieve Available Games

**Request Method:** `GET /games`

**Query Parameters:**

| Parameter   | Type   | Description              | Required |
|------------|--------|--------------------------|----------|
| `players`  | int    | Number of players        | ✅       |
| `genre`    | string | Game genre (optional)    | ❌       |
| `platform` | string | Game platform (optional) | ❌       |

**Example Request:**

```sh
curl "http://localhost:8080/games?players=2&genre=RPG&platform=PC"
```

**Example Response:**

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

## License

This project is licensed under the GPL-3.0 license. See the `LICENSE` file for details.

## Contributing

Feel free to submit Issues or Pull Requests to improve this project!

