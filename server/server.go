package server

import (
	"fmt"

	"github.com/Cloud-Pie/Passa/database"
	"github.com/Cloud-Pie/Passa/ymlparser"
	"github.com/gin-gonic/gin"
)

var config *ymlparser.Config
var stateChannel chan *ymlparser.State

//SetupServer setups the web interface server
func SetupServer(sc chan *ymlparser.State) *gin.Engine {
	r := gin.Default()
	stateChannel = sc //left: global, right: func param

	//r.LoadHTMLGlob(os.Getenv("GOPATH") + "/src/github.com/Cloud-Pie/Passa/server/templates/*") //FIXME: still needs a fix

	r.GET("/", func(ctx *gin.Context) {

		ctx.JSON(200, r.Routes())
	})

	r.GET("/ui/timeline", func(ctx *gin.Context) {

		//ctx.HTML(200, "timeline.html", config)
		ctx.JSON(200, gin.H{
			"timeline": "Not working will be fixed in v1.1",
		})
	})

	statesRest := r.Group("/api/states")
	{
		statesRest.POST("/", createState)
		statesRest.GET("/", getAllStates)
		statesRest.GET("/:name", getSingleState)
		statesRest.POST("/:name", updateState)
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

		stateChannel <- &newState
		c.JSON(200, gin.H{
			"data": "success",
		})
	}
}

func getAllStates(c *gin.Context) {
	fmt.Printf("%+v", database.ReadAllStates())
	c.JSON(200, database.ReadAllStates())
}
func getSingleState(c *gin.Context) {
	name := c.Params.ByName("name")
	postToReturn, err := database.GetSingleState(name)
	if err != nil {
		c.JSON(422, gin.H{"error": "Not Found!"})

	} else {

		c.JSON(200, postToReturn)
	}

}
func updateState(c *gin.Context) {
	name := c.Params.ByName("name")
	var updatedState ymlparser.State
	c.BindJSON(&updatedState)
	fmt.Printf("%v", updatedState)

	err := database.UpdateState(updatedState, name)
	if err != nil {
		c.JSON(422, gin.H{"error": "Not Found"})
	} else {

		c.JSON(200, updatedState)

	}

}
func deleteState(c *gin.Context) {
	name := c.Params.ByName("name")
	err := database.DeleteState(name)
	if err != nil { //Not Found
		c.JSON(422, gin.H{"error": "Not Found"})
	} else {

		c.JSON(200, gin.H{"data": "success"})
	}
}
