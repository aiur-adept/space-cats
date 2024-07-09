package scene

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/aiur-adept/sameriver/v7"

	"github.com/aiur-adept/space-cats/game/systems"
)

func (s *GameScene) buildWorld() {
	// construct world object
	s.w = sameriver.NewWorld(map[string]any{
		"width":  s.game.WindowSpec.Width,
		"height": s.game.WindowSpec.Height,
	})

	// register components must always be called before AddSystems()
	// since systems might want to create and listen on component bitarray
	// filters

	// add systems
	s.w.RegisterSystems(
		sameriver.NewPhysicsSystem(),
		sameriver.NewSpatialHashSystem(32, 32),
		sameriver.NewCollisionSystem(10*time.Millisecond),
		systems.NewCoinDespawnAtEdgeSystem(),
	)
	s.w.SetSystemSchedule("CollisionSystem", 16)
	// get updated entity list of coins
	s.coins = s.w.GetUpdatedEntityListByTag("coin")
	// add spawn random coin logic
	const COINS_PER_SEC = 50
	s.w.AddLogicWithSchedule("spawn-random-coin", s.spawnRandomCoin, 1000/COINS_PER_SEC)
	// add player coin collision logic
	s.w.AddLogic("player-collect-coin", s.playerCollectCoin)
	// add coin move logic
	s.w.AddLogic("move-coins", s.moveCoins)
}

func (s *GameScene) spawnInitialEntities() {
	mass := 1.0
	s.player = s.w.Spawn(map[string]any{
		"components": map[sameriver.ComponentID]any{
			sameriver.POSITION_:     sameriver.Vec2D{100, 100},
			sameriver.VELOCITY_:     sameriver.Vec2D{0, 0},
			sameriver.ACCELERATION_: sameriver.Vec2D{0, 0},
			sameriver.BOX_:          sameriver.Vec2D{10, 10},
			sameriver.MASS_:         mass,
			sameriver.RIGIDBODY_:    true,
		},
		"tags": []string{"player"},
	})
}

func (s *GameScene) SimpleEntityDraw(
	r *sdl.Renderer, e *sameriver.Entity, c sdl.Color) {

	box := s.w.GetVec2D(e, sameriver.BOX_)
	pos := s.w.GetVec2D(e, sameriver.POSITION_).ShiftedCenterToBottomLeft(*box)
	r.SetDrawColor(c.R, c.G, c.B, c.A)
	s.game.Screen.FillRect(r, &pos, box)
}

func (s *GameScene) playerHandleKeyboardState(kb []uint8) {
	v := s.w.GetVec2D(s.player, sameriver.VELOCITY_)
	// get player v1
	v.X = 0.4 * float64(
		int8(kb[sdl.SCANCODE_D]|kb[sdl.SCANCODE_RIGHT])-
			int8(kb[sdl.SCANCODE_A]|kb[sdl.SCANCODE_LEFT]))
	v.Y = 0.4 * float64(
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
	if s.coins.Length() < 500 {
		mass := 1.0
		s.w.Spawn(map[string]any{
			"tags": []string{"coin"},
			"components": map[sameriver.ComponentID]any{
				sameriver.POSITION_: sameriver.Vec2D{
					rand.Float64()*float64(s.w.Width/3) + float64(s.w.Width/3),
					rand.Float64()*float64(s.w.Height/3) + float64(s.w.Height/3),
				},
				sameriver.VELOCITY_:     sameriver.Vec2D{0, 0},
				sameriver.ACCELERATION_: sameriver.Vec2D{0, 0},
				sameriver.BOX_:          sameriver.Vec2D{4, 4},
				sameriver.MASS_:         mass,
				sameriver.RIGIDBODY_:    false,
			},
		})
	}
}

func (s *GameScene) moveCoins(dt_ms float64) {
	for _, c := range s.coins.GetEntities() {
		dist := s.w.GetVec2D(c, sameriver.POSITION_).Sub(*s.w.GetVec2D(s.player, sameriver.POSITION_))
		*s.w.GetVec2D(c, sameriver.VELOCITY_) = dist.Unit().Scale(0.1 * (1.0 - dist.Magnitude()/float64(s.w.Width)))
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
		s.growPlayer(0.8)
	}
}

func (s *GameScene) subscribeToPlayerCoinCollision() {
	s.playerCoinCollision = s.w.Events.Subscribe(
		sameriver.PredicateEventFilter(
			"collision",
			func(e sameriver.Event) bool {
				c := e.Data.(sameriver.CollisionData)
				return c.This == s.player &&
					s.w.GetTagList(c.Other, sameriver.GENERICTAGS_).Has("coin")
			}))
}

func (s *GameScene) augmentScore(x int) {
	s.score += x
	s.updateScoreTexture()
}

func (s *GameScene) growPlayer(increase float64) {
	playerBox := s.w.GetVec2D(s.player, sameriver.BOX_)
	if playerBox.X < 80 && playerBox.Y < 80 {
		playerBox.X += increase
		playerBox.Y += increase
	}
}
