package route

import (
	"github.com/gin-gonic/gin"
	"leadboard/model"
	"net/http"
)

//TODO:完成下方的两个Handle Function，其中第一个作为中间件使用，第二个处理投票逻辑

func CheckUserAgent(g *gin.Context) {
	//用于检查User Agent的中间件
	userAgent := g.Request.UserAgent()
	//TODO:在这里完成判断User Agent的逻辑，最简单的方法是判断User Agent是否为空字符串
	if userAgent == "" {
		g.JSON(http.StatusForbidden, gin.H{
			"msg": "No Robots!",
		})
		g.Abort()
	} else {
		g.Next()
	}
}

func HandleVote(g *gin.Context) {
	type VoteForm struct {
		UserName string `json:"user"`
	}
	var form VoteForm
	if err := g.ShouldBindJSON(&form); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid Form",
		})
	} else {
		//TODO:完成投票数+1这一操作，注意不要出现并发上的问题，加油 qwq
		//推荐自己完成，也可以使用model/user.go中给出的方法
		err = model.AddVoteForUser(form.UserName)
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{
				"code": -1,
			})
		} else {
			g.JSON(http.StatusAccepted, gin.H{
				"code": 0,
				"data": gin.H{
					"leaderboard": model.GetLeaderBoard(),
				},
			})
		}
	}
}
