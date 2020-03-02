package utils

import (
	"log"
	"net/url"
)

func FixURL(urlInput string) string {
	urlOut, err := url.Parse(urlInput)
	if err != nil {
		log.Println(err)
		return urlInput
	}

	return urlOut.String()
}
