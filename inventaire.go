package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type InventaireGUI struct {
	player   *Personnage
	open     bool
	keyPrevP bool
}

func NewInventaireGUI(p *Personnage) *InventaireGUI {
	return &InventaireGUI{player: p}
}

func (inv *InventaireGUI) Update() {
	// Toggle inventaire avec P
	keyP := ebiten.IsKeyPressed(ebiten.KeyP)
	if keyP && !inv.keyPrevP {
		inv.open = !inv.open
	}
	inv.keyPrevP = keyP
}

func (inv *InventaireGUI) Draw(screen *ebiten.Image) {
	if !inv.open {
		return
	}

	// Taille de l’inventaire
	screenW, screenH := screen.Size()
	width, height := screenW/2, screenH/2
	x := (screenW - width) / 2
	y := (screenH - height) / 2

	// Fond semi-transparent
	bg := ebiten.NewImage(width, height)
	bg.Fill(color.RGBA{0, 0, 0, 180}) // noir transparent
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(bg, opts)

	// Titre inventaire
	ebitenutil.DebugPrintAt(screen, "🎒 INVENTAIRE", x+20, y+20)

	// Afficher argent
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("💰 Or: %d", inv.player.Money), x+20, y+50)

	// Grille des items
	colSize := 4 // 4 items par ligne
	cellW, cellH := 100, 40
	startY := y + 80

	if len(inv.player.Inventory) == 0 {
		ebitenutil.DebugPrintAt(screen, "(vide)", x+20, startY)
		return
	}

	for i, item := range inv.player.Inventory {
		col := i % colSize
		row := i / colSize
		itemX := x + 20 + col*cellW
		itemY := startY + row*cellH
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("• %s", item), itemX, itemY)
	}
}
