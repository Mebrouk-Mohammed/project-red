package source

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ----------------- Structure Monstre -----------------
type Monster struct {
	Name       string          // Nom du monstre
	X, Y       float64         // Position
	Sprites    []*ebiten.Image // Images pour l'animation
	Index      int             // Frame actuelle
	LastUpdate time.Time       // Dernière mise à jour animation
	Speed      float64         // Vitesse du monstre
	DirX, DirY float64         // Direction du mouvement
	Health     int             // Points de vie du monstre
}

// Liste des monstres
var monsters []*Monster
var combatMessage string // Message affiché si combat

// Police par défaut intégrée
var combatFont = basicfont.Face7x13 // pas besoin de fichier .ttf

// ----------------- Initialisation des monstres -----------------
func InitMonsters() {
	serpent := &Monster{
		Name:       "Serpent",
		X:          1300,
		Y:          75,
		Sprites:    loadAndScale([]string{"source/assets/serpent1.png"}, 0.07),
		Speed:      1.5,
		LastUpdate: time.Now(),
		Health:     200,
	}

	scorpion := &Monster{
		Name:       "Scorpion",
		X:          220,
		Y:          350,
		Sprites:    loadAndScale([]string{"source/assets/scorpion1.png"}, 0.20),
		Speed:      2,
		LastUpdate: time.Now(),
		Health:     100,
	}

	hyene := &Monster{
		Name:       "Hyène",
		X:          350,
		Y:          650,
		Sprites:    loadAndScale([]string{"source/assets/hyene1.png"}, 0.20),
		Speed:      1,
		LastUpdate: time.Now(),
		Health:     400,
	}

	monsters = []*Monster{serpent, scorpion, hyene}
}

// ----------------- Mise à jour des monstres -----------------
func UpdateMonsters() {
	for _, m := range monsters {
		// Déplacement
		m.X += m.DirX * m.Speed
		m.Y += m.DirY * m.Speed

		// Rebondir sur les bords de la map
		if m.X < 0 || m.X > 1920 {
			m.DirX *= -1
		}
		if m.Y < 0 || m.Y > 1080 {
			m.DirY *= -1
		}

		// Animation
		if len(m.Sprites) > 1 && time.Since(m.LastUpdate) > 200*time.Millisecond {
			m.Index++
			if m.Index >= len(m.Sprites) {
				m.Index = 0
			}
			m.LastUpdate = time.Now()
		}
	}

	// Vérifier collisions avec le joueur
	checkCollisionWithPlayer()
}

// ----------------- Dessin des monstres -----------------
func DrawMonsters(screen *ebiten.Image) {
	for _, m := range monsters {
		if len(m.Sprites) > 0 {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(m.X, m.Y)
			screen.DrawImage(m.Sprites[m.Index%len(m.Sprites)], opts)
		}
	}
}

// ----------------- Détection des collisions -----------------
func checkCollisionWithPlayer() {
	playerW, playerH := 64.0, 64.0 // taille approximative du joueur

	for _, m := range monsters {
		if len(m.Sprites) == 0 {
			continue
		}

		monsterW, monsterH := m.Sprites[0].Size()

		if playerX < m.X+float64(monsterW) &&
			playerX+playerW > m.X &&
			playerY < m.Y+float64(monsterH) &&
			playerY+playerH > m.Y {
			// Collision détectée → afficher le message
			combatMessage = "DÉBUT DU COMBAT avec " + m.Name
			return
		}
	}

	combatMessage = "" // pas de combat
}

// ----------------- Dessin de la fenêtre combat -----------------
func DrawCombatMessage(screen *ebiten.Image) {
	if combatMessage == "" {
		return
	}

	screenW, screenH := screen.Size()
	winW, winH := 700, 250
	x := (screenW - winW) / 2
	y := (screenH - winH) / 2

	// Fenêtre semi-transparente couleur désert
	win := ebiten.NewImage(winW, winH)
	win.Fill(color.RGBA{237, 201, 175, 200}) // sable clair semi-transparent
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(win, opts)

	// Bordure marron foncé
	border := ebiten.NewImage(winW, winH)
	border.Fill(color.RGBA{139, 69, 19, 255}) // marron
	borderOpts := &ebiten.DrawImageOptions{}
	borderOpts.GeoM.Translate(float64(x-4), float64(y-4))
	screen.DrawImage(border, borderOpts)

	// Texte centré en rouge
	bounds := text.BoundString(combatFont, combatMessage)
	textW := bounds.Dx()
	textH := bounds.Dy()
	tx := x + (winW-textW)/2
	ty := y + (winH-textH)/2 + textH

	text.Draw(screen, combatMessage, combatFont, tx, ty, color.RGBA{255, 0, 0, 255})
}
