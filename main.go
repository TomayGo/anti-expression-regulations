package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

var replaceMap = map[string]string{
	"original string": "replacement string",

	// Add more key-value pairs here
}

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Add event handler for message create events
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	// Wait here until CTRL-C or other termination signal is received
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	<-make(chan struct{})
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the bot was mentioned in the message
	if strings.Contains(m.Content, s.State.User.Mention()) {
		fmt.Print("Message received: ", m.Content, "\n")

		// Remove the mention from the message
		m.Content = strings.ReplaceAll(m.Content, s.State.User.Mention(), "")

		// Create a map for reverse replacements
		reverseReplaceMap := make(map[string]string, len(replaceMap))
		for old, new := range replaceMap {
			reverseReplaceMap[new] = old
		}

		// Check if the message contains a string that needs to be replaced
		for old := range replaceMap {
			if strings.Contains(m.Content, old) {
				// Create a slice of strings to hold the key-value pairs
				oldNewPairs := make([]string, 0, len(replaceMap)*2)
				for old, new := range replaceMap {
					oldNewPairs = append(oldNewPairs, old, new)
				}

				// Create a replacer object
				replacer := strings.NewReplacer(oldNewPairs...)

				// Replace the desired string
				newMessage := replacer.Replace(m.Content)

				// Send the modified message back
				_, err := s.ChannelMessageSend(m.ChannelID, newMessage)
				if err != nil {
					fmt.Println("Error sending message:", err)
				}

				return
			}
		}

		// Check if the message contains a string that needs to be reversed
		for new := range reverseReplaceMap {
			if strings.Contains(m.Content, new) {
				// Create a slice of strings to hold the key-value pairs for reverse replacements
				oldNewPairs := make([]string, 0, len(reverseReplaceMap)*2)
				for old, new := range reverseReplaceMap {
					oldNewPairs = append(oldNewPairs, old, new)
				}

				// Create a replacer object for reverse replacements
				replacer := strings.NewReplacer(oldNewPairs...)

				// Replace the desired string
				originalMessage := replacer.Replace(m.Content)

				// Send the original message back
				_, err := s.ChannelMessageSend(m.ChannelID, originalMessage)
				if err != nil {
					fmt.Println("Error sending message:", err)
				}

				return
			}
		}
	}
}
