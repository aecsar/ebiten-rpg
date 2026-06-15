package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/ldtkgo"
	renderer "github.com/solarlune/ldtkgo/renderer/ebitengine"
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
	characterSprite *ebiten.Image

	characterAnimationCurrentFrame = 0
)

const (
	characterSpriteSize    = 16
	characterSpriteSpacing = 32
	characterSpritePadding = 16

	characterSpeed                = 3
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
		g.player.x += math.Round((dx / length) * characterSpeed)
		g.player.y += math.Round((dy / length) * characterSpeed)
	}

	g.counter++

	characterAnimationCurrentFrame = (g.counter / 10) % characterAnimationFramesCount
	if characterAnimationCurrentFrame >= characterAnimationFramesCount {
		characterAnimationCurrentFrame = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	charOp := &ebiten.DrawImageOptions{}

	// charOp.GeoM.Translate(g.player.x, g.player.y)
	charOp.GeoM.Translate(
		(tileSize*visibleHorizontalTiles)/2,
		(tileSize*visibleVerticalTiles)/2,
	)

	mapOp := &ebiten.DrawImageOptions{}

	mapOp.GeoM.Translate(
		-g.player.x+((tileSize*visibleHorizontalTiles)/2),
		-g.player.y+((tileSize*visibleVerticalTiles)/2),
	)

	level := ldtkProject.Levels[0]

	if err := ebitenRenderer.Render(level, screen, &renderer.DrawOptions{
		LayerDrawOptions: mapOp,
	}); err != nil {
		log.Fatal(err)
	}

	for _, layer := range level.Layers {
		if layer.Type != ldtkgo.LayerTypeEntity {
			continue
		}

		for _, entity := range layer.Entities {
			// entity.Position is [x, y] in world space
			ex := float64(entity.Position[0])
			ey := float64(entity.Position[1])

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(
				ex-g.player.x+(float64(tileSize*visibleHorizontalTiles)/2),
				ey-g.player.y+(float64(tileSize*visibleVerticalTiles)/2),
			)

			// Draw a placeholder rect or a real sprite based on entity.Identifier
			ebitenutil.DrawRect(screen,
				ex-g.player.x+(float64(tileSize*visibleHorizontalTiles)/2),
				ey-g.player.y+(float64(tileSize*visibleVerticalTiles)/2),
				float64(entity.Width), float64(entity.Height),
				color.RGBA{255, 0, 0, 180},
			)
		}
	}

	// move vertical spritesheet index based on direction
	currentMovOffset := ((characterSpriteSize + characterSpriteSpacing) * int(g.player.direction))

	// move frames are last two frames of spritesheet and idle frames are first two ones
	currentFrameOffset := (characterSpriteSize + characterSpriteSpacing) * characterAnimationCurrentFrame
	if g.player.isMoving {
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

var ldtkProject *ldtkgo.Project
var ebitenRenderer *renderer.Renderer

func main() {
	projFile, err := os.ReadFile("learning.ldtk")
	if err != nil {
		log.Fatal(err)
	}

	ldtkProject, err = ldtkgo.Read(projFile)
	if err != nil {
		log.Fatal(err)
	}

	ebitenRenderer, err = renderer.New(os.DirFS("."), ldtkProject)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowTitle("Ebiten RPG")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(true)

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
