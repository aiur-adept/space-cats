package scene

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/aiur-adept/sameriver/v4"

	"github.com/aiur-adept/space-cats/game/systems"
)

func (s *GameScene) buildWorld() {
	// construct world object
	s.w = sameriver.NewWorld(map[string]any{
		"Width":  s.game.WindowSpec.Width,
		"Height": s.game.WindowSpec.Height,
	})

	// register components must always be called before AddSystems()
	// since systems might want to create and listen on component bitarray
	// filters

	// add systems
	s.w.RegisterSystems(
		sameriver.NewPhysicsSystem(),
		sameriver.NewSpatialHashSystem(32, 32),
		sameriver.NewCollisionSystem(100*time.Millisecond),
		systems.NewCoinDespawnAtEdgeSystem(),
	)
	s.w.SetSystemSchedule("CollisionSystem", 16)
	// get updated entity list of coins
	s.coins = s.w.GetUpdatedEntityList(sameriver.NewEntityFilter("coin", func(e *sameriver.Entity) bool {
		return e.GetTagList(sameriver.GENERICTAGS).Has("coin")
	}))
	// add spawn random coin logic
	const COINS_PER_SEC = 50
	s.w.AddWorldLogicWithSchedule("spawn-random-coin", s.spawnRandomCoin, 1000/COINS_PER_SEC)
	// add player coin collision logic
	s.w.AddWorldLogic("player-collect-coin", s.playerCollectCoin)
}

func (s *GameScene) spawnInitialEntities() {
	mass := 1.0
	s.player = s.w.Spawn(map[string]any{
		"components": map[sameriver.ComponentID]any{
			sameriver.POSITION: sameriver.Vec2D{50, 50},
			sameriver.VELOCITY: sameriver.Vec2D{0, 0},
			sameriver.BOX:      sameriver.Vec2D{2, 2},
			sameriver.MASS:     mass,
		},
		"tags": []string{"player"},
	})
}

func (s *GameScene) SimpleEntityDraw(
	r *sdl.Renderer, e *sameriver.Entity, c sdl.Color) {

	box := e.GetVec2D(sameriver.BOX)
	pos := e.GetVec2D(sameriver.POSITION).ShiftedCenterToBottomLeft(*box)
	r.SetDrawColor(c.R, c.G, c.B, c.A)
	s.game.Screen.FillRect(r, &pos, box)
}

func (s *GameScene) playerHandleKeyboardState(kb []uint8) {
	v := s.player.GetVec2D(sameriver.VELOCITY)
	// get player v1
	v.X = 0.2 * float64(
		int8(kb[sdl.SCANCODE_D]|kb[sdl.SCANCODE_RIGHT])-
			int8(kb[sdl.SCANCODE_A]|kb[sdl.SCANCODE_LEFT]))
	v.Y = 0.2 * float64(
		int8(kb[sdl.SCANCODE_W]|kb[sdl.SCANCODE_UP])-
			int8(kb[sdl.SCANCODE_S]|kb[sdl.SCANCODE_DOWN]))
}

func (s *GameScene) updateScoreTexture() {
	if s.scoreSurface != nil {
		s.scoreSurface.Free()
	}
	if s.scoreTexture != nil {
		s.scoreTexture.Destroy()
	}
	// render message ("press space") surface
	score_msg := fmt.Sprintf("%d", s.score)
	var err error
	s.scoreSurface, err = s.UIFont.RenderUTF8Solid(
		score_msg,
		sdl.Color{255, 255, 255, 255})
	if err != nil {
		panic(err)
	}
	// create the texture
	s.scoreTexture, err = s.game.Renderer.CreateTextureFromSurface(s.scoreSurface)
	if err != nil {
		panic(err)
	}
	// set the width of the texture on screen
	w, h, err := s.UIFont.SizeUTF8(score_msg)
	if err != nil {
		panic(err)
	}
	s.scoreRect = sdl.Rect{10, 10, int32(w), int32(h)}
}

func (s *GameScene) spawnRandomCoin(dt_ms float64) {
	if s.coins.Length() < 1000 {
		mass := 1.0
		c := s.w.Spawn(map[string]any{
			"tags": []string{"coin"},
			"components": map[sameriver.ComponentID]any{
				sameriver.POSITION: sameriver.Vec2D{
					rand.Float64()*float64(s.w.Width/3) + float64(s.w.Width/3),
					rand.Float64()*float64(s.w.Height/3) + float64(s.w.Height/3),
				},
				sameriver.VELOCITY: sameriver.Vec2D{0, 0},
				sameriver.BOX:      sameriver.Vec2D{4, 4},
				sameriver.MASS:     mass,
			},
		})
		c.AddLogic("coin-logic", s.coinLogic(c))
	}
}

func (s *GameScene) coinLogic(c *sameriver.Entity) func(e *sameriver.Entity, dt_ms float64) {
	return func(e *sameriver.Entity, dt_ms float64) {
		dist := c.GetVec2D(sameriver.POSITION).Sub(*s.player.GetVec2D(sameriver.POSITION))
		*c.GetVec2D(sameriver.VELOCITY) = dist.Unit().Scale(0.1 * (1.0 - dist.Magnitude()/float64(s.w.Width)))
	}
}

func (s *GameScene) playerCollectCoin(dt_ms float64) {
	if s.playerCoinCollision == nil {
		s.subscribeToPlayerCoinCollision()
	}
	for len(s.playerCoinCollision.C) > 0 {
		e := <-s.playerCoinCollision.C
		coin := e.Data.(sameriver.CollisionData).Other
		s.w.Despawn(coin)
		s.augmentScore(10)
		s.growPlayer(0.5)
	}
}

func (s *GameScene) subscribeToPlayerCoinCollision() {
	s.playerCoinCollision = s.w.Events.Subscribe(
		sameriver.PredicateEventFilter(
			"coin-collision",
			func(e sameriver.Event) bool {
				c := e.Data.(sameriver.CollisionData)
				return c.This == s.player &&
					c.Other.GetTagList(sameriver.GENERICTAGS).Has("coin")
			}))
}

func (s *GameScene) augmentScore(x int) {
	s.score += x
	s.updateScoreTexture()
}

func (s *GameScene) growPlayer(increase float64) {
	playerBox := s.player.GetVec2D(sameriver.BOX)
	if playerBox.X < 50 && playerBox.Y < 50 {
		playerBox.X += increase
		playerBox.Y += increase
	}
}
