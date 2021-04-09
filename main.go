package main

import (
	"bufio"
	"os"
	"strconv"

	core "github.com/DigNZ/goinvaders/core"
	"github.com/DigNZ/goinvaders/system"
)

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
	var filesize uint16 = uint16(len(buffer))
	var pc uint16 = 0
	if len(os.Args) > 2 && os.Args[2] == "dasm" {
		for {
			if pc >= filesize {
				break
			}
			pc += core.Disassemble8080Op(buffer, pc)
		}
	} else {
		m := core.Machine{}
		m.Init()
		if len(os.Args) > 2 {
			addr, err := strconv.ParseInt(os.Args[2], 16, 16)
			if err != nil {
				panic(err)
			}
			m.Cpu.InitWithDataAt(buffer, int(addr))
		} else {
			m.Cpu.InitWithData(buffer)
		}
		system := system.System{Machine: &m}
		system.Start()
	}

}
