package leads

import (
	"github.com/gin-gonic/gin"
)

// ApplyRoutes leads
func ApplyRoutes(r *gin.RouterGroup) {
	leads := r.Group("/leads")
	{
		leads.GET("/page/:page", getLeads)
		leads.GET("/detail/:id", getLead)
	}
}
