package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	tilesize = 150.0
)

type Player struct {
	Speed         float64
	HP            int
	Inventory     []Item
	IHaveBoat     bool
	x, y          float64
	h, w          float64
	BoatFramesCNT int
	lifes         int
	GameOverFCNT  int

	color color.Color
}

type Item struct {
	Name  string
	Type  string // "key", "potion", "relic"
	Value int
}

type Enemy struct {
	x, y            float64
	HP              int
	Speed           float64
	HitcooldownFCNT int
}

type Tile struct {
	Type  string
	Solid bool
	Color color.Color
	x, y  float64
}

type Game struct {
	XMap, YMap    float64
	Player        Player
	Enemies       []*Enemy
	Items         []Item
	TimeOfDay     float64
	DayDuration   float64
	currentScreen string

	MapTiles  [][]Tile
	MapTiles2 [][]Tile
}

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	rand.Seed(time.Now().UnixNano())
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

var Map1 = [][]string{
	{"W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W"},
	{"W", "G", "G", "G", "G", "G", "S", "S", "S", "G", "G", "G", "R", "R", "G", "G", "G", "S", "S", "W"},
	{"W", "G", "S", "S", "S", "S", "S", "S", "G", "G", "R", "R", "R", "G", "S", "S", "S", "S", "S", "W"},
	{"W", "G", "S", "S", "S", "G", "G", "G", "W", "W", "G", "S", "S", "S", "S", "G", "G", "G", "S", "W"},
	{"W", "G", "S", "S", "S", "G", "G", "W", "W", "W", "G", "G", "S", "S", "S", "S", "S", "G", "S", "W"},
	{"W", "G", "S", "S", "S", "S", "S", "G", "W", "W", "W", "G", "S", "S", "S", "S", "S", "S", "G", "W"},
	{"W", "G", "G", "S", "S", "S", "S", "G", "G", "G", "G", "G", "S", "S", "S", "G", "R", "R", "G", "W"},
	{"W", "G", "G", "S", "S", "S", "S", "R", "R", "R", "R", "G", "S", "S", "S", "S", "R", "R", "G", "W"},
	{"W", "G", "S", "S", "S", "S", "S", "R", "S", "S", "G", "G", "S", "S", "S", "S", "S", "S", "G", "W"},
	{"W", "G", "S", "S", "G", "G", "G", "R", "S", "S", "R", "G", "G", "G", "S", "S", "S", "S", "G", "W"},
	{"W", "G", "S", "S", "G", "R", "R", "R", "R", "R", "R", "G", "S", "S", "S", "S", "S", "S", "G", "W"},
	{"W", "G", "G", "S", "S", "S", "S", "S", "G", "G", "G", "G", "S", "S", "S", "S", "S", "S", "G", "W"},
	{"W", "G", "S", "S", "S", "S", "S", "S", "S", "G", "S", "S", "S", "S", "S", "S", "S", "G", "G", "W"},
	{"W", "G", "S", "S", "S", "S", "G", "G", "G", "G", "G", "S", "S", "S", "S", "S", "S", "S", "G", "W"},
	{"W", "G", "S", "S", "G", "G", "R", "R", "G", "S", "S", "S", "S", "S", "S", "G", "G", "R", "R", "W"},
	{"W", "G", "G", "S", "S", "G", "R", "R", "R", "G", "S", "S", "S", "S", "G", "G", "R", "R", "G", "W"},
	{"W", "S", "S", "S", "S", "G", "G", "G", "G", "S", "S", "S", "S", "S", "G", "S", "S", "S", "S", "W"},
	{"W", "S", "S", "S", "S", "S", "S", "S", "G", "S", "S", "S", "S", "S", "S", "S", "S", "S", "G", "W"},
	{"W", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "W"},
	{"W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W"},
}
var Map2 = [][]string{
	{"W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W"},
	{"W", "D", "D", "D", "S", "S", "S", "L", "L", "L", "L", "S", "S", "D", "D", "D", "R", "R", "G", "W"},
	{"W", "D", "S", "S", "S", "S", "S", "L", "L", "L", "L", "L", "S", "S", "S", "S", "S", "R", "G", "W"},
	{"W", "D", "S", "S", "D", "D", "S", "S", "L", "L", "L", "S", "S", "S", "S", "R", "R", "R", "G", "W"},
	{"W", "S", "S", "S", "S", "D", "D", "S", "S", "S", "S", "S", "S", "S", "S", "R", "G", "G", "G", "W"},
	{"W", "S", "S", "S", "S", "S", "S", "D", "D", "D", "S", "S", "S", "S", "S", "S", "G", "G", "G", "W"},
	{"W", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "G", "G", "G", "W", "W"},
	{"W", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "G", "R", "R", "W", "W"},
	{"W", "L", "L", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "G", "R", "R", "G", "W"},
	{"W", "L", "L", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "G", "R", "R", "G", "W"},
	{"W", "L", "L", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "G", "R", "R", "G", "W"},
	{"W", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "G", "R", "R", "G", "W"},
	{"W", "D", "D", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "R", "R", "R", "R", "W"},
	{"W", "D", "D", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "R", "R", "R", "R", "W"},
	{"W", "D", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "R", "R", "R", "R", "W"},
	{"W", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "R", "R", "R", "R", "W"},
	{"W", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "R", "R", "R", "R", "W"},
	{"W", "G", "G", "G", "G", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "R", "R", "R", "R", "W"},
	{"W", "G", "G", "G", "G", "S", "S", "S", "S", "S", "S", "S", "S", "S", "S", "R", "R", "R", "R", "W"},
	{"W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W"},
}

func (t *Tile) Draw(screen *ebiten.Image, xMap, yMap float64) {
	ebitenutil.DrawRect(screen, t.x+xMap, t.y+yMap, tilesize, tilesize, t.Color)
}

func (p *Player) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.x, p.y, p.w, p.h, p.color)
}

func (e *Enemy) Draw(screen *ebiten.Image, xMap float64, yMap float64) {
	ebitenutil.DrawRect(screen, e.x+xMap, e.y+yMap, 50, 50, color.RGBA{R: 255, G: 50, B: 50, A: 255})
}
func NewGame() *Game {
	g := &Game{
		Player: Player{
			x:     725,
			y:     425,
			h:     50,
			w:     50,
			lifes: 5,
			color: color.RGBA{34, 139, 34, 255},
			Speed: 7,
		},
		Enemies: []*Enemy{
			{x: 450, y: 150, Speed: 1, HP: 1},
			{x: 450, y: 900, Speed: 2, HP: 1},
			{x: 600, y: 1500, Speed: 3, HP: 1},
			{x: 750, y: 2100, Speed: 4, HP: 1},
			{x: 300, y: 2250, Speed: 5, HP: 1},
		},
		TimeOfDay:     0.0,
		DayDuration:   60,
		currentScreen: "Menu",
	}
	g.XMap = 237
	g.YMap = 28

	// Prépare les slices pour chaque map
	g.MapTiles = make([][]Tile, len(Map1))
	g.MapTiles2 = make([][]Tile, len(Map2))

	// Remplir MapTiles (Map1)
	for i, row := range Map1 {
		g.MapTiles[i] = make([]Tile, len(row))
		for j, typ := range row {
			x := float64(j) * tilesize
			y := float64(i) * tilesize

			switch typ {
			case "G":
				g.MapTiles[i][j] = Tile{Type: "G", Solid: false, Color: color.RGBA{100, 255, 100, 255}, x: x, y: y}
			case "W":
				g.MapTiles[i][j] = Tile{Type: "W", Solid: true, Color: color.RGBA{0, 0, 255, 255}, x: x, y: y}
			case "S":
				g.MapTiles[i][j] = Tile{Type: "S", Solid: false, Color: color.RGBA{194, 178, 128, 255}, x: x, y: y}
			case "R":
				g.MapTiles[i][j] = Tile{Type: "R", Solid: true, Color: color.RGBA{112, 128, 144, 255}, x: x, y: y}
			default:
				// Sécurité : tuile par défaut (ex: sol)
				g.MapTiles[i][j] = Tile{Type: typ, Solid: false, Color: color.RGBA{150, 150, 150, 255}, x: x, y: y}
			}
		}
	}

	// Remplir MapTiles2 (Map2) — séparé, et on gère "D"
	for i, row := range Map2 {
		g.MapTiles2[i] = make([]Tile, len(row))
		for j, typ := range row {
			x := float64(j) * tilesize
			y := float64(i) * tilesize

			switch typ {
			case "G":
				g.MapTiles2[i][j] = Tile{Type: "G", Solid: false, Color: color.RGBA{100, 255, 100, 255}, x: x, y: y}
			case "W":
				g.MapTiles2[i][j] = Tile{Type: "W", Solid: true, Color: color.RGBA{0, 0, 255, 255}, x: x, y: y}
			case "S":
				g.MapTiles2[i][j] = Tile{Type: "S", Solid: false, Color: color.RGBA{194, 178, 128, 255}, x: x, y: y}
			case "R":
				g.MapTiles2[i][j] = Tile{Type: "R", Solid: true, Color: color.RGBA{112, 128, 144, 255}, x: x, y: y}
			case "L":
				g.MapTiles2[i][j] = Tile{Type: "L", Solid: true, Color: color.RGBA{255, 69, 0, 255}, x: x, y: y}
			case "D": // <<-- Ajouté : "D" (dirt / terre)
				g.MapTiles2[i][j] = Tile{Type: "D", Solid: false, Color: color.RGBA{150, 111, 51, 255}, x: x, y: y}
			default:
				// sécurité
				g.MapTiles2[i][j] = Tile{Type: typ, Solid: false, Color: color.RGBA{150, 150, 150, 255}, x: x, y: y}
			}
		}
	}

	return g
}

func (p *Player) Touch(t Tile, xMap, yMap float64) bool {
	if p.x < t.x+xMap+tilesize &&
		p.x+p.w > t.x+xMap &&
		p.y < t.y+yMap+tilesize &&
		p.y+p.h > t.y+yMap {
		return true
	}

	return false
}
func (e *Enemy) Touch(t Tile, xMap, yMap float64) bool {
	if e.x < t.x+xMap+tilesize &&
		e.x+50 > t.x+xMap &&
		e.y < t.y+yMap+tilesize &&
		e.y+50 > t.y+yMap {
		return true
	}

	return false
}

// Retourne toutes les tuiles sur lesquelles le Player touche actuellement
func (g *Game) currentTile(xMap, yMap float64) []Tile {
	ret := make([]Tile, 0)
	if g.currentScreen == "world 1" {
		for _, row := range g.MapTiles {
			for _, t := range row {
				if g.Player.Touch(t, xMap, yMap) {
					ret = append(ret, t)
				}
			}
		}
	} else if g.currentScreen == "world 2" {
		for _, row := range g.MapTiles2 {
			for _, t := range row {
				if g.Player.Touch(t, xMap, yMap) {
					ret = append(ret, t)
				}
			}
		}
	}
	return ret
} // <-- Ajoute cette accolade fermante ici
func RectRectCollision(x1, y1, w1, h1 float64, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 &&
		x1+w1 > x2 &&
		y1 < y2+h2 &&
		y1+h1 > y2
}
func (g *Game) currentTileE(newX, newY float64) []Tile {
	ret := make([]Tile, 0)
	tiles := g.MapTiles
	if g.currentScreen == "world 2" {
		tiles = g.MapTiles2
	}

	// newX et newY sont les nouvelles coordonnées monde de l'ennemi
	for _, row := range tiles {
		for _, t := range row {
			// Vérifier collision entre la nouvelle position de l'ennemi et la tuile
			if newX < t.x+tilesize &&
				newX+50 > t.x &&
				newY < t.y+tilesize &&
				newY+50 > t.y {
				if t.Solid {
					ret = append(ret, t)
				}
			}
		}
	}
	return ret
}

func (g *Game) PlayerX() float64 {
	return g.Player.x - g.XMap
}

func (g *Game) PlayerY() float64 {
	return g.Player.y - g.YMap
}

func (g *Game) Update() error {
	if g.Player.GameOverFCNT > 0 {
		g.Player.GameOverFCNT--
	}
	if g.currentScreen == "Menu" {
		if ebiten.IsKeyPressed(ebiten.Key1) {
			g.currentScreen = "world 1"
		}
		if ebiten.IsKeyPressed(ebiten.Key2) {
			g.currentScreen = "world 2"
		}
		return nil
	}
	// Temps du jour/nuit
	g.TimeOfDay += 1.0 / (g.DayDuration * 60) // avance à chaque frame (~60 FPS)
	if g.TimeOfDay > 1.0 {
		g.TimeOfDay = 0.0 // boucle infinie
	}
	if g.currentScreen != "Menu" {
		if g.Player.lifes <= 0 {
			g.Player.GameOverFCNT = 120
			g.currentScreen = "Menu"
		}
		for _, e := range g.Enemies {
			ex, ey := 0.0, 0.0
			playerX := g.PlayerX()
			playerY := g.PlayerY()

			if e.x > playerX {
				ex -= e.Speed
			}
			if e.x < playerX {
				ex += e.Speed
			}
			if e.y > playerY {
				ey -= e.Speed
			}
			if e.y < playerY {
				ey += e.Speed
			}
			if e.HitcooldownFCNT > 0 {
				e.HitcooldownFCNT--
			}
			if RectRectCollision(e.x, e.y, 50, 50, playerX, playerY, 50, 50) && e.HitcooldownFCNT <= 0 {
				g.Player.lifes--
				e.HitcooldownFCNT = 60
			}
			if ey != 0 || ex != 0 {
				canMove := true // Flag pour savoir si on peut bouger
				for _, t := range g.currentTileE(e.x+ex, e.y+ey) {
					if t.Solid {
						canMove = false
						break
					}
				}

				// ✅ APPLIQUER LE DÉPLACEMENT SEULEMENT SI AUTORISÉ
				if canMove {
					e.x += ex
					e.y += ey
				}
				fmt.Println("Enemy moving:", e.x, e.y, "canMove:", canMove)
			}
		}

		if g.Player.BoatFramesCNT > 0 {
			g.Player.BoatFramesCNT--
		}
		var xMap float64
		var yMap float64
		fmt.Printf("XMap: %v", g.XMap)
		fmt.Printf("YMap: %v", g.YMap)
		if g.XMap < -570 && g.XMap > -630 && g.YMap < -870 && g.YMap > -930 {
			g.Player.IHaveBoat = true
			g.Player.BoatFramesCNT = 300
		}

		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			xMap = -g.Player.Speed
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			xMap = g.Player.Speed
		}

		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			yMap = g.Player.Speed
		}

		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			yMap = -g.Player.Speed
		}

		if yMap != 0 || xMap != 0 {
			canMove := true // Flag pour savoir si on peut bouger

			for _, t := range g.currentTile(g.XMap+xMap, g.YMap+yMap) {
				if t.Solid {
					if t.Type == "W" && !g.Player.IHaveBoat {
						canMove = false
						break
					}
					if t.Type == "R" || t.Type == "L" {
						canMove = false
						break
					}
				}
			}

			// ✅ APPLIQUER LE DÉPLACEMENT SEULEMENT SI AUTORISÉ
			if canMove {
				g.XMap += xMap
				g.YMap += yMap
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScreen == "world 1" {
		for _, s := range g.MapTiles {
			for _, tile := range s {
				tile.Draw(screen, g.XMap, g.YMap)
			}
		}
		for _, e := range g.Enemies {
			e.Draw(screen, g.XMap, g.YMap)
		}
		g.Player.Draw(screen)
		chestX := 8*tilesize + 115 // position X (selon ta demande)
		chestY := 9*tilesize - 20  // position Y
		chestW := 90.0
		chestH := 60.0
		if !g.Player.IHaveBoat {
			ebitenutil.DrawRect(screen, chestX+g.XMap+5, chestY+g.YMap+5, chestW, chestH, color.RGBA{90, 70, 0, 255})
			ebitenutil.DrawRect(screen, chestX+g.XMap, chestY+g.YMap, chestW, chestH, color.RGBA{218, 165, 32, 255})
			ebitenutil.DrawRect(screen, chestX+g.XMap, chestY+g.YMap, chestW, chestH/3, color.RGBA{184, 134, 11, 255})
			ebitenutil.DrawRect(screen, chestX+g.XMap+chestW/2-5, chestY+g.YMap, 10, chestH, color.RGBA{120, 120, 120, 255})
			ebitenutil.DrawRect(screen, chestX+g.XMap+chestW/2-6, chestY+g.YMap+chestH/2-5, 12, 12, color.RGBA{30, 30, 30, 255})
		}
		if g.Player.BoatFramesCNT > 0 {
			opBoat := &text.DrawOptions{}
			opBoat.GeoM.Translate(g.Player.x-10, g.Player.y-50)
			opBoat.ColorScale.ScaleWithColor(color.RGBA{0, 150, 255, 255})
			text.Draw(screen, "Boat", &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   23,
			}, opBoat)
		}

		var darkness float64
		if g.TimeOfDay <= 0.5 {
			darkness = g.TimeOfDay * 2 // jour → nuit
		} else {
			darkness = (1.0 - g.TimeOfDay) * 2 // nuit → jour
		}

		overlay := ebiten.NewImage(1500, 900)
		overlay.Fill(color.RGBA{0, 0, 0, uint8(200 * darkness)}) // 200 = noirceur max
		screen.DrawImage(overlay, nil)
	} else if g.currentScreen == "world 2" {
		for _, s := range g.MapTiles2 {
			for _, tile := range s {
				tile.Draw(screen, g.XMap, g.YMap)
			}
		}
		for _, e := range g.Enemies {
			e.Draw(screen, g.XMap, g.YMap)
		}
		g.Player.Draw(screen)
		var darkness float64
		if g.TimeOfDay <= 0.5 {
			darkness = g.TimeOfDay * 2 // jour → nuit
		} else {
			darkness = (1.0 - g.TimeOfDay) * 2 // nuit → jour
		}

		overlay := ebiten.NewImage(1500, 900)
		overlay.Fill(color.RGBA{0, 0, 0, uint8(200 * darkness)}) // 200 = noirceur max
		screen.DrawImage(overlay, nil)
	} else if g.currentScreen == "Menu" {
		opTitle := &text.DrawOptions{}
		opTitle.GeoM.Translate(100, 100)
		opTitle.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
		text.Draw(screen, "Choisis un monde :", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   28,
		}, opTitle)

		op1 := &text.DrawOptions{}
		op1.GeoM.Translate(100, 150)
		op1.ColorScale.ScaleWithColor(color.RGBA{180, 220, 255, 255})
		text.Draw(screen, "[1] Monde 1", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   24,
		}, op1)

		op2 := &text.DrawOptions{}
		op2.GeoM.Translate(100, 190)
		op2.ColorScale.ScaleWithColor(color.RGBA{255, 220, 180, 255})
		text.Draw(screen, "[2] Monde 2", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   24,
		}, op2)

	}
	if g.Player.GameOverFCNT > 0 {
		// Affichage du titre du gagnant
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(80), float64(300))
		op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
		text.Draw(screen, "Game Over", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   150,
		}, op)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Lives: %d", g.Player.lifes), 10, 10)
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
