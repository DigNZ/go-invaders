package main

import (
	"bufio"
	"os"
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
	filesize := len(buffer)
	pc := 0
	if len(os.Args) > 2 && os.Args[2] == "dasm" {
		for {
			if pc >= filesize {
				break
			}
			pc += Disassemble8080Op(buffer, pc)
		}
	} else {
		s := State8080{}

	}

}
