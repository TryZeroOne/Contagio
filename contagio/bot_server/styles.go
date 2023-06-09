package bot_server

import (
	"fmt"
	"strings"
)

type Attribute int

const (
	Black Attribute = iota + 90
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

var colors = []string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
}

func BlackColor() string { return fmt.Sprintf("\x1b[%dm", Black) }

func RedColor() string { return fmt.Sprintf("\x1b[%dm", Red) }

func GreenColor() string { return fmt.Sprintf("\x1b[%dm", Green) }

func YellowColor() string { return fmt.Sprintf("\x1b[%dm", Yellow) }

func BlueColor() string { return fmt.Sprintf("\x1b[%dm", Blue) }

func MagentaColor() string { return fmt.Sprintf("\x1b[%dm", Magenta) }

func CyanColor() string { return fmt.Sprintf("\x1b[%dm", Cyan) }

func WhiteColor() string { return fmt.Sprintf("\x1b[%dm", White) }

func GetColor(color string) string {

	switch color {
	case "black":
		return BlackColor()
	case "red":
		return RedColor()
	case "green":
		return GreenColor()
	case "yellow":
		return YellowColor()
	case "blue":
		return BlueColor()
	case "magenta":
		return MagentaColor()
	case "cyan":
		return CyanColor()
	case "white":
		return WhiteColor()
	}

	return ""
}

func GeneratePrompt(str string) string {

	var prompt = str
	for _, i := range colors {
		prompt = strings.ReplaceAll(prompt, "{"+i+"}", GetColor(i))
	}
	return prompt

}
