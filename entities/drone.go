package entities

import (
	"image/color"
	"math"

	"drone-delivery-rush/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Package struct {
	Weight    float64
	Type      string
	Value     int
	Integrity float64 // 0 to 100
}

type Drone struct {
	X, Y           float64
	VX, VY         float64
	BaseWeight     float64
	ThrustPower    float64
	Battery        float64
	MaxBattery     float64
	BoostMulti     float64
	CurrentPackage *Package
	
	// Stats from upgrades
	MotorLevel       int
	BatteryLevel     int
	AerodynamicLevel int
	ShieldLevel      int
}

func NewDrone(startX, startY float64) *Drone {
	return &Drone{
		X:           startX,
		Y:           startY,
		BaseWeight:  10.0,
		ThrustPower: 0.6,
		MaxBattery:  100.0,
		Battery:     100.0,
		BoostMulti:  1.8,
	}
}

func (d *Drone) TotalWeight() float64 {
	w := d.BaseWeight
	if d.CurrentPackage != nil {
		w += d.CurrentPackage.Weight
	}
	return w
}

func (d *Drone) Update(windX, windY float64) {
	// Base gravity
	gravity := 0.2
	
	thrust := d.ThrustPower
	isBoosting := ebiten.IsKeyPressed(ebiten.KeyShift)
	
	if isBoosting && d.Battery > 0 {
		thrust *= d.BoostMulti
		d.Battery -= 0.1 // Boost battery drain
	} else {
		d.Battery -= 0.01 // Normal idle drain
	}
	if d.Battery < 0 {
		d.Battery = 0
		thrust = 0 // Cannot fly without battery
	}

	// Calculate forces based on weight
	weightFactor := 10.0 / d.TotalWeight() // Heavier = less effect

	var ax, ay float64
	
	// Input mapping
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		ay -= thrust * weightFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		ay += thrust * weightFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		ax -= thrust * weightFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		ax += thrust * weightFactor
	}

	// Apply wind (affected by aerodynamics)
	aeroFactor := 1.0 - (float64(d.AerodynamicLevel) * 0.1)
	if aeroFactor < 0 {
		aeroFactor = 0
	}
	ax += (windX * weightFactor * aeroFactor)
	ay += (windY * weightFactor * aeroFactor)

	// Apply gravity
	ay += gravity * (d.TotalWeight() / 10.0) // Heavier falls slightly faster, though galileo says otherwise, game feel is different. Actually let's keep gravity constant and let thrust be the differentiating factor.
	// Actually standard physics: F = ma -> a = F/m. Gravity is a force mg, so a = g. Gravity acceleration is constant.
	ay -= gravity * (d.TotalWeight() / 10.0) // wait
	
	// Let's reset to standard physics
	ay = 0 // resetting for re-calc
	ay += gravity // Constant downward acceleration

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		ay -= thrust * weightFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		ay += thrust * weightFactor * 0.5 // Weaker thrust downwards
	}

	// Add acceleration to velocity
	d.VX += ax
	d.VY += ay

	// Apply drag (air resistance)
	drag := 0.95
	d.VX *= drag
	d.VY *= drag

	// Max velocity limits? (optional, drag usually handles this)
	maxV := 15.0
	if math.Abs(d.VX) > maxV {
		d.VX = math.Copysign(maxV, d.VX)
	}
	if math.Abs(d.VY) > maxV {
		d.VY = math.Copysign(maxV, d.VY)
	}

	// Update position
	d.X += d.VX
	d.Y += d.VY
}

func (d *Drone) Draw(screen *ebiten.Image, camX, camY float64) {
	screenX := d.X - camX
	screenY := d.Y - camY
	
	if assets.DroneImage != nil {
		op := &ebiten.DrawImageOptions{}
		
		w, h := assets.DroneImage.Size()
		scale := 0.08
		
		// Center origin, flip horizontally to face left, scale, then translate to screen pos
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Scale(-scale, scale)
		op.GeoM.Translate(screenX, screenY)
		
		// Use Lighter blending to make the solid black background transparent
		op.Blend = ebiten.BlendLighter
		
		screen.DrawImage(assets.DroneImage, op)
	} else {
		// Fallback simple rectangle
		ebitenutil.DrawRect(screen, screenX-15, screenY-10, 30, 20, color.RGBA{0, 200, 255, 255})
	}
	
	// Draw thrust particle/flame if W is pressed
	if ebiten.IsKeyPressed(ebiten.KeyW) && d.Battery > 0 {
		flameLength := 10.0
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			flameLength = 20.0
		}
		ebitenutil.DrawRect(screen, screenX-5, screenY+15, 10, flameLength, color.RGBA{255, 100, 0, 255})
	}
	
	// Draw package if carrying one
	if d.CurrentPackage != nil {
		if assets.PackageImg != nil {
			op := &ebiten.DrawImageOptions{}
			scaleP := 0.04
			op.GeoM.Scale(scaleP, scaleP)
			op.Blend = ebiten.BlendLighter
			pw, _ := assets.PackageImg.Size()
			op.GeoM.Translate(screenX-(float64(pw)*scaleP/2), screenY+15)
			screen.DrawImage(assets.PackageImg, op)
		} else {
			ebitenutil.DrawRect(screen, screenX-10, screenY+25, 20, 15, color.RGBA{200, 150, 50, 255})
		}
	}
}
