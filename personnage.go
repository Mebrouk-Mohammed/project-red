package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Personnage représente le joueur
type Personnage struct {
	PosX   float64
	PosY   float64
	Width  float64
	Height float64

	Name      string
	Life      int
	MaxLife   int
	Shield    int
	MaxShield int
	Strength  int
	Money     int
	Inventory []string
}

// AjouterItem ajoute un item à l’inventaire et applique ses effets
func (p *Personnage) AjouterItem(item string) {
	 p.Inventory = append(p.Inventory, item)
	 fmt.Printf("%s a ajouté %s à son inventaire.\n", p.Name, item)
	 // Les effets sont appliqués uniquement lors de l'utilisation dans l'inventaire
}

// RetirerItem retire un item de l’inventaire
func (p *Personnage) RetirerItem(item string) {
	for i, v := range p.Inventory {
		if v == item {
			p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
			return
		}
	}
}

// AjouterShield ajoute des points de shield
func (p *Personnage) AjouterShield(amount int) {
	p.Shield += amount
	if p.Shield > p.MaxShield {
		p.Shield = p.MaxShield
	}
	fmt.Printf("%s a gagné %d points de shield. Shield: %d/%d\n", p.Name, amount, p.Shield, p.MaxShield)
}

// PrendreDegats applique des dégâts au shield et à la vie
func (p *Personnage) PrendreDegats(damage int) {
	if p.Shield > 0 {
		if damage <= p.Shield {
			p.Shield -= damage
			damage = 0
		} else {
			damage -= p.Shield
			p.Shield = 0
		}
	}

	p.Life -= damage
	if p.Life < 0 {
		p.Life = 0
	}

	fmt.Printf("%s a pris %d points de dégâts. Vie: %d/%d, Shield: %d/%d\n",
		p.Name, damage, p.Life, p.MaxLife, p.Shield, p.MaxShield)

	if p.Life == 0 {
		fmt.Printf("%s est mort!\n", p.Name)
	}
}

// Soigner soigne le joueur
func (p *Personnage) Soigner(heal int) {
	p.Life += heal
	if p.Life > p.MaxLife {
		p.Life = p.MaxLife
	}
	fmt.Printf("%s a été soigné de %d points. Vie: %d/%d\n", p.Name, heal, p.Life, p.MaxLife)
}

// AfficherInventaire affiche l’inventaire
func (p *Personnage) AfficherInventaire() {
	fmt.Println("Inventory:")
	if len(p.Inventory) == 0 {
		fmt.Println("  (vide)")
	}
	for _, item := range p.Inventory {
		fmt.Printf("  - %s\n", item)
	}
}

// AfficherStatut affiche les informations du joueur
func (p *Personnage) AfficherStatut() {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Vie: %d/%d\n", p.Life, p.MaxLife)
	fmt.Printf("Shield: %d/%d\n", p.Shield, p.MaxShield)
}

// DrawBars dessine les barres de vie et de shield
func (p *Personnage) DrawBars(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()
	barWidth := 400
	barHeight := 25
	padding := 10
	x := (screenWidth - barWidth) / 2
	y := screenHeight - barHeight*2 - padding*2

	// Vie
	lifeRatio := float64(p.Life) / float64(p.MaxLife)
	if lifeRatio < 0 {
		lifeRatio = 0
	}
	lifeWidth := int(float64(barWidth) * lifeRatio)
	if lifeWidth < 1 {
		lifeWidth = 1
	}
	drawRectBar(screen, x, y, lifeWidth, barHeight, color.RGBA{200, 0, 0, 255})

	// Shield
	shieldRatio := float64(p.Shield) / float64(p.MaxShield)
	if shieldRatio < 0 {
		shieldRatio = 0
	}
	shieldWidth := int(float64(barWidth) * shieldRatio)
	if shieldWidth < 1 {
		shieldWidth = 1
	}
	drawRectBar(screen, x, y+barHeight+padding, shieldWidth, barHeight, color.RGBA{0, 128, 255, 200})
}

// drawRectBar dessine un rectangle simple (fonction renommée pour éviter conflit)
func drawRectBar(screen *ebiten.Image, x, y, w, h int, col color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	img := ebiten.NewImage(w, h)
	img.Fill(col)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}
