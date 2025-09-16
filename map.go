package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	mapImage         *ebiten.Image
	playerX, playerY float64 = 1240, 600

	playerSpeed float64 = 3
	X, Y        float64
	Zoom        float64

	// Sprites par direction
	upSprites    []*ebiten.Image
	downSprites  []*ebiten.Image
	leftSprites  []*ebiten.Image
	rightSprites []*ebiten.Image

	currentSprites []*ebiten.Image
	index          int
	lastUpdate     time.Time
)

func loadAndScale(paths []string, factor float64) []*ebiten.Image {
	images := make([]*ebiten.Image, len(paths))
	for i, path := range paths {
		imgFile, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Fatal(err)
		}

		// Redimensionner l'image
		w, h := imgFile.Size()
		newImg := ebiten.NewImage(int(float64(w)*factor), int(float64(h)*factor))
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(factor, factor)
		newImg.DrawImage(imgFile, op)

		images[i] = newImg
	}
	return images
}

// Nouvelle fonction pour charger tous les sprites
func initSprites() {
	upSprites = loadAndScale([]string{
		"assets/perso/back-step1-.png",
		"assets/perso/back-step2.png",
		"assets/perso/back-fotor.png",
	}, 0.25)

	downSprites = loadAndScale([]string{
		"assets/perso/fromt-step1.png",
		"assets/perso/front-step2.png",
		"assets/perso/front-step3.png",
	}, 0.25)

	leftSprites = loadAndScale([]string{
		"assets/perso/left-step1.png",
		"assets/perso/left-step2.png",
		"assets/perso/left-step3.png",
	}, 0.25)

	rightSprites = loadAndScale([]string{
		"assets/perso/right-step1-.png",
		"assets/perso/right-step3.png",
	}, 0.25)

	// Par défaut, face vers le bas
	currentSprites = downSprites
	lastUpdate = time.Now()
}

func LoadMap() {
	// Charger la map
	img, _, err := ebitenutil.NewImageFromFile("assets/mapz.png")
	if err != nil {
		log.Fatal(err)
	}
	mapImage = img

	// Charger les sprites par direction
	upSprites = loadAndScale([]string{
		"assets/perso/back-step1-.png",
		"assets/perso/back-step2.png",
		"assets/perso/back-fotor.png",
	}, 0.25)

	downSprites = loadAndScale([]string{
		"assets/perso/fromt-step1.png",
		"assets/perso/front-step2.png",
		"assets/perso/front-step3.png",
	}, 0.25)

	leftSprites = loadAndScale([]string{
		"assets/perso/left-step1.png",
		"assets/perso/left-step2.png",
		"assets/perso/left-step3.png",
	}, 0.25)

	rightSprites = loadAndScale([]string{
		"assets/perso/right-step1-.png",
		"assets/perso/right-step3.png",
	}, 0.25)

}

// Fonction pour réduire un tableau d'images
func scaleImages(images []*ebiten.Image, factor float64) []*ebiten.Image {
	scaled := make([]*ebiten.Image, len(images))
	for i, img := range images {
		w, h := img.Size()
		newImg := ebiten.NewImage(int(float64(w)*factor), int(float64(h)*factor))
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(factor, factor)
		newImg.DrawImage(img, op)
		scaled[i] = newImg
	}
	return scaled
}

// Réduire toutes les images à 50%

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
		screenWidth, screenHeight := screen.Size()
		mapWidth, mapHeight := mapImage.Size()
		scaleX := float64(screenWidth) / float64(mapWidth)
		scaleY := float64(screenHeight) / float64(mapHeight)
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		// Appliquer l'échelle
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)

		// Dessiner la map
		screen.DrawImage(mapImage, op)
	}

	// Dessiner le personnage
	if len(currentSprites) > 0 {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(playerX, playerY)
		screen.DrawImage(currentSprites[index], opts)

	}

}
