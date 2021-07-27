package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const defaultPort = ":80"

var port = defaultPort

func init() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	} else if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
}

func main() {
	r := gin.Default()

	r.GET("/cc", func(ctx *gin.Context) {
		url := ctx.Query("url")
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			ctx.String(http.StatusBadRequest, "url %s error")
			return
		}

		oriString, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "read response error")
			return
		}

		replacer := strings.NewReplacer(
			"Proxy:", "proxies:",
			"Proxy Group:", "proxy-groups:",
			"Rule:", "rules:",
		)
		ctx.String(http.StatusOK, replacer.Replace(string(oriString)))
	})
	r.Run(port)
}
