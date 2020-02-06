package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

var db *gorm.DB

//历史事件
type Story struct {
	// ID is story's id
	ID            int    `json:"id"`
	TimeStamp     int64  `json:"timeStamp"`
	Name          string `json:"name"`
	StoryDescribe string `json:"storyDescribe"`
}

func ReadStories(c *gin.Context) {
	var stories []Story
	db.Find(&stories)
	c.JSON(200, stories)
}

func ReadStory(c *gin.Context) {
	id := c.Params.ByName("id")
	var story Story
	db.First(&story, id)
	if story.ID == 0 {
		c.JSON(404, gin.H{"message": "Story not found"})
		return
	}
	c.JSON(200, story)
}

func CreateStory(c *gin.Context) {
	var story Story

	if err := c.BindJSON(&story); err != nil {
		log.Println(err)
		c.JSON(400, "Not a Story")
		return
	}

	if story.ID != 0 {
		c.JSON(400, gin.H{"message": "Pass id in body is not allowed"})
		return
	}
	db.Create(&story)
	c.JSON(200, story)
}

func UpdateStory(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "your id is not a number"})
		return
	}
	var story Story
	db.First(&story, id)
	if story.ID == 0 {
		c.JSON(404, gin.H{"message": "Story not found"})
		return
	}
	err = c.ShouldBindJSON(&story)
	if err != nil {
		log.Println(err)
	}
	if story.ID != id {
		c.JSON(400, gin.H{"message": "Pass id in body is not allowed"})
		return
	}
	db.Save(&story)
	c.JSON(200, story)

}

func DeleteStory(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "your id is not a number"})
		return
	}
	var story Story
	db.First(&story, id)
	if story.ID == 0 {
		c.JSON(404, gin.H{"message": "Story not found"})
		return
	}
	db.Delete(&story)
	c.JSON(200, gin.H{"message": "delete success"})
}

func main() {
	engine := gin.Default()

	var err error
	db, err = gorm.Open("sqlite3", "./Wizz-Home-Page.Database")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Story{}) //数据库自动根据结构体建表

	storyGroup := engine.Group("/api/stories")
	storyGroup.GET("", ReadStories)
	storyGroup.GET("/:id", ReadStory)

	storyGroup.POST("", CreateStory)
	storyGroup.PUT("/:id", UpdateStory)
	storyGroup.DELETE("/:id", DeleteStory)

	err = engine.Run()
	if err != nil {
		log.Fatal(err)
	}
}
