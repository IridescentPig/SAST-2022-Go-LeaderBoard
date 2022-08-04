package model

import (
	"errors"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Submission struct {
	ID        uint    `gorm:"not null;autoIncrement"`
	UserName  string  `gorm:"type:varchar(255);"`
	Avatar    string  //头像base64，也可以是一个头像链接
	CreatedAt int64   //提交时间
	Score     float64 //评测成绩
	Subscore1 float64 //评测小分
	Subscore2 float64 //评测小分
	Subscore3 float64 //评测小分
}

//GetSub struct to get data from database
type GetSub struct {
	UserName  string  `json:"user"`
	Avatar    string  `json:"avatar"`
	CreatedAt int64   `json:"time"`
	Score     float64 `json:"score"`
	Subscore1 float64
	Subscore2 float64
	Subscore3 float64
}

//ReturnSub struct to return
type ReturnSub struct {
	UserName  string    `json:"user"`
	Avatar    string    `json:"avatar"`
	CreatedAt int64     `json:"time"`
	Score     float64   `json:"score"`
	UserVotes uint      `json:"votes"`
	Subscores []float64 `json:"subs"`
}

//CreateSubmission func to create submission
func CreateSubmission(name string, avatar string, content string) error {
	Score, Subscore1, Subscore2, Subscore3, err := GetScore(content)
	if err != nil {
		return err
	}
	err, _ = GetUserByName(name)
	if err != nil {
		errCu, _ := CreateUser(name)
		if errCu != nil {
			return errCu
		}
	}
	submission := Submission{
		UserName:  name,
		Avatar:    avatar,
		CreatedAt: time.Now().Unix(),
		Score:     Score,
		Subscore1: Subscore1,
		Subscore2: Subscore2,
		Subscore3: Subscore3,
	}
	tx := DB.Create(&submission)
	return tx.Error
}

//UserSubmission struct used for func GetUserSubmissions
type UserSubmission struct {
	Score     float64   `json:"score"`
	SubScore  []float64 `json:"subs"`
	CreatedAt int64     `json:"time"`
}

//GetUserSubmissions func to get user's all submissions
func GetUserSubmissions(username string) (error, []UserSubmission) {
	//返回某一用户的所有提交
	//在查询时可以使用.Order()来控制结果的顺序，详见https://gorm.io/zh_CN/docs/query.html#Order
	//当然，也可以查询后在这个函数里手动完成排序
	var submissions []Submission
	var ReturnSubmissions []UserSubmission
	tx := DB.Model(&Submission{}).Where("user_name=?", username).Order("created_at desc").Find(&submissions)
	for _, sub := range submissions {
		var subs []float64
		subs = append(subs, sub.Subscore1)
		subs = append(subs, sub.Subscore2)
		subs = append(subs, sub.Subscore3)
		temp := UserSubmission{
			Score:     sub.Score,
			SubScore:  subs,
			CreatedAt: sub.CreatedAt,
		}
		ReturnSubmissions = append(ReturnSubmissions, temp)
	}
	return tx.Error, ReturnSubmissions
}

//GetLeaderBoard func to get leaderboard
func GetLeaderBoard() []ReturnSub {
	//一个可行的思路，先全部选出submission，然后手动选出每个用户的最后一次提交
	var AllSub []GetSub
	var ChooseSub []ReturnSub
	DB.Model(&Submission{}).Where("1=1").Order("score desc, created_at").Find(&AllSub)
	//check if user has been picked
	user_exist := make(map[string]bool)
	for _, sub := range AllSub {
		_, flag := user_exist[sub.UserName]
		if !flag {
			user_exist[sub.UserName] = true
			var subs []float64
			subs = append(subs, sub.Subscore1)
			subs = append(subs, sub.Subscore2)
			subs = append(subs, sub.Subscore3)
			_, UserTemp := GetUserByName(sub.UserName)
			temp := ReturnSub{
				UserName:  sub.UserName,
				Avatar:    sub.Avatar,
				CreatedAt: sub.CreatedAt,
				Score:     sub.Score,
				UserVotes: UserTemp.Votes,
				Subscores: subs,
			}
			ChooseSub = append(ChooseSub, temp)
		}
	}
	return ChooseSub
}

//GetScore func to judge user's content
func GetScore(content string) (score, subscore1, subscore2, subscore3 float64, err error) {
	//pre-process ground_truth and content
	ground_truth_file, _ := os.Open("./model/ground_truth.txt")
	temp, _ := ioutil.ReadAll(ground_truth_file)
	ground_truth := strings.Split(string(temp), "\n")[1:]
	ground_truth = ground_truth[:len(ground_truth)-1]
	content_lines := strings.Split(content, "\n")
	if content_lines[len(content_lines)-1] == "" {
		content_lines = content_lines[:len(content_lines)-1]
	}
	//check invalid content
	if len(content_lines) != 1000 {
		var err error
		err = errors.New("you should evaluate all the data")
		//fmt.Println(len(content_lines))
		return 0, 0, 0, 0, err
	}
	//get correct rate
	subscore := [3]float64{0.0, 0.0, 0.0}
	for index, content_line := range content_lines {
		ground_truth_line := ground_truth[index]
		content_line_buf := strings.Split(content_line, ",")
		ground_truth_line_buf := strings.Split(ground_truth_line, ",")[1:]
		//check invalid content
		if len(content_line_buf) != 3 {
			var err error
			err = errors.New("you should evaluate all the data")
			return 0, 0, 0, 0, err
		}
		for i := 0; i < 3; i++ {
			content_flag, _ := strconv.ParseBool(content_line_buf[i])
			ground_truth_flag, _ := strconv.ParseBool(ground_truth_line_buf[i])
			if content_flag == ground_truth_flag {
				subscore[i] += 1.0 / 1000.0
			}
		}
	}
	return calc_score(subscore[:]), subscore[0], subscore[1], subscore[2], nil
}

//calc_score func to calc main_score
func calc_score(subscore []float64) float64 {
	var mean_score float64
	mean_score = 0.0
	for _, subscore := range subscore {
		mean_score += subscore
	}
	mean_score /= 3
	return 55*interpolate(0.5, 0.8, 0, 1, mean_score) + 15*interpolate(0.5, 0.7, 0, 1, subscore[0]) + 15*interpolate(0.5, 0.9, 0, 1, subscore[1]) + 15*interpolate(0.5, 0.75, 0, 1, subscore[2])
}

//interpolate func to evaluate subscores
func interpolate(x1, x2, y1, y2, x float64) float64 {
	if x < x1 {
		return y1
	}
	if x > x2 {
		return y2
	}
	return math.Sqrt((x-x1)/(x2-x1))*(y2-y1) + y1
}
