package systems

import (
	"github.com/aiur-adept/sameriver/v7"
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
		return s.w.GetTagList(e, sameriver.GENERICTAGS_).Has("coin")
	}))
}

func (s *CoinDespawnAtEdgeSystem) Update(dt_ms float64) {
	// despawn at up/down
	for y := 0; y <= s.sh.Hasher.GridY-1; y += (s.sh.Hasher.GridY - 1) {
		for x := 0; x < s.sh.Hasher.GridX; x++ {
			cell := s.sh.Hasher.Table[x][y]
			for _, e := range cell {
				if s.w.GetTagList(e, sameriver.GENERICTAGS_).Has("coin") {
					pos := s.w.GetVec2D(e, sameriver.POSITION_)
					box := s.w.GetVec2D(e, sameriver.BOX_)
					if pos.Y < box.Y || (s.w.Height-pos.Y) < box.Y {
						s.w.Despawn(e)
					}
				}
			}
		}
	}

	// despawn at L/R
	for x := 0; x <= s.sh.Hasher.GridX-1; x += (s.sh.Hasher.GridX - 1) {
		for y := 0; y < s.sh.Hasher.GridY; y++ {
			cell := s.sh.Hasher.Table[x][y]
			for _, e := range cell {
				if s.w.GetTagList(e, sameriver.GENERICTAGS_).Has("coin") {
					pos := s.w.GetVec2D(e, sameriver.POSITION_)
					box := s.w.GetVec2D(e, sameriver.BOX_)
					if pos.X < box.X || (s.w.Width-pos.X) < box.X {
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

func (s *CoinDespawnAtEdgeSystem) GetComponentDeps() []any {
	return []any{
		sameriver.POSITION_, sameriver.VEC2D, "POSITION",
		sameriver.BOX_, sameriver.VEC2D, "BOX",
	}
}
