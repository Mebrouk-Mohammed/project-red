package source

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type InventaireGUI struct {
	player   *Personnage
	open     bool
	keyPrevP bool
	message  string
	msgTime  time.Time
}

func NewInventaireGUI(player *Personnage) *InventaireGUI {
	return &InventaireGUI{
		player: player,
		open:   false,
	}
}

func (inv *InventaireGUI) Update() {
	keyP := ebiten.IsKeyPressed(ebiten.KeyP)
	if keyP && !inv.keyPrevP {
		inv.open = !inv.open
	}
	inv.keyPrevP = keyP

	if !inv.open {
		return
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		colSize := 5
		cellW, cellH := 110, 50
		screenW, screenH := ebiten.ScreenSizeInFullscreen()
		width, height := screenW*3/5, screenH*2/5
		x := (screenW - width) / 2
		y := (screenH - height) / 2
		startX := x + 20
		startY := y + 90

		for i, item := range inv.player.Inventory {
			col := i % colSize
			row := i / colSize
			itemX := startX + col*cellW
			itemY := startY + row*cellH

			if float64(mx) >= float64(itemX) && float64(mx) <= float64(itemX+cellW-10) &&
				float64(my) >= float64(itemY) && float64(my) <= float64(itemY+cellH-10) {

				switch item {
				case "Plante curative":
					inv.player.Soigner(50)
					inv.message = fmt.Sprintf("%s utilise %s ! Vie: %d/%d", inv.player.Name, item, inv.player.Life, inv.player.MaxLife)
				case "Potion magique":
					inv.player.AjouterShield(10)
					inv.message = fmt.Sprintf("%s utilise %s ! Shield: %d/%d", inv.player.Name, item, inv.player.Shield, inv.player.MaxShield)
				case "Armure":
					inv.player.MaxShield += 30
					inv.message = fmt.Sprintf("%s utilise %s ! MaxShield: %d", inv.player.Name, item, inv.player.MaxShield)
				case "Botte":
					inv.player.MaxShield += 20
					inv.message = fmt.Sprintf("%s utilise %s ! MaxShield: %d", inv.player.Name, item, inv.player.MaxShield)
				case "Chapeau":
					inv.player.MaxShield += 10
					inv.message = fmt.Sprintf("%s utilise %s ! MaxShield: %d", inv.player.Name, item, inv.player.MaxShield)
				default:
					inv.message = fmt.Sprintf("%s ne peut pas utiliser %s", inv.player.Name, item)
				}

				// Retirer l'item après usage
				inv.player.Inventory = append(inv.player.Inventory[:i], inv.player.Inventory[i+1:]...)
				inv.msgTime = time.Now()
				break
			}
		}
	}
}

func (inv *InventaireGUI) Draw(screen *ebiten.Image) {
	if !inv.open {
		return
	}

	screenW, screenH := screen.Size()
	width, height := screenW*3/5, screenH*2/5
	x := (screenW - width) / 2
	y := (screenH - height) / 2
	radius := 15

	// Ombre et fond
	drawRoundedRect(screen, x+5, y+5, width, height, radius, color.RGBA{120, 80, 30, 180})
	drawRoundedRect(screen, x, y, width, height, radius, color.RGBA{210, 180, 140, 230})

	face := basicfont.Face7x13

	// Titre
	title := "🏜️ Inventaire du Désert"
	tW := text.BoundString(face, title).Dx()
	text.Draw(screen, title, face, x+width/2-tW/2, y+30, color.RGBA{101, 67, 33, 255})

	// Argent joueur
	money := fmt.Sprintf("💰 Or: %d", inv.player.Money)
	tW = text.BoundString(face, money).Dx()
	text.Draw(screen, money, face, x+width/2-tW/2, y+50, color.RGBA{139, 69, 19, 255})

	// Grille des items
	colSize := 5
	cellW, cellH := 110, 50
	startX := x + 20
	startY := y + 90
	slotRadius := 10

	if len(inv.player.Inventory) == 0 {
		tW = text.BoundString(face, "(vide)").Dx()
		text.Draw(screen, "(vide)", face, x+width/2-tW/2, startY, color.RGBA{101, 67, 33, 255})
		return
	}

	for i, item := range inv.player.Inventory {
		col := i % colSize
		row := i / colSize
		itemX := startX + col*cellW
		itemY := startY + row*cellH

		slotColor := color.RGBA{184, 134, 11, 200}

		mx, my := ebiten.CursorPosition()
		if float64(mx) >= float64(itemX) && float64(mx) <= float64(itemX+cellW-10) &&
			float64(my) >= float64(itemY) && float64(my) <= float64(itemY+cellH-10) {
			slotColor = color.RGBA{218, 165, 32, 230}
		}

		drawRoundedRect(screen, itemX, itemY, cellW-10, cellH-10, slotRadius, slotColor)

		tW := text.BoundString(face, item).Dx()
		tH := text.BoundString(face, item).Dy()
		text.Draw(screen, item, face, itemX+(cellW-10)/2-tW/2, itemY+(cellH-10)/2+tH/2, color.RGBA{101, 67, 33, 255})
	}

	// Message temporaire
	if inv.message != "" && time.Since(inv.msgTime).Seconds() < 2 {
		msgW := text.BoundString(face, inv.message).Dx()
		text.Draw(screen, inv.message, face, x+width/2-msgW/2, y+height-20, color.RGBA{255, 0, 0, 255})
	}
}

// Fonctions utilitaires graphiques
func drawRect(screen *ebiten.Image, x, y, w, h int, fill color.RGBA) {
	img := ebiten.NewImage(w, h)
	img.Fill(fill)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}

func drawCircle(screen *ebiten.Image, cx, cy, r int, fill color.RGBA) {
	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			if dx*dx+dy*dy <= r*r {
				screen.Set(cx+dx, cy+dy, fill)
			}
		}
	}
}

func drawRoundedRect(screen *ebiten.Image, x, y, w, h, radius int, fill color.RGBA) {
	drawRect(screen, x+radius, y, w-2*radius, h, fill)
	drawRect(screen, x, y+radius, w, h-2*radius, fill)
	drawCircle(screen, x+radius, y+radius, radius, fill)
	drawCircle(screen, x+w-radius-1, y+radius, radius, fill)
	drawCircle(screen, x+radius, y+h-radius-1, radius, fill)
	drawCircle(screen, x+w-radius-1, y+h-radius-1, radius, fill)
}
