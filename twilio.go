package main

import (
	"flag"
	"fmt"
	"os"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func sendNotification(post *Post) (err error) {
	flag.Parse()
	client := twilio.NewRestClient()
	params := &openapi.CreateMessageParams{}
	params.SetTo(os.Getenv("TWILIO_TO_NUMBER"))
	params.SetFrom(os.Getenv("TWILIO_FROM_NUMBER"))
	params.SetBody(fmt.Sprintf("Openheavens for today has been scraped successfully.\n\n%s\n\nLink: https://myopenheavens.onrender.com/posts/%s.txt\n\nPlease wait 5-10 minutes for the website to be deployed.", post.Title, *date))
	_, err = client.Api.CreateMessage(params)
	return
}
