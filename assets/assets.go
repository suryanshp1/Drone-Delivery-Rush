package assets

import (
	"log"
	_ "image/png"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	DroneImage *ebiten.Image
	BgImage    *ebiten.Image
	PackageImg *ebiten.Image
	PickupImg  *ebiten.Image
	DeliveryImg *ebiten.Image
)

func LoadAssets() {
	var err error
	DroneImage, _, err = ebitenutil.NewImageFromFile("assets/drone.png")
	if err != nil {
		log.Fatal(err)
	}

	BgImage, _, err = ebitenutil.NewImageFromFile("assets/bg.png")
	if err != nil {
		log.Fatal(err)
	}

	PackageImg, _, err = ebitenutil.NewImageFromFile("assets/package.png")
	if err != nil {
		log.Println("Warning: package.png not found")
	}

	PickupImg, _, err = ebitenutil.NewImageFromFile("assets/pickup.png")
	if err != nil {
		log.Println("Warning: pickup.png not found")
	}

	DeliveryImg, _, err = ebitenutil.NewImageFromFile("assets/delivery.png")
	if err != nil {
		log.Println("Warning: delivery.png not found")
	}
}
