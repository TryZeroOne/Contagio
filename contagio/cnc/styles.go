package cnc

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/iskaa02/qalam/gradient"
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

	"reset",
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
	case "reset":
		return "\033[0m"
	}

	return ""
}

func GeneratePrompt(str string) string {

	oldStr := "{custom"

	substrings := strings.Split(str, oldStr)

	for i := 1; i < len(substrings); i++ {
		endIndex := strings.Index(substrings[i], "}")
		if endIndex == -1 {
			continue
		}
		newStr := CustomColor("custom{" + substrings[i])
		substrings[i] = newStr + substrings[i][endIndex+1:]

	}

	finalStr := strings.Join(substrings, "")

	var prompt = finalStr

	for _, i := range colors {
		prompt = strings.ReplaceAll(prompt, "{"+i+"}", GetColor(i))
	}
	return rainbow(prompt)

}

func (c *Connection) createCaptchaPrompt(code string) string {

	prompt := c.config.Auth.CaptchaPrompt

	prompt = GeneratePrompt(prompt)

	prompt = strings.ReplaceAll(prompt, "{code}", code)

	return prompt
}

func CustomColor(str string) string {
	// {custom(fg=30 bg=30 fgstyle=7)}

	var fgvar, bgvar, fgstylevar string
	var temp string
	var temparr = make([]string, 0)

	re := regexp.MustCompile(`fg=(\d+)`)
	fgMatches := re.FindStringSubmatch(str)
	if len(fgMatches) >= 2 {
		fgvar = fgMatches[1]
	}

	re = regexp.MustCompile(`bg=(\d+)`)
	bgMatches := re.FindStringSubmatch(str)
	if len(bgMatches) >= 2 {
		bgvar = bgMatches[1]
	}

	re = regexp.MustCompile(`fgstyle=(\d+)`)
	fgStyleMatches := re.FindStringSubmatch(str)
	if len(fgStyleMatches) >= 2 {
		fgstylevar = fgStyleMatches[1]

	}

	if fgstylevar == "" && bgvar == "" && fgvar == "" {
		fmt.Println("[Custom Color] Error. Invalid syntax")
		return ""
	}

	if fgstylevar != "" {
		temparr = append(temparr, fgstylevar)
	}
	if fgvar != "" {
		temparr = append(temparr, fgvar)
	}
	if bgvar != "" {
		temparr = append(temparr, bgvar)
	}
	for x, i := range temparr {
		if i != "" && x != len(temparr)-1 {
			temp += i + ";"
		}
		if x == len(temparr)-1 {
			temp += i
		}
	}

	temp = fmt.Sprintf("\x1b[%sm", temp)

	return temp

}

func (c *Connection) Cls() {
	c.conn.Write([]byte("\x1B[2J\x1B[H"))
}

func rainbow(text string) string {

	re := regexp.MustCompile(`{rainbow\((.*?)\)}`)

	matches := re.FindAllStringSubmatch(text, -1)

	if matches == nil {
		return text
	}

	var newtext string
	for _, match := range matches {
		newtext = re.ReplaceAllString(text, gradient.Rainbow().Apply(match[1]))
	}

	return newtext
}
