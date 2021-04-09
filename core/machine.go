package core

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
	shift0, shift1, shift_offset, port1, port2 uint8
	Cpu                                        *State8080
}

func (m *Machine) Init() {
	m.Cpu = &State8080{}
	m.Cpu.Init(m)
}

func (m *Machine) MachineIN(port uint8) uint8 {

	var a uint8
	switch port {
	case 0:
		return 1
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
	case 4:
		m.shift0 = m.shift1
		m.shift1 = value
	}
}

func (m *Machine) KeyDown() {

}

func (m *Machine) KeyUp() {

}
