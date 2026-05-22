package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
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
- 1-on-1 coaching with Ace
- 24/7 direct access
- Live trade breakdowns
- NQ-specific strategy
- Coaching until you're profitable

**The Guarantee**
You don't pay past $1,000 until you make your money back.
No ghosting. No upsells. One coach, one outcome.

👉 Apply: https://acevenen.github.io/avt-website/apply`)

	case "!faq":
		s.ChannelMessageSend(m.ChannelID, `**Frequently Asked Questions**

**What market do you trade?**
NQ Futures — Nasdaq 100 futures contracts.

**Do I need experience?**
No. Ace coaches all levels from zero experience upward.

**What do I need to start?**
A funded account or prop firm account (Topstep, Lucid, etc.) and the willingness to put in the work.

**How long is the program?**
Coaching continues until you're profitable. No arbitrary end date.

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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	dg.AddHandler(messageHandler)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	fmt.Println("AVT Bot is running. Press CTRL+C to stop.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	dg.Close()
}
