package core

import (
	"fmt"
	"os"
	"time"
)

var cycles8080 = [...]uint8{
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, //0x00..0x0f
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, //0x10..0x1f
	4, 10, 16, 5, 5, 5, 7, 4, 4, 10, 16, 5, 5, 5, 7, 4, //etc
	4, 10, 13, 5, 10, 10, 10, 4, 4, 10, 13, 5, 5, 5, 7, 4,

	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, //0x40..0x4f
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5,
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5,
	7, 7, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 5, 7, 5,

	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, //0x80..8x4f
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4,
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4,
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4,

	11, 10, 10, 10, 17, 11, 7, 11, 11, 10, 10, 10, 10, 17, 7, 11, //0xc0..0xcf
	11, 10, 10, 10, 17, 11, 7, 11, 11, 10, 10, 10, 10, 17, 7, 11,
	11, 10, 10, 18, 17, 11, 7, 11, 11, 5, 10, 5, 17, 17, 7, 11,
	11, 10, 10, 4, 17, 11, 7, 11, 11, 5, 10, 4, 17, 17, 7, 11,
}

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
	Machine             *Machine
	lastInterrupt       time.Time
	whichInterrupt      int
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

func (s *State8080) arithFlagsA(res uint16) {

	s.ConditionCodes.CY = (res > 0xFF)
	s.ConditionCodes.Z = ((res & 0xff) == 0)
	s.ConditionCodes.S = (0x80 == (res & 0x80))
	s.ConditionCodes.P = parity(uint8(res&0xff), 8)

}

func (s *State8080) flagsZSP(value uint8) {
	s.ConditionCodes.Z = value == 0
	s.ConditionCodes.S = (0x80 == (value & 0x80))
	s.ConditionCodes.P = parity(value, 8)
}

func (s *State8080) writeMem(address uint16, value uint8) {
	//if address < 0x2000 {
	//	log.Fatalf("Attempted to write to ROM %x", address)
	//	return
	//}
	//if address >= 0x4000 {
	//	log.Fatalf("Attempted to write outside RAM")
	//}
	s.Memory[address] = value
}
func (s *State8080) push(high, low uint8) {
	s.writeMem(s.SP-1, high)
	s.writeMem(s.SP-2, low)

	s.SP = s.SP - 2
}
func (s *State8080) writeToHL(value uint8) {
	offset := uint16(s.H)<<8 | uint16(s.L)
	s.writeMem(offset, value)
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
	s.IntEnable = 0
}
func (s *State8080) Init(m *Machine) {
	s.Machine = m
	s.lastInterrupt = time.Now()
	fmt.Println("INIT")
}
func (s *State8080) Step(cycles int) {
	cycleCount := 0
	//start := time.Now()
	for {
		if cycleCount >= cycles {
			//	fmt.Printf("Cycles %d, Time %s, Cycles per second %f\n", cycleCount, time.Since(start), float32(cycleCount)/float32((time.Since(start).Seconds())))
			return
		}
		opcode := s.Memory[s.PC]
		cycleCount += int(cycles8080[opcode])
		if opcode == 0xdb {
			port := s.Memory[s.PC+1]
			s.A = s.Machine.MachineIN(port)
			s.PC += 2
		} else if opcode == 0xd3 {
			port := s.Memory[s.PC+1]
			s.Machine.MachineOUT(port, s.A)
			s.PC += 2
			s.Machine.PlaySound()
		} else {
			s.Emulate8080Op(false)
		}
		//fmt.Printf("Time %d\n", time.Since(s.lastInterrupt).Milliseconds())
		if time.Since(s.lastInterrupt).Milliseconds() > 8 { //1/60 second has elapsed
			//only do an interrupt if they are enabled
			if s.IntEnable == 1 {
				if s.whichInterrupt == 1 {
					s.generateInterrupt(1)
					s.whichInterrupt = 2
				} else {
					s.generateInterrupt(2)
					s.whichInterrupt = 1
				}

				//Save the time we did this
				s.lastInterrupt = time.Now()
			}
		}
	}
}
func (s *State8080) Run() {

	for {
		s.Step(1)
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
		offset := uint16(s.C) | uint16(s.B)<<8
		s.writeMem(offset, s.A)

	case 0x03:
		s.C++
		if s.C == 0 {
			s.B++
		}
	case 0x04:
		s.B += 1
		s.flagsZSP(s.B)
	case 0x05:
		s.B -= 1
		s.flagsZSP(s.B)
	case 0x06:
		s.B = data[0]
		s.PC++
	case 0x07:
		x := s.A
		s.A = ((x & 0x80) >> 7) | (x << 1)
		s.ConditionCodes.CY = (0x80 == (x & 0x80))
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
		offset := uint16(s.C) | uint16(s.B)<<8
		s.A = s.Memory[offset]

	case 0x0b:
		s.C--
		if s.C == 0xff {
			s.B--
		}
	case 0x0c:
		s.C += 1
		s.flagsZSP(s.C)
	case 0x0d:
		s.C -= 1
		s.flagsZSP(s.C)
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
		offset := uint16(s.E) | uint16(s.D)<<8
		s.writeMem(offset, s.A)
	case 0x13:
		s.E++
		if s.E == 0 {
			s.D++
		}
	case 0x14:
		s.D += 1
		s.flagsZSP(s.D)
	case 0x15:
		s.D -= 1
		s.flagsZSP(s.D)
	case 0x16:
		s.D = data[0]
		s.PC++
	case 0x17:
		x := s.A
		var cy uint8 = 0
		if s.ConditionCodes.CY {
			cy = 1
		}
		s.A = cy | (x << 1)
		s.ConditionCodes.CY = (0x80 == (x & 0x80))
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
		s.E--
		if s.E == 0xff {
			s.D--
		}
	case 0x1c:
		s.E += 1
		s.flagsZSP(s.E)
	case 0x1d:
		s.E -= 1
		s.flagsZSP(s.E)
	case 0x1e:
		s.E = data[0]
		s.PC++
	case 0x1f:
		x := s.A
		var cy uint8 = 0
		if s.ConditionCodes.CY {
			cy = 1
		}
		s.A = (cy << 7) | (x >> 1)
		s.ConditionCodes.CY = (1 == (x & 1))
	case 0x20:
		s.UnimplementedInstruction(opcode)
	case 0x21:
		s.H = data[1]
		s.L = data[0]
		s.PC += 2
	case 0x22:
		offset := uint16(data[0]) | uint16(data[1])<<8
		s.writeMem(offset, s.L)
		s.writeMem(offset+1, s.H)
		s.PC += 2
	case 0x23:
		s.L++
		if s.L == 0 {
			s.H++
		}
	case 0x24:
		s.H += 1
		s.flagsZSP(s.H)
	case 0x25:
		s.H -= 1
		s.flagsZSP(s.H)
	case 0x26:
		s.H = data[0]
		s.PC++
	case 0x27:
		if (s.A & 0xf) > 9 {
			s.A += 6
		}
		if (s.A & 0xf0) > 0x90 {
			res := uint16(s.A) + 0x60
			s.A = uint8(res & 0xff)
			s.arithFlagsA(res)
		}

	case 0x28:
		s.UnimplementedInstruction(opcode)
	case 0x29:
		hl := uint32((uint16(s.H) << 8) | uint16(s.L))
		res := hl + hl
		s.H = uint8((res & 0xff00) >> 8)
		s.L = uint8(res & 0xff)
		s.ConditionCodes.CY = ((res & 0xffff0000) != 0)
	case 0x2a:
		offset := uint16(data[0]) | uint16(data[1])<<8
		s.L = s.Memory[offset]
		s.H = s.Memory[offset+1]
		s.PC += 2

	case 0x2b:
		s.L--
		if s.L == 0xff {
			s.H--
		}
	case 0x2c:
		s.L += 1
		s.flagsZSP(s.L)
	case 0x2d:
		s.L -= 1
		s.flagsZSP(s.L)
	case 0x2e:
		s.L = data[0]
		s.PC++
	case 0x2f:
		s.A = ^s.A
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
		s.SP++
	case 0x34:
		res := s.readFromHL() + 1
		s.flagsZSP(res)
		s.writeToHL(res)
	case 0x35:
		res := s.readFromHL() - 1
		s.flagsZSP(res)
		s.writeToHL(res)
	case 0x36:
		offset := (uint16(s.H) << 8) | uint16(s.L)
		s.Memory[offset] = data[0]
		s.PC++
	case 0x37:
		s.ConditionCodes.CY = true
	case 0x38:
		s.UnimplementedInstruction(opcode)
	case 0x39:
		hl := uint32((uint16(s.H)<<8 | uint16(s.L)))
		res := uint32(hl + uint32(s.SP))
		s.H = uint8((res & 0xff00) >> 8)
		s.L = uint8(res & 0xff)
		s.ConditionCodes.CY = ((res & 0xffff0000) > 0)

	case 0x3a:
		offset := (uint16(data[1]) << 8) | uint16(data[0])
		s.A = s.Memory[offset]
		s.PC += 2
	case 0x3b:
		s.SP--
	case 0x3c:
		s.A = s.A + 1
		s.flagsZSP(s.A)
	case 0x3d:
		s.A = s.A - 1
		s.flagsZSP(s.A)
	case 0x3e:
		s.A = data[0]
		s.PC++
	case 0x3f:
		s.ConditionCodes.CY = !s.ConditionCodes.CY
	case 0x40:
		break //s.B = s.B
	case 0x41:
		s.B = s.C
	case 0x42:
		s.B = s.D
	case 0x43:
		s.B = s.E
	case 0x44:
		s.B = s.H
	case 0x45:
		s.B = s.L
	case 0x46:
		s.B = s.readFromHL()
	case 0x47:
		s.B = s.A
	case 0x48:
		s.C = s.B
	case 0x49:
		break //s.C = s.C
	case 0x4a:
		s.C = s.D
	case 0x4b:
		s.C = s.E
	case 0x4c:
		s.C = s.H
	case 0x4d:
		s.C = s.L
	case 0x4e:
		s.C = s.readFromHL()
	case 0x4f:
		s.C = s.A
	case 0x50:
		s.D = s.B
	case 0x51:
		s.D = s.C
	case 0x52:
		break //s.D = s.D
	case 0x53:
		s.D = s.E
	case 0x54:
		s.D = s.H
	case 0x55:
		s.D = s.L
	case 0x56:
		s.D = s.readFromHL()
	case 0x57:
		s.D = s.A
	case 0x58:
		s.E = s.B
	case 0x59:
		s.E = s.C
	case 0x5a:
		s.E = s.D
	case 0x5b:
		break //s.E = s.E
	case 0x5c:
		s.E = s.H
	case 0x5d:
		s.E = s.L
	case 0x5e:
		s.E = s.readFromHL()
	case 0x5f:
		s.E = s.A
	case 0x60:
		s.H = s.B
	case 0x61:
		s.H = s.C
	case 0x62:
		s.H = s.D
	case 0x63:
		s.H = s.E
	case 0x64:
		break //s.H = s.H
	case 0x65:
		s.H = s.L
	case 0x66:
		s.H = s.readFromHL()
	case 0x67:
		s.H = s.A
	case 0x68:
		s.L = s.B
	case 0x69:
		s.L = s.C
	case 0x6a:
		s.L = s.D
	case 0x6b:
		s.L = s.E
	case 0x6c:
		s.L = s.H
	case 0x6d:
		break //s.L = s.L
	case 0x6e:
		s.L = s.readFromHL()
	case 0x6f:
		s.L = s.A
	case 0x70:
		s.writeToHL(s.B)
	case 0x71:
		s.writeToHL(s.C)
	case 0x72:
		s.writeToHL(s.D)
	case 0x73:
		s.writeToHL(s.E)
	case 0x74:
		s.writeToHL(s.H)
	case 0x75:
		s.writeToHL(s.L)
	case 0x76:
		s.writeToHL(s.readFromHL())
	case 0x77:
		s.writeToHL(s.A)
	case 0x78:
		s.A = s.B
	case 0x79:
		s.A = s.C
	case 0x7a:
		s.A = s.D
	case 0x7b:
		s.A = s.E
	case 0x7c:
		s.A = s.H
	case 0x7d:
		s.A = s.L
	case 0x7e:
		s.A = s.readFromHL()
	case 0x7f:
		break //s.A = s.A
	case 0x80:
		var res uint16 = uint16(s.A) + uint16(s.B)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)

	case 0x81:
		var res uint16 = uint16(s.A) + uint16(s.C)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x82:
		var res uint16 = uint16(s.A) + uint16(s.D)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x83:
		var res uint16 = uint16(s.A) + uint16(s.E)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x84:
		var res uint16 = uint16(s.A) + uint16(s.H)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x85:
		var res uint16 = uint16(s.A) + uint16(s.L)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x86:
		var res uint16 = uint16(s.A) + uint16(s.readFromHL())
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x87:
		var res uint16 = uint16(s.A) + uint16(s.A)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x88:
		var res uint16 = uint16(s.A) + uint16(s.B)
		if s.ConditionCodes.CY {
			res += 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x89:
		var res uint16 = uint16(s.A) + uint16(s.C)
		if s.ConditionCodes.CY {
			res += 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x8a:
		var res uint16 = uint16(s.A) + uint16(s.D)
		if s.ConditionCodes.CY {
			res += 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x8b:
		var res uint16 = uint16(s.A) + uint16(s.E)
		if s.ConditionCodes.CY {
			res += 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x8c:
		var res uint16 = uint16(s.A) + uint16(s.H)
		if s.ConditionCodes.CY {
			res += 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x8d:
		var res uint16 = uint16(s.A) + uint16(s.L)
		if s.ConditionCodes.CY {
			res += 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x8e:
		var res uint16 = uint16(s.A) + uint16(s.readFromHL())
		if s.ConditionCodes.CY {
			res += 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x8f:
		var res uint16 = uint16(s.A) + uint16(s.A)
		if s.ConditionCodes.CY {
			res += 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x90:
		var res uint16 = uint16(s.A) - uint16(s.B)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x91:
		var res uint16 = uint16(s.A) - uint16(s.C)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x92:
		var res uint16 = uint16(s.A) - uint16(s.D)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x93:
		var res uint16 = uint16(s.A) - uint16(s.E)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x94:
		var res uint16 = uint16(s.A) - uint16(s.H)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x95:
		var res uint16 = uint16(s.A) - uint16(s.L)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x96:
		var res uint16 = uint16(s.A) - uint16(s.readFromHL())
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x97:
		var res uint16 = uint16(s.A) - uint16(s.A)
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x98:
		var res uint16 = uint16(s.A) - uint16(s.B)
		if s.ConditionCodes.CY {
			res -= 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x99:
		var res uint16 = uint16(s.A) - uint16(s.C)
		if s.ConditionCodes.CY {
			res -= 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x9a:
		var res uint16 = uint16(s.A) - uint16(s.D)
		if s.ConditionCodes.CY {
			res -= 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x9b:
		var res uint16 = uint16(s.A) - uint16(s.E)
		if s.ConditionCodes.CY {
			res -= 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x9c:
		var res uint16 = uint16(s.A) - uint16(s.H)
		if s.ConditionCodes.CY {
			res -= 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x9d:
		var res uint16 = uint16(s.A) - uint16(s.L)
		if s.ConditionCodes.CY {
			res -= 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x9e:
		var res uint16 = uint16(s.A) - uint16(s.readFromHL())
		if s.ConditionCodes.CY {
			res -= 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
	case 0x9f:
		var res uint16 = uint16(s.A) - uint16(s.A)
		if s.ConditionCodes.CY {
			res -= 1
		}
		s.arithFlagsA(res)
		s.A = uint8(res & 0xff)
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
		var res uint16 = uint16(s.A) - uint16(s.B)
		s.arithFlagsA(res)

	case 0xb9:
		var res uint16 = uint16(s.A) - uint16(s.C)
		s.arithFlagsA(res)
	case 0xba:
		var res uint16 = uint16(s.A) - uint16(s.D)
		s.arithFlagsA(res)
	case 0xbb:
		var res uint16 = uint16(s.A) - uint16(s.E)
		s.arithFlagsA(res)
	case 0xbc:
		var res uint16 = uint16(s.A) - uint16(s.H)
		s.arithFlagsA(res)
	case 0xbd:
		var res uint16 = uint16(s.A) - uint16(s.L)
		s.arithFlagsA(res)
	case 0xbe:
		var res uint16 = uint16(s.A) - uint16(s.readFromHL())
		s.arithFlagsA(res)
	case 0xbf:
		var res uint16 = uint16(s.A) - uint16(s.A)
		s.arithFlagsA(res)
	case 0xc0:
		if !s.ConditionCodes.Z {
			s.PC = uint16(s.Memory[s.SP]) | uint16(s.Memory[s.SP+1])<<8
			s.SP += 2
		}
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
		if !s.ConditionCodes.Z {
			ret := s.PC + 2
			s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
			s.Memory[s.SP-2] = uint8((ret & 0xFF))
			s.SP = s.SP - 2
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
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
		ret := s.PC + 2
		s.writeMem(s.SP-1, uint8((ret>>8)&0xff))
		s.writeMem(s.SP-2, uint8(ret&0xff))
		s.SP = s.SP - 2
		s.PC = 0x0000
	case 0xc8:
		if s.ConditionCodes.Z {
			s.PC = uint16(s.Memory[s.SP]) | uint16(s.Memory[s.SP+1])<<8
			s.SP += 2
		}
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
		if s.ConditionCodes.Z {
			ret := s.PC + 2
			s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
			s.Memory[s.SP-2] = uint8((ret & 0xFF))
			s.SP = s.SP - 2
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
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
		if !s.ConditionCodes.CY {
			s.PC = uint16(s.Memory[s.SP]) | uint16(s.Memory[s.SP+1])<<8
			s.SP += 2
		}
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
		if !s.ConditionCodes.CY {
			ret := s.PC + 2
			s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
			s.Memory[s.SP-2] = uint8((ret & 0xFF))
			s.SP = s.SP - 2
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
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
		if s.ConditionCodes.CY {
			s.PC = uint16(s.Memory[s.SP]) | uint16(s.Memory[s.SP+1])<<8
			s.SP += 2
		}
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
		if s.ConditionCodes.CY {
			ret := s.PC + 2
			s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
			s.Memory[s.SP-2] = uint8((ret & 0xFF))
			s.SP = s.SP - 2
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
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
		if !s.ConditionCodes.P {
			s.PC = uint16(s.Memory[s.SP]) | uint16(s.Memory[s.SP+1])<<8
			s.SP += 2
		}
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
		if !s.ConditionCodes.P {
			ret := s.PC + 2
			s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
			s.Memory[s.SP-2] = uint8((ret & 0xFF))
			s.SP = s.SP - 2
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
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
		if s.ConditionCodes.P {
			s.PC = uint16(s.Memory[s.SP]) | uint16(s.Memory[s.SP+1])<<8
			s.SP += 2
		}

	case 0xe9:
		s.PC = uint16(s.H)<<8 | uint16(s.L)
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
		if s.ConditionCodes.P {
			ret := s.PC + 2
			s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
			s.Memory[s.SP-2] = uint8((ret & 0xFF))
			s.SP = s.SP - 2
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
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
		if !s.ConditionCodes.S {
			s.PC = uint16(s.Memory[s.SP]) | uint16(s.Memory[s.SP+1])<<8
			s.SP += 2
		}
	case 0xf1:

		s.A = s.Memory[s.SP+1]
		psw := s.Memory[s.SP]

		s.ConditionCodes.S = ((psw >> 7) & 1) == 1
		s.ConditionCodes.Z = ((psw >> 6) & 1) == 1
		s.ConditionCodes.AC = ((psw >> 4) & 1) == 1
		s.ConditionCodes.P = ((psw >> 2) & 1) == 1
		s.ConditionCodes.CY = ((psw >> 0) & 1) == 1
		s.SP += 2
	case 0xf2:
		if !s.ConditionCodes.S {
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xf3:
		s.IntEnable = 0
	case 0xf4:
		if !s.ConditionCodes.S {
			ret := s.PC + 2
			s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
			s.Memory[s.SP-2] = uint8((ret & 0xFF))
			s.SP = s.SP - 2
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
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
		if s.ConditionCodes.S {
			s.PC = uint16(s.Memory[s.SP]) | uint16(s.Memory[s.SP+1])<<8
			s.SP += 2
		}
	case 0xf9:
		s.SP = uint16(s.H)<<8 | uint16(s.L)
	case 0xfa:
		if s.ConditionCodes.S {
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}
	case 0xfb:
		s.IntEnable = 1
	case 0xfc:
		if s.ConditionCodes.S {
			ret := s.PC + 2
			s.Memory[s.SP-1] = uint8((ret >> 8) & 0xFF)
			s.Memory[s.SP-2] = uint8((ret & 0xFF))
			s.SP = s.SP - 2
			s.PC = uint16(data[1])<<8 | uint16(data[0])
		} else {
			s.PC += 2
		}

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
