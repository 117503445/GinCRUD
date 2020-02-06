# GinCRUD

基于 Gin 和 Gorm 的 CRUD 示例项目

## 前言

使用 Gin + Gorm 实现最基础的 CRUD,部分参考

<https://blog.csdn.net/u010129985/article/details/81744357>

并修改了其中的一些问题

## 代码解析

先定义结构体

```go
//历史事件
type Story struct {
    // ID is story's id
    ID            int    `json:"id"`
    TimeStamp     int64  `json:"timeStamp"`
    Name          string `json:"name"`
    StoryDescribe string `json:"storyDescribe"`
}
```

数据库变量

```go
    var db *gorm.DB
```

main函数中初始化数据库

```go
    var err error
    db, err = gorm.Open("sqlite3", "./Wizz-Home-Page.Database")
    if err != nil {
        log.Fatal(err)
    }
    db.AutoMigrate(&Story{}) //数据库自动根据结构体建表
```

读取所有

```go
func ReadStories(c *gin.Context) {
    var stories []Story
    db.Find(&stories)
    c.JSON(200, stories)
}
```

读取单个

```go
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
```

创建

```go
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
```

修改

```go
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
```

删除

```go
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
```

添加路由

```go
    storyGroup := engine.Group("/api/stories")
    storyGroup.GET("", ReadStories)
    storyGroup.GET("/:id", ReadStory)
    storyGroup.POST("", CreateStory)
    storyGroup.PUT("/:id", UpdateStory)
    storyGroup.DELETE("/:id", DeleteStory)
```

运行服务端

```go
    err = engine.Run()
    if err != nil {
        log.Fatal(err)
    }
```
