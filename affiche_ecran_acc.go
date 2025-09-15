package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	audioCtx *audio.Context
	player   *audio.Player
)

// Game représente l'état du jeu
type Game struct {
	frames        []*ebiten.Image
	index         int
	inMenu        bool // true = menu, false = jeu
	lastFrameTime time.Time
	frameDelay    time.Duration
	videoEnded    bool
	inventaire    *InventaireGUI
	player        *Personnage
}

// NewGame charge les frames de la vidéo
func NewGame() *Game {
	player := &Personnage{
		Name:      "Héros",
		Money:     100,
		Inventory: []string{"Épée", "Potion"},
	}

	g := &Game{
		frameDelay: time.Millisecond * 42,
		inMenu:     true,

		player:     player,
		inventaire: NewInventaireGUI(player),
	}
	// Chargement des frames vidéo
	totalFrames := 150
	for i := 1; i <= totalFrames; i++ {
		path := fmt.Sprintf("assets/video_frames/frame%03d.png", i)
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

// Lancer la musique en boucle
func playMusic() {
	audioCtx = audio.NewContext(44100)

	data, err := os.ReadFile("assets/menu.mp3") // <-- mets ton mp3 ici
	if err != nil {
		log.Fatal(err)
	}

	stream, err := mp3.Decode(audioCtx, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	loop := audio.NewInfiniteLoop(stream, stream.Length())
	player, err = audio.NewPlayer(audioCtx, loop)
	if err != nil {
		log.Fatal(err)
	}

	player.Play()
}

// Update gère la logique du jeu
func (g *Game) Update() error {
	if g.inventaire != nil {
		g.inventaire.Update()
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if g.inMenu {
			// START NEW GAME
			if x >= 90 && x <= 210 && y >= 520 && y <= 640 {
				fmt.Println("🎮 Start New Game !")
				g.inMenu = false // passe à la map
			}

			// LEAVE
			if x >= 240 && x <= 360 && y >= 520 && y <= 640 {
				fmt.Println("👋 Quitter le jeu...")
				os.Exit(0)
			}
		}
	}

	if !g.inMenu {
		UpdatePlayer()
	}

	if g.inMenu && !g.videoEnded {
		now := time.Now()
		if now.Sub(g.lastFrameTime) >= g.frameDelay {
			g.index++
			if g.index >= len(g.frames) {
				g.index = len(g.frames) - 1
				g.videoEnded = true
			}
			g.lastFrameTime = now
		}
	}
	g.inventaire.Update()

	return nil
}

// Draw affiche le menu ou la map
func (g *Game) Draw(screen *ebiten.Image) {
	if g.inMenu {
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
			ebitenutil.DebugPrint(screen, "SAHARA DEFENDER\nVidéo en cours...")
		} else {
			ebitenutil.DebugPrint(screen, "SAHARA DEFENDER\nFin de la vidéo.\nClique Start ou Leave")
		}
	} else {
		// Quand on quitte le menu -> afficher la map
		DrawMap(screen)
		g.inventaire.Draw(screen)

	}

}

// Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	// Charger la map dès le départ
	LoadMap()

	game := NewGame()

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("SAHARA DEFENDER")

	go func() {
		time.Sleep(500 * time.Millisecond)
		playMusic()
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
