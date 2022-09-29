package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	"time"
)

func login(c *gin.Context) {
	name := c.DefaultQuery("name", "jack")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

func submit(c *gin.Context) {
	name := c.DefaultQuery("name", "lily")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

func noRoute(c *gin.Context) {
	c.String(http.StatusOK, "not found route")
}

type Login struct {
	User   string `json:"user" form:"username" uri:"user" xml:"user" binding:"required"`
	Passwd string `json:"password" form:"password" uri:"password" xml:"password" binding:"required"`
}

type Login1 struct {
	User    string `uri:"user" validate:"checkName"`
	Pssword string `uri:"password"`
}

// 自定义验证函数
func checkName(fl validator.FieldLevel) bool {
	if fl.Field().String() != "root" {
		return false
	}
	return true
}

func Register() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello gin")
		//c.JSON(http.StatusOK, "hello gin")
	})

	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		//截取/
		action = strings.Trim(action, "/")

		//url参数获取
		a := c.DefaultQuery("a", "a1")
		b := c.DefaultQuery("b", "b1")
		d := c.Query("d")
		c.String(http.StatusOK, name+" is "+action+" "+a+" "+b)
		c.String(http.StatusOK, d)
	})

	r.POST("/form", func(c *gin.Context) {
		user := c.PostForm("user")
		age := c.PostForm("age")
		types := c.DefaultPostForm("type", "post_action")
		c.String(http.StatusOK, fmt.Sprintf("user:%s,age:%s,types:%s", user, age, types))
	})

	//限制最大上传8MB,8 * 2^10 * 2^10
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(c *gin.Context) {
		//file, err := c.FormFile("file")
		_, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.String(500, "上传文件出错")
		}
		if fileHeader.Size > 1024*1024*4 {
			fmt.Println("图片太大了")
			return
		}

		//保存上传文件到本地
		//c.SaveUploadedFile(file, file.Filename)
		c.String(http.StatusOK, fileHeader.Filename)
	})

	v1Group := r.Group("/v1")
	{
		v1Group.GET("/login", login)
		v1Group.GET("/submit", submit)
	}

	v2Group := r.Group("v2")
	{
		v2Group.GET("/login", login)
		v2Group.GET("/submit", submit)
	}

	r.NoRoute(noRoute)

	//json绑定
	r.POST("loginJson", func(c *gin.Context) {
		var jsonParam Login
		if err := c.ShouldBind(&jsonParam); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if jsonParam.User != "root" || jsonParam.Passwd != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "304"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 200})
	})

	r.GET("loginUrl", func(c *gin.Context) {
		var jsonParam Login
		if err := c.ShouldBind(&jsonParam); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if jsonParam.User != "root" || jsonParam.Passwd != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "304"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 200})
	})

	r.GET("/someXml", func(c *gin.Context) {
		c.XML(200, gin.H{"message": "xml"})
	})

	r.GET("/someYaml", func(c *gin.Context) {
		c.YAML(200, gin.H{"message": "yaml"})
	})

	r.GET("/someProtobuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "label"
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		c.ProtoBuf(http.StatusOK, data)
	})

	r.GET("/render", func(c *gin.Context) {
		//r.LoadHTMLGlob("render.*")
		r.LoadHTMLFiles("render.html")
		c.HTML(http.StatusOK, "render.html", gin.H{
			"title": "模板渲染",
			"ce":    "11111",
		})
	})

	r.GET("/baidu", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	})

	//r.Use(MiddleWare)
	//{
	r.GET("/middleWare", ceMiddle, func(c *gin.Context) {
		req, _ := c.Get("request")
		fmt.Println(req)
		c.JSON(http.StatusOK, gin.H{"request": req})
	})
	//}

	//r.Use(MiddleWare1())
	// {}为了代码规范
	//{
	r.GET("/ce", func(c *gin.Context) {
		// 取值
		req, _ := c.Get("request")
		fmt.Println("request:", req)
		// 页面接收
		c.JSON(200, gin.H{"request": req})
	})
	//}

	r.GET("/login", func(c *gin.Context) {
		c.SetCookie("abc", "123", 60, "/", "localhost", false, true)
		c.String(200, "Login Success!")
	})

	r.GET("/home", AuthMiddleWare, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "home"})
	})

	r.GET("/:user/:password", func(c *gin.Context) {
		var login1 Login1
		//注册自定义函数，与struct tag关联起来
		validate := validator.New()
		err := validate.RegisterValidation("checkName", checkName)
		if err := c.ShouldBindUri(&login1); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = validate.Struct(login1)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println(err)
			}
			return
		}
		fmt.Println("success")
	})

	r.Run(":8081")
}

func AuthMiddleWare(c *gin.Context) {
	if cookie, err := c.Cookie("abc"); err == nil {
		if cookie == "123" {
			c.Next()
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
	c.Abort()
	return
}

func MiddleWare(c *gin.Context) {
	t := time.Now()
	fmt.Println("中间件开始执行了")
	// 设置变量到Context的key中，可以通过Get()取
	c.Set("request", "中间件")

	c.Next()

	status := c.Writer.Status()
	fmt.Println("中间件执行完毕", status)
	t2 := time.Since(t)
	fmt.Println("time:", t2)
}

func ceMiddle(c *gin.Context) {
	fmt.Println("ce middle")
	//c.Negotiate(negotiate.AbortErrorWithStatus(c, errors.New("验签失败"), http.StatusUnauthorized))
	//c.Next()
}

func MiddleWare1() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Println("中间件1开始执行了")
		// 设置变量到Context的key中，可以通过Get()取
		c.Set("request", "中间件1")

		c.Next()

		status := c.Writer.Status()
		fmt.Println("中间件1执行完毕", status)
		t2 := time.Since(t)
		fmt.Println("time:", t2)
	}
}
