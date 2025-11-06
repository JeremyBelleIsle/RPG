package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tilesize = 150.0
)

type Player struct {
	Speed     float64
	HP        int
	Inventory []Item
}

type Item struct {
	Name  string
	Type  string // "key", "potion", "relic"
	Value int
}

type Enemy struct {
	X, Y   float64
	HP     int
	Speed  float64
	Active bool
}

type Tile struct {
	Type  string
	Solid bool
	Color color.Color
}

type Game struct {
	XMap, YMap float64
	Player     Player
	Enemies    []Enemy
	Items      []Item
	Map        [][]string

	Grass Tile
	Water Tile
	Sand  Tile
	Rock  Tile
}

func (t *Tile) Draw(screen *ebiten.Image, x, y float64) {
	ebitenutil.DrawRect(screen, x, y, tilesize, tilesize, t.Color)
}

func NewGame() *Game {
	g := &Game{
		Grass: Tile{
			Type:  "G",
			Solid: false,
			Color: color.RGBA{100, 255, 100, 255},
		},
		Water: Tile{
			Type:  "W",
			Solid: true,
			Color: color.RGBA{0, 0, 255, 255},
		},
		Sand: Tile{
			Type:  "S",
			Solid: true,
			Color: color.RGBA{194, 178, 128, 255},
		},
		Rock: Tile{
			Type:  "R",
			Solid: true,
			Color: color.RGBA{112, 128, 144, 255},
		},
		Map: [][]string{
			{"G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G"},
			{"G", "G", "S", "S", "S", "W", "W", "W", "G", "G", "G", "G"},
			{"G", "G", "S", "S", "S", "W", "W", "W", "W", "G", "G", "G"},
			{"R", "R", "S", "S", "S", "G", "G", "G", "W", "W", "G", "G"},
			{"R", "S", "S", "S", "S", "G", "G", "G", "W", "W", "G", "G"},
			{"R", "R", "R", "G", "G", "G", "G", "G", "G", "G", "G", "G"},
			{"R", "R", "S", "S", "S", "S", "G", "G", "G", "G", "G", "G"},
			{"R", "S", "S", "S", "S", "S", "G", "W", "W", "G", "G", "G"},
			{"G", "G", "G", "G", "S", "S", "S", "W", "W", "G", "G", "G"},
			{"G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G"},
		},
	}
	fmt.Println(g.Map)

	return g
}
func (g *Game) Update() error {
	x := tilesize + g.XMap
	y := tilesize + g.YMap
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.Map[int(5+x/125)][int(6+y/125)] == "W" || g.Map[int(5+x/125)][int(6+y/125)] == "R" {
			g.XMap -= 7
		} else {
			g.XMap += 7
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.Map[int(5+x/125)][int(6+y/125)] == "W" || g.Map[int(5+x/125)][int(6+y/125)] == "R" {
			g.XMap += 7
		} else {
			g.XMap -= 7
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if g.Map[int(5+x/125)][int(6+y/125)] == "W" || g.Map[int(5+x/125)][int(6+y/125)] == "R" {
			g.YMap += 7
		} else {
			g.YMap -= 7
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if g.Map[int(5+x/125)][int(6+y/125)] == "W" || g.Map[int(5+x/125)][int(6+y/125)] == "R" {
			g.YMap -= 7
		} else {
			g.YMap += 7
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, s := range g.Map {
		for j := range s {
			x := float64(j)*tilesize + g.XMap
			y := float64(i)*tilesize + g.YMap
			switch g.Map[i][j] {
			case "G":
				g.Grass.Draw(screen, x, y)
			case "W":
				g.Water.Draw(screen, x, y)
			case "S":
				g.Sand.Draw(screen, x, y)
			case "R":
				g.Rock.Draw(screen, x, y)
			}
		}
	}
	ebitenutil.DrawRect(screen, 725, 425, 50, 50, color.RGBA{34, 139, 34, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1500, 900
}

func main() {
	ebiten.SetWindowSize(1500, 900)
	ebiten.SetWindowTitle("Hello, World!")
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
