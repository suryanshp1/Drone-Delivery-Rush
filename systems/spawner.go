package systems

import (
	"math"
	"math/rand"

	"drone-delivery-rush/entities"
)

type Spawner struct {
	GeneratedChunks map[int]bool // Track which X-chunks have been generated
	ChunkWidth      float64
}

func NewSpawner() *Spawner {
	return &Spawner{
		GeneratedChunks: make(map[int]bool),
		ChunkWidth:      1000,
	}
}

// GenerateChunk generates obstacles for a given chunk index if not already generated.
// X ranges from chunkIdx*ChunkWidth to (chunkIdx+1)*ChunkWidth
func (s *Spawner) GenerateChunk(chunkIdx int) []entities.Obstacle {
	if s.GeneratedChunks[chunkIdx] {
		return nil
	}
	s.GeneratedChunks[chunkIdx] = true

	var obs []entities.Obstacle
	startX := float64(chunkIdx) * s.ChunkWidth

	// Ground Level
	groundY := 400.0

	// Safe starting zone for chunk 0
	if chunkIdx == 0 {
		// Just spawn a single safe starting pad on the ground and no buildings
		obs = append(obs, entities.NewDeliveryZone(startX+s.ChunkWidth/2, groundY-10))
		return obs
	}

	absChunk := int(math.Abs(float64(chunkIdx)))

	// Difficulty Scaling
	// Base buildings: 2 to 5. Scales up to 8 max around chunk 30.
	maxBuildings := 5 + (absChunk / 5)
	if maxBuildings > 10 {
		maxBuildings = 10
	}
	numBuildings := rand.Intn(maxBuildings-2) + 2

	for i := 0; i < numBuildings; i++ {
		// Gap calculation based on numBuildings
		w := 80.0 + rand.Float64()*(200.0-float64(absChunk)) // Buildings get slightly narrower
		if w < 50 {
			w = 50
		}
		
		// Height scales up slightly
		baseHeight := 100.0 + float64(absChunk)*2.0
		h := baseHeight + rand.Float64()*200.0
		if h > 350 {
			h = 350 // Leave room at top
		}

		x := startX + (float64(i) * (s.ChunkWidth / float64(numBuildings))) + rand.Float64()*30.0
		y := groundY - h
		
		obs = append(obs, entities.NewBuilding(x, y, w, h))

		// Decide if this building gets a Pad
		hasPad := false
		if rand.Float64() < 0.25 {
			hasPad = true
			if rand.Float64() < 0.5 {
				obs = append(obs, entities.NewPackagePickup(x+w/2-10, y-30)) // Raised higher safely
			} else {
				obs = append(obs, entities.NewDeliveryZone(x+w/2-30, y-30))
			}
		}

		// Only add antennas if there is NO pad, preventing overlapping traps
		if !hasPad && rand.Float64() < 0.4 {
			ah := 30.0 + rand.Float64()*80.0
			ax := x + w/2 - 2.5
			ay := y - ah
			obs = append(obs, entities.NewAntenna(ax, ay, ah))
		}
	}

	return obs
}
