package io

type Machine struct {
	shift0, shift1, shift_offset uint8
}

func (m *Machine) MachineIN(port uint8) uint8 {

	var a uint8
	switch port {
	case 3:
		v := (uint16(m.shift1) << 8) | uint16(m.shift0)
		a = uint8((v >> (8 - m.shift_offset)) & 0xff)
	}
	return a

}

func (m *Machine) MachineOUT(port uint8) {

}
