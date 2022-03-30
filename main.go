/************************
Created by: Mitchell Noun
Date created: 3/25/22
Class: COMP415 Emerging Languages
Assignment: Project 3
*************************/
package main

import (
	"embed"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/png"
	"log"
	"math/rand"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

const (
	GameWidth   = 1000
	GameHeight  = 1000
	PlayerSpeed = 2
)

type Sprite struct {
	pict *ebiten.Image
	xloc int
	yloc int
	dX   int
	dY   int
}

type Game struct {
	player   Sprite
	enemies  []Sprite
	score    int
	drawOps  ebiten.DrawImageOptions
	collided bool
}

func (g *Game) Update() error {
	processPlayerInput(g)
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	g.drawOps.GeoM.Reset()
	g.drawOps.GeoM.Translate(float64(g.player.xloc), float64(g.player.yloc))
	screen.DrawImage(g.player.pict, &g.drawOps)
	for i := 0; i <= 10; i++ {
		if !g.collided {
			g.drawOps.GeoM.Reset()
			g.drawOps.GeoM.Translate(float64(g.enemies[i].xloc), float64(g.enemies[i].yloc))
			screen.DrawImage(g.enemies[i].pict, &g.drawOps)
		}
	}
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

func main() {
	ebiten.SetWindowSize(GameWidth, GameHeight)
	ebiten.SetWindowTitle("Project 3: Ebiten Game")
	game := Game{score: 0}
	game.player = Sprite{
		pict: loadPNGImageFromEmbedded("PlayerSprite.png"),
		xloc: 200,
		yloc: 300,
		dX:   0,
		dY:   0,
	}
	game.enemies = []Sprite{}
	width, height := (*ebiten.Image).Size(loadPNGImageFromEmbedded("EnemySprite.png"))
	for i := 0; i <= 10; i++ {
		s := Sprite{
			pict: loadPNGImageFromEmbedded("EnemySprite.png"),
			xloc: rand.Intn(GameWidth - width),
			yloc: rand.Intn(GameHeight - height),
			dX:   0,
			dY:   0,
		}
		game.enemies = append(game.enemies, s)
	}
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal("Oh no! something terrible happened and the game crashed", err)
	}
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("assets")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("assets/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func processPlayerInput(theGame *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		theGame.player.dY = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		theGame.player.dY = PlayerSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		theGame.player.dY = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		theGame.player.dX = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		theGame.player.dX = PlayerSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		theGame.player.dX = 0
	}
	theGame.player.xloc += theGame.player.dX
	theGame.player.yloc += theGame.player.dY
	if theGame.player.yloc <= 0 {
		theGame.player.dY = 0
		theGame.player.yloc = 0
	} else if theGame.player.yloc+theGame.player.pict.Bounds().Size().Y > GameHeight {
		theGame.player.dY = 0
		theGame.player.yloc = GameHeight - theGame.player.pict.Bounds().Size().Y
	}
	if theGame.player.xloc <= 0 {
		theGame.player.dX = 0
		theGame.player.xloc = 0
	} else if theGame.player.xloc+theGame.player.pict.Bounds().Size().X > GameWidth {
		theGame.player.dX = 0
		theGame.player.xloc = GameWidth - theGame.player.pict.Bounds().Size().X
	}
}

/*
func isColliding(player, enemy Sprite) bool {
	enemyWidth, enemyHeight := enemy.pict.Size()
	playerWidth, playerHeight := player.pict.Size()
	if player.xLoc < enemy.xLoc+enemyWidth &&
		player.xLoc+playerWidth > enemy.xLoc &&
		player.yLoc < gold.yLoc+enemyHeight &&
		player.yLoc+playerHeight > enemy.yLoc {
		return true
	}
	return false
}
*/
