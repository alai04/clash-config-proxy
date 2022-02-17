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
		oriString, err := getOrigin(ctx)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "get url error")
			return
		}

		replacer := strings.NewReplacer(
			"Proxy:", "proxies:",
			"Proxy Group:", "proxy-groups:",
			"Rule:", "rules:",
		)
		ctx.String(http.StatusOK, replacer.Replace(string(oriString)))
	})

	r.GET("/ori", func(ctx *gin.Context) {
		oriString, err := getOrigin(ctx)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "get url error")
			return
		}
		ctx.String(http.StatusOK, string(oriString))
	})

	r.GET("/test", func(ctx *gin.Context) {
		const testResponse = `Proxy: "This is a test"`
		ctx.String(http.StatusOK, testResponse)
	})

	r.Run(port)
}

func getOrigin(ctx *gin.Context) (oriString []byte, err error) {
	var resp *http.Response
	url := ctx.Query("url")
	resp, err = http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
