package home

// SmsHomeRecommendSubject 首页推荐专题表
type SmsHomeRecommendSubject struct {
	Id              int64  `json:"id" gorm:"id"`
	SubjectId       int64  `json:"subjectId" gorm:"subject_id"`
	SubjectName     string `json:"subjectName" gorm:"subject_name"`
	RecommendStatus int    `json:"recommendStatus" gorm:"recommend_status"`
	Sort            int    `json:"sort" gorm:"sort"`
}
type RecommendSubjectList []SmsHomeRecommendSubject

func (SmsHomeRecommendSubject) TableName() string {
	return "sms_home_recommend_subject"
}
