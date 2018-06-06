package page

import (
	"fmt"
	"log"
	"time"
	cf "goBoss/config"

	"github.com/fedesog/webdriver"
	"strings"
)

type Message struct {
	Driver  *webdriver.ChromeDriver
	Session *webdriver.Session
	MsgList []map[string]string
	Latest  string
}

func (m *Message) Listen() {
	m.EnterMessage()
	m.Receive()

}

func (m *Message) Receive() {
	for {
		msgList, latest := m.GetMsgList()
		if len(msgList) > 0 {
			if msgList[0]["bossName"] == m.MsgList[0]["bossName"] && latest == m.Latest {
				// 没有新boss消息
				m.ReFetch()
			}
			fmt.Printf("您的最新职位为: %+v\n消息为: %s\n", msgList[0], latest)
		} else {
			m.ReFetch()
		}
		break //TODO 暂时只读一次
	}
}

func (m *Message) ReFetch() {
	//没有新消息或者没有消息
	m.Session.Refresh()
	time.Sleep(time.Duration(cf.Config.Delay) * time.Second) // 延迟Delay秒刷新
}

func (m *Message) EnterMessage() {
	time.Sleep(3 * time.Second)
	err := GetElement("首页", "消息").Click(m.Session)
	Assert(err)
}

func (m *Message) GetMsgList() ([]map[string]string, string) {
	var latest string
	time.Sleep(3 * time.Second)
	// 获取消息列表
	messageList, e := GetElement("消息页面", "消息列表").GetElements(m.Session)
	Assert(e)
	if len(messageList) == 0 {
		// 消息列表为空， 持续检查
		log.Printf("消息列表为空, 请检查!")
		return []map[string]string{}, ""
	}
	msgList := make([]map[string]string, 0)
	for i, ms := range messageList[:10] {
		ms.Click()
		// 输出前10条最新消息
		time.Sleep(2 * time.Second)
		if i == 0 {
			eles, _ := ms.FindElements("css selector", ".text p")
			if len(eles) > 0 && m.Latest == "" {
				latest, _ = eles[len(eles)-1].Text()
				m.Latest = latest
			}
		}
		msgList = append(msgList, m.getInfo())
		time.Sleep(1 * time.Second)
	}
	for _, msg := range msgList {
		fmt.Printf("%+v\n", msg)
	}
	if len(m.MsgList) == 0 {
		m.MsgList = msgList
	}
	return msgList, latest
}

func (m *Message) getInfo() map[string]string {
	info := make(map[string]string)
	bossEle, err := GetElement("消息页面", "Boss信息").GetElements(m.Session)
	Assert(err)
	info["bossName"], _ = bossEle[0].Text()
	info["company"], _ = bossEle[1].Text()
	info["bossTitle"], _ = bossEle[2].Text()
	jobEle, err := GetElement("消息页面", "职位信息").GetElements(m.Session)
	Assert(err)
	info["position"], _ = jobEle[1].Text()
	info["money"], _ = jobEle[2].Text()
	info["base"], _ = jobEle[3].Text()
	for k, v := range info {
		info[k] = strings.Replace(v, " ", "", -1)
	}
	return info
}
