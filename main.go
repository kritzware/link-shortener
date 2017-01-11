package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/redis.v5"
)

type Config struct {
	Port  int
	Redis string
}

type link_struct struct {
	Link string
}

func main() {
	router := gin.Default()

	router.StaticFile("/", "./static/index.html")
	router.POST("/link/", shortenLink)
	router.GET("/s/:token", getLink)

	router.Run(":7821")
}

func getLink(c *gin.Context) {
	token := c.Param("token")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(pong)

	link, err := client.Get("link:" + token).Result()
	if err != nil {
		panic(err)
	}
	c.Redirect(http.StatusMovedPermanently, link)
}

func shortenLink(c *gin.Context) {

	var body link_struct
	c.BindJSON(&body)
	link := body.Link

	token := randToken()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error Occurred. Please Try Again Later.",
		})
	}

	err = client.Set("link:"+token, link, 0).Err()
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"error": "Internal Server Error Occurred. Please Try Again Later.",
		})
	}

	c.JSON(200, gin.H{
		"URL":   link,
		"Token": token,
	})
}

func randToken() string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
