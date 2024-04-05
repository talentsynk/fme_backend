package middleware

import (
	
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"

)



func ExtractMdaID(c *gin.Context){
	mdaID := extractMdaIDFromRequest(c.Request)
	c.Set("mdaID", mdaID)
	c.Next()
}

func extractMdaIDFromRequest(req *http.Request) uint{
	mdaIDParam := req.URL.Query().Get("mdaID")
	mdaID, err := strconv.ParseUint(mdaIDParam, 10 ,64)
	if err != nil {
		return 0 
	}

	return uint(mdaID)
}

