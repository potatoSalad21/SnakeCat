package main

import (
    "fmt"

    rl "github.com/gen2brain/raylib-go/raylib"
)

const (
    screenWidth = 720
    screenHeight = 480
)

type Cat struct {
    posX int32
    posY int32
}

func input(cat *Cat) {
    if rl.IsKeyDown(rl.KeyW) {
        cat.posY -= 10
    }
    if rl.IsKeyDown(rl.KeyS) {
        cat.posY += 10
    }
    if rl.IsKeyDown(rl.KeyA) {
        cat.posX -= 10
    }
    if rl.IsKeyDown(rl.KeyD) {
        cat.posX += 10
    }
}

func update() {

}

func drawScene() {

}

func render(c *Cat) {
    rl.BeginDrawing()

    // draw background
    rl.ClearBackground(rl.NewColor(147, 211, 196, 255))

    rl.DrawRectangle(c.posX, c.posY, 64, 64, rl.NewColor(255, 0, 0, 255))

    rl.EndDrawing()
}

func main() {
    fmt.Println("Peak gameplay")
    rl.InitWindow(screenWidth, screenHeight, "DEMO")
    defer rl.CloseWindow()

    rl.SetExitKey(0)
    rl.SetTargetFPS(60)

    cat := new(Cat)

    for !rl.WindowShouldClose() {
        input(cat)
        update()
        render(cat)
    }
}
