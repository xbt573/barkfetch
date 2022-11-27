package console

import (
	"fmt"
	"strings"
)

type Console struct {
	lastLine    int
	offset      int
	framebuffer [][]byte
}

func (c *Console) PushChar(x, y int, char byte) {
	c.framebuffer[y][x] = char
}

func (c *Console) PushAt(x, y int, str string) {
	x += c.offset
	for _, line := range strings.Split(str, "\n") {
		for _, char := range line {
			c.PushChar(x, y, []byte(string(char))[0])
			x++
		}
		y++
		x = c.offset
	}

	c.lastLine = y - 1
}

func (c *Console) SetOffset(offset int) {
	c.offset = offset
}

func (c *Console) Render() {
	for _, line := range c.framebuffer {
		for _, char := range line {
			if char == ' ' {
				fmt.Print(" ")
				continue
			}

			fmt.Print(string(char))
		}

		fmt.Println()
	}
}

func (c *Console) Clear() {
	for i, line := range c.framebuffer {
		for k := range line {
			c.framebuffer[i][k] = ' '
		}
	}
}

func (c *Console) GetLastLine() int {
	return c.lastLine
}

func (c *Console) SetLastLine(l int) {
	c.lastLine = l
}

func NewConsole(x, y int) Console {
	arr := make([][]byte, y)
	for i := range arr {
		arr[i] = make([]byte, x)
	}

	return Console{
		framebuffer: arr,
	}
}
