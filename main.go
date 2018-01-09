package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var VERIFY_TOKEN string

func init() {
	VERIFY_TOKEN = os.Getenv("VERIFY_TOKEN")
}

/**
 * Hook GET request
 */

type HookVerifyQuery struct {
	VerifyToken string `form:"hub.verify_token"`
	Challenge   string `form:"hub.challenge"`
	Mode        string `form:"hub.mode"`
}

func hookVerify(c *gin.Context) {
	h := HookVerifyQuery{}
	if err := c.Bind(&h); err != nil {
		log.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	log.Info(h)
	if h.VerifyToken != VERIFY_TOKEN {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	if h.Mode != "subscribe" {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.String(200, h.Challenge)
}

/**
 * Hook POST request
 */

type HookMessage struct {
	Message string `json:"message"`
	Sender  struct {
		Id string `json:"id"`
	}
	Recipient struct {
		Id string `json:"id"`
	}
}

type HookEntry struct {
	Id        string        `json:"id"`
	Time      time.Time     `json:"time"`
	Messaging []HookMessage `json:"messaging"`
}

type Hook struct {
	Object string      `json:"object"`
	Entry  []HookEntry `json:"entry"`
}

func hook(c *gin.Context) {
	h := Hook{}
	if err := c.Bind(&h); err != nil {
		log.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	log.Info(h)
	c.String(200, "EVENT_RECEIVED")
}

/**
 * main
 */

func main() {
	r := gin.Default()
	r.GET("/hook", hookVerify)
	r.POST("/hook", hook)
	r.Run() // listen and serve on 0.0.0.0:8080
}
