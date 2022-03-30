/************************
Created by: Mitchell Noun
Date created: 3/25/22
Class: COMP415 Emerging Languages
Assignment: Project 3
*************************/
package main

import (
	"embed"
	"fmt"
	_ "github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	_ "github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	_ "github.com/blizzy78/ebitenui/widget"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font/basicfont"
	Image "image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
)

//go:embed assets/*
var EmbeddedAssets embed.FS
var score = 0
var playbutton *widget.Button
var quitbutton *widget.Button

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
	for num, enemy := range g.enemies {
		if isColliding(g.player, enemy) == true {
			removeEnemy(g.enemies, num)
			score += 1
		} else {
			g.drawOps.GeoM.Reset()
			g.drawOps.GeoM.Translate(float64(enemy.xloc), float64(enemy.yloc))
			screen.DrawImage(enemy.pict, &g.drawOps)
		}
	}
	if score == 10 {
		playbutton.SetLocation(Image.Rectangle{
			Min: Image.Point{
				X: 0,
				Y: 0},
			Max: Image.Point{
				X: 400,
				Y: 500,
			},
		})
		playbutton.Text().Label = "Play Again"
		playbutton.Text().SetLocation(Image.Rectangle{
			Min: Image.Point{
				X: 0,
				Y: 0},
			Max: Image.Point{
				X: 400,
				Y: 500,
			},
		})
		quitbutton.SetLocation(Image.Rectangle{
			Min: Image.Point{
				X: 0,
				Y: 0},
			Max: Image.Point{
				X: 600,
				Y: 500,
			},
		})
		quitbutton.Text().Label = "Quit"
		quitbutton.Text().SetLocation(Image.Rectangle{
			Min: Image.Point{
				X: 0,
				Y: 0},
			Max: Image.Point{
				X: 600,
				Y: 500,
			},
		})
	}
	message := fmt.Sprintf("Score:%d", score)
	ebitenutil.DebugPrint(screen, message)
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
	game.enemies = populateEnemy(game, width, height)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal("Oh no! something terrible happened and the game crashed", err)
	}
	idle, err := loadImageNineSlice("graphics/button-idle.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	hover, err := loadImageNineSlice("graphics/button-hover.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	pressed, err := loadImageNineSlice("graphics/button-pressed.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	disabled, err := loadImageNineSlice("graphics/button-disabled.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	buttonImage := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
	quitbutton = widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Quit", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		widget.ButtonOpts.ClickedHandler(quit),
	)
	playbutton = widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Play Again", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.ClickedHandler(playAgain),
	)
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

func isColliding(player, enemy Sprite) bool {
	enemyWidth, enemyHeight := enemy.pict.Size()
	playerWidth, playerHeight := player.pict.Size()
	if player.xloc < enemy.xloc+enemyWidth &&
		player.xloc+playerWidth > enemy.xloc &&
		player.yloc < enemy.yloc+enemyHeight &&
		player.yloc+playerHeight > enemy.yloc {
		return true
	}
	return false
}

func populateEnemy(g Game, width, height int) []Sprite {
	for i := 0; i <= 10; i++ {
		s := Sprite{
			pict: loadPNGImageFromEmbedded("EnemySprite.png"),
			xloc: rand.Intn(GameWidth - width),
			yloc: rand.Intn(GameHeight - height),
			dX:   0,
			dY:   0,
		}
		g.enemies = append(g.enemies, s)
	}
	return g.enemies
}

func removeEnemy(s []Sprite, index int) []Sprite {
	return append(s[:index], s[index+1:]...)
}

func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i := loadPNGImageFromEmbedded(path)

	w, h := i.Size()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}

func playAgain(args *widget.ButtonClickedEventArgs) {
	score = 0
	g := Game{}
	width, height := (*ebiten.Image).Size(loadPNGImageFromEmbedded("EnemySprite.png"))
	g.enemies = populateEnemy(g, width, height)
}

func quit(args *widget.ButtonClickedEventArgs) {
	os.Exit(1)
}
