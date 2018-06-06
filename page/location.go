package page

import (
	"encoding/json"
	"errors"
	"fmt"
	"goBoss/config"
	"io/ioutil"
	"log"
	"os"
	"time"

	dr "github.com/fedesog/webdriver"
)

type Element struct {
	Method dr.FindElementStrategy `json:"method,omitempty"`
	Value  string                 `json:"value,omitempty"`
}

type Method interface {
	GetEle(s *dr.Session) (dr.WebElement, error)
	GetElements(s *dr.Session) ([]dr.WebElement, error)
	SendKeys(s *dr.Session, text string) error
	Click(s *dr.Session) error
	Attr(s *dr.Session, attribute string) (string, error)
	Text(s *dr.Session) (string, error)
	WaitAndGet(now time.Time, s *dr.Session) (dr.WebElement, error)
}

var Page = map[string]map[string]Element{}

func GetElement(root, name string) *Element {
	ele, ok := Page[root][name]
	if !ok {
		log.Panicf("page/element.json未找到root: [%s] key: [%s]", root, name)
	}
	return &ele
}

func (e *Element) GetEle(s *dr.Session) (dr.WebElement, error) {
	ele, err := s.FindElement(e.Method, e.Value)
	return ele, err
}

func (e *Element) GetElements(s *dr.Session) ([]dr.WebElement, error) {
	ele, err := s.FindElements(e.Method, e.Value)
	return ele, err
}

func (e *Element) Click(s *dr.Session) error {
	ele, err := e.GetEle(s)

	// ele, err := e.WaitAndGet(time.Now(), s)
	if err != nil {
		return err
	}
	return ele.Click()
}

func (e *Element) SendKeys(s *dr.Session, text string) error {
	ele, err := e.GetEle(s)
	if err != nil {
		return err
	}
	return ele.SendKeys(text)
}

func (e *Element) Attr(s *dr.Session, attribute string) (string, error) {
	ele, err := e.GetEle(s)
	if err != nil {
		return "", err
	}
	return ele.GetAttribute(attribute)
}

func (e *Element) Text(s *dr.Session) (string, error) {
	ele, err := e.GetEle(s)
	if err != nil {
		return "", err
	}
	return ele.Text()
}

func (e *Element) WaitAndGet(now time.Time, s *dr.Session) (dr.WebElement, error) {
	tm := config.Config.WebTimeout
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("开始时间: %v 现在时间: %v元素暂时还未找到：Error: %v\n", now, time.Now(), err)
			e.WaitAndGet(now, s)
		}
	}()
	for {
		ele, err := e.GetEle(s)
		if err == nil {
			return ele, nil
		}
		if ses := time.Now().Sub(now).Seconds(); ses > float64(tm) {
			fmt.Printf("耗时: %v", ses)
			return dr.WebElement{}, errors.New(fmt.Sprintf("等待元素[%+v]超时: ", e))
		}
		time.Sleep(1 * time.Second) // 等待1秒后继续寻找
	}

}

func TearDown(w *Login) {
	// 截图screen
	// defer w.Close()
	pic, _ := w.Session.Screenshot()
	filename := fmt.Sprintf("%s_error.png", time.Now().Format("2006_01_02_15_04_05"))
	f, err := os.OpenFile(fmt.Sprintf("%s/picture/%s", config.Environ.Root, filename), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Printf("截图失败!Error: %s", err.Error())
	}
	f.Write(pic)
	// defer f.Close()
}

func init() {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/page/element.json", config.Environ.Root))
	if err != nil {
		log.Panicf("打开定位配置文件失败! Error: %s", err.Error())
	}
	err = json.Unmarshal(data, &Page)
	if err != nil {
		log.Panicf("解析用户配置文件data.json失败!Error: %s", err.Error())
	}
}
