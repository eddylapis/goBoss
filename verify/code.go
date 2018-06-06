package verify

import (
	"encoding/base64"
	"fmt"
	cf "goBoss/config"
	"goBoss/utils"
	"log"
	"net/url"
	"strings"

	sp "github.com/bitly/go-simplejson"
)

func GetCode(src string) string {
	token := GetBaiduToken()
	bs := GetPic(src)
	bs64 := Encode(base64.StdEncoding, bs)
	return SendPic(token, bs64)

}

func Encode(enc *base64.Encoding, bt []byte) string {
	// 编码
	encStr := enc.EncodeToString(bt)
	return encStr
}

func GetPic(src string) []byte {
	code_url := utils.Request{
		Url:    src,
		Method: "GET",
	}
	res := code_url.Http()
	bt_data := GetData(res)
	return bt_data
}

func SendPic(token, bs64 string) string {
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	values := url.Values{}
	values.Add("image", bs64)
	values.Add("language_type", "ENG")
	code_req := utils.Request{
		Url:     fmt.Sprintf("%s?access_token=%s", cf.Config.BaiduOcrUrl, token),
		Method:  "POST",
		Headers: header,
		Data:    values.Encode(),
	}
	res := code_req.Http()
	bt_data := GetData(res)
	data, _ := sp.NewJson(bt_data)
	word, e := data.GetPath("words_result").Array()
	if e != nil {
		log.Printf("调用百度API解析图片验证码失败, Error: %s", e.Error())
		return ""
	}
	if len(word) > 0 {
		wd, _ := word[0].(map[string]interface{})
		code := wd["words"]
		str := code.(string)
		str = strings.Replace(str, " ", "", -1)
		fmt.Printf("本次识别的验证码为: %s\n", str)
		return str
	} else {
		return ""
	}
}

func GetBaiduToken() string {
	token_request := utils.Request{
		Url:    cf.Environ.BaiduTokenUrl,
		Method: "GET",
	}
	res := token_request.Http()
	bt_data := GetData(res)
	data, _ := sp.NewJson(bt_data)
	token, err := data.Get("access_token").String()
	if err != nil {
		log.Panicf("获取百度api token出错， msg: %s", err.Error())
	}
	return token
}

func GetData(res utils.H) []byte {
	if !res["status"].(bool) {
		log.Panicf("Http请求出错, Msg: %s", fmt.Sprintf("%v", res["result"]))
	}
	return res["result"].([]byte)
}
