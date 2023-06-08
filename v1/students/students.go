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
	var apiKey []structs.ApiKey
	return func(c *gin.Context) {
		if err := c.BindJSON(&apiKey); err != nil {
			c.IndentedJSON(http.StatusNotFound,
				gin.H{
					"Message":    "Record not found",
					"Method":     http.MethodPost,
					"StatusCode": http.StatusBadRequest,
				})
			return
		}
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
		results := db.Raw("SELECT reg_students.std_id,reg_students.std_id_str,reg_students.national_id,reg_students.full_name,reg_student_nls.full_nls_name,reg_students.acc_avg,   reg_majors.major_name,reg_majors_nls.major_name AS nls_major_name,reg_major_plans.academic_level_desc,reg_major_plans.academic_level_desc_ar,reg_student_personal_info.phone,   reg_student_statuses.status_name,fin_fund_inst.fund_inst_name,reg_students.block   FROM reg_students   LEFT JOIN reg_student_nls ON reg_student_nls.std_id = reg_students.std_id   LEFT JOIN reg_majors ON reg_majors.major_id = reg_students.major_id   LEFT JOIN reg_majors_nls ON reg_majors_nls.major_id = reg_students.major_id   LEFT JOIN reg_major_plans ON reg_major_plans.plan_id = reg_students.plan_id   LEFT JOIN reg_student_personal_info ON reg_student_personal_info.std_id = reg_students.std_id   LEFT JOIN reg_student_statuses ON reg_student_statuses.status_id = reg_students.status_id   LEFT JOIN reg_student_semester_summary ON reg_students.last_semester_summary_id = reg_student_semester_summary.std_mark_cache_id    LEFT JOIN fin_student_funds ON reg_student_semester_summary.student_fund_id = fin_student_funds.fund_id   LEFT JOIN fin_fund_types ON fin_student_funds.fund_type_id = fin_fund_types.fund_type_id   LEFT JOIN fin_fund_inst ON fin_fund_types.fund_inst_id = fin_fund_inst.fund_inst_id    WHERE reg_students.std_id_str = ?", paramID).Scan(&students)
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
