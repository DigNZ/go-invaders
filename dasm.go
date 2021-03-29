package main

import (
	"bufio"
	"fmt"
	"os"
)

func Disassemble8080Op(buffer []byte, pc int) int {
	code := buffer[pc]
	opbytes := 1
	fmt.Printf("%02X", code)
	switch code {
	case 0x00:
		fmt.Print("\tNOP")
	case 0x01:
		fmt.Printf("\tLXI\tB,#$%02X%02X", buffer[pc+2], buffer[pc+1])
		opbytes = 3
	case 0x02:
		fmt.Print("\tSTAX B")
	case 0x03:
		fmt.Print("\tINX B")
	case 0x04:
		fmt.Print("\tINR B")
	case 0x05:
		fmt.Print("\tDCR B")
	case 0x06:
		fmt.Printf("\tMVI B, #$%02X", buffer[pc+1])
		opbytes = 2
	case 0x07:
		fmt.Print("\tRLC")
	case 0x08:
		fmt.Print("")
	case 0x09:
		fmt.Print("\tDAD B")
	case 0x0a:
		fmt.Print("\tLDAX B")
	case 0x0b:
		fmt.Print("\tLDAX B")
	case 0x0c:
		fmt.Print("\tINR C")
	case 0x0d:
		fmt.Print("\tDCR C")
	case 0x0e:
		fmt.Printf("\tMVI C, #$%02X", buffer[pc+1])
		opbytes = 2
	case 0x0f:
		fmt.Print("\tRRC")
	case 0x10:
		fmt.Print("")
	case 0x11:
		fmt.Printf("\tLXI D, #$%02X%02X", buffer[pc+2], buffer[pc+1])
		opbytes = 3
	case 0x12:
		fmt.Print("\tSTAX D")
	case 0x13:
		fmt.Print("\tINX D")
	case 0x14:
		fmt.Print("\tINR D")
	case 0x15:
		fmt.Print("\tDCR D")
	case 0x16:
		fmt.Printf("\tMVI D, #$%02X", buffer[pc+1])
		opbytes = 2
	case 0x17:
		fmt.Print("\tRAL")
	case 0x18:
		fmt.Print("")
	case 0x19:
		fmt.Print("\tDAD D")
	case 0x1a:
		fmt.Print("\tLDAX D")
	case 0x1b:
		fmt.Print("\tDCX D")
	case 0x1c:
		fmt.Print("\tINR E")
	case 0x1d:
		fmt.Print("\tDCR E")
	case 0x1e:
		fmt.Printf("\tMVI E, #$%02X", buffer[pc+1])
		opbytes = 2
	case 0x1f:
		fmt.Print("\tRAR")
	case 0x20:
		fmt.Print("\tRIM")
	case 0x21:
		fmt.Printf("\tLXI H, #$%02X%02X", buffer[pc+2], buffer[pc+1])
		opbytes = 3
	case 0x22:
		fmt.Printf("\tSHLD #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0x23:
		fmt.Print("\tINX H")
	case 0x24:
		fmt.Print("\tINR H")
	case 0x25:
		fmt.Print("\tDCR H")
	case 0x26:
		fmt.Printf("\tMVI H, #$%02X", buffer[pc+1])
		opbytes = 2
	case 0x27:
		fmt.Print("\tDAA")
	case 0x28:
		fmt.Print("")
	case 0x29:
		fmt.Print("\tDAD H")
	case 0x2a:
		fmt.Printf("\tLHLD #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0x2b:
		fmt.Printf("\tDCX H")
	case 0x2c:
		fmt.Print("\tCMA")
	case 0x2d:
		fmt.Print("\tDCR L")
	case 0x2e:
		fmt.Printf("\tMVI L, #$%02X", buffer[pc+1])
		opbytes = 2
	case 0x2f:
		fmt.Print("\tCMA")
	case 0x30:
		fmt.Print("\tSIM")
	case 0x31:
		fmt.Printf("\tLXI SP, #$%02X%02X", buffer[pc+2], buffer[pc+1])
		opbytes = 3
	case 0x32:
		//TODO: work this one out
		fmt.Printf("\tSTA #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0x33:
		fmt.Print("\tINX SP")
	case 0x34:
		fmt.Print("\tINR M")
	case 0x35:
		fmt.Print("\tDCR M")
	case 0x36:
		fmt.Printf("\tMVI M, #$%02X", buffer[pc+1])
		opbytes = 2
	case 0x37:
		fmt.Print("\tSTC")
	case 0x38:
		fmt.Print("")
	case 0x39:
		fmt.Print("\tDAD SP")
	case 0x3a:
		fmt.Printf("\tLDA #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0x3b:
		fmt.Print("\tDCX SP")
	case 0x3c:
		fmt.Print("\tINR A")
	case 0x3d:
		fmt.Print("\tDCR A")
	case 0x3e:
		fmt.Printf("\tMVI A, #$%02X", buffer[pc+1])
		opbytes = 2
	case 0x3f:
		fmt.Print("\tCMC")
	case 0x40:
		fmt.Printf("\tMOV B,B")
	case 0x41:
		fmt.Printf("\tMOV B,C")
	case 0x42:
		fmt.Printf("\tMOV B,D")
	case 0x43:
		fmt.Printf("\tMOV B,E")
	case 0x44:
		fmt.Printf("\tMOV B,H")
	case 0x45:
		fmt.Printf("\tMOV B,L")
	case 0x46:
		fmt.Printf("\tMOV B,M")
	case 0x47:
		fmt.Printf("\tMOV B,A")
	case 0x48:
		fmt.Printf("\tMOV C,B")
	case 0x49:
		fmt.Printf("\tMOV C,C")
	case 0x4a:
		fmt.Printf("\tMOV C,D")
	case 0x4b:
		fmt.Printf("\tMOV C,E")
	case 0x4c:
		fmt.Printf("\tMOV C,H")
	case 0x4d:
		fmt.Printf("\tMOV C,L")
	case 0x4e:
		fmt.Printf("\tMOV C,M")
	case 0x4f:
		fmt.Printf("\tMOV C,A")
	case 0x50:
		fmt.Printf("\tMOV D,B")
	case 0x51:
		fmt.Printf("\tMOV D,C")
	case 0x52:
		fmt.Printf("\tMOV D,D")
	case 0x53:
		fmt.Printf("\tMOV D,E")
	case 0x54:
		fmt.Printf("\tMOV D,H")
	case 0x55:
		fmt.Printf("\tMOV D,L")
	case 0x56:
		fmt.Printf("\tMOV D,M")
	case 0x57:
		fmt.Printf("\tMOV D,A")
	case 0x58:
		fmt.Printf("\tMOV E,B")
	case 0x59:
		fmt.Printf("\tMOV E,C")
	case 0x5a:
		fmt.Printf("\tMOV E,D")
	case 0x5b:
		fmt.Printf("\tMOV E,E")
	case 0x5c:
		fmt.Printf("\tMOV E,H")
	case 0x5d:
		fmt.Printf("\tMOV E,L")
	case 0x5e:
		fmt.Printf("\tMOV E,M")
	case 0x5f:
		fmt.Printf("\tMOV E,A")
	case 0x60:
		fmt.Printf("\tMOV H,B")
	case 0x61:
		fmt.Printf("\tMOV H,C")
	case 0x62:
		fmt.Printf("\tMOV H,D")
	case 0x63:
		fmt.Printf("\tMOV H,E")
	case 0x64:
		fmt.Printf("\tMOV H,H")
	case 0x65:
		fmt.Printf("\tMOV H,L")
	case 0x66:
		fmt.Printf("\tMOV H,M")
	case 0x67:
		fmt.Printf("\tMOV H,A")
	case 0x68:
		fmt.Printf("\tMOV L,B")
	case 0x69:
		fmt.Printf("\tMOV L,C")
	case 0x6a:
		fmt.Printf("\tMOV L,D")
	case 0x6b:
		fmt.Printf("\tMOV L,E")
	case 0x6c:
		fmt.Printf("\tMOV L,H")
	case 0x6d:
		fmt.Printf("\tMOV L,L")
	case 0x6e:
		fmt.Printf("\tMOV L,M")
	case 0x6f:
		fmt.Printf("\tMOV L,A")
	case 0x70:
		fmt.Printf("\tMOV M,B")
	case 0x71:
		fmt.Printf("\tMOV M,C")
	case 0x72:
		fmt.Printf("\tMOV M,D")
	case 0x73:
		fmt.Printf("\tMOV M,E")
	case 0x74:
		fmt.Printf("\tMOV M,H")
	case 0x75:
		fmt.Printf("\tMOV M,L")
	case 0x76:
		fmt.Printf("\tHLT")
	case 0x77:
		fmt.Printf("\tMOV M,A")
	case 0x78:
		fmt.Printf("\tMOV A,B")
	case 0x79:
		fmt.Printf("\tMOV A,C")
	case 0x7a:
		fmt.Printf("\tMOV A,D")
	case 0x7b:
		fmt.Printf("\tMOV A,E")
	case 0x7c:
		fmt.Printf("\tMOV A,H")
	case 0x7d:
		fmt.Printf("\tMOV A,L")
	case 0x7e:
		fmt.Printf("\tMOV A,M")
	case 0x7f:
		fmt.Printf("\tMOV A,A")
	case 0x80:
		fmt.Printf("\tADD B")
	case 0x81:
		fmt.Printf("\tADD C")
	case 0x82:
		fmt.Printf("\tADD D")
	case 0x83:
		fmt.Printf("\tADD E")
	case 0x84:
		fmt.Printf("\tADD H")
	case 0x85:
		fmt.Printf("\tADD L")
	case 0x86:
		fmt.Printf("\tADD M")
	case 0x87:
		fmt.Printf("\tADD A")
	case 0x88:
		fmt.Printf("\tADC B")
	case 0x89:
		fmt.Printf("\tADC C")
	case 0x8a:
		fmt.Printf("\tADC D")
	case 0x8b:
		fmt.Printf("\tADC E")
	case 0x8c:
		fmt.Printf("\tADC H")
	case 0x8d:
		fmt.Printf("\tADC L")
	case 0x8e:
		fmt.Printf("\tADC M")
	case 0x8f:
		fmt.Printf("\tADC A")
	case 0x90:
		fmt.Printf("\tSUB B")
	case 0x91:
		fmt.Printf("\tSUB C")
	case 0x92:
		fmt.Printf("\tSUB D")
	case 0x93:
		fmt.Printf("\tSUB E")
	case 0x94:
		fmt.Printf("\tSUB H")
	case 0x95:
		fmt.Printf("\tSUB L")
	case 0x96:
		fmt.Printf("\tSUB M")
	case 0x97:
		fmt.Printf("\tSUB A")
	case 0x98:
		fmt.Printf("\tSBB B")
	case 0x99:
		fmt.Printf("\tSBB C")
	case 0x9a:
		fmt.Printf("\tSBB D")
	case 0x9b:
		fmt.Printf("\tSBB E")
	case 0x9c:
		fmt.Printf("\tSBB H")
	case 0x9d:
		fmt.Printf("\tSBB L")
	case 0x9e:
		fmt.Printf("\tSBB M")
	case 0x9f:
		fmt.Printf("\tSBB A")
	case 0xa0:
		fmt.Printf("\tANA B")
	case 0xa1:
		fmt.Printf("\tANA C")
	case 0xa2:
		fmt.Printf("\tANA D")
	case 0xa3:
		fmt.Printf("\tANA E")
	case 0xa4:
		fmt.Printf("\tANA H")
	case 0xa5:
		fmt.Printf("\tANA L")
	case 0xa6:
		fmt.Printf("\tANA M")
	case 0xa7:
		fmt.Printf("\tANA A")
	case 0xa8:
		fmt.Printf("\tXRA B")
	case 0xa9:
		fmt.Printf("\tXRA C")
	case 0xaa:
		fmt.Printf("\tXRA D")
	case 0xab:
		fmt.Printf("\tXRA E")
	case 0xac:
		fmt.Printf("\tXRA H")
	case 0xad:
		fmt.Printf("\tXRA L")
	case 0xae:
		fmt.Printf("\tXRA M")
	case 0xaf:
		fmt.Printf("\tXRA A")
	case 0xB0:
		fmt.Printf("\tORA B")
	case 0xB1:
		fmt.Printf("\tORA C")
	case 0xB2:
		fmt.Printf("\tORA D")
	case 0xB3:
		fmt.Printf("\tORA E")
	case 0xB4:
		fmt.Printf("\tORA H")
	case 0xB5:
		fmt.Printf("\tORA L")
	case 0xB6:
		fmt.Printf("\tORA M")
	case 0xB7:
		fmt.Printf("\tORA A")
	case 0xB8:
		fmt.Printf("\tCMP B")
	case 0xB9:
		fmt.Printf("\tCMP C")
	case 0xBa:
		fmt.Printf("\tCMP D")
	case 0xBb:
		fmt.Printf("\tCMP E")
	case 0xBc:
		fmt.Printf("\tCMP H")
	case 0xBd:
		fmt.Printf("\tCMP L")
	case 0xBe:
		fmt.Printf("\tCMP M")
	case 0xBf:
		fmt.Printf("\tCMP A")
	case 0xc0:
		fmt.Printf("\tRNZ")
	case 0xc1:
		fmt.Printf("\tPOP B")
	case 0xc2:
		fmt.Printf("\tJNZ #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xc3:
		fmt.Printf("\tJMP #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xc4:
		fmt.Printf("\tCNZ #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xc5:
		fmt.Printf("\tPUSH B")
	case 0xc6:
		fmt.Printf("\tADI #$%02X", buffer[pc+1])
		opbytes = 2
	case 0xc7:
		fmt.Printf("\tRST 0")
	case 0xc8:
		fmt.Printf("\tRZ")
	case 0xc9:
		fmt.Printf("\tRET")
	case 0xca:
		fmt.Printf("\tJZ #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xcb:
		fmt.Printf("")
	case 0xcc:
		fmt.Printf("\tCZ #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xcd:
		fmt.Printf("\tCALL #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xce:
		fmt.Printf("\tACI #$%02X", buffer[pc+1])
		opbytes = 2
	case 0xcf:
		fmt.Printf("\tRST 1")
	case 0xd0:
		fmt.Printf("\tRNC")
	case 0xd1:
		fmt.Printf("\tPOP D")
	case 0xd2:
		fmt.Printf("\tJNC #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xd3:
		fmt.Printf("\tOUT #$%02X", buffer[pc+1])
		opbytes = 2
	case 0xd4:
		fmt.Printf("\tCNC #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xd5:
		fmt.Printf("\tPUSH D")
	case 0xd6:
		fmt.Printf("\tSUI #$%02X", buffer[pc+1])
		opbytes = 2
	case 0xd7:
		fmt.Printf("\tRST 2")
	case 0xd8:
		fmt.Printf("\tRC")
	case 0xd9:
		fmt.Printf("")
	case 0xda:
		fmt.Printf("\tJC #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xdb:
		fmt.Printf("\tIN #$%02X", buffer[pc+1])
		opbytes = 2
	case 0xdc:
		fmt.Printf("\tCC #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xdd:
		fmt.Printf("")
	case 0xde:
		fmt.Printf("\tSBI #$%02X", buffer[pc+1])
		opbytes = 2
	case 0xdf:
		fmt.Printf("\tRST 3")
	case 0xe0:
		fmt.Printf("\tRPO")
	case 0xe1:
		fmt.Printf("\tPOP H")
	case 0xe2:
		fmt.Printf("\tJPO #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xe3:
		fmt.Printf("\tXHTL")
	case 0xe4:
		fmt.Printf("\tCPO #$%02X%02X", buffer[pc+1], buffer[pc+2])
	case 0xe5:
		fmt.Printf("\tPUSH H")
	case 0xe6:
		fmt.Printf("\tANI #$%02X", buffer[pc+1])
		opbytes = 2
	case 0xe7:
		fmt.Printf("\tRST 4")
	case 0xe8:
		fmt.Printf("\tRPE")
	case 0xe9:
		fmt.Printf("\tPCHL")
	case 0xea:
		fmt.Printf("\tJPE #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xeb:
		fmt.Printf("\tXCHG")
	case 0xec:
		fmt.Printf("\tCPE #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xee:
		fmt.Printf("\tXRI #$%02X", buffer[pc+1])
		opbytes = 2
	case 0xef:
		fmt.Printf("\tRST 5")
	case 0xf0:
		fmt.Printf("\tRP")
	case 0xf1:
		fmt.Print("\tPOP PSW")
	case 0xf2:
		fmt.Printf("\tJP #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xf3:
		fmt.Printf("\tDI")
	case 0xf4:
		fmt.Printf("\tCP #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xf5:
		fmt.Printf("\tPUSH PSW")
	case 0xf6:
		fmt.Printf("\tORI #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xf7:
		fmt.Printf("\tRST 6")
	case 0xf8:
		fmt.Printf("\tRM")
	case 0xf9:
		fmt.Printf("\tSPHL")
	case 0xfa:
		fmt.Printf("\tJM #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xfb:
		fmt.Printf("\tEI")
	case 0xfc:
		fmt.Printf("\tCM #$%02X%02X", buffer[pc+1], buffer[pc+2])
		opbytes = 3
	case 0xfd:
		fmt.Printf("")
	case 0xfe:
		fmt.Printf("\tCPI #$%02X", buffer[pc+1])
		opbytes = 2
	case 0xff:
		fmt.Printf("\tRST 7")
	default:
		panic("\tUNIMP")
	}
	fmt.Println("")
	return opbytes
}
func RetrieveROM(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}
func main() {
	if len(os.Args) < 2 {
		panic("Invalid use")
	}
	filename := os.Args[1]
	buffer, err := RetrieveROM(filename)
	if err != nil {
		panic("Cannot load file")
	}
	filesize := len(buffer)
	pc := 0

	for {
		if pc >= filesize {
			break
		}
		pc += Disassemble8080Op(buffer, pc)
	}
}
