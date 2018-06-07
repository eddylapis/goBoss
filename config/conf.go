package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Config  UserConfig
	Environ Env
)

type UserConfig struct {
	Directory     string   `json:"directory"`
	User          string   `json:"user"`
	Password      string   `json:"password"`
	Receiver      string   `json:"receiver"`
	Sender        string   `json:"sender"`
	SenderPwd     string   `json:"sender_pwd"`
	AppID         string   `json:"app_id"`
	APIKey        string   `json:"api_key"`
	SecretKey     string   `json:"secret_key"`
	BlackList     []string `json:"black_list"`
	Retry         int      `json:"retry"`
	Delay         int      `json:"delay"`
	AutoResume    bool     `json:"auto_resume"`
	LoginURL      string   `json:"login_url"`
	Host          string   `json:"host"`
	LoginJSON     string   `json:"login_json"`
	MsgPage       string   `json:"msg_page"`
	JobJSON       string   `json:"job_json"`
	HisMsg        string   `json:"his_msg"`
	ResumeURL     string   `json:"resume_url"`
	WebTimeout    int      `json:"web_timeout"`
	BaiduTokenUrl string   `json:"baidu_token_url"`
	BaiduOcrUrl   string   `json:"baidu_ocr_url"`
	Headless      bool     `json:"headless"`
	StarCompany   []string `json:"star_company"`
	StarReply     string   `json:"star_reply"`
	BlackReply    string   `json:"black_reply"`
	CommonReply   string   `json:"common_reply"`
}

type Env struct {
	Root          string
	BaiduTokenUrl string
	BaiduCode     string
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func init() {
	Environ.Root, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	// Environ.Root = "/Users/wuranxu/go/src/goBoss"
	// Environ.Root, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	//Environ.Root = "C:/Users/Woody/go/src/goBoss"
	// 解析json
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/config/data.json", Environ.Root))
	if err != nil {
		log.Panicf("打开用户配置文件失败! Error: %s", err.Error())
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Panicf("解析用户配置文件data.json失败!Error: %s", err.Error())
	}

	Environ.BaiduTokenUrl = fmt.Sprintf(`%s&client_id=%s&client_secret=%s`,
		Config.BaiduTokenUrl, Config.APIKey, Config.SecretKey)
	Environ.BaiduCode = `https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic`

}
