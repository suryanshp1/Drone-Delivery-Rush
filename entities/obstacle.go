package entities

import (
	"image/color"

	"drone-delivery-rush/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ObstacleType int

const (
	Building ObstacleType = iota
	Antenna
	PowerLine
	PackagePickup
	DeliveryZone
)

type Obstacle struct {
	Type   ObstacleType
	X, Y   float64
	W, H   float64
	Color  color.RGBA
	Active bool
}

func NewBuilding(x, y, w, h float64) Obstacle {
	return Obstacle{
		Type:   Building,
		X:      x,
		Y:      y,
		W:      w,
		H:      h,
		Color:  color.RGBA{20, 20, 30, 255}, // Dark outline building
		Active: true,
	}
}

func NewAntenna(x, y, h float64) Obstacle {
	return Obstacle{
		Type:   Antenna,
		X:      x,
		Y:      y,
		W:      5,
		H:      h,
		Color:  color.RGBA{100, 100, 100, 255},
		Active: true,
	}
}

func NewPackagePickup(x, y float64) Obstacle {
	return Obstacle{
		Type:   PackagePickup,
		X:      x,
		Y:      y,
		W:      20,
		H:      20,
		Color:  color.RGBA{255, 200, 0, 255},
		Active: true,
	}
}

func NewDeliveryZone(x, y float64) Obstacle {
	return Obstacle{
		Type:   DeliveryZone,
		X:      x,
		Y:      y,
		W:      60,
		H:      10,
		Color:  color.RGBA{0, 255, 0, 150},
		Active: true,
	}
}

func (o *Obstacle) Draw(screen *ebiten.Image, camX, camY float64) {
	if !o.Active {
		return
	}
	screenX := o.X - camX
	screenY := o.Y - camY

	if o.Type == PackagePickup && assets.PickupImg != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 0.08
		op.GeoM.Scale(scale, scale)
		op.Blend = ebiten.BlendLighter
		w, _ := assets.PickupImg.Size()
		op.GeoM.Translate(screenX+o.W/2-(float64(w)*scale/2), screenY-20)
		screen.DrawImage(assets.PickupImg, op)
		return
	}

	if o.Type == DeliveryZone && assets.DeliveryImg != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 0.08
		op.GeoM.Scale(scale, scale)
		op.Blend = ebiten.BlendLighter
		w, _ := assets.DeliveryImg.Size()
		op.GeoM.Translate(screenX+o.W/2-(float64(w)*scale/2), screenY-20)
		screen.DrawImage(assets.DeliveryImg, op)
		return
	}

	ebitenutil.DrawRect(screen, screenX, screenY, o.W, o.H, o.Color)
	
	// Details
	if o.Type == Building {
		// Draw some windows
		for wy := screenY + 10; wy < screenY+o.H-10; wy += 20 {
			for wx := screenX + 10; wx < screenX+o.W-10; wx += 20 {
				ebitenutil.DrawRect(screen, wx, wy, 10, 10, color.RGBA{255, 255, 150, 100})
			}
		}
	}
}

// Collides returns true if the given rect overlaps with this obstacle
func (o *Obstacle) Collides(dx, dy, dw, dh float64) bool {
	if !o.Active {
		return false
	}
	return dx < o.X+o.W && dx+dw > o.X && dy < o.Y+o.H && dy+dh > o.Y
}
