package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"net/http"
)

// GenCaptcha 获取服务器生成的验证码
type GenCaptcha struct {
	ID   string `json:"id"`
	B64s string `json:"b64s"`
}

var store = new(RedisStore).NewRedisClient()

func GenerateCaptcha(c *gin.Context) {
	gc := new(GenCaptcha)
	var err error
	//bgColor := color.RGBA{3, 102, 214, 125}
	//bgColor := color.RGBA{R: 255, G: 255, B: 255, A: 1}
	bgColor := color.RGBA{R: 40, G: 52, B: 67, A: 1}
	resp := make(gin.H, 0)
	driver := base64Captcha.NewDriverMath(
		47,
		200,
		0,
		base64Captcha.OptionShowHollowLine,
		&bgColor,
		nil,
		[]string{"chromohv.ttf"},
	)

	// driver 生成captcha id及answer，存入store中，此处使用的是我们自定义的redis store
	captcha := base64Captcha.NewCaptcha(driver, store)
	// 生成验证码
	gc.ID, gc.B64s, err = captcha.Generate()
	if err != nil {
		resp = gin.H{
			"code": 1,
			"data": gc,
			"msg":  err.Error(),
		}
	} else {
		resp = gin.H{
			"code": 0,
			"data": gc,
			"msg":  "success",
		}
	}
	c.SecureJSON(http.StatusOK, resp)
}

// AnswerCaptcha 获取用户提交的验证码
type AnswerCaptcha struct {
	ID     string `json:"id"`
	Answer string `json:"answer"`
	Clear  bool   `json:"clear"`
}

// CaptchaVerify 验证用户提交的验证码
func (answer *AnswerCaptcha) CaptchaVerify() bool {
	// clear设置为false，否则会直接清理answer，用户由一次提交答案的机会
	result := store.Verify(answer.ID, answer.Answer, answer.Clear)
	return result
}
