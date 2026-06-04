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
	
	titleText := "DRONE DELIVERY RUSH\n\nPress ENTER to Start"
	controlsText := `HOW TO PLAY:
- W/A/S/D: Move Drone
- SHIFT: Boost (Uses more battery!)

Pick up YELLOW packages.
Drop them at GREEN delivery zones.
Avoid Buildings and Antennas!
Heavier packages change your flight physics.`

	ebitenutil.DebugPrintAt(screen, titleText, 10, 10)
	ebitenutil.DebugPrintAt(screen, controlsText, 10, 80)
}
