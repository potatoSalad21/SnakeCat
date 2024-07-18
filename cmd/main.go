package main

import (
	"fmt"
    "math/rand/v2"
    "time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var tileDest rl.Rectangle

const (
	screenWidth  = 960
	screenHeight = 960
	tileSize     = 48
	tileNum      = 20
)

type Cat struct {
    src       rl.Rectangle
    dest      rl.Rectangle
	direction rl.Vector2
	texture   rl.Texture2D
}

type Food struct {
	src     rl.Rectangle
	dest    rl.Rectangle
	texture rl.Texture2D
}

func spawnFood(f *Food) {
    row := rand.IntN(tileNum)
    col := rand.IntN(tileNum)

    f.dest.X = float32(row) * tileSize
    f.dest.Y = float32(col) * tileSize
}

func (c *Cat) checkOutOfBounds() {
    if c.dest.X > screenWidth || c.dest.X < 0 || c.dest.Y > screenHeight || c.dest.Y < 0 {
        fmt.Println("+ PLAYER DIED")
        // TODO: display death screen
        c.respawn()
    }
}

func (c *Cat) respawn() {
    c.dest.X = tileSize * (tileNum / 2)
    c.dest.Y = tileSize * (tileNum / 2)
    c.direction = rl.Vector2{X: 1, Y: 0}
    // TODO: reset cat size
    // TODO: reset game score
}

func (c *Cat) move() {
    dir := c.direction

    switch {
    case dir.X == 1 && dir.Y == 0:
        c.dest.X += tileSize
    case dir.X == -1 && dir.Y == 0:
        c.dest.X -= tileSize
    case dir.X == 0 && dir.Y == 1:
        c.dest.Y += tileSize
    case dir.X == 0 && dir.Y == -1:
        c.dest.Y -= tileSize
    }
    c.checkOutOfBounds()
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
	cat.dest = rl.NewRectangle(tileSize * (tileNum / 2), tileSize * (tileNum / 2), 48, 48)
	cat.texture = rl.LoadTexture("./assets/Block.png")
	defer rl.UnloadTexture(cat.texture)

	food := new(Food)
	food.src = rl.NewRectangle(0, 0, 48, 32)
    food.dest.Width = tileSize
    food.dest.Height = tileSize
	food.texture = rl.LoadTexture("./assets/fish.png")
	defer rl.UnloadTexture(food.texture)

    ticker := time.NewTicker(200 * time.Millisecond)
    done := make(chan bool)
    go func() {
        for {
            select {
            case <- done:
                return
            case <- ticker.C:
                cat.move()
            }
        }
    }()

    spawnFood(food)
	for !rl.WindowShouldClose() {
		handleMovement(cat)
		render(cat, food, grassSprite, tileSrc)
	}
    ticker.Stop()
}
