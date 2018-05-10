package server

import (
	"github.com/gin-gonic/gin"
	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

//SetupServer setups the web interface server
func SetupServer(c *ymlparser.Config) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/ui/states", func(ctx *gin.Context) {
		ctx.JSON(200, c)
	})

	r.GET("/", func(ctx *gin.Context) {

		ctx.HTML(200, "index.html", gin.H{
			"Links": r.Routes(),
		})
	})

	r.GET("/ui/timeline", func(ctx *gin.Context) {

		ctx.HTML(200, "timeline.html", c)
	})
	return r
}

func main() {
	c := ymlparser.ParseStatesfile("../test/passa-states-test.yml")
	r := SetupServer(c)
	r.Run()
}
