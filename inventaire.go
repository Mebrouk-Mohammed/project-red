package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// InventaireGUI gère l'affichage graphique de l'inventaire
type InventaireGUI struct {
	Personnage *Personnage
	Visible    bool
	keyPrevP   bool // pour détecter la touche juste pressée

}

// NewInventaireGUI crée une nouvelle interface d'inventaire
func NewInventaireGUI(p *Personnage) *InventaireGUI {
	return &InventaireGUI{
		Personnage: p,
		Visible:    false,
		keyPrevP:   false,
	}
}

// Update détecte l'appui sur la touche P pour ouvrir/fermer l'inventaire
func (inv *InventaireGUI) Update() {
	keyP := ebiten.IsKeyPressed(ebiten.KeyP)

	// Détecte la touche juste pressée
	if keyP && !inv.keyPrevP {
		inv.Visible = !inv.Visible
	}

	inv.keyPrevP = keyP
}

// Draw affiche l'inventaire sur l'écran
func (inv *InventaireGUI) Draw(screen *ebiten.Image) {
	if !inv.Visible {
		return
	}

	// Fond semi-transparent
	screen.Fill(color.RGBA{0, 0, 0, 200})

	// Titre
	ebitenutil.DebugPrintAt(screen, "INVENTAIRE", 20, 20)

	// Liste des items
	y := 50
	if len(inv.Personnage.Inventory) == 0 {
		ebitenutil.DebugPrintAt(screen, "(vide)", 20, y)
	} else {
		itemCount := make(map[string]int)
		for _, item := range inv.Personnage.Inventory {
			itemCount[item]++
		}
		for item, count := range itemCount {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("- %s x%d", item, count), 20, y)
			y += 20
		}
	}

	// Argent
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Argent: %d", inv.Personnage.Money), 20, y+20)
}
