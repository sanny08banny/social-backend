package routes

import (
	"social-backend/utils"

	"github.com/gin-gonic/gin"
)

func ImageRoutes(r *gin.Engine) {
	r.POST("/images/upload", utils.UploadImage)
	r.GET("/images/view/:imageName", utils.ViewImage)
	r.DELETE("/images/delete/:imageName", utils.DeleteImage)
}
