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
				color := rl.White
				if (256-x)*3 > 512 {
					color = rl.Green
				}
				rl.DrawRectangle(y*3, (256-x)*3, 3, 3, color)
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
func (s *System) updateInput() {
	//keys down
	if rl.IsKeyDown(rl.KeyC) {
		s.Machine.Port1 |= 0x01
	}
	if rl.IsKeyDown(rl.KeyOne) {
		s.Machine.Port1 |= 0b00000100
	}
	if rl.IsKeyDown(rl.KeyRight) {
		s.Machine.Port1 |= 0b01000000
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		s.Machine.Port1 |= 0b00100000
	}
	if rl.IsKeyDown(rl.KeySpace) {
		s.Machine.Port1 |= 0b00010000
	}

	//Keys up
	if rl.IsKeyUp(rl.KeyC) {
		var bit uint8 = 0x01
		s.Machine.Port1 &= ^bit
	}
	if rl.IsKeyUp(rl.KeyOne) {
		var bit uint8 = 0b00000100
		s.Machine.Port1 &= ^bit
	}
	if rl.IsKeyUp(rl.KeyRight) {
		var bit uint8 = 0b01000000
		s.Machine.Port1 &= ^bit
	}
	if rl.IsKeyUp(rl.KeyLeft) {
		var bit uint8 = 0b00100000
		s.Machine.Port1 &= ^bit
	}
	if rl.IsKeyUp(rl.KeySpace) {
		var bit uint8 = 0b00010000
		s.Machine.Port1 &= ^bit
	}
}
func (s *System) Start() {
	cycles := 33334
	rl.InitWindow(448+224, 512+256, "Go Invaders")
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		s.Machine.Cpu.Step(cycles)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		s.drawScreen()
		s.updateInput()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
