package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
    "strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
    tileDest rl.Rectangle
    score    int
)

const (
	screenWidth  = 960
	screenHeight = 960
	tileSize     = 48
	tileNum      = 20
)

type Cat struct {
	dead    bool
	blocks  []rl.Vector2
	src     rl.Rectangle
	dir     rl.Vector2
	texture rl.Texture2D
}

type Food struct {
	src     rl.Rectangle
	dest    rl.Rectangle
	texture rl.Texture2D
}

func (f *Food) spawnFood() {
	row := rand.IntN(tileNum) + 1
	col := rand.IntN(tileNum) + 1

	f.dest.X = float32(row) * tileSize
	f.dest.Y = float32(col) * tileSize
}

func (c *Cat) checkOutOfBounds() {
	head := c.blocks[0]
	if head.X > screenWidth || head.X <= 0 || head.Y > screenHeight || head.Y <= 0 {
		fmt.Println("+ PLAYER DIED")
		c.dead = true
		time.Sleep(2 * time.Second)
		c.spawnCat()
	}
}

func (c *Cat) checkCollision() {
    headRec := rl.NewRectangle(c.blocks[0].X, c.blocks[0].Y, tileSize, tileSize)
    for _, block := range c.blocks[1:] {
        if rl.CheckCollisionPointRec(block, headRec) {
            c.dead = true
            time.Sleep(2 * time.Second)
            c.spawnCat()
        }
    }
}

func (c *Cat) spawnCat() {
	c.dir = rl.Vector2{X: 1, Y: 0} // default direction: right
	var mapMiddle float32 = tileSize * (tileNum / 2)
	c.blocks = []rl.Vector2{
		{X: mapMiddle, Y: mapMiddle},
		{X: mapMiddle - tileSize, Y: mapMiddle},
		{X: mapMiddle - 2*tileSize, Y: mapMiddle}}

    score = 0
	c.dead = false
}

func (c *Cat) grow() {
    score++
	tail := c.blocks[len(c.blocks)-1]
	c.blocks = append(c.blocks,
		rl.Vector2{X: tail.X + c.dir.X*float32(tileSize), Y: tail.Y + c.dir.Y*float32(tileSize)})
}

func (c *Cat) move() {
	dir := c.dir
	c.blocks = c.blocks[:len(c.blocks)-1]
	head := c.blocks[0]
	c.blocks = slices.Insert(c.blocks, 0, rl.Vector2{X: head.X + dir.X*tileSize, Y: head.Y + dir.Y*tileSize})

	c.checkOutOfBounds()
    c.checkCollision()
}

func handleMovement(c *Cat) {
	if rl.IsKeyDown(rl.KeyW) && c.dir.Y != 1 {
		c.dir.X = 0
		c.dir.Y = -1
	}
	if rl.IsKeyDown(rl.KeyS) && c.dir.Y != -1 {
		c.dir.X = 0
		c.dir.Y = 1
	}
	if rl.IsKeyDown(rl.KeyA) && c.dir.X != 1 {
		c.dir.X = -1
		c.dir.Y = 0
	}
	if rl.IsKeyDown(rl.KeyD) && c.dir.X != -1 {
		c.dir.X = 1
		c.dir.Y = 0
	}
}

func render(c *Cat, f *Food, grassSprite rl.Texture2D, tileSrc rl.Rectangle) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(147, 211, 196, 255))

	for row := 0; row < tileNum+1; row++ {
		for col := 0; col < tileNum+1; col++ {
			tileDest.X = float32(col) * tileSrc.Width
			tileDest.Y = float32(row) * tileSrc.Height
			rl.DrawTexturePro(grassSprite, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
		}
	}
	if c.dead {
		rl.DrawText("YOU DIED", tileSize*(tileNum/2), tileSize*(tileNum/2), 72, rl.NewColor(255, 0, 0, 255))
	}

	head := c.blocks[0]
	if rl.CheckCollisionRecs(rl.NewRectangle(head.X, head.Y, tileSize, tileSize), f.dest) {
		f.spawnFood()
		c.grow()
	}

	for _, block := range c.blocks {
		dest := rl.NewRectangle(block.X, block.Y, tileSize, tileSize)
		rl.DrawTexturePro(c.texture, c.src, dest, rl.NewVector2(tileSize, tileSize), 0, rl.White)
	}
	rl.DrawTexturePro(f.texture, f.src, f.dest, rl.NewVector2(tileSize, tileSize), 0, rl.White)
    rl.DrawText(strconv.Itoa(score), 10, 0, 64, rl.NewColor(200, 100, 20, 255))

	rl.EndDrawing()
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Stretchy Cats")
	defer rl.CloseWindow()
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	tileDest.Width = 48
	tileDest.Height = 48
	tileSrc := rl.NewRectangle(0, 0, 48, 48)
	grassSprite := rl.LoadTexture("./assets/grass.png")
	defer rl.UnloadTexture(grassSprite)

	cat := new(Cat)
	cat.src = rl.NewRectangle(0, 0, 40, 40)
	cat.spawnCat()
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
			case <-done:
				return
			case <-ticker.C:
				cat.move()
			}
		}
	}()

	food.spawnFood()
	for !rl.WindowShouldClose() {
		handleMovement(cat)
		render(cat, food, grassSprite, tileSrc)
	}
	ticker.Stop()
}
