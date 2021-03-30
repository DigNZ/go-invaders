package core

import "fmt"

type ConditionCodes struct {
	Z, S, P, CY, AC bool
	PAD             uint8
}
type State8080 struct {
	A, B, C, D, E, H, L uint8
	SP, PC              uint16
	Memory              [0x10000]uint8
	ConditionCodes      ConditionCodes
	IntEnable           uint8
}

func (s *State8080) UnimplementedInstruction(opcode byte) {
	panic(fmt.Sprintf("Error: Unimplemented Instruction: %02X", opcode))
}

func (s *State8080) InitWithData(data []byte) {
	for idx, d := range data {
		s.Memory[idx] = d
	}
}
func parity(x, size uint8) bool {
	var i uint8
	p := 0
	x = (x & ((1 << size) - 1))
	for i = 0; i < size; i++ {
		if (x & 0x1) != 0 {
			p++
		}
		x = x >> 1
	}
	return (0 == (p & 0x1))

}
func (s *State8080) Emulate8080Op() {
	opcode := s.Memory[s.PC]
	data := s.Memory[s.PC+1 : s.PC+3]
	Disassemble8080Op(s.Memory[:], s.PC)
	//Remember that we've advanced the PC BEFORE executing the opcode.
	s.PC += 1
	switch opcode {
	case 0x00:
		// NOP
	case 0x01:
		s.UnimplementedInstruction(opcode)
	case 0x02:
		s.UnimplementedInstruction(opcode)
	case 0x03:
		s.UnimplementedInstruction(opcode)
	case 0x04:
		s.UnimplementedInstruction(opcode)
	case 0x05:
		res := s.B - 1
		s.ConditionCodes.Z = res == 0
		s.ConditionCodes.S = 0x80 == (res & 0x80)
		s.ConditionCodes.P = parity(res, 8)
		s.B = res
	case 0x06:
		s.B = s.Memory[s.PC]
		s.PC++
	case 0x07:
		s.UnimplementedInstruction(opcode)
	case 0x08:
		s.UnimplementedInstruction(opcode)
	case 0x09:
		s.UnimplementedInstruction(opcode)
	case 0x0a:
		s.UnimplementedInstruction(opcode)
	case 0x0b:
		s.UnimplementedInstruction(opcode)
	case 0x0c:
		s.UnimplementedInstruction(opcode)
	case 0x0d:
		s.UnimplementedInstruction(opcode)
	case 0x0e:
		s.UnimplementedInstruction(opcode)
	case 0x0f:
		s.UnimplementedInstruction(opcode)
	case 0x10:
		s.UnimplementedInstruction(opcode)
	case 0x11:
		s.D = data[1]
		s.E = data[0]
		s.PC += 2
	case 0x12:
		s.UnimplementedInstruction(opcode)
	case 0x13:
		s.E++
		if s.E == 0 {
			s.D++
		}
	case 0x14:
		s.UnimplementedInstruction(opcode)
	case 0x15:
		s.UnimplementedInstruction(opcode)
	case 0x16:
		s.UnimplementedInstruction(opcode)
	case 0x17:
		s.UnimplementedInstruction(opcode)
	case 0x18:
		s.UnimplementedInstruction(opcode)
	case 0x19:
		s.UnimplementedInstruction(opcode)
	case 0x1a:
		offset := (uint16(s.D) << 8) | uint16(s.E)
		s.A = s.Memory[offset]
	case 0x1b:
		s.UnimplementedInstruction(opcode)
	case 0x1c:
		s.UnimplementedInstruction(opcode)
	case 0x1d:
		s.UnimplementedInstruction(opcode)
	case 0x1e:
		s.UnimplementedInstruction(opcode)
	case 0x1f:
		s.UnimplementedInstruction(opcode)
	case 0x20:
		s.UnimplementedInstruction(opcode)
	case 0x21:
		s.H = data[1]
		s.L = data[0]
		s.PC += 2
	case 0x22:
		s.UnimplementedInstruction(opcode)
	case 0x23:
		s.L++
		if s.L == 0 {
			s.H++
		}
	case 0x24:
		s.UnimplementedInstruction(opcode)
	case 0x25:
		s.UnimplementedInstruction(opcode)
	case 0x26:
		s.UnimplementedInstruction(opcode)
	case 0x27:
		s.UnimplementedInstruction(opcode)
	case 0x28:
		s.UnimplementedInstruction(opcode)
	case 0x29:
		s.UnimplementedInstruction(opcode)
	case 0x2a:
		s.UnimplementedInstruction(opcode)
	case 0x2b:
		s.UnimplementedInstruction(opcode)
	case 0x2c:
		s.UnimplementedInstruction(opcode)
	case 0x2d:
		s.UnimplementedInstruction(opcode)
	case 0x2e:
		s.UnimplementedInstruction(opcode)
	case 0x2f:
		s.UnimplementedInstruction(opcode)
	case 0x30:
		s.UnimplementedInstruction(opcode)
	case 0x31:
		s.SP = (uint16(data[1])<<8 | uint16(data[0]))
		s.PC += 2
	case 0x32:
		s.UnimplementedInstruction(opcode)
	case 0x33:
		s.UnimplementedInstruction(opcode)
	case 0x34:
		s.UnimplementedInstruction(opcode)
	case 0x35:
		s.UnimplementedInstruction(opcode)
	case 0x36:
		offset := (uint16(s.H) << 8) | uint16(s.L)
		s.Memory[offset] = data[0]
		s.PC++
	case 0x37:
		s.UnimplementedInstruction(opcode)
	case 0x38:
		s.UnimplementedInstruction(opcode)
	case 0x39:
		s.UnimplementedInstruction(opcode)
	case 0x3a:
		s.UnimplementedInstruction(opcode)
	case 0x3b:
		s.UnimplementedInstruction(opcode)
	case 0x3c:
		s.UnimplementedInstruction(opcode)
	case 0x3d:
		s.UnimplementedInstruction(opcode)
	case 0x3e:
		s.UnimplementedInstruction(opcode)
	case 0x3f:
		s.UnimplementedInstruction(opcode)
	case 0x40:
		s.UnimplementedInstruction(opcode)
	case 0x41:
		s.UnimplementedInstruction(opcode)
	case 0x42:
		s.UnimplementedInstruction(opcode)
	case 0x43:
		s.UnimplementedInstruction(opcode)
	case 0x44:
		s.UnimplementedInstruction(opcode)
	case 0x45:
		s.UnimplementedInstruction(opcode)
	case 0x46:
		s.UnimplementedInstruction(opcode)
	case 0x47:
		s.UnimplementedInstruction(opcode)
	case 0x48:
		s.UnimplementedInstruction(opcode)
	case 0x49:
		s.UnimplementedInstruction(opcode)
	case 0x4a:
		s.UnimplementedInstruction(opcode)
	case 0x4b:
		s.UnimplementedInstruction(opcode)
	case 0x4c:
		s.UnimplementedInstruction(opcode)
	case 0x4d:
		s.UnimplementedInstruction(opcode)
	case 0x4e:
		s.UnimplementedInstruction(opcode)
	case 0x4f:
		s.UnimplementedInstruction(opcode)
	case 0x50:
		s.UnimplementedInstruction(opcode)
	case 0x51:
		s.UnimplementedInstruction(opcode)
	case 0x52:
		s.UnimplementedInstruction(opcode)
	case 0x53:
		s.UnimplementedInstruction(opcode)
	case 0x54:
		s.UnimplementedInstruction(opcode)
	case 0x55:
		s.UnimplementedInstruction(opcode)
	case 0x56:
		s.UnimplementedInstruction(opcode)
	case 0x57:
		s.UnimplementedInstruction(opcode)
	case 0x58:
		s.UnimplementedInstruction(opcode)
	case 0x59:
		s.UnimplementedInstruction(opcode)
	case 0x5a:
		s.UnimplementedInstruction(opcode)
	case 0x5b:
		s.UnimplementedInstruction(opcode)
	case 0x5c:
		s.UnimplementedInstruction(opcode)
	case 0x5d:
		s.UnimplementedInstruction(opcode)
	case 0x5e:
		s.UnimplementedInstruction(opcode)
	case 0x5f:
		s.UnimplementedInstruction(opcode)
	case 0x60:
		s.UnimplementedInstruction(opcode)
	case 0x61:
		s.UnimplementedInstruction(opcode)
	case 0x62:
		s.UnimplementedInstruction(opcode)
	case 0x63:
		s.UnimplementedInstruction(opcode)
	case 0x64:
		s.UnimplementedInstruction(opcode)
	case 0x65:
		s.UnimplementedInstruction(opcode)
	case 0x66:
		s.UnimplementedInstruction(opcode)
	case 0x67:
		s.UnimplementedInstruction(opcode)
	case 0x68:
		s.UnimplementedInstruction(opcode)
	case 0x69:
		s.UnimplementedInstruction(opcode)
	case 0x6a:
		s.UnimplementedInstruction(opcode)
	case 0x6b:
		s.UnimplementedInstruction(opcode)
	case 0x6c:
		s.UnimplementedInstruction(opcode)
	case 0x6d:
		s.UnimplementedInstruction(opcode)
	case 0x6e:
		s.UnimplementedInstruction(opcode)
	case 0x6f:
		s.UnimplementedInstruction(opcode)
	case 0x70:
		s.UnimplementedInstruction(opcode)
	case 0x71:
		s.UnimplementedInstruction(opcode)
	case 0x72:
		s.UnimplementedInstruction(opcode)
	case 0x73:
		s.UnimplementedInstruction(opcode)
	case 0x74:
		s.UnimplementedInstruction(opcode)
	case 0x75:
		s.UnimplementedInstruction(opcode)
	case 0x76:
		s.UnimplementedInstruction(opcode)
	case 0x77:
		offset := (uint16(s.H) << 8) | uint16(s.L)
		s.Memory[offset] = s.A

	case 0x78:
		s.UnimplementedInstruction(opcode)
	case 0x79:
		s.UnimplementedInstruction(opcode)
	case 0x7a:
		s.UnimplementedInstruction(opcode)
	case 0x7b:
		s.UnimplementedInstruction(opcode)
	case 0x7c:
		s.A = s.H
	case 0x7d:
		s.UnimplementedInstruction(opcode)
	case 0x7e:
		s.UnimplementedInstruction(opcode)
	case 0x7f:
		s.UnimplementedInstruction(opcode)
	case 0x80:
		s.UnimplementedInstruction(opcode)
	case 0x81:
		s.UnimplementedInstruction(opcode)
	case 0x82:
		s.UnimplementedInstruction(opcode)
	case 0x83:
		s.UnimplementedInstruction(opcode)
	case 0x84:
		s.UnimplementedInstruction(opcode)
	case 0x85:
		s.UnimplementedInstruction(opcode)
	case 0x86:
		s.UnimplementedInstruction(opcode)
	case 0x87:
		s.UnimplementedInstruction(opcode)
	case 0x88:
		s.UnimplementedInstruction(opcode)
	case 0x89:
		s.UnimplementedInstruction(opcode)
	case 0x8a:
		s.UnimplementedInstruction(opcode)
	case 0x8b:
		s.UnimplementedInstruction(opcode)
	case 0x8c:
		s.UnimplementedInstruction(opcode)
	case 0x8d:
		s.UnimplementedInstruction(opcode)
	case 0x8e:
		s.UnimplementedInstruction(opcode)
	case 0x8f:
		s.UnimplementedInstruction(opcode)
	case 0x90:
		s.UnimplementedInstruction(opcode)
	case 0x91:
		s.UnimplementedInstruction(opcode)
	case 0x92:
		s.UnimplementedInstruction(opcode)
	case 0x93:
		s.UnimplementedInstruction(opcode)
	case 0x94:
		s.UnimplementedInstruction(opcode)
	case 0x95:
		s.UnimplementedInstruction(opcode)
	case 0x96:
		s.UnimplementedInstruction(opcode)
	case 0x97:
		s.UnimplementedInstruction(opcode)
	case 0x98:
		s.UnimplementedInstruction(opcode)
	case 0x99:
		s.UnimplementedInstruction(opcode)
	case 0x9a:
		s.UnimplementedInstruction(opcode)
	case 0x9b:
		s.UnimplementedInstruction(opcode)
	case 0x9c:
		s.UnimplementedInstruction(opcode)
	case 0x9d:
		s.UnimplementedInstruction(opcode)
	case 0x9e:
		s.UnimplementedInstruction(opcode)
	case 0x9f:
		s.UnimplementedInstruction(opcode)
	case 0xa0:
		s.UnimplementedInstruction(opcode)
	case 0xa1:
		s.UnimplementedInstruction(opcode)
	case 0xa2:
		s.UnimplementedInstruction(opcode)
	case 0xa3:
		s.UnimplementedInstruction(opcode)
	case 0xa4:
		s.UnimplementedInstruction(opcode)
	case 0xa5:
		s.UnimplementedInstruction(opcode)
	case 0xa6:
		s.UnimplementedInstruction(opcode)
	case 0xa7:
		s.UnimplementedInstruction(opcode)
	case 0xa8:
		s.UnimplementedInstruction(opcode)
	case 0xa9:
		s.UnimplementedInstruction(opcode)
	case 0xaa:
		s.UnimplementedInstruction(opcode)
	case 0xab:
		s.UnimplementedInstruction(opcode)
	case 0xac:
		s.UnimplementedInstruction(opcode)
	case 0xad:
		s.UnimplementedInstruction(opcode)
	case 0xae:
		s.UnimplementedInstruction(opcode)
	case 0xaf:
		s.UnimplementedInstruction(opcode)
	case 0xb0:
		s.UnimplementedInstruction(opcode)
	case 0xb1:
		s.UnimplementedInstruction(opcode)
	case 0xb2:
		s.UnimplementedInstruction(opcode)
	case 0xb3:
		s.UnimplementedInstruction(opcode)
	case 0xb4:
		s.UnimplementedInstruction(opcode)
	case 0xb5:
		s.UnimplementedInstruction(opcode)
	case 0xb6:
		s.UnimplementedInstruction(opcode)
	case 0xb7:
		s.UnimplementedInstruction(opcode)
	case 0xb8:
		s.UnimplementedInstruction(opcode)
	case 0xb9:
		s.UnimplementedInstruction(opcode)
	case 0xba:
		s.UnimplementedInstruction(opcode)
	case 0xbb:
		s.UnimplementedInstruction(opcode)
	case 0xbc:
		s.UnimplementedInstruction(opcode)
	case 0xbd:
		s.UnimplementedInstruction(opcode)
	case 0xbe:
		s.UnimplementedInstruction(opcode)
	case 0xbf:
		s.UnimplementedInstruction(opcode)
	case 0xc0:
		s.UnimplementedInstruction(opcode)
	case 0xc1:
		s.UnimplementedInstruction(opcode)
	case 0xc2:
		if !s.ConditionCodes.Z {
			s.PC = (uint16(data[1]) << 8) | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xc3:
		s.PC = (uint16(data[1]) << 8) | uint16(data[0])
	case 0xc4:
		s.UnimplementedInstruction(opcode)
	case 0xc5:
		s.UnimplementedInstruction(opcode)
	case 0xc6:
		s.UnimplementedInstruction(opcode)
	case 0xc7:
		s.UnimplementedInstruction(opcode)
	case 0xc8:
		s.UnimplementedInstruction(opcode)
	case 0xc9:
		s.PC = uint16(s.Memory[s.SP]) | (uint16(s.Memory[s.SP+1]) << 8)
		s.SP += 2
	case 0xca:
		s.UnimplementedInstruction(opcode)
	case 0xcb:
		s.UnimplementedInstruction(opcode)
	case 0xcc:
		s.UnimplementedInstruction(opcode)
	case 0xcd:
		ret := s.PC + 2
		s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
		s.Memory[s.SP-2] = uint8((ret & 0xFF))
		s.SP = s.SP - 2
		s.PC = (uint16(data[1])<<8 | uint16(data[0]))

	case 0xce:
		s.UnimplementedInstruction(opcode)
	case 0xcf:
		s.UnimplementedInstruction(opcode)
	case 0xd0:
		s.UnimplementedInstruction(opcode)
	case 0xd1:
		s.UnimplementedInstruction(opcode)
	case 0xd2:
		s.UnimplementedInstruction(opcode)
	case 0xd3:
		s.UnimplementedInstruction(opcode)
	case 0xd4:
		s.UnimplementedInstruction(opcode)
	case 0xd5:
		s.UnimplementedInstruction(opcode)
	case 0xd6:
		s.UnimplementedInstruction(opcode)
	case 0xd7:
		s.UnimplementedInstruction(opcode)
	case 0xd8:
		s.UnimplementedInstruction(opcode)
	case 0xd9:
		s.UnimplementedInstruction(opcode)
	case 0xda:
		s.UnimplementedInstruction(opcode)
	case 0xdb:
		s.UnimplementedInstruction(opcode)
	case 0xdc:
		s.UnimplementedInstruction(opcode)
	case 0xdd:
		s.UnimplementedInstruction(opcode)
	case 0xde:
		s.UnimplementedInstruction(opcode)
	case 0xdf:
		s.UnimplementedInstruction(opcode)
	case 0xe0:
		s.UnimplementedInstruction(opcode)
	case 0xe1:
		s.UnimplementedInstruction(opcode)
	case 0xe2:
		s.UnimplementedInstruction(opcode)
	case 0xe3:
		s.UnimplementedInstruction(opcode)
	case 0xe4:
		s.UnimplementedInstruction(opcode)
	case 0xe5:
		s.UnimplementedInstruction(opcode)
	case 0xe6:
		s.UnimplementedInstruction(opcode)
	case 0xe7:
		s.UnimplementedInstruction(opcode)
	case 0xe8:
		s.UnimplementedInstruction(opcode)
	case 0xe9:
		s.UnimplementedInstruction(opcode)
	case 0xea:
		s.UnimplementedInstruction(opcode)
	case 0xeb:
		s.UnimplementedInstruction(opcode)
	case 0xec:
		s.UnimplementedInstruction(opcode)
	case 0xed:
		s.UnimplementedInstruction(opcode)
	case 0xee:
		s.UnimplementedInstruction(opcode)
	case 0xef:
		s.UnimplementedInstruction(opcode)
	case 0xf0:
		s.UnimplementedInstruction(opcode)
	case 0xf1:
		s.UnimplementedInstruction(opcode)
	case 0xf2:
		s.UnimplementedInstruction(opcode)
	case 0xf3:
		s.UnimplementedInstruction(opcode)
	case 0xf4:
		s.UnimplementedInstruction(opcode)
	case 0xf5:
		s.UnimplementedInstruction(opcode)
	case 0xf6:
		s.UnimplementedInstruction(opcode)
	case 0xf7:
		s.UnimplementedInstruction(opcode)
	case 0xf8:
		s.UnimplementedInstruction(opcode)
	case 0xf9:
		s.UnimplementedInstruction(opcode)
	case 0xfa:
		s.UnimplementedInstruction(opcode)
	case 0xfb:
		s.UnimplementedInstruction(opcode)
	case 0xfc:
		s.UnimplementedInstruction(opcode)
	case 0xfd:
		s.UnimplementedInstruction(opcode)
	case 0xfe:
		x := s.A - data[0]
		s.ConditionCodes.Z = (x == 0)
		s.ConditionCodes.S = (0x80 == (x & 0x80))
		s.ConditionCodes.P = parity(x, 8)
		s.ConditionCodes.CY = (s.A < data[0])
		s.PC++
	case 0xff:
		s.UnimplementedInstruction(opcode)
	default:
		s.UnimplementedInstruction(opcode)
	}

}
