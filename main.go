package main

import (
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

const (
	tileSize               = 16
	visibleVerticalTiles   = 18
	visibleHorizontalTiles = 32
)

type playerDirection int

const (
	playerDirectionDown playerDirection = iota
	playerDirectionUp
	playerDirectionLeft
	playerDirectionRight
)

type player struct {
	x, y      float64
	direction playerDirection
	isMoving  bool
}

var (
	mapImage        *ebiten.Image
	characterSprite *ebiten.Image

	characterAnimationCurrentFrame = 0
)

const (
	characterSpriteSize    = 16
	characterSpriteSpacing = 32
	characterSpritePadding = 16

	characterSpeed                = 1
	characterAnimationFramesCount = 2
)

type Game struct {
	player  player
	counter int
}

func (g *Game) Update() error {
	dx, dy := 0.0, 0.0

	g.player.isMoving = false
	g.player.direction = playerDirectionDown

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx -= 1
		g.player.direction = playerDirectionLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx += 1
		g.player.direction = playerDirectionRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy -= 1
		g.player.direction = playerDirectionUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy += 1
		g.player.direction = playerDirectionDown
	}

	if dx != 0 || dy != 0 {
		g.player.isMoving = true

		length := math.Sqrt(dx*dx + dy*dy)
		g.player.x += (dx / length) * characterSpeed
		g.player.y += (dy / length) * characterSpeed
	}

	g.counter++

	characterAnimationCurrentFrame = (g.counter / 10) % characterAnimationFramesCount
	if characterAnimationCurrentFrame >= characterAnimationFramesCount {
		characterAnimationCurrentFrame = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(mapImage, nil)

	charOp := &ebiten.DrawImageOptions{}
	charOp.GeoM.Translate(float64(g.player.x), float64(g.player.y))

	currentMovOffset := ((characterSpriteSize + characterSpriteSpacing) * int(g.player.direction))

	currentFrameOffset := (characterSpriteSize + characterSpriteSpacing) * characterAnimationCurrentFrame
	if g.player.isMoving {
		// move frames are last two frames of spritesheet and idle frames are first two ones
		currentFrameOffset += (characterSpriteSize + characterSpriteSpacing) * 2
	}

	screen.DrawImage(characterSprite.SubImage(image.Rect(
		characterSpritePadding+currentFrameOffset,
		characterSpritePadding+currentMovOffset,
		characterSpritePadding+characterSpriteSize+currentFrameOffset,
		characterSpritePadding+characterSpriteSize+currentMovOffset,
	)).(*ebiten.Image), charOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return tileSize * visibleHorizontalTiles, tileSize * visibleVerticalTiles
}

func main() {
	gameMap, err := tiled.LoadFile("maps/map.tmx")
	if err != nil {
		log.Fatal(err)
	}

	renderer, err := render.NewRenderer(gameMap)
	if err != nil {
		log.Fatal(err)
	}

	err = renderer.RenderVisibleLayersAndObjectGroups()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(true)

	mapImage = ebiten.NewImageFromImage(renderer.Result)
	characterSprite, _, err = ebitenutil.NewImageFromFile("sprites/character.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&Game{
		player: player{
			x:         tileSize * 8,
			y:         tileSize * 8,
			direction: playerDirectionDown,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
