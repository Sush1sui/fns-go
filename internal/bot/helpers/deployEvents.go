package helpers

import (
	"log"

	"github.com/Sush1sui/fns-go/internal/bot/events"
	"github.com/bwmarrin/discordgo"
)


var EventHandlers = []any{
	events.OnSticky,
	events.OnAutoReact,
	events.OnGifAndAttachment,
	events.OnKakClaim,
	events.OnMemberBoost,
	events.OnMemberJoin,
	events.OnMemberLeave,
	events.OnStealEmoji,
	// Add more event handlers here, e.g.:
	// Go doesn't support dynamic runtime imports
	// You have to manually add each event handler
}

func DeployEvents(sess *discordgo.Session) {
	for _, handler := range EventHandlers {
		sess.AddHandler(handler)
	}
	log.Println("Event handlers deployed successfully.")
}