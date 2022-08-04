package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	//Middleware
	r.Use(CORS)
	//Routes
	r.GET("/leaderboard", HandleGetBoard)
	r.GET("/history/:username", HandleUserHistory)
	r.POST("/submit", HandleSubmit)
	r.POST("/vote", CheckUserAgent, HandleVote)
	//TODO:register your route here
	//for example:
	//r.POST("/create-user",HandleCreateUser)，这个接口不是要求中的，仅仅作为示例
	//r.GET("/leaderboard",HandleGetBoard)
	return r
}

//CORS Options request process
func CORS(g *gin.Context) {
	g.Header("Access-Control-Allow-Origin", "*")
	g.Header("Access-Control-Allow-Headers", "content-type")
	if g.Request.Method == "OPTIONS" {
		g.Status(http.StatusNoContent)
	}
}
