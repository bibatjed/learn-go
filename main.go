package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"net/http"
)

type User struct {
	Username string `json:"username" validate:"min=6,max=10"`
	Age      uint8  `json:"age" validate:"gte=1,lte=20"`
	Sex      string `json:"sex" validate:"oneof=female male"`
}

var trans ut.Translator

func transInit(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //chinese
		enT := en.New() //english
		uni := ut.New(enT, zhT, enT)

		var o bool
		trans, o = uni.GetTranslator(local)
		if !o {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		//register translate
		//Register translator
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

func main() {
	if err := transInit("en"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/users", func(c *gin.Context) {
		validate := validator.New()
		var requestBody User
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "error"})
			return
		}

		if err := validate.Struct(requestBody); err != nil {
			errs, _ := err.(validator.ValidationErrors)
			fmt.Println(errs)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errs.Translate(trans),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"email": requestBody.Username})
	})

	r.Run()
}
