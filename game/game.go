package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	screenWidth  int
	screenHeight int
	currentScene Scene
	nextScene    Scene
}

func NewGame(width, height int) *Game {
	g := &Game{
		screenWidth:  width,
		screenHeight: height,
	}
	g.SwitchScene(NewTitleScene(g))
	return g
}

func (g *Game) SwitchScene(s Scene) {
	g.nextScene = s
}

func (g *Game) Update() error {
	if g.nextScene != nil {
		g.currentScene = g.nextScene
		g.currentScene.Init()
		g.nextScene = nil
	}
	if g.currentScene != nil {
		return g.currentScene.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}
