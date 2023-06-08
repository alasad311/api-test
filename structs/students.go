package structs

type RegStudents struct {
	StdId               uint    `gorm:" PRIMARY_KEY" json:"std_id"`
	StdIdStr            uint    `json:"std_id_str"`
	NationalId          uint    `json:"national_id"`
	FullName            string  `json:"full_name"`
	FullNlsName         string  `json:"full_nls_name"`
	AccAvg              float32 `json:"acc_avg"`
	MajorName           string  `json:"accmajor_name_avg"`
	NlsMajorName        string  `json:"nls_major_name"`
	AcademicLevelDesc   string  `json:"academic_level_desc"`
	AcademicLevelDescAr string  `json:"academic_level_desc_ar"`
	Phone               uint    `json:"phone"`
	StatusName          string  `json:"status_name"`
	FundInstName        string  `json:"fund_inst_name"`
	Block               int     `json:"block"`
}
