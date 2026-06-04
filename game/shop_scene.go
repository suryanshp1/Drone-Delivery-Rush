package game

import (
	"fmt"
	"image/color"

	"drone-delivery-rush/systems"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ShopScene struct {
	game *Game
}

func NewShopScene(g *Game) *ShopScene {
	return &ShopScene{game: g}
}

func (s *ShopScene) Init() {
}

func (s *ShopScene) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.game.SwitchScene(NewTitleScene(s.game))
	}

	cost := 500

	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		if s.game.Save.Coins >= cost {
			s.game.Save.Coins -= cost
			s.game.Save.Upgrades["battery"]++
			systems.SaveGame(s.game.Save)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		if s.game.Save.Coins >= cost {
			s.game.Save.Coins -= cost
			s.game.Save.Upgrades["motor"]++
			systems.SaveGame(s.game.Save)
		}
	}

	return nil
}

func (s *ShopScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 40, 20, 255}) // Dark green for shop

	batLvl := s.game.Save.Upgrades["battery"]
	motLvl := s.game.Save.Upgrades["motor"]

	msg := fmt.Sprintf(`GARAGE (UPGRADES)

COINS: %d

[B] Upgrade Battery (Lvl %d) - Cost: 500
[M] Upgrade Motors (Lvl %d) - Cost: 500


[ESC] Back to Title`, s.game.Save.Coins, batLvl, motLvl)

	ebitenutil.DebugPrintAt(screen, msg, s.game.screenWidth/2-100, s.game.screenHeight/2-100)
}
