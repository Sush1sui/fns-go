package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sush1sui/fns-go/internal/bot/events"
	"github.com/Sush1sui/fns-go/internal/bot/helpers"
	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/Sush1sui/fns-go/internal/config"
	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func StartBot() {
	// Load configuration
	cfg, err := config.New()
	if err != nil{
		panic(err)
	}


	// create new discord session
	if cfg.DiscordToken == "" {
		fmt.Println("Bot token not found")
	}
	sess, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildPresences | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
    s.UpdateStatusComplex(discordgo.UpdateStatusData{
        Status: "idle",
        Activities: []*discordgo.Activity{
            {
                Name: "Do it with Finesse!",
            },
        },
    })
	})

	err = sess.Open()
	if err != nil {
		log.Fatalf("error opening connection to Discord: %v", err)
	}
	defer sess.Close()

	// Deploy commands
	helpers.DeployCommands(sess)

	// Deploy events
	helpers.DeployEvents(sess)

	// periodic vanity scan
	go func() {
		for {
			helpers.ScanForVanityLinks(sess)
			time.Sleep(time.Hour) // sleep for 1 hour
		}
	}()

	// initialize nickname requests
	go func() {
		err := repository.NicknameRequestService.DBClient.InitializeNicknameRequests(sess)
		if err != nil {
			fmt.Println("Error initializing nickname requests:", err)
		}
	}()

	go func() {
		err = repository.StickyService.DBClient.InitializeStickyChannels()
		if err != nil {
			fmt.Println("Error initializing sticky channels:", err)
		}
	}()
	if err != nil {
		log.Fatalf("error initializing sticky channels: %v", err)
	}

	// initialize boost cache
	go events.SyncMemberBoostCache(sess, common.GuildID)

	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

