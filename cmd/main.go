package main

import (
	"fmt"
    "math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var tileDest rl.Rectangle

const (
	screenWidth  = 960
	screenHeight = 960
	tileSize     = 32
	tileNum      = 30
	vel          = 5.0
)

type Cat struct {
	posX      float32
	posY      float32
	direction rl.Vector2
	src       rl.Rectangle
	dest      rl.Rectangle
	texture   rl.Texture2D
}

type Food struct {
	posX    float32
	posY    float32
	src     rl.Rectangle
	dest    rl.Rectangle
	texture rl.Texture2D
}

func handleMovement(c *Cat) {
	if rl.IsKeyDown(rl.KeyW) {
		c.direction.X = 0
		c.direction.Y = -1
	}
	if rl.IsKeyDown(rl.KeyS) {
		c.direction.X = 0
		c.direction.Y = 1
	}
	if rl.IsKeyDown(rl.KeyA) {
		c.direction.X = -1
		c.direction.Y = 0
	}
	if rl.IsKeyDown(rl.KeyD) {
		c.direction.X = 1
		c.direction.Y = 0
	}

	c.dest.X += c.direction.X * vel
	c.dest.Y += c.direction.Y * vel
}

func spawnFood(f *Food) {
    row := rand.IntN(31)
    col := rand.IntN(31)

    f.dest.X = float32(row) * f.dest.Width
    f.dest.Y = float32(col) * f.dest.Height
}

func render(c *Cat, food *Food, grassSprite rl.Texture2D, tileSrc rl.Rectangle) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(147, 211, 196, 255))

	for row := 0; row < tileNum+1; row++ {
		for col := 0; col < tileNum+1; col++ {
			tileDest.X = float32(col) * tileSrc.Width
			tileDest.Y = float32(row) * tileSrc.Height
			rl.DrawTexturePro(grassSprite, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
		}
	}

	rl.DrawTexturePro(c.texture, c.src, c.dest, rl.NewVector2(c.dest.Width, c.dest.Height), 0, rl.White)
	rl.DrawTexturePro(food.texture, food.src, food.dest, rl.NewVector2(food.dest.Width, food.dest.Height), 0, rl.White)

	rl.EndDrawing()
}

func main() {
	fmt.Println("Peak gameplay")
	rl.InitWindow(screenWidth, screenHeight, "DEMO")
	defer rl.CloseWindow()
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	tileDest.Width = 48
	tileDest.Height = 48
	tileSrc := rl.NewRectangle(0, 0, 48, 48)
	grassSprite := rl.LoadTexture("./assets/grass.png")
	defer rl.UnloadTexture(grassSprite)

	cat := new(Cat)
	cat.direction = rl.Vector2{X: 1, Y: 0} // default direction: right
	cat.src = rl.NewRectangle(0, 0, 40, 40)
	cat.dest = rl.NewRectangle(200, 200, 48, 48)
	cat.texture = rl.LoadTexture("./assets/Block.png")
	defer rl.UnloadTexture(cat.texture)

	food := new(Food)
	food.src = rl.NewRectangle(0, 0, 48, 32)
    food.dest.Width = 48
    food.dest.Height = 32
	food.texture = rl.LoadTexture("./assets/fish.png")
	defer rl.UnloadTexture(food.texture)

    spawnFood(food)
	for !rl.WindowShouldClose() {
		handleMovement(cat)
		render(cat, food, grassSprite, tileSrc)
	}
}
