//Package server provides routes for web interface
package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

//StartServer starts the web interface server
func StartServer(c *ymlparser.Config) *gin.Engine {
	r := gin.Default()
	r.GET("/ui/states", func(ctx *gin.Context) {
		ctx.JSON(200, c)
	})

	r.GET("/", func(ctx *gin.Context) {
		htmlToSend := "<html><body><ul>"
		for _, rt := range r.Routes() {
			htmlToSend += fmt.Sprintf("<li><a href=\"%s\">%s</a>", rt.Path, rt.Path)
		}

		htmlToSend += "</ul></body></html>"
		ctx.Data(200, "text/html; charset=utf-8", []byte(htmlToSend))
	})

	go r.Run() //NOTE: there is also a go here

	return r
}
