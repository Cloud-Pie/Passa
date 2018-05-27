package server

import (
	"fmt"
	"os"

	"github.com/Cloud-Pie/Passa/database"
	"github.com/Cloud-Pie/Passa/ymlparser"
	"github.com/gin-gonic/gin"
)

var config *ymlparser.Config

//SetupServer setups the web interface server
func SetupServer(c *ymlparser.Config) *gin.Engine {
	r := gin.Default()
	//d, _ := os.Getwd()

	goPath := os.Getenv("GOPATH")
	r.LoadHTMLGlob(goPath + "/src/github.com/Cloud-Pie/Passa/server/templates/*")
	config = c

	r.GET("/", func(ctx *gin.Context) {

		ctx.HTML(200, "index.html", r.Routes())
	})

	r.GET("/ui/timeline", func(ctx *gin.Context) {

		ctx.HTML(200, "timeline.html", config)
	})

	statesRest := r.Group("/api/states")
	{
		statesRest.POST("/", createState)
		statesRest.GET("/", getAllStates)
		statesRest.GET("/:name", getSingleState)
		statesRest.PUT("/:name", updateState)
		statesRest.DELETE("/:name", deleteState)
	}

	r.POST("/test", func(c *gin.Context) {
		var myState ymlparser.State
		c.BindJSON(&myState)
		fmt.Printf("%v", myState)
		c.JSON(200, myState)
	})
	return r
}

func createState(c *gin.Context) {
	var newState ymlparser.State
	c.BindJSON(&newState)

	if newState.ISODate.IsZero() || newState.Services == nil { //input validation
		c.JSON(422, gin.H{
			"error": "Time or service field is empty",
		})
	} else {
		config.States = append(config.States, newState)
		c.JSON(200, gin.H{
			"data": "success",
		})
	}
}

func getAllStates(c *gin.Context) {
	fmt.Printf("%+v", config.States)
	c.JSON(200, config.States)
}
func getSingleState(c *gin.Context) {
	name := c.Params.ByName("name")
	postToReturn := database.SearchQuery(config.States, name)
	if postToReturn == -1 {
		c.JSON(422, gin.H{"error": "Not Found!"})

	} else {

		c.JSON(200, config.States[postToReturn])
	}

}
func updateState(c *gin.Context) {
	name := c.Params.ByName("name")
	var updatedState ymlparser.State
	c.BindJSON(&updatedState)
	fmt.Printf("%v", updatedState)
	posToUpdate := database.SearchQuery(config.States, name)
	if posToUpdate == -1 {
		c.JSON(422, gin.H{"error": "Not Found"})
	} else {

		config.States[posToUpdate] = updatedState

		c.JSON(200, config.States[posToUpdate])

	}

}
func deleteState(c *gin.Context) {
	name := c.Params.ByName("name")
	postToDelete := database.SearchQuery(config.States, name)
	if postToDelete == -1 { //Not Found
		c.JSON(422, gin.H{"error": "Not Found"})
	} else {
		config.States[postToDelete] = config.States[len(config.States)-1]
		config.States[len(config.States)-1] = ymlparser.State{}
		config.States = config.States[:len(config.States)-1]
		c.JSON(200, gin.H{"data": "success"})
	}
}
