package source

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// MenuMarchand gère l'interface du marchand
// MenuMarchand gère l'interface du marchand et les interactions d'achat
type MenuMarchand struct {
	player      *Personnage // Référence au joueur
	open        bool        // Menu ouvert ou fermé
	shopItems   []ShopItem  // Liste des objets en vente
	message     string      // Message temporaire
	messageTime time.Time   // Temps d'affichage du message

	shopZoneX        float64 // Position X du marchand
	shopZoneY        float64 // Position Y du marchand
	shopZoneW        float64 // Largeur de la zone du marchand
	shopZoneH        float64 // Hauteur de la zone du marchand
	lastMousePressed bool    // Pour détecter le front du clic
}

// ShopItem représente un objet à vendre
// ShopItem représente un objet à vendre
type ShopItem struct {
	Name  string // Nom de l'objet
	Price int    // Prix de l'objet
}

// NewMenuMarchand initialise le marchand
// Initialise le menu du marchand avec les objets disponibles
func NewMenuMarchand(p *Personnage) *MenuMarchand {
	items := []ShopItem{
		{"Plante curative", 50},
		{"Potion magique", 25},
		{"Épée", 50},
		{"Épée améliorée", 150},
		{"Armure", 50},
		{"Botte", 50},
		{"Chapeau", 50},
	}
	return &MenuMarchand{
		player:    p,
		shopItems: items,
		shopZoneX: 193, // coordonnées du marchand sur la map
		shopZoneY: 9,
		shopZoneW: 120, // largeur du sprite du marchand
		shopZoneH: 120, // hauteur du sprite
	}
}

// Update gère l'ouverture automatique et les achats
// Met à jour l'état du menu marchand et gère les achats
func (m *MenuMarchand) Update() {
	playerX := m.player.PosX
	playerY := m.player.PosY
	playerW := m.player.Width
	playerH := m.player.Height

	// Détecte collision joueur <-> zone du marchand
	if playerX < m.shopZoneX+m.shopZoneW &&
		playerX+playerW > m.shopZoneX &&
		playerY < m.shopZoneY+m.shopZoneH &&
		playerY+playerH > m.shopZoneY {
		m.open = true
	} else {
		m.open = false
	}

	if !m.open {
		return
	}

	// Achat sur front du clic gauche
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if mousePressed && !m.lastMousePressed {
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
					m.player.AjouterItem(item.Name) // applique effets automatiquement
					m.message = fmt.Sprintf("Vous avez acheté %s pour %d pièces !", item.Name, item.Price)
					m.messageTime = time.Now()
				} else {
					m.message = "Pas assez d'or !"
					m.messageTime = time.Now()
				}
			}
		}
	}
	m.lastMousePressed = mousePressed
}

// Draw affiche le menu marchand
func (m *MenuMarchand) Draw(screen *ebiten.Image) {
	if !m.open {
		return
	}

	screenW, screenH := screen.Size()
	width, height := screenW*3/5, screenH*2/5
	x := (screenW - width) / 2
	y := (screenH - height) / 2
	radius := 15

	// Fond et ombre
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

	// Affiche les items
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
			slotColor = color.RGBA{218, 165, 32, 230}
		}

		drawRoundedRect(screen, itemX, itemY, cellW-10, cellH-10, slotRadius, slotColor)

		textStr := fmt.Sprintf("%s (%d)", item.Name, item.Price)
		tW := text.BoundString(face, textStr).Dx()
		tH := text.BoundString(face, textStr).Dy()
		text.Draw(screen, textStr, face, itemX+(cellW-10)/2-tW/2, itemY+(cellH-10)/2+tH/2, color.RGBA{101, 67, 33, 255})
	}

	// Message achat ou erreur
	if m.message != "" && time.Since(m.messageTime).Seconds() < 2 {
		msgW := text.BoundString(face, m.message).Dx()
		text.Draw(screen, m.message, face, x+width/2-msgW/2, y+height-20, color.RGBA{255, 0, 0, 255})
	}
}
