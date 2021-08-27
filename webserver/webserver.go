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

	router := gin.New()

	setupRouter(router)

	setupRoutes(router)

	router.Run(":" + port)
}

func setupRouter(router *gin.Engine) {
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
}

func setupRoutes(router *gin.Engine) {
	router.GET("/", handleIndexRoute)
	router.GET("/aram", handleAramRoute)
	router.POST("/aram", handleAramRoute)
	router.GET("/players", handlePlayersRoute)
	router.GET("/database", handleDatabaseRoute)
	router.POST("/database", handleDatabaseRoute)
	router.GET("/auth/login", handleAuthLoginRoute)
	router.GET("/auth/callback", handleAuthCallbackRoute)
	router.GET("/auth/logout", handleAuthLogoutRoute)
}

func getSession(c *gin.Context) *sessions.Session {
	session, _ := CookieStore.Get(c.Request, "les-randoms-cookie")
	return session
}

func isAuthentified(s *sessions.Session) bool {
	auth, ok := s.Values["authenticated"].(bool)
	return !ok || !auth
}
