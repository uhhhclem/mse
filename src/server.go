package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"mse"
)

var games map[string]*mse.Game


func apiNewGame(w http.ResponseWriter, r *http.Request) {
	g := mse.NewGame()
	go g.Run()
	
	if games == nil {
		games = make(map[string]*mse.Game)
	}
	games[g.ID] = g

	resp := struct{
		ID string
	}{
		ID:g.ID,
	}

	if b, err := json.Marshal(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

type apiHandler func(*mse.Game, http.ResponseWriter, *http.Request)

func apiWrapper(h apiHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.URL)
		id := r.FormValue("ID")
		if game, ok := games[id]; !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Game ID %s not found.", id)))
		} else {
			h(game, w, r)
		}
	}
}

func apiGetStatus(game *mse.Game, w http.ResponseWriter, r *http.Request) {
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

func apiGetBoard(game *mse.Game, w http.ResponseWriter, r *http.Request) {
	<-game.Ready

	if b, err := json.Marshal(game.GetBoard()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func apiGetPrompt(game *mse.Game, w http.ResponseWriter, r *http.Request) {
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

func apiPostChoice(w http.ResponseWriter, r *http.Request) {
	req := struct{ 
		ID string
		Key string 
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	game, ok := games[req.ID]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Game %s not found.", req.ID)))
	} else if err := game.MakeChoice(req.Key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	http.HandleFunc("/api/newGame", apiNewGame)
	http.HandleFunc("/api/choice", apiPostChoice)

	handlers := []struct{
		url string
		handler apiHandler
	}{
		{"/api/board", apiGetBoard},
		{"/api/status", apiGetStatus},
		{"/api/prompt", apiGetPrompt},
	}
	for _, h := range handlers {
		http.HandleFunc(h.url, apiWrapper(h.handler))
	}
		
	http.Handle("/", http.FileServer(http.Dir("./..")))

	http.ListenAndServe(":8080", nil)
}
