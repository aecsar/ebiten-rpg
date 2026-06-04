package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

type Game struct {
	mapImage *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(.25, .25)

	screen.DrawImage(g.mapImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 400
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

	ebiten.SetWindowSize(600, 400)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(&Game{
		mapImage: ebiten.NewImageFromImage(renderer.Result),
	}); err != nil {
		log.Fatal(err)
	}
}
