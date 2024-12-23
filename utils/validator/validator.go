package validator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"gomall/utils/response"
	"regexp"
)

var ValidTrans ut.Translator
var ValidObj *validator.Validate

func init() {
	ValidObj = validator.New()
	english := en.New()
	chinese := zh.New()
	uni := ut.New(chinese, english)
	ValidTrans, _ = uni.GetTranslator("zh")
	_ = zhTranslations.RegisterDefaultTranslations(ValidObj, ValidTrans)
}

func CheckParams(ctx *gin.Context, err error) {
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range err.(validator.ValidationErrors) {
				msg, _ := ValidTrans.T(fieldError.Tag(), fieldError.Field(), fieldError.Param())
				response.Error(ctx, msg)
				return
			}
		} else {
			response.TypeError(ctx, err.Error())
			return
		}
	}
}

// VerifyMobileFormat 手机号验证
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func VerifyEmailFormat(email string) bool {
	regex := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	return regex.MatchString(email)
}

func CheckVideoSuffix(suffix string) error {
	switch suffix {
	case ".jpg", ".jpeg", ".png", ".ico", ".gif", ".wbmp", ".bmp", ".svg", ".webp", ".mp4":
		return nil
	default:
		return fmt.Errorf("非法后缀！")
	}
}
