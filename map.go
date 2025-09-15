package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	mapImage         *ebiten.Image
	playerImage      *ebiten.Image
	playerX, playerY float64 = 400, 300 // position initiale (centre écran)
	playerSpeed      float64 = 3
)

// Charger la map PNG
func LoadMap() {
	img, _, err := ebitenutil.NewImageFromFile("map_v1.png") // <-- mets ton PNG ici
	if err != nil {
		log.Fatal(err)
	}
	mapImage = img

	png, _, err := ebitenutil.NewImageFromFile("player.png") // 👉 ton perso
	if err != nil {
		log.Fatal(err)
	}
	playerImage = png
}

func UpdatePlayer() {
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyZ) { // Z ou W = avancer
		playerY -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		playerY += playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyQ) { // Q ou A = gauche
		playerX -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		playerX += playerSpeed
	}
}

// Dessiner la map
func DrawMap(screen *ebiten.Image) {
	if mapImage != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(mapImage, op)
	}
	if playerImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(playerX, playerY)
		screen.DrawImage(playerImage, op)
	}
}
