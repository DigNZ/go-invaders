package core

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Player2Start = "PLAYER2START"
	Player1Start = "PLAYER1START"
	Player1Left  = "PLAYER1LEFT"
	Player1Right = "PLAYER1RIGHT"
	Player2Left  = "PLAYER2LEFT"
	Player2Right = "PLAYER2RIGHT"
	Player1Shoot = "PLAYER1SHOOT"
	Player2Shoot = "PLAYER2SHOOT"
)

type Machine struct {
	shift0, shift1, shift_offset, Port1, Port2, outport3, outport5, lastOutport3, lastOutport5 uint8
	Cpu                                                                                        *State8080
	ufo                                                                                        rl.Music
	audioPath                                                                                  string
	samples                                                                                    []rl.Sound
}

func (m *Machine) Init(path string) {
	for i := 0; i < 8; i++ {
		snd := rl.LoadSound(fmt.Sprintf("%s/%d.wav", path, i))
		m.samples = append(m.samples, snd)
	}
	m.Cpu = &State8080{}
	m.Cpu.Init(m)
}

func (m *Machine) MachineIN(port uint8) uint8 {

	var a uint8
	switch port {
	case 0:
		return 1
	case 1:
		return m.Port1
	case 2:
		return m.Port2
	case 3:
		v := (uint16(m.shift1) << 8) | uint16(m.shift0)
		a = uint8((v >> (8 - m.shift_offset)) & 0xff)
	}
	return a

}

func (m *Machine) MachineOUT(port, value uint8) {
	switch port {
	case 2:
		m.shift_offset = value & 0x7
	case 3:
		m.outport3 = value
	case 4:
		m.shift0 = m.shift1
		m.shift1 = value
	case 5:
		m.outport5 = value
	}
}
func (m *Machine) playSample(id int) {
	fmt.Printf("Play Sample %d", id)
	rl.PlaySound(m.samples[id])
}
func (m *Machine) PlaySound() {
	if m.outport3 != m.lastOutport3 {
		if (m.outport3&0x1) == 1 && (m.lastOutport3&0x1) == 0 {
			//start UFO
			m.ufo = rl.LoadMusicStream(fmt.Sprintf("%s/0.wav", m.audioPath))
			rl.PlayMusicStream(m.ufo)
		} else if (m.outport3&0x1) == 0 && (m.lastOutport3&0x1) == 1 {
			rl.StopMusicStream(m.ufo)
		}

		if (m.outport3&0x2) == 1 && (m.lastOutport3&0x2) == 0 {
			m.playSample(1)
		}

		if (m.outport3&0x4) == 1 && (m.lastOutport3&0x4) == 0 {
			m.playSample(2)
		}

		if (m.outport3&0x8) == 1 && (m.lastOutport3&0x8) == 0 {
			m.playSample(3)
		}

		m.lastOutport3 = m.outport3
	}
	if m.outport5 != m.lastOutport5 {
		if (m.outport5&0x1) == 1 && (m.lastOutport5&0x1) == 0 {
			m.playSample(4)
		}

		if (m.outport5&0x2) == 1 && (m.lastOutport5&0x2) == 0 {
			m.playSample(5)
		}

		if (m.outport5&0x4) == 1 && (m.lastOutport5&0x4) == 0 {
			m.playSample(6)
		}

		if (m.outport5&0x8) == 1 && (m.lastOutport5&0x8) == 0 {
			m.playSample(7)
		}

		if (m.outport5&0x10) == 1 && (m.lastOutport5&0x10) == 0 {
			m.playSample(8)
		}

		m.lastOutport5 = m.outport5
	}

}
