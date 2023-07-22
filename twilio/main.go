package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var date *string

func init() {
	date = flag.String("date", "", "date of scraped post")
}

func main() {
	flag.Parse()
	client := twilio.NewRestClient()
	params := &openapi.CreateMessageParams{}
	params.SetTo(os.Getenv("TWILIO_TO_NUMBER"))
	params.SetFrom(os.Getenv("TWILIO_FROM_NUMBER"))
	params.SetBody(fmt.Sprintf("Openheavens for today has been scraped.\n\nLink: https://myopenheavens.onrender.com/%s.txt", *date))
	_, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Message sent successfully!")
	}
}
