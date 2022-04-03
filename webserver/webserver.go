package webserver

import (
	"context"
	"errors"
	"fmt"
	"les-randoms/database"
	"les-randoms/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var CookieStore *sessions.CookieStore
var Router *gin.Engine
var server http.Server
var Conf *oauth2.Config
var upgrader *websocket.Upgrader

func StartWebServer() {
	upgrader = &websocket.Upgrader{}

	CookieStore = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

	Conf = &oauth2.Config{
		RedirectURL:  os.Getenv("WEBSITE_URL") + "/auth/callback",
		ClientID:     os.Getenv("DISCORD_CLIENTID"),
		ClientSecret: os.Getenv("DISCORD_CLIENTSECRET"),
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}

	port := os.Getenv("PORT")
	if port == "" {
		utils.HandlePanicError(errors.New("$PORT must be set"))
	}

	gin.DefaultWriter = &utils.Logger
	gin.DefaultErrorWriter = &utils.Logger

	Router = gin.New()

	Router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Gin error recovered: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	setupRouter()

	setupRoutes()

	utils.LogSuccess("Webserver successfully started")
	var err error
	server, err = Router.Run(":" + port)
	if err != nil {
		utils.HandlePanicError(errors.New("An error happened while the server was running : " + err.Error()))
	}

}

func StopServer() {
	utils.LogNotNilError(server.Shutdown(context.Background()))
	utils.LogSuccess("Webserver successfully shutdowned")
}

func setupRouter() {
	Router.Use(gin.Logger())
	Router.LoadHTMLGlob("templates/*.tmpl.html")
	Router.Static("/static", "static")
}

func setupRoutes() {
	Router.GET("/", handleIndexRoute)

	lol := Router.Group("/lol")

	aram := lol.Group("/aram")
	aram.GET("", handleAramRoute)
	aram.POST("", handleAramRoute)
	aram.GET("/:subNavItem", handleAramRoute)
	aram.POST("/:subNavItem", handleAramRoute)

	players := lol.Group("/players")
	players.GET("", handlePlayersRoute)
	players.POST("", handlePlayersRoute)
	players.GET("/:subNavItem", handlePlayersRoute)
	players.POST("/:subNavItem", handlePlayersRoute)
	players.GET("/:subNavItem/:param1", handlePlayersRoute)
	players.POST("/:subNavItem/:param1", handlePlayersRoute)

	discordbot := Router.Group("/discord-bot")

	music := discordbot.Group("/music")
	music.GET("", handleMusicRoute)
	music.POST("", handleMusicRoute)
	music.GET("/:subNavItem", handleMusicRoute)
	music.POST("/:subNavItem", handleMusicRoute)
	music.GET("/:subNavItem/:order", handleMusicRoute)
	music.POST("/:subNavItem/:order", handleMusicRoute)
	music.GET("/playing/ws", handlePlayingWs)

	database := Router.Group("/database")
	database.GET("", handleDatabaseRoute)
	database.POST("", handleDatabaseRoute)
	database.GET("/:subNavItem", handleDatabaseRoute)
	database.POST("/:subNavItem", handleDatabaseRoute)

	auth := Router.Group("/auth")
	auth.GET("/login", handleAuthLoginRoute)
	auth.GET("/callback", handleAuthCallbackRoute)
	auth.GET("/logout", handleAuthLogoutRoute)
}

func getSession(c *gin.Context) *sessions.Session {
	session, _ := CookieStore.Get(c.Request, "les-randoms-cookie")
	return session
}

func isNotAuthentified(s *sessions.Session) bool {
	auth, ok := s.Values["authenticated"].(bool)
	return !ok || !auth
}

func getUsername(s *sessions.Session) string {
	if isNotAuthentified(s) {
		return ""
	}
	return s.Values["username"].(string)
}

func getDiscordId(s *sessions.Session) string {
	if isNotAuthentified(s) {
		return ""
	}
	return s.Values["discordId"].(string)
}

func getAvatarId(s *sessions.Session) string {
	if isNotAuthentified(s) {
		return ""
	}
	return s.Values["avatarId"].(string)
}

func getUserId(s *sessions.Session) int {
	if isNotAuthentified(s) {
		return 0
	}
	return s.Values["userId"].(int)
}

func getAccessStatus(s *sessions.Session, path string) int {
	if isNotAuthentified(s) {
		return database.RightTypes.Hidden // Default right access value for non authentified users
	}
	accessRight, err := database.AccessRight_SelectFirst("WHERE userId=" + fmt.Sprint(getUserId(s)) + " AND path='" + path + "'")
	if err != nil {
		return database.RightTypes.Hidden // Default right access value for users without a specified access right
	}
	return accessRight.RightType
}
