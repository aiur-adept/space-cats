package systems

import (
	"github.com/aiur-adept/sameriver/v4"
)

type CoinDespawnAtEdgeSystem struct {
	w     *sameriver.World
	sh    *sameriver.SpatialHashSystem `sameriver-system-dependency:"-"`
	coins *sameriver.UpdatedEntityList
}

func NewCoinDespawnAtEdgeSystem() *CoinDespawnAtEdgeSystem {
	return &CoinDespawnAtEdgeSystem{}
}

func (s *CoinDespawnAtEdgeSystem) LinkWorld(w *sameriver.World) {
	s.w = w
	s.coins = s.w.GetUpdatedEntityList(sameriver.NewEntityFilter("coin", func(e *sameriver.Entity) bool {
		return e.GetTagList(sameriver.GENERICTAGS).Has("coin")
	}))
}

func (s *CoinDespawnAtEdgeSystem) Update(dt_ms float64) {
	for y := 0; y <= s.sh.Hasher.GridY-1; y += (s.sh.Hasher.GridY - 1) {
		for x := 0; x < s.sh.Hasher.GridX; x++ {
			cell := s.sh.Hasher.Table[x][y]
			for _, e := range cell {
				if e.GetTagList(sameriver.GENERICTAGS).Has("coin") {
					pos := e.GetVec2D(sameriver.POSITION)
					box := e.GetVec2D(sameriver.BOX)
					if pos.Y < box.Y || (s.w.Height-pos.Y) < box.Y {
						s.w.Despawn(e)
					}
				}
			}
		}
	}
}

func (s *CoinDespawnAtEdgeSystem) Expand(n int) {
	s.sh.Hasher.Expand(n)
}

func (s *CoinDespawnAtEdgeSystem) GetComponentDeps() []string {
	return []string{"TagList,GenericTags", "Vec2D,Position", "Vec2D,Box"}
}
