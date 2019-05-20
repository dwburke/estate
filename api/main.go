package api

import (
	"github.com/gin-gonic/gin"

	endpoint_key "github.com/dwburke/estate/api/key"
)

func SetupRoutes(r *gin.Engine) {
	endpoint_key.SetupRoutes(r)
}
