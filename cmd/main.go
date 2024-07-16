package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 720
	screenHeight = 480
    vel = 10.0
)

type Cat struct {
	posX    float32
	posY    float32
	src     rl.Rectangle
	dest    rl.Rectangle
	texture rl.Texture2D
}

func input(c *Cat) {
	if rl.IsKeyDown(rl.KeyW) {
		c.dest.Y -= vel
	}
	if rl.IsKeyDown(rl.KeyS) {
		c.dest.Y += vel
	}
	if rl.IsKeyDown(rl.KeyA) {
		c.dest.X -= vel
	}
	if rl.IsKeyDown(rl.KeyD) {
		c.dest.X += vel
	}
}

func render(c *Cat) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.NewColor(147, 211, 196, 255))
	rl.DrawTexturePro(c.texture, c.src, c.dest, rl.NewVector2(c.dest.Width, c.dest.Height), 0, rl.White)

	rl.EndDrawing()
}

func main() {
	fmt.Println("Peak gameplay")
	rl.InitWindow(screenWidth, screenHeight, "DEMO")
	defer rl.CloseWindow()

	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	cat := new(Cat)
	cat.texture = rl.LoadTexture("./assets/Block.png")
	defer rl.UnloadTexture(cat.texture)

	cat.src = rl.NewRectangle(0, 0, 40, 40)
	cat.dest = rl.NewRectangle(200, 200, 48, 48)

	for !rl.WindowShouldClose() {
		input(cat)
		render(cat)
	}
}
