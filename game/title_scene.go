package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TitleScene struct {
	game *Game
}

func NewTitleScene(g *Game) *TitleScene {
	return &TitleScene{game: g}
}

func (s *TitleScene) Init() {
	// Initialize scene variables
}

func (s *TitleScene) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.game.SwitchScene(NewGameScene(s.game))
	}
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 40, 255}) // Dark neo-tokyo purple/blue
	
	msg := "DRONE DELIVERY RUSH\n\nPress ENTER to Start\n\n\nHOW TO PLAY:\n[W] Thrust Up\n[A] Fly Left (Forward)\n[D] Fly Right (Brake)\n[SHIFT] Boost\n\nDeliver packages to recharge your battery!"
	ebitenutil.DebugPrintAt(screen, msg, s.game.screenWidth/2-100, s.game.screenHeight/2-50)
}
