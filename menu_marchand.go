package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// MenuMarchand gère l'interface du marchand
type MenuMarchand struct {
	player      *Personnage
	open        bool
	keyPrevG    bool
	shopItems   []ShopItem
	message     string
	messageTime time.Time
}

// ShopItem représente un objet à vendre
type ShopItem struct {
	Name  string
	Price int
}

// Création du menu marchand
func NewMenuMarchand(p *Personnage) *MenuMarchand {
	items := []ShopItem{
		{"Plante curative", 50},
		{"Épée améliorée", 200},
		{"Armure", 150},
		{"Arc", 120},
		{"Potion magique", 80},
	}
	return &MenuMarchand{player: p, shopItems: items}
}

// Mise à jour du menu : G pour toggle, Q pour fermer
func (m *MenuMarchand) Update() {
	keyG := ebiten.IsKeyPressed(ebiten.KeyG)
	keyQ := ebiten.IsKeyPressed(ebiten.KeyQ)

	// Toggle avec G
	if keyG && !m.keyPrevG {
		m.open = !m.open
	}
	// Fermer avec Q
	if keyQ && m.open {
		m.open = false
	}

	m.keyPrevG = keyG

	if !m.open {
		return
	}

	// Achat si clic gauche
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

		for i, item := range m.shopItems {
			col := i % colSize
			row := i / colSize
			itemX := startX + col*cellW
			itemY := startY + row*cellH

			if float64(mx) >= float64(itemX) && float64(mx) <= float64(itemX+cellW-10) &&
				float64(my) >= float64(itemY) && float64(my) <= float64(itemY+cellH-10) {
				if m.player.Money >= item.Price {
					m.player.Money -= item.Price
					m.player.Inventory = append(m.player.Inventory, item.Name)
					m.message = fmt.Sprintf("Vous avez acheté %s pour %d pièces !", item.Name, item.Price)
					m.messageTime = time.Now()
				} else {
					m.message = "Pas assez d'or !"
					m.messageTime = time.Now()
				}
			}
		}
	}
}

// Dessine le menu marchand
func (m *MenuMarchand) Draw(screen *ebiten.Image) {
	if !m.open {
		return
	}

	screenW, screenH := screen.Size()
	width, height := screenW*3/5, screenH*2/5
	x := (screenW - width) / 2
	y := (screenH - height) / 2
	radius := 15

	// Ombre et fond semi-transparent
	drawRoundedRect(screen, x+5, y+5, width, height, radius, color.RGBA{120, 80, 30, 180})
	drawRoundedRect(screen, x, y, width, height, radius, color.RGBA{210, 180, 140, 230})

	face := basicfont.Face7x13

	// Titre
	title := "🏜️ Marchand du Désert"
	tW := text.BoundString(face, title).Dx()
	text.Draw(screen, title, face, x+width/2-tW/2, y+30, color.RGBA{101, 67, 33, 255})

	// Argent joueur
	money := fmt.Sprintf("💰 Or: %d", m.player.Money)
	tW = text.BoundString(face, money).Dx()
	text.Draw(screen, money, face, x+width/2-tW/2, y+50, color.RGBA{139, 69, 19, 255})

	// Grille des items
	colSize := 5
	cellW, cellH := 110, 50
	startX := x + 20
	startY := y + 90
	slotRadius := 10

	for i, item := range m.shopItems {
		col := i % colSize
		row := i / colSize
		itemX := startX + col*cellW
		itemY := startY + row*cellH

		slotColor := color.RGBA{184, 134, 11, 200}
		mx, my := ebiten.CursorPosition()
		if float64(mx) >= float64(itemX) && float64(mx) <= float64(itemX+cellW-10) &&
			float64(my) >= float64(itemY) && float64(my) <= float64(itemY+cellH-10) {
			slotColor = color.RGBA{218, 165, 32, 230} // survol doré
		}

		drawRoundedRect(screen, itemX, itemY, cellW-10, cellH-10, slotRadius, slotColor)

		// Texte centré dans le slot
		textStr := fmt.Sprintf("%s (%d)", item.Name, item.Price)
		tW := text.BoundString(face, textStr).Dx()
		tH := text.BoundString(face, textStr).Dy()
		text.Draw(screen, textStr, face, itemX+(cellW-10)/2-tW/2, itemY+(cellH-10)/2+tH/2, color.RGBA{101, 67, 33, 255})
	}

	// Affichage du message (achat ou pas assez d'or) pendant 2 secondes
	if m.message != "" && time.Since(m.messageTime).Seconds() < 2 {
		msgW := text.BoundString(face, m.message).Dx()
		text.Draw(screen, m.message, face, x+width/2-msgW/2, y+height-20, color.RGBA{255, 0, 0, 255})
	}
}
