package main

import (
	"fmt"
	"log"
	"time"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game représente l'état du jeu
type Game struct {
	frames        []*ebiten.Image // images simulant la vidéo
	index         int             // index de la frame actuelle
	start         bool            // false = écran d'accueil, true = jeu
	lastFrameTime time.Time       // dernier temps de changement de frame
	frameDelay    time.Duration   // durée entre chaque frame
	videoEnded    bool            // vrai si la vidéo est terminée
}

// NewGame charge les images pour l'écran d'accueil
func NewGame() *Game {
	g := &Game{
		frameDelay: time.Millisecond * 42, // ~24 FPS
	}

	// ----------------------------
	// PLACE TES IMAGES ICI
	// 1. Convertis ta vidéo intro.mp4 en images PNG
	// 2. Mets-les dans le dossier "video_frames" à la racine du projet
	// ----------------------------
	totalFrames := 150 // change selon le nombre d'images que tu as
	for i := 1; i <= totalFrames; i++ {
		path := fmt.Sprintf("video_frames/frame%03d.png", i)
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Println("Impossible de charger l'image :", path, err)
			continue
		}
		g.frames = append(g.frames, img)
	}

	g.lastFrameTime = time.Now()
	return g
}

// Update gère la logique du jeu
func (g *Game) Update() error {
	// Quitter avec ESC
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if !g.videoEnded {
		// Écran d'accueil
		now := time.Now()
		if now.Sub(g.lastFrameTime) >= g.frameDelay {
			g.index++
			if g.index >= len(g.frames) {
				// La vidéo est terminée, rester sur la dernière frame
				g.index = len(g.frames) - 1
				g.videoEnded = true
			}
			g.lastFrameTime = now
		}
	}

	return nil
}

// Draw dessine l'écran
func (g *Game) Draw(screen *ebiten.Image) {
	if len(g.frames) > 0 && g.index < len(g.frames) {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(
			float64(screen.Bounds().Dx())/float64(g.frames[g.index].Bounds().Dx()),
			float64(screen.Bounds().Dy())/float64(g.frames[g.index].Bounds().Dy()),
		)
		screen.DrawImage(g.frames[g.index], opts)
	} else {
		screen.Fill(color.Black)
	}

	if !g.videoEnded {
		ebitenutil.DebugPrint(screen, "SAHARA DEFENDER\nVidéo en cours... ESC pour quitter")
	} else {
		ebitenutil.DebugPrint(screen, "SAHARA DEFENDER\nFin de la vidéo. ESC pour quitter ou ENTER pour continuer")
	}
}

// Layout retourne les dimensions de l'écran
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	game := NewGame()

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("SAHARA DEFENDER")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
