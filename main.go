package main

import(
	"github.com/JinHyeokOh01/go-crwl-server/crwl"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/cse", crwl.GetCSE)
	r.GET("/sw", crwl.GetSW)
	r.Run(":5000")
}