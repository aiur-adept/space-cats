package main

import (
	"fmt"

	"github.com/aiur-adept/sameriver/v7"
	"github.com/aiur-adept/space-cats/game/scene"
)

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
