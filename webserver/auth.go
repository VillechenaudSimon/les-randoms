package webserver

import (
	"context"
	"io/ioutil"
	"les-randoms/database"
	"les-randoms/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func handleAuthLoginRoute(c *gin.Context) {
	http.Redirect(c.Writer, c.Request, Conf.AuthCodeURL(""), http.StatusTemporaryRedirect)
}

func handleAuthCallbackRoute(c *gin.Context) {
	token, err := Conf.Exchange(context.Background(), c.Request.FormValue("code"))

	// Error handling
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	res, err := Conf.Client(context.Background(), token).Get("https://discordapp.com/api/users/@me")

	// Error handling
	if err != nil || res.StatusCode != 200 {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			c.Writer.Write([]byte(err.Error()))
		} else {
			c.Writer.Write([]byte(res.Status))
		}
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	// Error handling
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	session := getSession(c)

	session.Values["authenticated"] = true
	err = session.Save(c.Request, c.Writer)

	if err != nil {
		utils.LogError("Error while saving session (authenticated value) while logging in : " + err.Error())
		return
	}

	username := string(body)
	username = username[strings.Index(username, "\"username\": \"")+13:]
	username = username[0:strings.Index(username, "\"")]
	discordId := string(body)
	discordId = discordId[strings.Index(discordId, "\"id\": \"")+7:]
	discordId = discordId[0:strings.Index(discordId, "\"")]
	avatarId := string(body)
	avatarId = avatarId[strings.Index(avatarId, "\"avatar\": \"")+11:]
	avatarId = avatarId[0:strings.Index(avatarId, "\"")]
	user, err := database.User_SelectFirst("WHERE discordId=" + discordId)
	if err != nil {
		user, _, err = database.User_CreateNew(username, discordId)
		if err != nil {
			utils.LogError("Error while inserting a new user : " + err.Error())
			return
		}
	}

	session.Values["discordId"] = discordId
	session.Values["username"] = username
	session.Values["avatarId"] = avatarId
	session.Values["userId"] = user.Id
	err = session.Save(c.Request, c.Writer)

	if err != nil {
		utils.LogError("Error while saving session (discordId value) while logging in : " + err.Error())
		return
	}

	utils.LogClassic(user.Name + " successfully logged in with discord")
	redirectToIndex(c)
}

func handleAuthLogoutRoute(c *gin.Context) {
	session := getSession(c)
	session.Values["authenticated"] = false
	session.Save(c.Request, c.Writer)
	redirectToIndex(c)
}

func redirectToAuth(c *gin.Context) {
	c.Redirect(http.StatusFound, "/auth/login")
}
