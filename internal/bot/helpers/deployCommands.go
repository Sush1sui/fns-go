package helpers

import (
	"fmt"
	"log"

	"github.com/Sush1sui/fns-go/internal/bot/commands"
	"github.com/bwmarrin/discordgo"
)

// List all slash commands here
var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "sticky-get-list", // Use underscore, not space
		Description: "Replies with a list of sticky channels",
		Type:        discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	{
		Name:        "sticky-delete-all",
		Description: "Deletes all sticky channels",
		Type:        discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	{
		Name:        "sticky-delete",
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
		Name:        "sticky-create",
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
		Name:        "sticky-set-message",
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
	{
		Name:        "edit-kc",
		Description: "Sets a timer for kak/trash claim",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "timer",
				Description: "The number of seconds to set the timer for",
				Required:    true,
			},
		},
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionKickMembers); return &p }(),
	},
	{
		Name:        "vanity-add",
		Description: "Adds a user for vanity exemption",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to add for vanity exemption",
				Required:    true,
			},
		},
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	{
		Name:        "vanity-remove",
		Description: "Removes a user from vanity exemption",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to remove from vanity exemption",
				Required:    true,
			},
		},
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	{
		Name:        "default-perms-category-jtc",
		Description: "Sets default permissions for a category (for JTCs)",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "category",
				Description: "The category to set default permissions for",
				Required:    true,
			},
		},
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(), // Administrators only
	},
	// Add more commands here
}

// Map command names to handler functions
var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"sticky-create":    commands.StickyCreate,
	"sticky-get-list": commands.StickyGetAll,
	"sticky-delete":   commands.StickyDelete,
	"sticky-delete-all": commands.StickyDeleteAll,
	"sticky-set-message": commands.StickySetMessage,
	"edit-kc": commands.KakClaimSetTimer,
	"vanity-add": commands.VanityAdd,
	"vanity-remove": commands.VanityRemove,
	"default-perms-category-jtc": commands.SetDefaultPermsCategoryJTC,
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
					fmt.Printf("Unknown command: %s\n", i.ApplicationCommandData().Name)
					fmt.Printf("Available commands: %v\n", CommandHandlers)
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