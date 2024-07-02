/**
  *
  *
  *
  *
**/

package scene

import (
	"github.com/aiur-adept/sameriver/v4"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type LoadingScene struct {
	// used to make destroy() idempotent
	destroyed bool
	// used to make Init() idempotent
	initialized bool
	// the game
	game *sameriver.Game

	message_font    *ttf.Font
	message_surface *sdl.Surface
	message_texture *sdl.Texture

	// time accumulator for bouncing the word "loading"
	accum_5000 sameriver.TimeAccumulator
}
