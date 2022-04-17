package logic

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

const (
	musicCacheFolderPath       = "musics/"
	musicCacheYoutubeSubfolder = "youtube/"
	musicCacheSpotifySubfolder = "spotify/"
)

type DiscordBot struct {
	DiscordSession            *discordgo.Session
	Token                     string
	Id                        string
	LogChannelId              string
	Prefix                    string
	isMessageHandlerActivated bool
	commandParser             func(bot *DiscordBot, m *discordgo.MessageCreate) string
	commands                  map[string]func(bot *DiscordBot, m *discordgo.MessageCreate) error
	defaultCommand            func(bot *DiscordBot, m *discordgo.MessageCreate) error
	streamingSessions         map[string]*dca.StreamingSession // Mapped by Guild Id
	encodeSessions            map[string]*dca.EncodeSession    // Mapped by Guild Id
	musicQueues               map[string][]*MusicInfos         // Mapped by Guild Id
	queueAppender             map[string]chan []*MusicInfos    // Mapped by Guild Id
	queuePlayer               map[string]chan *MusicInfos      // Mapped by Guild Id
}

type MusicInfosSource string

var MusicInfosSources = struct {
	Youtube MusicInfosSource
	Spotify MusicInfosSource
}{
	Youtube: "Youtube",
	Spotify: "Spotify",
}

type MusicInfos struct {
	Id     string
	Title  string
	Url    string
	Source MusicInfosSource
}

func NewMusicInfos(id string, title string, url string, source MusicInfosSource) *MusicInfos {
	return &MusicInfos{Id: id, Title: title, Url: url, Source: source}
}

/*
type Command struct {
	Name           string
	Execute        func(bot *DiscordBot, m *discordgo.MessageCreate) error
	PossibleErrors map[string]error
}
*/

func New(prefix string, token string, logChannelId string) *DiscordBot {
	return &DiscordBot{
		Prefix:                    prefix,
		Token:                     token,
		isMessageHandlerActivated: false,
		commands:                  make(map[string]func(bot *DiscordBot, m *discordgo.MessageCreate) error),
		defaultCommand:            func(bot *DiscordBot, m *discordgo.MessageCreate) error { return nil },
		streamingSessions:         make(map[string]*dca.StreamingSession),
		encodeSessions:            make(map[string]*dca.EncodeSession),
		musicQueues:               make(map[string][]*MusicInfos),
		queueAppender:             make(map[string]chan []*MusicInfos),
		queuePlayer:               make(map[string]chan *MusicInfos),
	}
}

func (bot *DiscordBot) Start() error {
	var err error
	bot.DiscordSession, err = discordgo.New("Bot " + bot.Token)
	if err != nil {
		return err
	}

	u, err := bot.DiscordSession.User("@me")
	if err != nil {
		return err
	}
	bot.Id = u.ID

	err = bot.DiscordSession.Open()
	if err != nil {
		return err
	}

	return nil
}

func (bot *DiscordBot) IsMessageHandlerActivated() bool {
	return bot.isMessageHandlerActivated
}
