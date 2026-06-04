package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene interface that all game screens must implement
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	Init()
}
