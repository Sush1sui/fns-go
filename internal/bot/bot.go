package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sush1sui/fns-go/internal/bot/helpers"
	"github.com/Sush1sui/fns-go/internal/config"
	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func StartBot() {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	// create new discord session
	sess, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
    s.UpdateStatusComplex(discordgo.UpdateStatusData{
        Status: "idle",
        Activities: []*discordgo.Activity{
            {
                Name: "with Finesse!",
                Type: discordgo.ActivityTypeGame,
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

	err = repository.StickyService.DBClient.InitializeStickyChannels()
	if err != nil {
		log.Fatalf("error initializing sticky channels: %v", err)
	}

	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

