package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	frames []*ebiten.Image // Tableau des frames de la "vidéo"
	index  int             // Index de la frame actuelle
}

// Nouvelle instance du jeu
func NewGame() *Game {
	g := &Game{}

	for i := 1; i <= 10; i++ { // Exemple avec 10 images
		path := "video_frames/frame0" + string(rune(i+'0')) + ".png"
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Fatal(err) // Si l'image n'est pas trouvée
		}
		g.frames = append(g.frames, img)
	}

	return g
}

func (g *Game) Update() error {
	// Quitter le jeu si ESC est pressé
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Passer à la frame suivante pour simuler la vidéo
	g.index++
	if g.index >= len(g.frames) {
		g.index = 0 // Boucle
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if len(g.frames) == 0 {
		return
	}

	// Dessiner la frame actuelle
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(
		float64(screen.Bounds().Dx())/float64(g.frames[g.index].Bounds().Dx()),
		float64(screen.Bounds().Dy())/float64(g.frames[g.index].Bounds().Dy()),
	)
	screen.DrawImage(g.frames[g.index], opts)

	// Exemple : texte sur le fond
	ebitenutil.DebugPrint(screen, "SAHARA DEFENDER - Appuyez sur ESC pour quitter")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	game := NewGame()

	ebiten.SetFullscreen(true)               // Plein écran
	ebiten.SetWindowTitle("SAHARA DEFENDER") // Titre de la fenêtre

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
