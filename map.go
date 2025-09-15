package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	mapImage         *ebiten.Image
	playerX, playerY float64 = 400, 300
	playerSpeed      float64 = 3
	X, Y             float64
	Zoom             float64

	// Sprites par direction
	upSprites    []*ebiten.Image
	downSprites  []*ebiten.Image
	leftSprites  []*ebiten.Image
	rightSprites []*ebiten.Image

	currentSprites []*ebiten.Image
	index          int
	lastUpdate     time.Time
)

func LoadMap() {
	// Charger la map
	img, _, err := ebitenutil.NewImageFromFile("assets/mapdV.png")
	if err != nil {
		log.Fatal(err)
	}
	mapImage = img

	// Charger les sprites par direction
	upSprites = loadImages([]string{
		"assets/perso/back-step1.png",
		"assets/perso/back-step2.png",
		"assets/perso/back.png",
	})

	downSprites = loadImages([]string{
		"assets/perso/front.png",
		"assets/perso/front.png", // tu peux dupliquer pour plus de frames
	})

	//leftSprites = loadImages([]string{
	//	"assets/perso/left-step1.png",
	//	"assets/perso/left-step2.png",
	//})

	//rightSprites = loadImages([]string{
	//	"assets/perso/right-step1.png",
	//	"assets/perso/right-step2.png",
	//})

	// Par défaut, face vers le bas
	currentSprites = downSprites
	lastUpdate = time.Now()
}

func loadImages(paths []string) []*ebiten.Image {
	var imgs []*ebiten.Image
	for _, f := range paths {
		img, _, err := ebitenutil.NewImageFromFile(f)
		if err != nil {
			log.Fatal(err)
		}
		imgs = append(imgs, img)
	}
	return imgs
}

func UpdatePlayer() {
	moving := false

	// Déplacement et direction
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyZ) {
		playerY -= playerSpeed
		currentSprites = upSprites
		moving = true
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		playerY += playerSpeed
		currentSprites = downSprites
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		playerX -= playerSpeed
		currentSprites = leftSprites
		moving = true
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		playerX += playerSpeed
		currentSprites = rightSprites
		moving = true
	}

	// Animation : avancer seulement si le personnage bouge
	if moving && time.Since(lastUpdate) > 150*time.Millisecond {
		index++
		if index >= len(currentSprites) {
			index = 0
		}
		lastUpdate = time.Now()
	} else if !moving {
		// Reset sur la frame de repos quand il ne bouge pas
		index = 0
	}
}

func DrawMap(screen *ebiten.Image) {
	
	// Dessiner la map
	if mapImage != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(mapImage, op)
	}

	// Dessiner le personnage
	if len(currentSprites) > 0 {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(playerX, playerY)
		screen.DrawImage(currentSprites[index], opts)

	}
}
