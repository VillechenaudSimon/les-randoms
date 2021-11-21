package webserver

import (
	"errors"
	"les-randoms/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var CookieStore *sessions.CookieStore
var Router *gin.Engine
var Conf *oauth2.Config

func StartWebServer() {
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

	Router = gin.New()

	setupRouter()

	setupRoutes()

	Router.Run(":" + port)
}

func setupRouter() {
	Router.Use(gin.Logger())
	Router.LoadHTMLGlob("templates/*.tmpl.html")
	Router.Static("/static", "static")
}

func setupRoutes() {
	Router.GET("/", handleIndexRoute)

	aram := Router.Group("/aram")
	aram.GET("", handleAramRoute)
	aram.POST("", handleAramRoute)
	aram.GET("/:subNavItem", handleAramRoute)
	aram.POST("/:subNavItem", handleAramRoute)

	players := Router.Group("/players")
	players.GET("", handlePlayersRoute)
	players.POST("", handlePlayersRoute)
	players.GET("/:subNavItem", handlePlayersRoute)
	players.POST("/:subNavItem", handlePlayersRoute)
	players.GET("/:subNavItem/:param1", handlePlayersRoute)
	players.POST("/:subNavItem/:param1", handlePlayersRoute)

	Router.GET("/database", handleDatabaseRoute)
	Router.POST("/database", handleDatabaseRoute)

	auth := Router.Group("/auth")
	auth.GET("/login", handleAuthLoginRoute)
	auth.GET("/callback", handleAuthCallbackRoute)
	auth.GET("/logout", handleAuthLogoutRoute)
}

/*
func handlePOSTParameters(c *gin.Context) {
	Router.HandleContext(c)
}
*/

func getSession(c *gin.Context) *sessions.Session {
	session, _ := CookieStore.Get(c.Request, "les-randoms-cookie")
	return session
}

func isNotAuthentified(s *sessions.Session) bool {
	auth, ok := s.Values["authenticated"].(bool)
	return !ok || !auth
}

func isNotAdmin(s *sessions.Session) bool {
	if isNotAuthentified(s) {
		return true
	}
	discordId, ok := s.Values["discordId"].(string)
	return !ok || !(discordId == "178853941189148672") // Discord Id of Vemuni#4770
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
