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

func DiceString(dice *data.AttackDice) string {
	var strSlice []string

	if dice.White > 0 {
		strSlice = append(strSlice, WithTemplate("white: %d", dice.White))
	}

	if dice.Black > 0 {
		strSlice = append(strSlice, WithTemplate("black: %d", dice.Black))
	}

	if dice.Red > 0 {
		strSlice = append(strSlice, WithTemplate("red: %d", dice.Red))
	}

	return strings.Join(strSlice, ", ")
}
