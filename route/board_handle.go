package route

import (
	"github.com/gin-gonic/gin"
	"leadboard/model"
	"net/http"
)

//TODO:在这里完成handle function，返回所有的leader board内容
func HandleGetBoard(g *gin.Context) {
	g.JSON(http.StatusAccepted, model.GetLeaderBoard())
}

//TODO:在这里完成返回一个用户提交历史的Handle function
func HandleUserHistory(g *gin.Context) {
	username := g.Param("username")
	err, _ := model.GetUserByName(username)
	//check if user exist
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "User doesn't exist.",
		})
		return
	}
	//get user's all submission
	_, submissons := model.GetUserSubmissions(username)
	g.JSON(http.StatusAccepted, submissons)
}

//TODO:在这里完成接受提交内容，进行评判的handle function
func HandleSubmit(g *gin.Context) {
	type SubmitForm struct {
		Username string `json:"user"`
		Avatar   string `json:"avatar"`
		Content  string `json:"content"`
	}
	var form SubmitForm
	//check if all paras are provided
	if err := g.ShouldBindJSON(&form); err != nil {
		//log.Fatal(err)
		g.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "参数不全啊",
		})
		return
	}
	//check if username is valid
	if len(form.Username) > 255 {
		g.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "用户名太长了",
		})
		return
	}
	//check if avatar is valid
	if len(form.Avatar) > 1000000 {
		g.JSON(http.StatusBadRequest, gin.H{
			"code": -2,
			"msg":  "图像太大了",
		})
		return
	}
	//check if content is valid
	err := model.CreateSubmission(form.Username, form.Avatar, form.Content)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"code": -3,
			"msg":  "提交内容非法呜呜",
		})
		return
	}
	//submit successfully and return the new leaderboard
	g.JSON(http.StatusAccepted, gin.H{
		"code": 0,
		"msg":  "提交成功",
		"data": gin.H{
			"leaderboard": model.GetLeaderBoard(),
		},
	})
}
