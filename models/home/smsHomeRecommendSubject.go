package home

// SmsHomeRecommendSubject 首页推荐专题表
type SmsHomeRecommendSubject struct {
	Id              int64  `json:"id" gorm:"id"`
	SubjectId       int64  `json:"subject_id" gorm:"subject_id"`
	SubjectName     string `json:"subject_name" gorm:"subject_name"`
	RecommendStatus int    `json:"recommend_status" gorm:"recommend_status"`
	Sort            int    `json:"sort" gorm:"sort"`
}
type RecommendSubjectList []SmsHomeRecommendSubject

func (SmsHomeRecommendSubject) TableName() string {
	return "sms_home_recommend_subject"
}
