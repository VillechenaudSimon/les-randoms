package webserver

import (
	"context"
	"io/ioutil"
	"les-randoms/utils"
	"net/http"

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

	utils.LogDebug(string(body))
}

func handleAuthLogoutRoute(c *gin.Context) {
	session := getSession(c)
	session.Values["authenticated"] = false
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/")
}
