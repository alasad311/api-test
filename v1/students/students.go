package students

import (
	"net/http"
	"strconv"

	"api-test/structs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetStudent(db *gorm.DB) gin.HandlerFunc {
	var students []structs.RegStudents
	return func(c *gin.Context) {
		paramID, ParamErr := strconv.Atoi(c.Param("id"))
		if ParamErr != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"StatusCode": http.StatusBadRequest,
				"Method":     http.MethodPost,
				"Message":    "Record not found",
				"id":         c.Param("id"),
			})
			return
		}
		results := db.Where("std_id_str = ?", paramID).Find(&students)
		if results.RowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound,
				gin.H{
					"Message":    "Record not found",
					"Method":     http.MethodPost,
					"StatusCode": http.StatusBadRequest,
				})
			return
		} else {
			c.SecureJSON(http.StatusOK, gin.H{
				"Data":       students,
				"Message":    "success",
				"Method":     http.MethodGet,
				"StatusCode": http.StatusOK,
			})
			return
		}
	}
}
