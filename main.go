package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Game struct {
	Name       string   `json:"name"`
	MinPlayers int      `json:"min_players"`
	MaxPlayers int      `json:"max_players"`
	Platforms  []string `json:"platforms"`
	Genre      string   `json:"genre"`
	Online     bool     `json:"online"`
}

var games []Game

func loadGames() error {
	data, err := ioutil.ReadFile("games.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &games)
	if err != nil {
		return err
	}
	log.Println("Loaded", len(games), "games from games.json")
	return nil
}

type WSMessage struct {
	Type     string `json:"type"`
	Players  int    `json:"players"`
	Genre    string `json:"genre"`
	Platform string `json:"platform"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade failed:", err)
		return
	}
	defer conn.Close()
	log.Println("Client connected via WebSocket")

	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		switch msg.Type {
		case "search":
			results := filterGames(msg.Players, msg.Genre, msg.Platform)
			conn.WriteJSON(gin.H{
				"type":    "search",
				"results": results,
			})
		case "random":
			results := filterGames(msg.Players, msg.Genre, msg.Platform)
			if len(results) == 0 {
				conn.WriteJSON(gin.H{
					"type":  "random",
					"error": "No games found",
				})
			} else {
				randomGame := results[rand.Intn(len(results))]
				conn.WriteJSON(gin.H{
					"type": "random",
					"game": randomGame,
				})
			}
		default:
			conn.WriteJSON(gin.H{
				"error": "Unknown message type",
			})
		}
	}
}

func filterGames(players int, genre, platform string) []Game {
	genre = strings.ToLower(genre)
	platform = strings.ToLower(platform)

	var filtered []Game
	for _, game := range games {
		if !game.Online || players < game.MinPlayers || players > game.MaxPlayers {
			continue
		}
		if genre != "" && strings.ToLower(game.Genre) != genre {
			continue
		}
		if platform != "" {
			match := false
			for _, p := range game.Platforms {
				if strings.ToLower(p) == platform {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		filtered = append(filtered, game)
	}
	return filtered
}

func getRandomGame(c *gin.Context) {
	players, err := strconv.Atoi(c.Query("players"))
	if err != nil || players <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player count"})
		return
	}

	genreFilter := strings.ToLower(c.Query("genre"))
	platformFilter := strings.ToLower(c.Query("platform"))

	filteredGames := filterGames(players, genreFilter, platformFilter)

	if len(filteredGames) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No games found"})
		return
	}

	randomGame := filteredGames[rand.Intn(len(filteredGames))]
	c.JSON(http.StatusOK, randomGame)
}

func getGamesByPlayerCount(c *gin.Context) {
	players, err := strconv.Atoi(c.Query("players"))
	if err != nil || players <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player count"})
		return
	}

	genreFilter := strings.ToLower(c.Query("genre"))
	platformFilter := strings.ToLower(c.Query("platform"))

	filteredGames := filterGames(players, genreFilter, platformFilter)

	c.JSON(http.StatusOK, filteredGames)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := loadGames(); err != nil {
		log.Fatal("Failed to load games.json: ", err)
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	r.GET("/ws", serveWebSocket)

	r.GET("/games", getGamesByPlayerCount)
	r.GET("/randomGame", getRandomGame)

	log.Println("API Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
