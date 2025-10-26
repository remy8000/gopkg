package tools

import (
	"fmt"
)

type Color string

const (
	ERROR   Color = "error"
	WARNING Color = "warning"
	INFO    Color = "info"
	PINK    Color = "pink"
	BLUE    Color = "blue"
	CYAN    Color = "cyan"
	MAGENTA Color = "magenta"
	YELLOW  Color = "yellow"
)

func CLIprintColored(text string, color Color) string {
	var code string
	switch color {
	case ERROR:
		code = "\033[31m"
	case WARNING:
		code = "\033[38;5;208m" // orange (256-color, may not work in all terminals)
	case YELLOW:
		code = "\033[93m" // bright yellow
	case INFO:
		code = "\033[32m"
	case PINK:
		code = "\033[95m" // bright magenta (pink)
	case BLUE:
		code = "\033[94m" // bright blue
	case CYAN:
		code = "\033[96m" // bright cyan
	case MAGENTA:
		code = "\033[95m" // bright magenta
	default:
		code = "\033[0m" // default color
	}

	reset := "\033[0m"
	return fmt.Sprintf("%s%s%s", code, text, reset)
}
