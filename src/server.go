package main

import (
	"encoding/json"
	"net/http"

	"mse"
)

var game *mse.Game

func apiNewGame(w http.ResponseWriter, r *http.Request) {
	game := mse.NewGame()

	if b, err := json.Marshal(game.GetBoard()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func main() {
	http.HandleFunc("/api/newGame", apiNewGame)
	http.Handle("/", http.FileServer(http.Dir("./..")))

	http.ListenAndServe(":8080", nil)
}
