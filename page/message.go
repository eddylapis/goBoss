package page

import (
	"fmt"
	cf "goBoss/config"
	"log"
	"time"

	"strings"

	"github.com/fedesog/webdriver"
)

type Message struct {
	Driver    *webdriver.ChromeDriver
	Session   *webdriver.Session
	MsgList   []map[string]string
	ReplyList map[string]bool
}

func (m *Message) Listen() {
	m.EnterMessage()
	m.Receive()

}

func (m *Message) Receive() {
	for {
		fmt.Printf("[%s]---发送消息列表: %+v\n", time.Now().Format("2006-01-02 15:04:05"), m.ReplyList)
		msgList, latest := m.GetMsgList()
		if len(msgList) > 0 {
			if msgList[0]["bossName"] == m.MsgList[0]["bossName"] && latest == m.MsgList[0]["latest"] {
				// 没有新boss消息
				m.ReFetch()
				continue
			} else {
				bossName, company := msgList[0]["bossName"], msgList[0]["company"]
				star := m.IsStar(company)
				if status, ok := m.ReplyList[bossName]; !ok {
					// 发送消息
					m.SendMsg(star, bossName, company)
					m.ReFetch()
					continue
				} else {
					if star == "star" {
						// 回复包含简历且未发送过简历
						if strings.Contains(latest, "简历") && !status {
							// 发送简历
							m.SendInfo(bossName, company)
							m.ReplyList[bossName] = true
						}
					}
					// 非大厂不自动发送简历
				}
			}
			fmt.Printf("您的最新职位为: %+v\n消息为: %s\n", msgList[0], latest)
		} else {
			m.ReFetch()
			continue
		}
		// break //TODO 暂时只读一次
	}
}

func (m *Message) SendMsg(companyType, bossName, company string) {
	var reply string
	dialog := GetElement("消息页面", "消息对话框")
	switch {
	case companyType == "star":
		reply = fmt.Sprintf(cf.Config.StarReply, bossName, company)
	case companyType == "black":
		reply = cf.Config.BlackReply
	default:
		reply = cf.Config.CommonReply
	}
	err := dialog.SendKeys(m.Session, reply)
	Assert(err)
	err = GetElement("消息页面", "发送按钮").Click(m.Session)
	if err != nil {
		fmt.Printf("自动回复失败!内容: %s, 接受者公司: %s, 接受者: %s\n Error: %s\n", reply, company, bossName, err.Error())
	}
	fmt.Printf("自动回复成功!内容: %s, 接受者公司: %s, 接受者: %s\n", reply, company, bossName)
	m.ReplyList[bossName] = false
}

func (m *Message) IsStar(company string) string {
	// 判断是否是大厂
	stars := cf.Config.StarCompany
	black_list := cf.Config.BlackList
	for _, star := range stars {
		if strings.Contains(strings.ToUpper(company), strings.ToUpper(star)) {
			return "star"
		}
	}
	for _, black := range black_list {
		if strings.Contains(strings.ToUpper(company), strings.ToUpper(black)) {
			return "black"
		}
	}
	return "common"
}

func (m *Message) SendInfo(bossName, company string) {
	err := GetElement("消息页面", "发送简历").Click(m.Session)
	if err != nil {
		fmt.Printf("遇到问题: 发送简历给公司: %s Boss: %s 出错!Error: %s\n", company, bossName, err.Error())
	}
	time.Sleep(1 * time.Second)
	err = GetElement("消息页面", "发送简历确认").Click(m.Session)
	Assert(err)
	fmt.Printf("发送简历给公司: %s Boss: %s 成功!", company, bossName)
}

func (m *Message) ReFetch() {
	//没有新消息或者没有消息
	fmt.Println("暂时没有新的消息")
	m.Session.Refresh()
	time.Sleep(time.Duration(cf.Config.Delay) * time.Second) // 延迟Delay秒刷新
}

func (m *Message) EnterMessage() {
	time.Sleep(5 * time.Second)
	err := GetElement("首页", "消息").Click(m.Session)
	Assert(err)
}

func (m *Message) GetMsgList() ([]map[string]string, string) {
	var lt string
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
		info := m.getInfo()

		eles, _ := GetElement("消息页面", "聊天内容").GetElements(m.Session)
		var latest string
		if len(eles) > 0 {
			latest, _ = eles[len(eles)-1].Text()
		}
		if i == 0 {
			lt = latest
		}
		info["latest"] = latest
		msgList = append(msgList, info)
		time.Sleep(1 * time.Second)
	}
	// for _, msg := range msgList {
	// 	fmt.Printf("%+v\n", msg)
	// }
	if len(m.MsgList) == 0 {
		m.MsgList = msgList
		return msgList, lt
	}
	// 回到第一个对话
	messageList[0].Click()
	return msgList, lt
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
