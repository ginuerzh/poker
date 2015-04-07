package main

import (
	"labix.org/v2/mgo/bson"
)

const (
	MaxLevel = 40
)

func init() {
	initLevelScores()
}

type Props struct {
	Physical int64
	Literal  int64
	Mental   int64
	Wealth   int64
	Score    int64
}

type User struct {
	Id       string `bson:"_id"`
	Nickname string
	Profile  string
	Gender   string
	Chips    int64
	Props    Props
}

func (u *User) FindById(id string) error {
	return findOne("accounts", bson.M{"_id": id}, nil, u)
}

func (u *User) Level() int64 {
	return Score2Level(u.Props.Score)
}

var (
	levelScores = make([]int64, MaxLevel)
)

func scoreOfUpgrade(n int) int64 {
	difficult := func(n int) int {
		if n < 10 {
			return 0
		} else if n < 20 {
			return 1
		} else if n < 30 {
			return 3
		} else if n < 35 {
			return 6
		} else {
			return 5 * (n - 33)
		}
	}

	factor := func(n int) float64 {
		if n <= 10 {
			return 1
		} else if n < 30 {
			return (1.0 - float64(n-10)/100)
		} else {
			return 0.82
		}
	}

	s := int64(float64(2*n+difficult(n)) * float64(40+3*n) * factor(n))
	return s - s%10
}

func initLevelScores() {
	var total int64
	for i := 1; i < len(levelScores); i++ {
		total += scoreOfUpgrade(i)
		levelScores[i] = total
		//fmt.Println(i, total)
	}
}

func Score2Level(score int64) int64 {
	for i := 1; i < len(levelScores); i++ {
		if score < levelScores[i] {
			return int64(i)
		}
	}

	return int64(MaxLevel)
}
