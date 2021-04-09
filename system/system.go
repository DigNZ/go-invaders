package system

import (
	core "github.com/DigNZ/goinvaders/core"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type System struct {
	Machine *core.Machine
}

func (s *System) drawScreen() {
	var x, y int32
	for idx := 0x2400; idx < 0x3FFF; idx++ {
		b := s.Machine.Cpu.Memory[idx]

		for i := 0; i < 8; i++ {
			if (b & 0x1) != 0 {
				rl.DrawPixel(x, y, rl.White)
			}
			b = b >> 1
			x++
		}

		if x > 255 {
			x = 0
			y++
		}

	}
}
func (s *System) Start() {
	cycles := 2000000 / 60
	rl.InitWindow(256, 244, "raylib [core] example - basic window")

	for !rl.WindowShouldClose() {
		s.Machine.Cpu.Step(cycles)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		s.drawScreen()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
