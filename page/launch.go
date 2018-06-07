package page

import (
	"errors"
	"fmt"
	cf "goBoss/config"
	"goBoss/utils"
	"goBoss/verify"
	"log"
	"time"

	"github.com/fedesog/webdriver"
)

type Login struct {
	Driver  *webdriver.ChromeDriver
	Session *webdriver.Session
}

func Assert(err error) {
	if err != nil {
		log.Printf("Error: %s", err.Error())
		// panic("程序遇到问题啦, 请检查截图和日志...")
	}
}

func MaxWindow(w *Login) error {
	p := fmt.Sprintf(`{"windowHandle": "current", "sessionId": "%s"}`, w.Session.Id)
	req := utils.Request{
		Data:   p,
		Method: "POST",
		Url:    fmt.Sprintf("http://127.0.0.1:%d/session/%s/window/current/maximize", w.Driver.Port, w.Session.Id),
	}

	res := req.Http()
	if !res["status"].(bool) {
		log.Printf("response: %+v", res)
		return errors.New(fmt.Sprintf(`最大化窗口失败, 请检查!%+v`, res["msg"]))
	}
	return nil
}

func SetWindow(w *Login, width, height int) error {
	p := fmt.Sprintf(`{"windowHandle": "current", "sessionId": "%s", "height": %d, "width": %d}`, w.Session.Id, height, width)
	req := utils.Request{
		Data:   p,
		Method: "POST",
		Url:    fmt.Sprintf("http://127.0.0.1:%d/session/%s/window/current/size", w.Driver.Port, w.Session.Id),
	}

	res := req.Http()
	if !res["status"].(bool) {
		log.Printf("response: %+v", res)
		return errors.New(fmt.Sprintf(`设置浏览器窗口失败, 请检查!%+v`, res["msg"]))
	}
	return nil
}

func MaxWin(w *Login) error {
	args := make([]interface{}, 0)
	_, err := w.Session.ExecuteScript(`window.resizeTo( screen.availWidth, screen.availHeight );`, args)
	return err
}

func (w *Login) Start() {
	var err error
	w.Driver.Start()
	args := make([]string, 0)
	if cf.Config.Headless {
		args = append(args, "--headless")
	}
	desired := webdriver.Capabilities{
		"Platform":           "Mac",
		"goog:chromeOptions": map[string][]string{"args": args, "extensions": []string{}},
		"browserName":        "chrome",
		"version":            "",
		"platform":           "ANY",
	}
	required := webdriver.Capabilities{}
	w.Session, err = w.Driver.NewSession(desired, required)
	if err != nil {
		log.Printf("open browser failed: %s", err.Error())
	}

}

func (w *Login) OpenBrowser() {
	w.Session.Url(cf.Config.LoginURL)
	// err := MaxWindow(w)
	// err := MaxWin(w)
	err := SetWindow(w, 1600, 900)
	if err != nil {
		log.Panicf("最大化浏览器失败!!!Msg: %s", err.Error())
	}
	w.Session.SetTimeoutsImplicitWait(cf.Config.WebTimeout)
}

func (w *Login) sendCode() {
	// 识别验证码
	for {
		// image, err := w.Session.FindElement(lg["验证码"].Method, lg["验证码"].Value)
		image := GetElement("登录页面", "验证码")
		src, err := image.Attr(w.Session, "src")
		Assert(err)
		// src, _ := image.GetAttribute("src")
		code := verify.GetCode(src)
		if len(code) != 4 {
			// 验证码识别有误
			time.Sleep(3 * time.Second / 2)
			fmt.Println("验证码长度不为4, 重新获取!")
			image.Click(w.Session)
			continue
		} else {
			err = GetElement("登录页面", "验证码输入框").SendKeys(w.Session, code)
			Assert(err)

			err = GetElement("登录页面", "登录").Click(w.Session)
			Assert(err)
			time.Sleep(3 * time.Second / 2)
			text, _ := GetElement("登录页面", "验证码错误").Text(w.Session)
			// Assert(err)
			if text == "" {
				// 登录成功, break
				fmt.Println("恭喜您登录成功...")
				break
			} else {
				fmt.Println("验证码错误, 重新登录...")
				time.Sleep(3 * time.Second / 2)
				GetElement("登录页面", "验证码").Click(w.Session)
				continue
			}
		}

	}
}

func (w *Login) Login() {

	// 进入密码登录页面
	err := GetElement("登录页面", "密码登录").Click(w.Session)
	Assert(err)
	err = GetElement("登录页面", "用户名输入框").SendKeys(w.Session, cf.Config.User)
	Assert(err)
	err = GetElement("登录页面", "密码输入框").SendKeys(w.Session, cf.Config.Password)
	Assert(err)
	w.sendCode()

}

func (w *Login) Close() {
	w.Session.CloseCurrentWindow()
	w.Driver.Stop()

}
