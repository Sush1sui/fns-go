package commands

import (
	"github.com/bwmarrin/discordgo"
)

func SetDefaultPermsCategoryJTC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		return
	}

	category_id := i.ApplicationCommandData().GetOption("category_id").StringValue()
	if category_id == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must specify a category to set default permissions.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
    overwrites := []*discordgo.PermissionOverwrite{
        {
            ID:   i.GuildID, // @everyone
            Type: discordgo.PermissionOverwriteTypeRole,
            Deny: discordgo.PermissionSendMessages,
        },
        {
            ID:    "1299577480868528330", // music bots
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect | discordgo.PermissionVoiceSpeak,
        },
        {
            ID:    "1292473360114122784", // finest role
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionCreatePublicThreads | discordgo.PermissionSendMessages | discordgo.PermissionCreatePrivateThreads | discordgo.PermissionSendMessagesInThreads | discordgo.PermissionAddReactions | discordgo.PermissionManageThreads | discordgo.PermissionReadMessageHistory | discordgo.PermissionVoiceSpeak | discordgo.PermissionVoiceStreamVideo | discordgo.PermissionUseEmbeddedActivities | discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect,
        },
        {
            ID:    "1303998295911436309", // lvl50
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
        },
        {
            ID:    "1303998297538560060", // lvl60
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
        },
        {
            ID:    "1303998299031736393", // lvl70
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
        },
        {
            ID:    "1303998300671709186", // lvl80
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
        },
        {
            ID:    "1303998302785900544", // lvl90
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
        },
        {
            ID:    "1303998304710819940", // lvl100
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
        },
        {
            ID:    "1303916681692839956", // pioneers
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
        },
        {
            ID:    "1303924607555997776", // supporter
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
        },
        {
            ID:    "1292420325002448930", // booster
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles | discordgo.PermissionUseExternalEmojis | discordgo.PermissionUseExternalStickers,
        },
        {
            ID:    "1310186525606154340", // staff
            Type:  discordgo.PermissionOverwriteTypeRole,
            Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles | discordgo.PermissionUseExternalEmojis | discordgo.PermissionUseExternalStickers,
        },
    }

    _, err := s.ChannelEditComplex(category_id, &discordgo.ChannelEdit{
        PermissionOverwrites: overwrites,
    })
    if err != nil {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Failed to apply default permissions: " + err.Error(),
                Flags:   discordgo.MessageFlagsEphemeral,
            },
        })
        return
    }

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "Default permissions applied to the category.",
            Flags:   discordgo.MessageFlagsEphemeral,
        },
    })
}