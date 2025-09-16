package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ----------------- Variables de combat -----------------
var inCombat bool                    // vrai si le joueur est en combat
var combatMonster *Monster           // monstre actuel
var combatPlayerImage *ebiten.Image  // image du joueur pour le combat
var combatFonts = basicfont.Face7x13 // police intégrée

// Images fixes pour le combat
var combatBackground *ebiten.Image
var combatBorder *ebiten.Image

// ----------------- Initialisation du combat -----------------
func InitCombatGraphics() {
	winW, winH := 1000, 400
	combatBackground = ebiten.NewImage(winW, winH)
	combatBackground.Fill(color.RGBA{237, 201, 175, 230}) // fond sable semi-transparent

	combatBorder = ebiten.NewImage(winW, winH)
	combatBorder.Fill(color.RGBA{139, 69, 19, 255}) // bordure marron
}

// ----------------- Début du combat -----------------
func StartCombat(monster *Monster, playerImg *ebiten.Image) {
	if monster == nil || playerImg == nil {
		return
	}
	inCombat = true
	combatMonster = monster
	combatPlayerImage = playerImg
}

// ----------------- Fin du combat -----------------
func EndCombat() {
	inCombat = false
	combatMonster = nil
	combatPlayerImage = nil
}

// ----------------- Mise à jour du combat -----------------
func UpdateCombat() {
	if !inCombat {
		return
	}

	// Quitter le combat avec Escape
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		EndCombat()
	}
}

// ----------------- Dessin de la fenêtre de combat -----------------
func DrawCombatScreen(screen *ebiten.Image) {
	if !inCombat {
		return
	}

	screenW, screenH := screen.Size()
	winW, winH := 1000, 400
	x := (screenW - winW) / 2
	y := (screenH - winH) / 2

	// Fond semi-transparent
	if combatBackground != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(combatBackground, opts)
	}

	// Bordure
	if combatBorder != nil {
		borderOpts := &ebiten.DrawImageOptions{}
		borderOpts.GeoM.Translate(float64(x-4), float64(y-4))
		screen.DrawImage(combatBorder, borderOpts)
	}

	// Texte centré
	if combatMonster != nil {
		msg := "DÉBUT DU COMBAT avec " + combatMonster.Name
		bounds := text.BoundString(combatFonts, msg)
		textX := x + (winW-bounds.Dx())/2
		textY := y + 40
		text.Draw(screen, msg, combatFonts, textX, textY, color.RGBA{255, 0, 0, 255})
	}

	// Dessiner le monstre à gauche
	if combatMonster != nil && len(combatMonster.Sprites) > 0 {
		monsterImg := combatMonster.Sprites[combatMonster.Index%len(combatMonster.Sprites)]
		if monsterImg != nil {
			monsterOpts := &ebiten.DrawImageOptions{}
			monsterOpts.GeoM.Translate(float64(x+50), float64(y+100))
			screen.DrawImage(monsterImg, monsterOpts)
		}
	}

	// Dessiner le joueur à droite
	if combatPlayerImage != nil {
		playerOpts := &ebiten.DrawImageOptions{}
		playerOpts.GeoM.Translate(float64(x+winW-150), float64(y+100))
		screen.DrawImage(combatPlayerImage, playerOpts)
	}
}

// ----------------- Vérification collision et déclenchement combat -----------------
func CheckCollisionWithPlayerCombat() {
	if inCombat || len(currentSprites) == 0 {
		return // déjà en combat ou pas d'image joueur
	}

	playerW, playerH := 64.0, 64.0

	for _, m := range monsters {
		if len(m.Sprites) == 0 {
			continue
		}

		monsterW, monsterH := m.Sprites[0].Size()

		if playerX < m.X+float64(monsterW) &&
			playerX+playerW > m.X &&
			playerY < m.Y+float64(monsterH) &&
			playerY+playerH > m.Y {
			// Début combat → jeu en pause, monstres arrêtent de bouger et disparaissent
			StartCombat(m, currentSprites[index])
			return
		}
	}
}
