package main

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

const (
    screenWidth = 720
    screenHeight = 480
)

type Cat struct {
}

func input() {
}

func update() {

}

func drawScene() {

}

func render() {
    rl.BeginDrawing()

    // draw background
    rl.ClearBackground(rl.RayWhite)

    drawScene()

    rl.EndDrawing()
}

func main() {
    rl.InitWindow(screenWidth, screenHeight, "DEMO")
    defer rl.CloseWindow()

    for !rl.WindowShouldClose() {
        input()
        update()
        render()
    }

}
