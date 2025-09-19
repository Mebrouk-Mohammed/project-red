package source

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
	audioCtx     *audio.Context
	player       *audio.Player
	gameInstance *Game
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
	marchand      *MenuMarchand

	camera Camera
}

// Camera gère la position et le zoom de la vue
type Camera struct {
	X, Y float64 // Position de la caméra
	Zoom float64 // Facteur de zoom
}

// NewGame charge les frames de la vidéo
func NewGame() *Game {
	player := &Personnage{
		Name:      "Héros",
		Life:      100,
		MaxLife:   100,
		Shield:    0,
		MaxShield: 100, // valeur de base
		Strength:  10,
		Money:     10000,
		Inventory: []string{},
	}

	g := &Game{
		frameDelay: time.Millisecond * 42,
		inMenu:     true,

		player: player,
		inventaire: &InventaireGUI{
			player: player,
		},
		marchand: NewMenuMarchand(player),
		camera: Camera{
			X:    0,
			Y:    0,
			Zoom: 1.0, // zoom normal
		},
	}

	// Chargement des frames vidéo
	totalFrames := 150
	for i := 1; i <= totalFrames; i++ {
		path := fmt.Sprintf("source/assets/video_frames/frame%03d.png", i)
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

	data, err := os.ReadFile("source/assets/menu.mp3") // <-- mets ton mp3 ici
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
	// Gestion du combat
	if inCombat {
		UpdateCombat()
		return nil
	}

	// Gestion des entrées clavier/souris
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.inMenu {
		x, y := ebiten.CursorPosition()
		// Bouton Start
		if x >= 90 && x <= 210 && y >= 520 && y <= 640 {
			fmt.Println("🎮 Start New Game !")
			g.inMenu = false
		}
		// Bouton Quitter
		if x >= 240 && x <= 360 && y >= 520 && y <= 640 {
			fmt.Println("👋 Quitter le jeu...")
			os.Exit(0)
		}
	}

	// Mise à jour du joueur si hors menu
	if !g.inMenu {
		UpdatePlayer()
	}

	// Animation vidéo menu
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

	// Gestion du zoom
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyKPAdd), ebiten.IsKeyPressed(ebiten.KeyEqual):
		g.camera.Zoom += 0.01
	case ebiten.IsKeyPressed(ebiten.KeyKPSubtract), ebiten.IsKeyPressed(ebiten.KeyMinus):
		if g.camera.Zoom > 0.2 {
			g.camera.Zoom -= 0.01
		}
	}

	// Déplacement caméra
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.camera.Y -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.camera.Y += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.camera.X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.camera.X += 5
	}

	// Mise à jour des entités
	UpdateMonsters()
	CheckCollisionWithPlayerCombat()
	if g.inventaire != nil {
		g.inventaire.Update()
	}
	if g.marchand != nil {
		g.marchand.Update()
	}

	return nil
}

// Draw affiche le menu ou la map
func (g *Game) Draw(screen *ebiten.Image) {
	if g.inMenu {
		// Affichage du menu vidéo
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
		if inCombat {
			DrawCombatScreen(screen)
			return
		}
	} else {
		// Affichage principal hors menu
		DrawMap(screen)
		g.marchand.Draw(screen)
		g.player.DrawBars(screen)
		DrawMonsters(screen)
		DrawCombatMessage(screen)
		DrawCombatScreen(screen)
		g.inventaire.Draw(screen)
	}

}

// Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Layout de la fenêtre (conserve la taille demandée)
	return outsideWidth, outsideHeight
}

func Main() {
	// Initialisation du jeu
	LoadMap()      // Charge la map
	InitMonsters() // Initialise les monstres

	game := NewGame() // Crée l'instance principale
	gameInstance = game

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("SAHARA DEFENDER")

	// Lancement de la musique en arrière-plan
	go func() {
		time.Sleep(500 * time.Millisecond)
		playMusic()
	}()

	// Démarrage de la boucle principale Ebiten
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
