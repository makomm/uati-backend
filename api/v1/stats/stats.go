package stats

import (
	"github.com/gin-gonic/gin"
)

// ApplyRoutes stats
func ApplyRoutes(r *gin.RouterGroup) {
	stats := r.Group("/stats")
	{
		stats.GET("/", getStats)
	}
}
