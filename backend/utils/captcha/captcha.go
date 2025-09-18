package captcha

import (
	"errors"
	"fmt"
	"strings"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/mojocn/base64Captcha"
)

// 生成图片验证码
func NewCaptcha() (string, string, error) {
	driverString := base64Captcha.DriverString{
		Height:          50,
		Width:           120,
		NoiseCount:      0,
		ShowLineOptions: 0,
		Length:          4,
		Source:          "347ACDEFHJKMNPRTUVWXY",
		BgColor:         nil,
		Fonts: []string{
			"RitaSmith.ttf",
			"actionj.ttf",
			"chromohv.ttf",
		},
	}
	driver := driverString.ConvertFonts()
	id, content, answer := driver.GenerateIdQuestionAnswer()
	item, err := driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}
	key := fmt.Sprintf("captcha_%v", id)
	err = goi.Cache.Set(key, answer, 10*60)
	if err != nil {
		return "", "", err
	}
	base64 := item.EncodeB64string()
	return id, base64, nil
}

// 验证图片验证码
func VerifyCode(id string, code string) error {
	var answer string
	var err error

	if id == "" {
		return errors.New("验证码错误")
	}

	key := fmt.Sprintf("captcha_%v", id)
	if goi.Cache.Has(key) == false {
		return errors.New("验证码已过期")
	}
	err = goi.Cache.Get(key, &answer)
	if err != nil {
		return errors.New("验证码已过期")
	}
	if strings.EqualFold(answer, code) == false {
		return errors.New("验证码错误")
	}
	return nil
}
