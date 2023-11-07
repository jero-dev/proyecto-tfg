// Package services holds all the services that connects repositories into a business flow
package services

import (
	"regexp"
	"strconv"
	"strings"
)

// MessageProcessor is a interface that defines the rules around what a message
// processor has to be able to perform
type MessageProcessor interface {
	ParseMessage(message string) (string, string, string, float64)
}

type MessageProcessorService struct{}

// ParseMessage takes an offer message and returns the name, platform, link and price of the offer
func (messageProcessor *MessageProcessorService) ParseMessage(message string) (string, string, string, float64) {
	var name, platform, link string
	var price float64

	gamePattern := regexp.MustCompile(`⬇️ ([^#]+) #(\w+)`)
	pricePattern := regexp.MustCompile(`(?:BAJONAZO|FLASH).+?(\d{1,2}(?:,\d{1,2})?)€`)
	linkPattern := regexp.MustCompile(`https?://\S+`)

	gameMatches := gamePattern.FindStringSubmatch(message)
	priceMatches := pricePattern.FindStringSubmatch(message)
	linkMatches := linkPattern.FindAllString(message, -1)

	if len(gameMatches) >= 3 {
		name = gameMatches[1]
		platform = gameMatches[2]
	}

	if len(priceMatches) >= 2 {
		foundPrice := strings.Replace(priceMatches[1], ",", ".", -1)
		price, _ = strconv.ParseFloat(foundPrice, 64)
	}

	if len(linkMatches) > 1 {
		goodLinks := []string{}
		for _, foundLink := range linkMatches {
			if !strings.Contains(foundLink, ")") {
				goodLinks = append(goodLinks, foundLink)
			}
		}
		link = goodLinks[0]
	}

	return name, platform, link, price
}
