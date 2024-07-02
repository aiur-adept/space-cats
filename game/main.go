package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aiur-adept/sameriver/v4"
	"github.com/aiur-adept/space-cats/game/scene"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("space cats, motherfucker")
	sameriver.RunGame(sameriver.GameInitSpec{
		WindowSpec: sameriver.WindowSpec{
			Title:      "space cats",
			Width:      800,
			Height:     800,
			Fullscreen: false},
		LoadingScene: &scene.LoadingScene{},
		FirstScene:   &scene.GameScene{},
	})
}
