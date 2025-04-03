package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func getGamesByPlayerCount(w http.ResponseWriter, r *http.Request) {
	players, err := strconv.Atoi(r.URL.Query().Get("players"))
	if err != nil || players <= 0 {
		http.Error(w, "Invalid player count", http.StatusBadRequest)
		return
	}

	genreFilter := strings.ToLower(r.URL.Query().Get("genre"))
	platformFilter := strings.ToLower(r.URL.Query().Get("platform"))

	var filteredGames []Game
	for _, game := range games {
		if !game.Online || players < game.MinPlayers || players > game.MaxPlayers {
			continue
		}
		if genreFilter != "" && strings.ToLower(game.Genre) != genreFilter {
			continue
		}
		if platformFilter != "" {
			match := false
			for _, p := range game.Platforms {
				if strings.ToLower(p) == platformFilter {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		filteredGames = append(filteredGames, game)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredGames)
}

func main() {
	if err := loadGames(); err != nil {
		log.Fatal("Failed to load games.json: ", err)
	}
	log.Println("API Server running at http://localhost:8080")
	http.HandleFunc("/games", getGamesByPlayerCount)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
