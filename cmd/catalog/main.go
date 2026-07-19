package main

import "github.com/gin-gonic/gin"
import "log"

func main() {

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	if err := r.Run(":9001"); err != nil {
		log.Fatal(err)
	}
}
