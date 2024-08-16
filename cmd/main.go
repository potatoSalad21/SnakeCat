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
	blocks  []*CatBlock
	src     rl.Rectangle
	dir     rl.Vector2
	texture map[string]rl.Texture2D
}

type CatBlock struct {
	vec rl.Vector2
	dir rl.Vector2
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
	head := c.blocks[0].vec
	if head.X > screenWidth || head.X <= 0 || head.Y > screenHeight || head.Y <= 0 {
		fmt.Println("+ PLAYER DIED")
		c.dead = true
		time.Sleep(2 * time.Second)
		c.spawnCat()
	}
}

func (c *Cat) checkCollision() {
	headRec := rl.NewRectangle(c.blocks[0].vec.X, c.blocks[0].vec.Y, tileSize, tileSize)
	for _, block := range c.blocks[1:] {
		if rl.CheckCollisionPointRec(block.vec, headRec) {
			c.dead = true
			time.Sleep(2 * time.Second)
			c.spawnCat()
		}
	}
}

func (c *Cat) spawnCat() {
	c.blocks = nil
	var mapMiddle float32 = tileSize * (tileNum / 2)
	for i := 0; i < 3; i++ {
		block := CatBlock{
			vec: rl.Vector2{
				X: mapMiddle - float32(tileSize*i),
				Y: mapMiddle,
			},
			dir: rl.Vector2{ // default direction: right
				X: 1,
				Y: 0,
			},
		}
		c.blocks = append(c.blocks, &block)
	}

	score = 0
	c.dead = false
}

func (c *Cat) draw() {
	for i, block := range c.blocks {
		src := c.src
		dest := rl.NewRectangle(block.vec.X, block.vec.Y, tileSize, tileSize)
		var texture rl.Texture2D

		if i == 0 {
			texture = c.texture["head"]
		} else if i == len(c.blocks)-1 {
			if block.dir.X == 0 {
				texture = c.texture["tailV"]
			} else if block.dir.Y == 0 {
				texture = c.texture["tailH"]
			}
		} else {
			if block.dir.X == 0 {
				texture = c.texture["bodyV"]
			} else if block.dir.Y == 0 {
				texture = c.texture["bodyH"]
			}
		}

		if block.dir.Y == 0 {
			src.Width *= -block.dir.X
		} else {
			src.Height *= -block.dir.Y
		}

		rl.DrawTexturePro(texture, src, dest, rl.NewVector2(tileSize, tileSize), 0, rl.White)
	}
}

func (c *Cat) grow() {
	score++
	tail := c.blocks[len(c.blocks)-1]
	block := CatBlock{
		vec: rl.Vector2{
			X: tail.vec.X + tail.dir.X*float32(tileSize),
			Y: tail.vec.Y + tail.dir.Y*float32(tileSize),
		},
		dir: rl.Vector2{
			X: tail.dir.X,
			Y: tail.dir.Y,
		},
	}

	c.blocks = append(c.blocks, &block)
}

func (c *Cat) move() {
	c.blocks = c.blocks[:len(c.blocks)-1]
	head := c.blocks[0]
	block := CatBlock{
		vec: rl.Vector2{
			X: head.vec.X + head.dir.X*float32(tileSize),
			Y: head.vec.Y + head.dir.Y*float32(tileSize),
		},
		dir: rl.Vector2{
			X: head.dir.X,
			Y: head.dir.Y,
		},
	}
	c.blocks = slices.Insert(c.blocks, 0, &block)

	c.checkOutOfBounds()
	c.checkCollision()
}

func handleMovement(c *Cat) {
	head := c.blocks[0]
	if rl.IsKeyDown(rl.KeyW) && head.dir.Y != 1 {
		head.dir.X = 0
		head.dir.Y = -1
	}
	if rl.IsKeyDown(rl.KeyS) && head.dir.Y != -1 {
		head.dir.X = 0
		head.dir.Y = 1
	}
	if rl.IsKeyDown(rl.KeyA) && head.dir.X != 1 {
		head.dir.X = -1
		head.dir.Y = 0
	}
	if rl.IsKeyDown(rl.KeyD) && head.dir.X != -1 {
		head.dir.X = 1
		head.dir.Y = 0
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
	if rl.CheckCollisionRecs(rl.NewRectangle(head.vec.X, head.vec.Y, tileSize, tileSize), f.dest) {
		f.spawnFood()
		c.grow()
	}

	c.draw()
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

	cat.texture = make(map[string]rl.Texture2D)
	cat.texture["head"] = rl.LoadTexture("./assets/cathead.png")
	cat.texture["bodyV"] = rl.LoadTexture("./assets/catbodyV.png")
	cat.texture["bodyH"] = rl.LoadTexture("./assets/catbodyH.png")
	cat.texture["tailV"] = rl.LoadTexture("./assets/catbuttV.png")
	cat.texture["tailH"] = rl.LoadTexture("./assets/catbuttH.png")
	defer rl.UnloadTexture(cat.texture["head"])
	defer rl.UnloadTexture(cat.texture["bodyV"])
	defer rl.UnloadTexture(cat.texture["bodyH"])
	defer rl.UnloadTexture(cat.texture["tailV"])
	defer rl.UnloadTexture(cat.texture["tailH"])

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
	close(done)
}
