package server

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.lrz.de/ga53lis/PASSA/ymlparser"
)

var config *ymlparser.Config

//SetupServer setups the web interface server
func SetupServer(c *ymlparser.Config) *gin.Engine {
	r := gin.Default()
	//d, _ := os.Getwd()

	goPath := os.Getenv("GOPATH")
	r.LoadHTMLGlob(goPath + "/src/gitlab.lrz.de/ga53lis/PASSA/server/templates/*")
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
		myState.ISODate, _ = time.Parse(ymlparser.TimeLayout, myState.Time)
		fmt.Printf("%v", myState)
		c.JSON(200, gin.H{"ok": "ok"})
	})
	return r
}

func createState(c *gin.Context) {
	var newState ymlparser.State
	c.BindJSON(&newState)

	if newState.Time == "" || newState.Services == nil { //input validation
		c.JSON(422, gin.H{
			"error": "Fields are empty",
		})
	} else {
		isoTimeFormat, err := time.Parse(ymlparser.TimeLayout, newState.Time)
		if err != nil {
			defer c.JSON(422, gin.H{"error": "Cannot parse Time"})
			panic(err)

		}
		newState.ISODate = isoTimeFormat
		config.States = append(config.States, newState)
		c.JSON(200, config.States)
	}

}
func getAllStates(c *gin.Context) {
	c.JSON(200, config.States)
}
func getSingleState(c *gin.Context) {
	name := c.Params.ByName("name")
	postToReturn := searchQuery(config.States, name)
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
	posToUpdate := searchQuery(config.States, name)
	if posToUpdate == -1 {
		c.JSON(422, gin.H{"error": "Not Found"})
	} else {
		isoTimeFormat, err := time.Parse(ymlparser.TimeLayout, updatedState.Time)
		if err != nil {
			defer c.JSON(422, gin.H{"error": "Cannot parse Time"})
			panic(err)

		}
		updatedState.ISODate = isoTimeFormat
		config.States[posToUpdate] = updatedState

		c.JSON(200, config.States)
	}
}
func deleteState(c *gin.Context) {
	name := c.Params.ByName("name")
	postToDelete := searchQuery(config.States, name)
	if postToDelete == -1 { //Not Found
		c.JSON(422, gin.H{"error": "Not Found"})
	} else {
		config.States[postToDelete] = config.States[len(config.States)-1]
		config.States[len(config.States)-1] = ymlparser.State{}
		config.States = config.States[:len(config.States)-1]
		c.JSON(200, config.States)
	}
}

func searchQuery(currentStates []ymlparser.State, searchName string) int {

	for idx := range currentStates {
		if currentStates[idx].Name == searchName {
			return idx
		}
	}
	return -1
}
