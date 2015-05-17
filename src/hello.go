package main

import (
    "fmt"
    
    "mse"
)

func main() {
    var card int
    for deck := mse.NearSystemDeck; len(deck) > 0; {
        card, deck = mse.Draw(deck)
        fmt.Println(*mse.Systems[card])
    }
}