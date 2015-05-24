package main

import (
	"encoding/json"
	"net/http"

	"mse"
)

var game *mse.Game

func apiNewGame(w http.ResponseWriter, r *http.Request) {
	game = mse.NewGame()
	go game.Run()
}

func apiGetStatus(w http.ResponseWriter, r *http.Request) {
	s := <-game.NextStatus
	resp := mse.StatusResponse{End: s == nil}
	if s != nil {
		resp.Status = *s
	}
	if b, err := json.Marshal(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func apiGetBoard(w http.ResponseWriter, r *http.Request) {
	<-game.Ready

	if b, err := json.Marshal(game.GetBoard()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func apiGetPrompt(w http.ResponseWriter, r *http.Request) {
	p := <-game.NextPrompt
	resp := mse.PromptResponse{End: p == nil}
	if p != nil {
		resp.Prompt = *p
	}
	if b, err := json.Marshal(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func main() {
	http.HandleFunc("/api/newGame", apiNewGame)
	http.HandleFunc("/api/board", apiGetBoard)
	http.HandleFunc("/api/status", apiGetStatus)
	http.HandleFunc("/api/prompt", apiGetPrompt)
	http.Handle("/", http.FileServer(http.Dir("./..")))

	http.ListenAndServe(":8080", nil)
}
