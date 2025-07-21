package helpers

import (
	"log"

	"github.com/Sush1sui/fns-go/internal/bot/commands"
	"github.com/bwmarrin/discordgo"
)

// List all slash commands here
var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "sticky_get_list", // Use underscore, not space
		Description: "Replies with a list of sticky channels",
		Type:        discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	{
		Name:        "sticky_delete_all",
		Description: "Deletes all sticky channels",
		Type:        discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	{
		Name:        "sticky_delete",
		Description: "Deletes a sticky channel",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The channel to delete from sticky channels",
				Required:    true,
			},
		},
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	{
		Name:        "sticky_create",
		Description: "Creates a sticky channel",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The channel to create a sticky channel in",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The message to send in the sticky channel",
				Required:    false,
			},
		},
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	{
		Name:        "sticky_set_message",
		Description: "Sets the message for a sticky channel",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The channel to set the sticky message for",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The message to set for the sticky channel",
				Required:    true,
			},
		},
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	// Add more commands here
}

// Map command names to handler functions
var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"sticky_create":    commands.StickyCreate,
	"sticky_get_list": commands.StickyGetAll,
	"sticky_delete":   commands.StickyDelete,
	"sticky_delete_all": commands.StickyDeleteAll, // Match the command name exactly
	// Add more: "hello": commands.HelloCommand, etc.
}

func DeployCommands(sess *discordgo.Session) {
	// Remove all global commands
	globalCmds, err := sess.ApplicationCommands(sess.State.User.ID, "")
	if err == nil {
			for _, cmd := range globalCmds {
					err := sess.ApplicationCommandDelete(sess.State.User.ID, "", cmd.ID)
					if err != nil {
							log.Printf("Failed to delete global command %s: %v", cmd.Name, err)
					} else {
							log.Printf("Deleted global command: %s", cmd.Name)
					}
			}
	}

	// Bulk overwrite commands for each guild (this replaces all commands)
	guilds := sess.State.Guilds
	for _, guild := range guilds {
			_, err := sess.ApplicationCommandBulkOverwrite(sess.State.User.ID, guild.ID, SlashCommands)
			if err != nil {
					log.Fatalf("Cannot create slash commands for guild %s: %v", guild.ID, err)
			}
	}

	// Register handler for slash commands
	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if handler, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
					handler(s, i)
			} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
									Content: "Unknown command.",
									Flags:   discordgo.MessageFlagsEphemeral,
							},
					})
			}
	})

	log.Println("Slash commands deployed successfully.")
}