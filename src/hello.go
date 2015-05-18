package main

import (
	"encoding/json"
	"fmt"

	"mse"
)

func main() {
	g := mse.NewGame()

	b, err := json.Marshal(g.GetBoard())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
