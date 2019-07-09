package utils

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/data"
)

func WithTemplate(tmpl string, digits ...interface{}) (out string) {
	out = fmt.Sprintf(tmpl, digits...)
	return
}

func CleanName(in string) (out string) {
	out = strings.ToLower(in)
	out = JustAlphanumeric(out)
	return
}

func JustAlphanumeric(in string) (out string) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	out = reg.ReplaceAllString(in, "")

	return
}

func Courage(val int) interface{} {
	if val < 1 {
		return "-"
	}

	return val
}

func SurgeString(surge *data.Surge) string {
	var str []string

	if surge.Attack != "" {
		str = append(str, WithTemplate("attack: %s", surge.Attack))
	}

	if surge.Defense != "" {
		str = append(str, WithTemplate("defense: %s", surge.Defense))
	}

	return strings.Join(str, ", ")
}

func DiceString(dice *data.AttackDice) string {
	str := []string{}

	if dice.White > 0 {
		str = append(str, WithTemplate("white: %d", dice.White))
	}

	if dice.Black > 0 {
		str = append(str, WithTemplate("black: %d", dice.Black))
	}

	if dice.Red > 0 {
		str = append(str, WithTemplate("red: %d", dice.Red))
	}

	return strings.Join(str, ", ")
}