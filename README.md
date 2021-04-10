# Go Invaders

8080 Emulator for playing Space Invaders written in Go.

Following details from http://emulator101.com/

Uses Raylib via github.com/gen2brain/raylib-go/raylib


To test the 8080 emulation run

`go run main.go rom/cpudiag.bin 100`

Parameters are the path the the binary file and the offset (in hex) of where to load and start the file from.

To play space invaders you need a space invaders rom, online these are often distributed as 4 separate files so
they must be combined into a single file first.

`go run main.go roms/invaders`

Will get the game started.

Player 2 and SFX not implemented yet.


| Key    |   Action      |
| ------ | ------------- |
|  C     | Insert Coin   |
|  Space | Shoot         |
| Left   | Left          |
| Right  | Right         |
| 1      | Player1 Start | 

