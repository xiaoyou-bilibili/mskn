package main

import (
	"github.com/gin-gonic/gin"
	"mskn-server/web"
)

func main() {
	r := gin.Default()
	web.RegisterRouter(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
