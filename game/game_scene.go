package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"drone-delivery-rush/assets"
	"drone-delivery-rush/entities"
	"drone-delivery-rush/systems"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameScene struct {
	game      *Game
	drone     *entities.Drone
	cameraX   float64
	cameraY   float64
	windX     float64
	windY     float64
	score     int
	obstacles []entities.Obstacle
	spawner   *systems.Spawner
	gameOver  bool
}

func NewGameScene(g *Game) *GameScene {
	return &GameScene{
		game: g,
	}
}

func (s *GameScene) Init() {
	// Start drone on ground at the right side of chunk 0
	s.drone = entities.NewDrone(900.0, 390.0)
	s.cameraX = 0
	s.cameraY = 0
	
	// Initial wind
	s.windX = (rand.Float64() - 0.5) * 0.1
	s.windY = 0

	s.spawner = systems.NewSpawner()
	s.obstacles = []entities.Obstacle{}
	s.gameOver = false
}

func (s *GameScene) Update() error {
	if s.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			s.game.SwitchScene(NewGameScene(s.game)) // Restart
		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			s.game.SwitchScene(NewTitleScene(s.game))
		}
		return nil
	}

	// Random turbulence
	if rand.Float64() < 0.05 {
		s.windX += (rand.Float64() - 0.5) * 0.05
	}
	
	s.drone.Update(s.windX, s.windY)

	// Check chunk generation
	chunkIdx := int(math.Floor(s.drone.X / s.spawner.ChunkWidth))
	// Generate current, previous, and next chunks to be safe
	for i := chunkIdx - 1; i <= chunkIdx+1; i++ {
		if newObs := s.spawner.GenerateChunk(i); newObs != nil {
			s.obstacles = append(s.obstacles, newObs...)
		}
	}

	// Collision detection
	// Standard hitbox for pickups
	droneRectX := s.drone.X - 15
	droneRectY := s.drone.Y - 10
	droneRectW := 30.0
	droneRectH := 20.0

	// Forgiving hitbox for deadly obstacles
	deadlyRectX := s.drone.X - 8
	deadlyRectY := s.drone.Y - 5
	deadlyRectW := 16.0
	deadlyRectH := 10.0

	for i := range s.obstacles {
		obs := &s.obstacles[i]
		
		// Check pad interactions first (standard size)
		if obs.Type == entities.PackagePickup || obs.Type == entities.DeliveryZone {
			if obs.Collides(droneRectX, droneRectY, droneRectW, droneRectH) {
				if obs.Type == entities.PackagePickup {
					if s.drone.CurrentPackage == nil {
						s.drone.CurrentPackage = &entities.Package{
							Weight:    5.0 + rand.Float64()*15.0, // Random weight 5-20kg
							Value:     100,
							Integrity: 100.0,
						}
						obs.Active = false // Consume pickup
					}
				} else if obs.Type == entities.DeliveryZone {
					if s.drone.CurrentPackage != nil {
						s.score += s.drone.CurrentPackage.Value
						s.drone.CurrentPackage = nil
						
						// Recharge battery on successful delivery
						s.drone.Battery += 30.0
						if s.drone.Battery > s.drone.MaxBattery {
							s.drone.Battery = s.drone.MaxBattery
						}
					}
				}
			}
		} else {
			// Check deadly obstacles (forgiving size)
			if obs.Collides(deadlyRectX, deadlyRectY, deadlyRectW, deadlyRectH) {
				s.gameOver = true // Crash!
			}
		}
	}

	// Ceiling limit
	if s.drone.Y < -200 {
		s.drone.Y = -200
		if s.drone.VY < 0 {
			s.drone.VY = 0
		}
	}

	// Right Fence (Prevent flying backwards to the right infinitely)
	if s.drone.X > 950 {
		s.drone.X = 950
		if s.drone.VX > 0 {
			s.drone.VX = 0
		}
	}

	// Floor collision
	if s.drone.Y > 400 {
		s.drone.Y = 400
		s.drone.VY = 0
		// Optional: crash on hard landing
		if s.drone.VY > 5.0 {
			s.gameOver = true
		}
	}

	// Update camera to follow drone, keeping it centered
	targetCamX := s.drone.X - float64(s.game.screenWidth)/2.0
	targetCamY := s.drone.Y - float64(s.game.screenHeight)/2.0

	// Smooth camera follow
	s.cameraX += (targetCamX - s.cameraX) * 0.1
	s.cameraY += (targetCamY - s.cameraY) * 0.1

	// Return to title on ESC
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.game.SwitchScene(NewTitleScene(s.game))
	}

	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	// Background color (sky) fallback
	screen.Fill(color.RGBA{20, 20, 40, 255}) // Darker sky for cyberpunk

	if assets.BgImage != nil {
		bgW, _ := assets.BgImage.Size()
		// Parallax scrolling
		parallaxX := s.cameraX * 0.2
		offsetX := math.Mod(parallaxX, float64(bgW))
		
		op1 := &ebiten.DrawImageOptions{}
		op1.GeoM.Translate(-offsetX, 0)
		screen.DrawImage(assets.BgImage, op1)
		
		// Draw second instance for seamless looping
		if offsetX > 0 {
			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Translate(float64(bgW)-offsetX, 0)
			screen.DrawImage(assets.BgImage, op2)
		} else if offsetX < 0 {
			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Translate(-float64(bgW)-offsetX, 0)
			screen.DrawImage(assets.BgImage, op2)
		}
	}

	// Draw Ground
	groundY := 400.0 - s.cameraY
	if groundY < float64(s.game.screenHeight) {
		ebitenutil.DrawRect(screen, 0, groundY, float64(s.game.screenWidth), float64(s.game.screenHeight)-groundY, color.RGBA{30, 30, 40, 255}) // Cyberpunk ground
	}

	// Draw Obstacles
	for _, obs := range s.obstacles {
		obs.Draw(screen, s.cameraX, s.cameraY)
	}

	// Draw Tutorial in Chunk 0
	if s.drone.X > 0 && s.drone.X < 1000 {
		ebitenutil.DebugPrintAt(screen, "HOW TO PLAY:\n[W] Thrust Up\n[A] Fly Left (Forward)\n[D] Fly Right (Brake)\n[SHIFT] Boost (Drains Battery Fast)\n\nDeliver packages to recharge battery!", int(900-s.cameraX-100), int(390-s.cameraY-150))
	}

	// Draw Drone
	s.drone.Draw(screen, s.cameraX, s.cameraY)

	// HUD
	s.drawHUD(screen)

	if s.gameOver {
		ebitenutil.DebugPrintAt(screen, "CRASHED!\nPress ENTER to Restart", s.game.screenWidth/2-50, s.game.screenHeight/2)
	}
}

func (s *GameScene) drawHUD(screen *ebiten.Image) {
	// Battery bar
	ebitenutil.DrawRect(screen, 10, 10, 200, 20, color.RGBA{50, 50, 50, 255})
	
	batRatio := s.drone.Battery / s.drone.MaxBattery
	batColor := color.RGBA{0, 255, 0, 255}
	if batRatio < 0.2 {
		batColor = color.RGBA{255, 0, 0, 255}
	} else if batRatio < 0.5 {
		batColor = color.RGBA{255, 255, 0, 255}
	}
	ebitenutil.DrawRect(screen, 10, 10, 200*batRatio, 20, batColor)
	
	// HUD Text
	weightStr := "None"
	if s.drone.CurrentPackage != nil {
		weightStr = fmt.Sprintf("%.1f kg", s.drone.CurrentPackage.Weight)
	}

	hudText := fmt.Sprintf("Battery: %.0f%%\nSpeed: %.1f\nWind: %.2f\nPackage: %s\nScore: %d\nCoins: %d", 
		s.drone.Battery, s.drone.VX, s.windX, weightStr, s.score, s.score) // using score as coins for now
	ebitenutil.DebugPrintAt(screen, hudText, 10, 40)
}
