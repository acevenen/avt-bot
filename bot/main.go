package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/joho/godotenv"
    "github.com/jonas747/dca"
)

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmd := strings.ToLower(strings.TrimSpace(m.Content))

	switch cmd {

	case "!help":
		s.ChannelMessageSend(m.ChannelID, `**AVT Bot Commands**
!apply   → Apply for coaching
!price   → Pricing & guarantee
!faq     → Common questions
!results → Student results
!socials → Links & socials`)

	case "!apply":
		s.ChannelMessageSend(m.ChannelID, `**Ready to get to work?**
Applications are reviewed personally by Ace.
Limited spots available.

👉 Apply here: https://acevenen.github.io/avt-website/apply`)

	case "!price":
		s.ChannelMessageSend(m.ChannelID, `**AVT Coaching — $1,000**
- 1-on-1 coaching with Ace for 3 months
- 24/7 direct access
- Live trade breakdowns
- NQ-specific strategy
- Coaching until you pass your first funded account

No ghosting. No upsells. One outcome.

👉 Apply: https://acevenen.github.io/avt-website/apply`)

	case "!faq":
		s.ChannelMessageSend(m.ChannelID, `**Frequently Asked Questions**

**What market do you trade?**
NQ Futures — Nasdaq 100 futures contracts.

**Do I need experience?**
No. Ace coaches all levels from zero experience upward.

**What do I need to start?**
A funded or prop firm account (Topstep, Lucid, etc.) and the willingness to put in the work.

**How long is the program?**
3 months of 1-on-1 coaching. After that, AVT Group at $147/mo.

**How do I apply?**
Type !apply or visit the apply channel.`)

	case "!results":
		s.ChannelMessageSend(m.ChannelID, `**Student Results**

+$18,686 — Chris, 15 days part-time (Topstep combine)
+$3,037  — April combine, passed Topstep $50K
+$1,510  — Tuesday NQ scalp, 4 trades, 62 minutes

Real receipts posted in #student-results
👉 Apply: https://acevenen.github.io/avt-website/apply`)

	case "!socials":
		s.ChannelMessageSend(m.ChannelID, `**Ace Venen Trading**
🌐 Website: https://acevenen.github.io/avt-website
📸 Instagram: https://instagram.com/ace.venen
🐦 Twitter/X: https://x.com/AceVenen`)
	}
}

func playIntro(s *discordgo.Session, guildID string, channelID string) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return
	}
	defer vc.Disconnect()

	time.Sleep(250 * time.Millisecond)

	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96

	encodeSession, err := dca.EncodeFile("intro.dca", opts)
	if err != nil {
		return
	}
	defer encodeSession.Cleanup()

	done := make(chan error)
	dca.NewStream(encodeSession, vc, done)
	<-done
}

func main() {
	godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	dg.AddHandler(messageHandler)

	dg.AddHandler(func(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
		log.Printf("VoiceStateUpdate fired: userID=%s channelID=%s guildID=%s", vs.UserID, vs.ChannelID, vs.GuildID)
		if vs.ChannelID == "" || vs.UserID == s.State.User.ID {
			return
		}
		if vs.BeforeUpdate != nil && vs.BeforeUpdate.ChannelID != "" {
			return
		}
		go playIntro(s, vs.GuildID, vs.ChannelID)
	})
	
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	fmt.Println("AVT Bot is running.")

	// Keep-alive for Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		http.ListenAndServe(":"+port, nil)
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	dg.Close()
}
