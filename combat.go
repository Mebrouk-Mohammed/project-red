package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ----------------- Variables de combat -----------------
var inCombat bool
var combatMonster *Monster
var combatPlayerImage *ebiten.Image
var combatFonts = basicfont.Face7x13

var combatPlayerEntity *Entity
var combatMonsterEntity *Entity

var basicPunch = Weapon{Name: "Coup de poing", Damage: 10}
var sword = Weapon{Name: "Épée", Damage: 25}

// Tour par tour
var playerTurn bool = true
var bPressedLastFrame bool
var vPressedLastFrame bool
var shieldPotion int = 30 // valeur à adapter si besoin
var healPotion int = 50   // valeur à adapter si besoin

// Appui unique pour éviter multi-dégâts
var aPressedLastFrame bool
var ePressedLastFrame bool
var spacePressedLastFrame bool

// ----------------- Début du combat -----------------
func StartCombat(monster *Monster, playerImg *ebiten.Image) {
	if monster == nil || playerImg == nil {
		return
	}

	inCombat = true
	combatMonster = monster
	combatPlayerImage = playerImg
	playerTurn = true

	// PV joueur
	combatPlayerEntity = &Entity{Name: "Joueur", Health: 100}

	// PV monstre selon type
	hp := 50
	switch monster.Name {
	case "Serpent":
		hp = 100
	case "Scorpion":
		hp = 50
	case "Hyene":
		hp = 200
	}
	combatMonsterEntity = &Entity{Name: monster.Name, Health: hp}
}

// ----------------- Fin du combat -----------------
func EndCombat() {
	inCombat = false
	combatMonster = nil
	combatPlayerImage = nil
	combatPlayerEntity = nil
	combatMonsterEntity = nil
	aPressedLastFrame = false
	ePressedLastFrame = false
	spacePressedLastFrame = false
}

// ----------------- Mise à jour du combat -----------------
func UpdateCombat() {
	if !inCombat {
		return
	}

	// Quitter combat avec SPACE
	spacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	if spacePressed && !spacePressedLastFrame {
		EndCombat()
		return
	}
	spacePressedLastFrame = spacePressed

	if playerTurn {
		// Attaque simple "A"
		aPressed := ebiten.IsKeyPressed(ebiten.KeyA)
		if aPressed && !aPressedLastFrame && combatMonsterEntity.Health > 0 {
			combatMonsterEntity.TakeDamage(basicPunch.Damage)
			playerTurn = false // fin du tour → passe au monstre
		}
		aPressedLastFrame = aPressed

		// Attaque épée "E" ou épée améliorée
		ePressed := ebiten.IsKeyPressed(ebiten.KeyE)
		if ePressed && !ePressedLastFrame && combatMonsterEntity.Health > 0 {
			var hasSword, hasSwordPlus bool
			if gameInstance != nil && gameInstance.player != nil {
				for _, item := range gameInstance.player.Inventory {
					if item == "Épée" {
						hasSword = true
					}
					if item == "Épée améliorée" {
						hasSwordPlus = true
					}
				}
			}
			if hasSwordPlus {
				combatMonsterEntity.TakeDamage(50) // Dégâts épée améliorée
				playerTurn = false
			} else if hasSword {
				combatMonsterEntity.TakeDamage(sword.Damage)
				playerTurn = false
			} else {
				fmt.Println("Vous n'avez pas d'épée !")
			}
		}
		ePressedLastFrame = ePressed

		// Potion de shield "B"
		bPressed := ebiten.IsKeyPressed(ebiten.KeyB)
		if bPressed && !bPressedLastFrame {
			if gameInstance != nil && gameInstance.player != nil {
				gameInstance.player.Soigner(shieldPotion)
			}
			playerTurn = false
		}
		bPressedLastFrame = bPressed

		// Potion de soin "V"
		vPressed := ebiten.IsKeyPressed(ebiten.KeyV)
		if vPressed && !vPressedLastFrame {
			if gameInstance != nil && gameInstance.player != nil {
				gameInstance.player.Soigner(healPotion)
			}
			playerTurn = false
		}
		vPressedLastFrame = vPressed

	} else {
		// --- Tour du monstre ---
		if combatMonsterEntity.Health > 0 {
			damage := 20 // valeur par défaut, à adapter si besoin
			combatPlayerEntity.TakeDamage(damage)
			fmt.Printf("%s attaque le joueur et inflige %d dégâts !\n", combatMonster.Name, damage)
		}
		playerTurn = true // fin du tour → revient au joueur
	}

	// Fin combat si monstre mort
	if combatMonsterEntity.Health <= 0 {
		RemoveMonsterFromMap(combatMonster)
		EndCombat()
	}
}

// ----------------- Dessin de la fenêtre de combat -----------------
func DrawCombatScreen(screen *ebiten.Image) {
	if !inCombat || combatMonsterEntity == nil || combatPlayerEntity == nil {
		return
	}

	screenW, screenH := screen.Size()
	winW, winH := 1000, 400
	x := (screenW - winW) / 2
	y := (screenH - winH) / 2

	// Fond
	win := ebiten.NewImage(winW, winH)
	win.Fill(color.RGBA{237, 201, 175, 230})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(win, opts)

	// PV affichés
	text.Draw(screen, "Combat contre "+combatMonsterEntity.Name, combatFonts, x+20, y+40, color.Black)
	text.Draw(screen, "PV Joueur: "+itoa(combatPlayerEntity.Health), combatFonts, x+20, y+80, color.RGBA{0, 0, 255, 255})
	text.Draw(screen, "PV "+combatMonsterEntity.Name+": "+itoa(combatMonsterEntity.Health), combatFonts, x+20, y+120, color.RGBA{255, 0, 0, 255})

	// Monstre à gauche
	if combatMonster != nil && len(combatMonster.Sprites) > 0 {
		img := combatMonster.Sprites[combatMonster.Index%len(combatMonster.Sprites)]
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x+50), float64(y+150))
		screen.DrawImage(img, opts)
	}

	// Joueur à droite
	if combatPlayerImage != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x+winW-150), float64(y+150))
		screen.DrawImage(combatPlayerImage, opts)
	}

	// Instructions
	text.Draw(screen, "A = Attaque simple | E = Épée | SPACE = Fuir", combatFonts, x+20, y+winH-30, color.Black)
}

// ----------------- Collision pour lancer combat -----------------
func CheckCollisionWithPlayerCombat() {
	if inCombat || len(currentSprites) == 0 {
		return
	}

	playerW, playerH := 64.0, 64.0
	for _, m := range monsters {
		if len(m.Sprites) == 0 {
			continue
		}
		w, h := m.Sprites[0].Size()
		if playerX < m.X+float64(w) && playerX+playerW > m.X &&
			playerY < m.Y+float64(h) && playerY+playerH > m.Y {
			StartCombat(m, currentSprites[index])
			return
		}
	}
}

// ----------------- Supprimer monstre de la map -----------------
func RemoveMonsterFromMap(monster *Monster) {
	newList := []*Monster{}
	for _, m := range monsters {
		if m != monster {
			newList = append(newList, m)
		}
	}
	monsters = newList
}

// ----------------- Int -> string -----------------
func itoa(num int) string {
	return fmt.Sprintf("%d", num)
}
