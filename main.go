package main

import (
	"log"

	"drone-delivery-rush/assets"
	"drone-delivery-rush/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 540
)

func main() {
	assets.LoadAssets()

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Drone Delivery Rush")

	g := game.NewGame(ScreenWidth, ScreenHeight)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
