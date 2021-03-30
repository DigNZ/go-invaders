package main

type ConditionCodes struct {
	Z, S, P, CY, AC, PAD uint8
}
type State8080 struct {
	A, B, C, D, E, H, L uint8
	SP, PC              uint16
	Memory              []uint8
	ConditionCodes      ConditionCodes
	IntEnable           uint8
}

func (s *State8080) UnimplementedInstruction() {
	panic("Error: Unimplemented Instruction")
}

func (s *State8080) Emulate8080Op() {
	opcode := s.Memory[s.PC]
	switch opcode {
	case 0x00:
		s.UnimplementedInstruction()
	case 0x01:
		s.UnimplementedInstruction()
	case 0x02:
		s.UnimplementedInstruction()
	case 0x03:
		s.UnimplementedInstruction()
	case 0x04:
		s.UnimplementedInstruction()
	case 0x05:
		s.UnimplementedInstruction()
	case 0x06:
		s.UnimplementedInstruction()
	case 0x07:
		s.UnimplementedInstruction()
	case 0x08:
		s.UnimplementedInstruction()
	case 0x09:
		s.UnimplementedInstruction()
	case 0x0a:
		s.UnimplementedInstruction()
	case 0x0b:
		s.UnimplementedInstruction()
	case 0x0c:
		s.UnimplementedInstruction()
	case 0x0d:
		s.UnimplementedInstruction()
	case 0x0e:
		s.UnimplementedInstruction()
	case 0x0f:
		s.UnimplementedInstruction()
	case 0x10:
		s.UnimplementedInstruction()
	case 0x11:
		s.UnimplementedInstruction()
	case 0x12:
		s.UnimplementedInstruction()
	case 0x13:
		s.UnimplementedInstruction()
	case 0x14:
		s.UnimplementedInstruction()
	case 0x15:
		s.UnimplementedInstruction()
	case 0x16:
		s.UnimplementedInstruction()
	case 0x17:
		s.UnimplementedInstruction()
	case 0x18:
		s.UnimplementedInstruction()
	case 0x19:
		s.UnimplementedInstruction()
	case 0x1a:
		s.UnimplementedInstruction()
	case 0x1b:
		s.UnimplementedInstruction()
	case 0x1c:
		s.UnimplementedInstruction()
	case 0x1d:
		s.UnimplementedInstruction()
	case 0x1e:
		s.UnimplementedInstruction()
	case 0x1f:
		s.UnimplementedInstruction()
	case 0x20:
		s.UnimplementedInstruction()
	case 0x21:
		s.UnimplementedInstruction()
	case 0x22:
		s.UnimplementedInstruction()
	case 0x23:
		s.UnimplementedInstruction()
	case 0x24:
		s.UnimplementedInstruction()
	case 0x25:
		s.UnimplementedInstruction()
	case 0x26:
		s.UnimplementedInstruction()
	case 0x27:
		s.UnimplementedInstruction()
	case 0x28:
		s.UnimplementedInstruction()
	case 0x29:
		s.UnimplementedInstruction()
	case 0x2a:
		s.UnimplementedInstruction()
	case 0x2b:
		s.UnimplementedInstruction()
	case 0x2c:
		s.UnimplementedInstruction()
	case 0x2d:
		s.UnimplementedInstruction()
	case 0x2e:
		s.UnimplementedInstruction()
	case 0x2f:
		s.UnimplementedInstruction()
	case 0x30:
		s.UnimplementedInstruction()
	case 0x31:
		s.UnimplementedInstruction()
	case 0x32:
		s.UnimplementedInstruction()
	case 0x33:
		s.UnimplementedInstruction()
	case 0x34:
		s.UnimplementedInstruction()
	case 0x35:
		s.UnimplementedInstruction()
	case 0x36:
		s.UnimplementedInstruction()
	case 0x37:
		s.UnimplementedInstruction()
	case 0x38:
		s.UnimplementedInstruction()
	case 0x39:
		s.UnimplementedInstruction()
	case 0x3a:
		s.UnimplementedInstruction()
	case 0x3b:
		s.UnimplementedInstruction()
	case 0x3c:
		s.UnimplementedInstruction()
	case 0x3d:
		s.UnimplementedInstruction()
	case 0x3e:
		s.UnimplementedInstruction()
	case 0x3f:
		s.UnimplementedInstruction()
	case 0x40:
		s.UnimplementedInstruction()
	case 0x41:
		s.UnimplementedInstruction()
	case 0x42:
		s.UnimplementedInstruction()
	case 0x43:
		s.UnimplementedInstruction()
	case 0x44:
		s.UnimplementedInstruction()
	case 0x45:
		s.UnimplementedInstruction()
	case 0x46:
		s.UnimplementedInstruction()
	case 0x47:
		s.UnimplementedInstruction()
	case 0x48:
		s.UnimplementedInstruction()
	case 0x49:
		s.UnimplementedInstruction()
	case 0x4a:
		s.UnimplementedInstruction()
	case 0x4b:
		s.UnimplementedInstruction()
	case 0x4c:
		s.UnimplementedInstruction()
	case 0x4d:
		s.UnimplementedInstruction()
	case 0x4e:
		s.UnimplementedInstruction()
	case 0x4f:
		s.UnimplementedInstruction()
	case 0x50:
		s.UnimplementedInstruction()
	case 0x51:
		s.UnimplementedInstruction()
	case 0x52:
		s.UnimplementedInstruction()
	case 0x53:
		s.UnimplementedInstruction()
	case 0x54:
		s.UnimplementedInstruction()
	case 0x55:
		s.UnimplementedInstruction()
	case 0x56:
		s.UnimplementedInstruction()
	case 0x57:
		s.UnimplementedInstruction()
	case 0x58:
		s.UnimplementedInstruction()
	case 0x59:
		s.UnimplementedInstruction()
	case 0x5a:
		s.UnimplementedInstruction()
	case 0x5b:
		s.UnimplementedInstruction()
	case 0x5c:
		s.UnimplementedInstruction()
	case 0x5d:
		s.UnimplementedInstruction()
	case 0x5e:
		s.UnimplementedInstruction()
	case 0x5f:
		s.UnimplementedInstruction()
	case 0x60:
		s.UnimplementedInstruction()
	case 0x61:
		s.UnimplementedInstruction()
	case 0x62:
		s.UnimplementedInstruction()
	case 0x63:
		s.UnimplementedInstruction()
	case 0x64:
		s.UnimplementedInstruction()
	case 0x65:
		s.UnimplementedInstruction()
	case 0x66:
		s.UnimplementedInstruction()
	case 0x67:
		s.UnimplementedInstruction()
	case 0x68:
		s.UnimplementedInstruction()
	case 0x69:
		s.UnimplementedInstruction()
	case 0x6a:
		s.UnimplementedInstruction()
	case 0x6b:
		s.UnimplementedInstruction()
	case 0x6c:
		s.UnimplementedInstruction()
	case 0x6d:
		s.UnimplementedInstruction()
	case 0x6e:
		s.UnimplementedInstruction()
	case 0x6f:
		s.UnimplementedInstruction()
	case 0x70:
		s.UnimplementedInstruction()
	case 0x71:
		s.UnimplementedInstruction()
	case 0x72:
		s.UnimplementedInstruction()
	case 0x73:
		s.UnimplementedInstruction()
	case 0x74:
		s.UnimplementedInstruction()
	case 0x75:
		s.UnimplementedInstruction()
	case 0x76:
		s.UnimplementedInstruction()
	case 0x77:
		s.UnimplementedInstruction()
	case 0x78:
		s.UnimplementedInstruction()
	case 0x79:
		s.UnimplementedInstruction()
	case 0x7a:
		s.UnimplementedInstruction()
	case 0x7b:
		s.UnimplementedInstruction()
	case 0x7c:
		s.UnimplementedInstruction()
	case 0x7d:
		s.UnimplementedInstruction()
	case 0x7e:
		s.UnimplementedInstruction()
	case 0x7f:
		s.UnimplementedInstruction()
	case 0x80:
		s.UnimplementedInstruction()
	case 0x81:
		s.UnimplementedInstruction()
	case 0x82:
		s.UnimplementedInstruction()
	case 0x83:
		s.UnimplementedInstruction()
	case 0x84:
		s.UnimplementedInstruction()
	case 0x85:
		s.UnimplementedInstruction()
	case 0x86:
		s.UnimplementedInstruction()
	case 0x87:
		s.UnimplementedInstruction()
	case 0x88:
		s.UnimplementedInstruction()
	case 0x89:
		s.UnimplementedInstruction()
	case 0x8a:
		s.UnimplementedInstruction()
	case 0x8b:
		s.UnimplementedInstruction()
	case 0x8c:
		s.UnimplementedInstruction()
	case 0x8d:
		s.UnimplementedInstruction()
	case 0x8e:
		s.UnimplementedInstruction()
	case 0x8f:
		s.UnimplementedInstruction()
	case 0x90:
		s.UnimplementedInstruction()
	case 0x91:
		s.UnimplementedInstruction()
	case 0x92:
		s.UnimplementedInstruction()
	case 0x93:
		s.UnimplementedInstruction()
	case 0x94:
		s.UnimplementedInstruction()
	case 0x95:
		s.UnimplementedInstruction()
	case 0x96:
		s.UnimplementedInstruction()
	case 0x97:
		s.UnimplementedInstruction()
	case 0x98:
		s.UnimplementedInstruction()
	case 0x99:
		s.UnimplementedInstruction()
	case 0x9a:
		s.UnimplementedInstruction()
	case 0x9b:
		s.UnimplementedInstruction()
	case 0x9c:
		s.UnimplementedInstruction()
	case 0x9d:
		s.UnimplementedInstruction()
	case 0x9e:
		s.UnimplementedInstruction()
	case 0x9f:
		s.UnimplementedInstruction()
	case 0xa0:
		s.UnimplementedInstruction()
	case 0xa1:
		s.UnimplementedInstruction()
	case 0xa2:
		s.UnimplementedInstruction()
	case 0xa3:
		s.UnimplementedInstruction()
	case 0xa4:
		s.UnimplementedInstruction()
	case 0xa5:
		s.UnimplementedInstruction()
	case 0xa6:
		s.UnimplementedInstruction()
	case 0xa7:
		s.UnimplementedInstruction()
	case 0xa8:
		s.UnimplementedInstruction()
	case 0xa9:
		s.UnimplementedInstruction()
	case 0xaa:
		s.UnimplementedInstruction()
	case 0xab:
		s.UnimplementedInstruction()
	case 0xac:
		s.UnimplementedInstruction()
	case 0xad:
		s.UnimplementedInstruction()
	case 0xae:
		s.UnimplementedInstruction()
	case 0xaf:
		s.UnimplementedInstruction()
	case 0xb0:
		s.UnimplementedInstruction()
	case 0xb1:
		s.UnimplementedInstruction()
	case 0xb2:
		s.UnimplementedInstruction()
	case 0xb3:
		s.UnimplementedInstruction()
	case 0xb4:
		s.UnimplementedInstruction()
	case 0xb5:
		s.UnimplementedInstruction()
	case 0xb6:
		s.UnimplementedInstruction()
	case 0xb7:
		s.UnimplementedInstruction()
	case 0xb8:
		s.UnimplementedInstruction()
	case 0xb9:
		s.UnimplementedInstruction()
	case 0xba:
		s.UnimplementedInstruction()
	case 0xbb:
		s.UnimplementedInstruction()
	case 0xbc:
		s.UnimplementedInstruction()
	case 0xbd:
		s.UnimplementedInstruction()
	case 0xbe:
		s.UnimplementedInstruction()
	case 0xbf:
		s.UnimplementedInstruction()
	case 0xc0:
		s.UnimplementedInstruction()
	case 0xc1:
		s.UnimplementedInstruction()
	case 0xc2:
		s.UnimplementedInstruction()
	case 0xc3:
		s.UnimplementedInstruction()
	case 0xc4:
		s.UnimplementedInstruction()
	case 0xc5:
		s.UnimplementedInstruction()
	case 0xc6:
		s.UnimplementedInstruction()
	case 0xc7:
		s.UnimplementedInstruction()
	case 0xc8:
		s.UnimplementedInstruction()
	case 0xc9:
		s.UnimplementedInstruction()
	case 0xca:
		s.UnimplementedInstruction()
	case 0xcb:
		s.UnimplementedInstruction()
	case 0xcc:
		s.UnimplementedInstruction()
	case 0xcd:
		s.UnimplementedInstruction()
	case 0xce:
		s.UnimplementedInstruction()
	case 0xcf:
		s.UnimplementedInstruction()
	case 0xd0:
		s.UnimplementedInstruction()
	case 0xd1:
		s.UnimplementedInstruction()
	case 0xd2:
		s.UnimplementedInstruction()
	case 0xd3:
		s.UnimplementedInstruction()
	case 0xd4:
		s.UnimplementedInstruction()
	case 0xd5:
		s.UnimplementedInstruction()
	case 0xd6:
		s.UnimplementedInstruction()
	case 0xd7:
		s.UnimplementedInstruction()
	case 0xd8:
		s.UnimplementedInstruction()
	case 0xd9:
		s.UnimplementedInstruction()
	case 0xda:
		s.UnimplementedInstruction()
	case 0xdb:
		s.UnimplementedInstruction()
	case 0xdc:
		s.UnimplementedInstruction()
	case 0xdd:
		s.UnimplementedInstruction()
	case 0xde:
		s.UnimplementedInstruction()
	case 0xdf:
		s.UnimplementedInstruction()
	case 0xe0:
		s.UnimplementedInstruction()
	case 0xe1:
		s.UnimplementedInstruction()
	case 0xe2:
		s.UnimplementedInstruction()
	case 0xe3:
		s.UnimplementedInstruction()
	case 0xe4:
		s.UnimplementedInstruction()
	case 0xe5:
		s.UnimplementedInstruction()
	case 0xe6:
		s.UnimplementedInstruction()
	case 0xe7:
		s.UnimplementedInstruction()
	case 0xe8:
		s.UnimplementedInstruction()
	case 0xe9:
		s.UnimplementedInstruction()
	case 0xea:
		s.UnimplementedInstruction()
	case 0xeb:
		s.UnimplementedInstruction()
	case 0xec:
		s.UnimplementedInstruction()
	case 0xed:
		s.UnimplementedInstruction()
	case 0xee:
		s.UnimplementedInstruction()
	case 0xef:
		s.UnimplementedInstruction()
	case 0xf0:
		s.UnimplementedInstruction()
	case 0xf1:
		s.UnimplementedInstruction()
	case 0xf2:
		s.UnimplementedInstruction()
	case 0xf3:
		s.UnimplementedInstruction()
	case 0xf4:
		s.UnimplementedInstruction()
	case 0xf5:
		s.UnimplementedInstruction()
	case 0xf6:
		s.UnimplementedInstruction()
	case 0xf7:
		s.UnimplementedInstruction()
	case 0xf8:
		s.UnimplementedInstruction()
	case 0xf9:
		s.UnimplementedInstruction()
	case 0xfa:
		s.UnimplementedInstruction()
	case 0xfb:
		s.UnimplementedInstruction()
	case 0xfc:
		s.UnimplementedInstruction()
	case 0xfd:
		s.UnimplementedInstruction()
	case 0xfe:
		s.UnimplementedInstruction()
	case 0xff:
		s.UnimplementedInstruction()
	default:
		s.UnimplementedInstruction()
	}
	s.PC += 1
}
