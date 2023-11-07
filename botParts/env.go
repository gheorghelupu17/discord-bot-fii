package botParts

import (
	"github.com/diamondburned/arikawa/v3/session"
	"log"
	"os"
)

// GetSession is used to create a session
func GetSession() (*session.Session, error) {
	// the environment variables that are needed for the bot to work
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN provided.")
	}
	// create a new session using the given Discord bot token
	s, err := session.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return s, nil
}

