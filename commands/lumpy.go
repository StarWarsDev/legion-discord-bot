package commands

import (
	"math/rand"
	"time"
)

// Lumpy handles the !lumpy command
func Lumpy() string {
	urls := []string{
		"https://media.giphy.com/media/PSnTTyw6BdjHy/giphy.gif",
		"https://media.giphy.com/media/10QqGj0eqGOWIw/giphy.gif",
		"https://media.giphy.com/media/FY5dT7KDV2i0o/giphy.gif",
		"https://media.giphy.com/media/hffHBmxUSfHlm/giphy.gif",
	}
	rand.Seed(time.Now().Unix())
	// initialize global pseudo random generator
	randomLumpyURL := urls[rand.Intn(len(urls))]
	return randomLumpyURL
}
