package core

import (
	"fmt"
	"log"
	"os"
	"time"

	io "github.com/DigNZ/goinvaders/io"
)

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
func (s *State8080) InitWithDataAt(data []byte, addr int) {
	for idx, d := range data {
		s.Memory[idx+addr] = d
	}
	s.PC = uint16(addr)

	//cpudiag hacks
	s.Memory[368] = 0x7

	//Skip DAA test
	s.Memory[0x59c] = 0xc3 //JMP
	s.Memory[0x59d] = 0xc2
	s.Memory[0x59e] = 0x05
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
func (s *State8080) logicFlagsA() {

	s.ConditionCodes.CY = false
	s.ConditionCodes.AC = false
	s.ConditionCodes.Z = s.A == 0
	s.ConditionCodes.S = (0x80 == (s.A & 0x80))
	s.ConditionCodes.P = parity(s.A, 8)
}

func (s *State8080) flagsZSP(value uint8) {
	s.ConditionCodes.Z = value == 0
	s.ConditionCodes.S = (0x80 == (value & 0x80))
	s.ConditionCodes.P = parity(value, 8)
}

func (s *State8080) writeMem(address uint16, value uint8) {
	if address < 0x2000 {
		//	log.Fatalf("Attempted to write to ROM %x", address)
		return
	}
	if address >= 0x4000 {
		log.Fatalf("Attempted to write outside RAM")
	}
	s.Memory[address] = value
}
func (s *State8080) push(high, low uint8) {
	s.writeMem(s.SP-1, high)
	s.writeMem(s.SP-2, low)

	s.SP = s.SP - 2
}

func (s *State8080) readFromHL() uint8 {
	offset := uint16(s.H)<<8 | uint16(s.L)
	return s.Memory[offset]
}

func (s *State8080) generateInterrupt(interrupt_num uint16) {
	//perform "PUSH PC"
	s.push(uint8((s.PC&0xFF00)>>8), uint8(s.PC&0xff))

	//Set the PC to the low memory vector.
	//This is identical to an "RST interrupt_num" instruction.
	s.PC = 8 * interrupt_num
}
func (s *State8080) Run() {
	m := io.Machine{}
	var lastInterrupt time.Time
	for {
		opcode := s.Memory[s.PC]
		if opcode == 0xdb {
			fmt.Print("0XDB")
			port := s.Memory[s.PC+1]
			s.A = m.MachineIN(port)
			s.PC += 2
		} else if opcode == 0xd3 {
			port := s.Memory[s.PC+1]
			m.MachineOUT(port)
			s.PC += 2
		} else {
			s.Emulate8080Op(true)
		}
		if (time.Now().Nanosecond())-lastInterrupt.Nanosecond() > 16666666 { //1/60 second has elapsed
			//only do an interrupt if they are enabled
			if s.IntEnable == 1 {
				s.generateInterrupt(2) //interrupt 2

				//Save the time we did this
				lastInterrupt = time.Now()
			}
		}

	}
}
func (s *State8080) Emulate8080Op(dasm bool) {
	opcode := s.Memory[s.PC]
	var data []uint8
	if s.PC+3 > s.PC+1 {
		data = s.Memory[s.PC+1 : s.PC+3]
	}
	if dasm {
		Disassemble8080Op(s.Memory[:], s.PC)
	}
	//Remember that we've advanced the PC BEFORE executing the opcode.
	s.PC += 1
	switch opcode {
	case 0x00:
		// NOP
	case 0x01:
		s.B = data[1]
		s.C = data[0]
		s.PC += 2
	case 0x02:
		s.UnimplementedInstruction(opcode)
	case 0x03:
		s.UnimplementedInstruction(opcode)
	case 0x04:
		s.UnimplementedInstruction(opcode)
	case 0x05:
		s.B -= 1
		s.flagsZSP(s.B)
	case 0x06:
		s.B = s.Memory[s.PC]
		s.PC++
	case 0x07:
		s.UnimplementedInstruction(opcode)
	case 0x08:
		s.UnimplementedInstruction(opcode)
	case 0x09:
		hl := uint32((uint16(s.H) << 8) | uint16(s.L))
		bc := uint32((uint16(s.B) << 8) | uint16(s.C))
		res := hl + bc
		s.H = uint8((res & 0xff00) >> 8)
		s.L = uint8(res & 0xff)
		s.ConditionCodes.CY = ((res & 0xffff0000) != 0)
	case 0x0a:
		s.UnimplementedInstruction(opcode)
	case 0x0b:
		s.UnimplementedInstruction(opcode)
	case 0x0c:
		s.UnimplementedInstruction(opcode)
	case 0x0d:
		res := s.C - 1
		s.ConditionCodes.Z = (res == 0)
		s.ConditionCodes.S = (0x80 == (res & 0x80))
		s.ConditionCodes.P = parity(res, 8)
		s.C = res
	case 0x0e:
		s.C = data[0]
		s.PC++
	case 0x0f:
		x := s.A
		s.A = ((x & 1) << 7) | (x >> 1)
		s.ConditionCodes.CY = (1 == (x & 1))
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
		hl := uint32((uint16(s.H) << 8) | uint16(s.L))
		de := uint32((uint16(s.D) << 8) | uint16(s.E))
		res := hl + de
		s.H = uint8((res & 0xff00) >> 8)
		s.L = uint8(res & 0xff)
		s.ConditionCodes.CY = ((res & 0xffff0000) != 0)
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
		s.H = data[0]
		s.PC++
	case 0x27:
		s.UnimplementedInstruction(opcode)
	case 0x28:
		s.UnimplementedInstruction(opcode)
	case 0x29:
		hl := uint32((uint16(s.H) << 8) | uint16(s.L))
		res := hl + hl
		s.H = uint8((res & 0xff00) >> 8)
		s.L = uint8(res & 0xff)
		s.ConditionCodes.CY = ((res & 0xffff0000) != 0)
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
		offset := (uint16(data[1]) << 8) | uint16(data[0])
		s.Memory[offset] = s.A
		s.PC += 2
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
		offset := (uint16(data[1]) << 8) | uint16(data[0])
		s.A = s.Memory[offset]
		s.PC += 2
	case 0x3b:
		s.UnimplementedInstruction(opcode)
	case 0x3c:
		s.UnimplementedInstruction(opcode)
	case 0x3d:
		s.UnimplementedInstruction(opcode)
	case 0x3e:
		s.A = data[0]
		s.PC++
	case 0x3f:
		s.UnimplementedInstruction(opcode)
	case 0x40:
		s.UnimplementedInstruction(opcode)
	case 0x41:
		s.UnimplementedInstruction(opcode)
	case 0x42:
		s.UnimplementedInstruction(opcode)
	case 0x43:
		s.B = s.E
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
		s.C = s.L
	case 0x4e:
		s.UnimplementedInstruction(opcode)
	case 0x4f:
		s.C = s.A
	case 0x50:
		s.UnimplementedInstruction(opcode)
	case 0x51:
		s.UnimplementedInstruction(opcode)
	case 0x52:
		s.D = s.D
	case 0x53:
		s.D = s.E
	case 0x54:
		s.UnimplementedInstruction(opcode)
	case 0x55:
		s.UnimplementedInstruction(opcode)
	case 0x56:
		offset := (uint16(s.H) << 8) | uint16(s.L)
		s.D = s.Memory[offset]
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
		offset := (uint16(s.H) << 8) | uint16(s.L)
		s.E = s.Memory[offset]
	case 0x5f:
		s.E = s.A
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
		offset := (uint16(s.H) << 8) | uint16(s.L)
		s.H = s.Memory[offset]
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
		s.L = s.A
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
		s.A = s.D
	case 0x7b:
		s.A = s.E
	case 0x7c:
		s.A = s.H
	case 0x7d:
		s.A = s.L
	case 0x7e:
		offset := (uint16(s.H) << 8) | uint16(s.L)
		s.A = s.Memory[offset]
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
		s.A = s.A & s.B
		s.logicFlagsA()
	case 0xa1:
		s.A = s.A & s.C
		s.logicFlagsA()
	case 0xa2:
		s.A = s.A & s.D
		s.logicFlagsA()
	case 0xa3:
		s.A = s.A & s.E
		s.logicFlagsA()
	case 0xa4:
		s.A = s.A & s.H
		s.logicFlagsA()
	case 0xa5:
		s.A = s.A & s.L
		s.logicFlagsA()
	case 0xa6:
		s.A = s.A & s.readFromHL()
		s.logicFlagsA()
	case 0xa7:
		s.A = s.A & s.A
		s.logicFlagsA()
	case 0xa8:
		s.A = s.A ^ s.B
		s.logicFlagsA()
	case 0xa9:
		s.A = s.A ^ s.C
		s.logicFlagsA()
	case 0xaa:
		s.A = s.A ^ s.D
		s.logicFlagsA()
	case 0xab:
		s.A = s.A ^ s.E
		s.logicFlagsA()
	case 0xac:
		s.A = s.A ^ s.H
		s.logicFlagsA()
	case 0xad:
		s.A = s.A ^ s.L
		s.logicFlagsA()
	case 0xae:
		s.A = s.A ^ s.readFromHL()
		s.logicFlagsA()
	case 0xaf:
		s.A = s.A ^ s.A
		s.logicFlagsA()
	case 0xb0:
		s.A = s.A | s.B
		s.logicFlagsA()
	case 0xb1:
		s.A = s.A | s.C
		s.logicFlagsA()
	case 0xb2:
		s.A = s.A | s.D
		s.logicFlagsA()
	case 0xb3:
		s.A = s.A | s.E
		s.logicFlagsA()
	case 0xb4:
		s.A = s.A | s.H
		s.logicFlagsA()
	case 0xb5:
		s.A = s.A | s.L
		s.logicFlagsA()
	case 0xb6:
		s.A = s.A | s.readFromHL()
		s.logicFlagsA()
	case 0xb7:
		s.A = s.A | s.A
		s.logicFlagsA()
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
		s.C = s.Memory[s.SP]
		s.B = s.Memory[s.SP+1]
		s.SP += 2
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
		s.Memory[s.SP-1] = s.B
		s.Memory[s.SP-2] = s.C
		s.SP -= 2
	case 0xc6:
		x := uint16(s.A) + uint16(data[0])

		s.ConditionCodes.Z = ((x & 0xff) == 0)
		s.ConditionCodes.S = (0x80 == (x & 0x80))
		s.ConditionCodes.P = parity(uint8(x&0xff), 8)
		s.ConditionCodes.CY = (x > 0xff)
		s.A = uint8(x)
		s.PC++
	case 0xc7:
		s.UnimplementedInstruction(opcode)
	case 0xc8:
		s.UnimplementedInstruction(opcode)
	case 0xc9:
		s.PC = uint16(s.Memory[s.SP]) | (uint16(s.Memory[s.SP+1]) << 8)
		s.SP += 2
	case 0xca:
		if s.ConditionCodes.Z {
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xcb:
		s.UnimplementedInstruction(opcode)
	case 0xcc:
		s.UnimplementedInstruction(opcode)
	case 0xcd:
		if 5 == uint16(data[1])<<8|uint16(data[0]) {
			if s.C == 9 {
				offset := uint16(s.D)<<8 | uint16(s.E)
				var idx uint16 = 3
				str := s.Memory[offset+3]
				for {
					if str == '$' {
						fmt.Println("")
						os.Exit(0)
						break
					}
					fmt.Printf("%c", str)
					idx++
					str = s.Memory[offset+idx]
				}

			} else if s.C == 2 {
				fmt.Println("This got called")
			}
		} else if 0 == uint16(data[1])<<8|uint16(data[0]) {
			os.Exit(0)
		} else {
			ret := s.PC + 2
			s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
			s.Memory[s.SP-2] = uint8((ret & 0xFF))
			s.SP = s.SP - 2
			s.PC = (uint16(data[1])<<8 | uint16(data[0]))
		}

	case 0xce:
		x := uint16(s.A + data[0])
		if s.ConditionCodes.CY {
			x += 1
		}
		s.flagsZSP(uint8(x & 0xff))

		s.ConditionCodes.CY = (x > 0xff)
		s.A = uint8(x & 0xff)
		s.PC++
	case 0xcf:
		s.UnimplementedInstruction(opcode)
	case 0xd0:
		s.UnimplementedInstruction(opcode)
	case 0xd1:
		s.E = s.Memory[s.SP]
		s.D = s.Memory[s.SP+1]
		s.SP += 2
	case 0xd2:
		if !s.ConditionCodes.CY {
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xd3:
		//Special?
		s.PC++
	case 0xd4:
		s.UnimplementedInstruction(opcode)
	case 0xd5:
		s.Memory[s.SP-1] = s.D
		s.Memory[s.SP-2] = s.E
		s.SP -= 2
	case 0xd6:
		x := s.A - data[0]
		s.flagsZSP(x & 0xff)
		s.ConditionCodes.CY = (s.A < data[0])
		s.A = x
		s.PC++
	case 0xd7:
		s.UnimplementedInstruction(opcode)
	case 0xd8:
		s.UnimplementedInstruction(opcode)
	case 0xd9:
		s.UnimplementedInstruction(opcode)
	case 0xda:
		if s.ConditionCodes.CY {
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xdb:
		s.UnimplementedInstruction(opcode)
	case 0xdc:
		s.UnimplementedInstruction(opcode)
	case 0xdd:
		s.UnimplementedInstruction(opcode)
	case 0xde:
		var x uint16 = uint16(s.A) - uint16(data[0])
		if s.ConditionCodes.CY {
			x -= 1
		}
		s.flagsZSP(uint8(x & 0xff))
		s.ConditionCodes.CY = (x > 0xff)
		s.A = uint8(x & 0xFF)
		s.PC++

	case 0xdf:
		s.UnimplementedInstruction(opcode)
	case 0xe0:
		s.UnimplementedInstruction(opcode)
	case 0xe1:
		s.L = s.Memory[s.SP]
		s.H = s.Memory[s.SP+1]
		s.SP += 2
	case 0xe2:
		if !s.ConditionCodes.P {
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xe3:
		h := s.H
		l := s.L
		s.L = s.Memory[s.SP]
		s.H = s.Memory[s.SP+1]
		s.writeMem(s.SP, l)
		s.writeMem(s.SP+1, h)

	case 0xe4:
		s.UnimplementedInstruction(opcode)
	case 0xe5:
		s.Memory[s.SP-1] = s.H
		s.Memory[s.SP-2] = s.L
		s.SP -= 2
	case 0xe6:
		s.A = s.A & data[0]
		s.logicFlagsA()
		s.PC++
	case 0xe7:
		s.UnimplementedInstruction(opcode)
	case 0xe8:
		s.UnimplementedInstruction(opcode)
	case 0xe9:
		s.UnimplementedInstruction(opcode)
	case 0xea:
		if s.ConditionCodes.P {
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xeb:
		save1 := s.D
		save2 := s.E
		s.D = s.H
		s.E = s.L
		s.H = save1
		s.L = save2
	case 0xec:
		s.UnimplementedInstruction(opcode)
	case 0xed:
		s.UnimplementedInstruction(opcode)
	case 0xee:
		x := s.A ^ data[0]
		s.flagsZSP(x)
		s.ConditionCodes.CY = false
		s.A = x
		s.PC++
	case 0xef:
		s.UnimplementedInstruction(opcode)
	case 0xf0:
		s.UnimplementedInstruction(opcode)
	case 0xf1:

		s.A = s.Memory[s.SP+1]
		psw := s.Memory[s.SP]

		s.ConditionCodes.S = ((psw >> 7) & 1) == 1
		s.ConditionCodes.Z = ((psw >> 6) & 1) == 1
		s.ConditionCodes.AC = ((psw >> 4) & 1) == 1
		s.ConditionCodes.P = ((psw >> 2) & 1) == 1
		s.ConditionCodes.CY = ((psw >> 0) & 1) == 1
		/*s.ConditionCodes.Z = (0x01 == (psw & 0x01))
		s.ConditionCodes.S = (0x02 == (psw & 0x02))
		s.ConditionCodes.P = (0x04 == (psw & 0x04))
		s.ConditionCodes.CY = (0x05 == (psw & 0x08))
		s.ConditionCodes.AC = (0x10 == (psw & 0x10))*/
		s.SP += 2
	case 0xf2:
		if !s.ConditionCodes.S {
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xf3:
		s.UnimplementedInstruction(opcode)
	case 0xf4:
		s.UnimplementedInstruction(opcode)
	case 0xf5:
		s.Memory[s.SP-1] = s.A
		var z, s1, p, cy, ac uint8
		if s.ConditionCodes.Z {
			z = 1
		}
		if s.ConditionCodes.S {
			s1 = 1
		}
		if s.ConditionCodes.P {
			p = 1
		}
		if s.ConditionCodes.CY {
			cy = 1
		}
		if s.ConditionCodes.AC {
			ac = 1
		}
		var psw uint8 = 0
		psw |= s1 << 7
		psw |= z << 6
		psw |= ac << 4
		psw |= p << 2
		psw |= 1 << 1 // bit 1 is always 1
		psw |= cy << 0

		s.Memory[s.SP-2] = uint8(psw)
		s.SP -= 2

	case 0xf6:
		x := s.A | data[0]
		s.flagsZSP(x)
		s.ConditionCodes.CY = false
		s.A = x
		s.PC++

	case 0xf7:
		s.UnimplementedInstruction(opcode)
	case 0xf8:
		s.UnimplementedInstruction(opcode)
	case 0xf9:
		s.UnimplementedInstruction(opcode)
	case 0xfa:
		if s.ConditionCodes.S {
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xfb:
		s.IntEnable = 1
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
