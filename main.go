package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

const (
	tileSize               = 16
	visibleVerticalTiles   = 18
	visibleHorizontalTiles = 32
)

type Game struct {
	mapImage *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.mapImage, nil)
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

	if err := ebiten.RunGame(&Game{
		mapImage: ebiten.NewImageFromImage(renderer.Result),
	}); err != nil {
		log.Fatal(err)
	}
}
