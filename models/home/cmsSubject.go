package home

import "time"

type CmsSubject struct {
	Id              int64     `json:"id" gorm:"id"`
	CategoryId      int64     `json:"category_id" gorm:"category_id"`
	Title           string    `json:"title" gorm:"title"`
	Pic             string    `json:"pic" gorm:"pic"`
	ProductCount    int       `json:"product_count" gorm:"product_count"`
	RecommendStatus int       `json:"recommend_status" gorm:"recommend_status"`
	CreatTime       time.Time `json:"create_time" gorm:"create_time"`
	CollectCount    int       `json:"collect_count" gorm:"collect_count"`
	ReadCount       int       `json:"read_count" gorm:"read_count"`
	CommentCount    int       `json:"comment_count" gorm:"comment_count"`
	AlbumPics       string    `json:"album_pics" gorm:"album_pics"`
	Description     string    `json:"description" gorm:"description"`
	ShowStatus      int       `json:"show_status" gorm:"show_status"`
	Content         string    `json:"content" gorm:"content"`
	ForwardCount    int       `json:"forward_count" gorm:"forward_count"`
	CategoryName    string    `json:"category_name" gorm:"category_name"`
}

type CmsSubjectList []CmsSubject

func (CmsSubject) TableName() string {
	return "cms_subject"
}
