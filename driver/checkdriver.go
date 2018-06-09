package driver

import (
	"goBoss/config"
	"fmt"

	"log"
	"strings"
	"os/exec"
	"bytes"
	"goBoss/utils"
	"regexp"
	"strconv"
	"os"
	"archive/zip"
	"io"
)

func SetDriver() {
	ver := getChromeVer()
	getDriver(ver)
}

func getDriver(ver string) {
	drVer := getDriverVer(ver)
	setDriverName(drVer)
	status, err := checkDriver()
	if !status {
		if err != nil {
			log.Fatal("查找driver目录下的", config.Environ.DriverName, "失败!")
		} else {
			downloadDriver(drVer)
			return
		}
	}
	log.Println("驱动已存在, 无需重新下载...")
}

func setDriverName(drVer string) {
	switch config.Environ.Sys {
	case "windows":
		config.Environ.DriverName = fmt.Sprintf("chromedriver%s.exe", drVer)
		config.Environ.DriverZip = "chromedriver_win32.zip"
	case "darwin":
		config.Environ.DriverName = fmt.Sprintf("chromedriver%s", drVer)
		config.Environ.DriverZip = "chromedriver_mac64.zip"
	default:
		config.Environ.DriverName = fmt.Sprintf("chromedriver%s", drVer)
		config.Environ.DriverZip = "chromedriver_linux64.zip"
	}
}

func downloadDriver(s string) {
	log.Println("正在下载chromedriver驱动, 版本: ", s)
	zipfileName := fmt.Sprintf("%s/driver/%s", config.Environ.Root, config.Environ.DriverZip)
	req := utils.Request{
		Url: fmt.Sprintf("%s%s/%s", config.Config.DriverUrl, s, config.Environ.DriverZip),
		Method: "GET",
	}
	res := req.Http()
	if result := res["result"].([]byte); !res["status"].(bool) {
		log.Panicf("下载浏览器驱动版本失败, 请检查Url是否更换: %s!%s", req.Url, string(result))
	} else {
		// 下载浏览器驱动zip
		f, _ := os.OpenFile(zipfileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		f.Write(result)
		f.Close()
	}
	// 解压文件
	r, err := zip.OpenReader(zipfileName)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// 迭代压缩文件
	fmt.Println("正在解压", config.Environ.DriverZip)
	for _, fl := range r.File {
		if strings.Contains(fl.Name, "chromedriver") {
			rc, err := fl.Open()
			if err != nil {
				log.Fatal(err)
			}
			f, err := os.OpenFile(fmt.Sprintf("%s/driver/%s", config.Environ.Root, config.Environ.DriverName),
				os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			io.CopyN(f, rc, int64(fl.UncompressedSize64))
			rc.Close()
			f.Close()
			fmt.Println("浏览器驱动准备就绪...")
			return
		}
	}
}

func checkDriver() (bool, error) {
	_, err := os.Stat(fmt.Sprintf("%s/driver/%s", config.Environ.Root, config.Environ.DriverName))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getDriverVer(number string) string {
	var file_vr string
	req := utils.Request{
		Url:    fmt.Sprintf("%sLATEST_RELEASE", config.Config.DriverUrl),
		Method: "GET",
	}
	res := req.Http()
	var latest string
	if result := string(res["result"].([]byte)); !res["status"].(bool) {
		log.Panicf("获取浏览器驱动版本失败, 请检查Url是否更换: %s!%s", req.Url, result)
	} else {
		latest = result
	}
	req.Url = fmt.Sprintf("%s%s/notes.txt", config.Config.DriverUrl, latest)
	res = req.Http()
	info := res["result"].([]byte)
	reg, _ := regexp.Compile(`-+ChromeDriver\s+v(\d+\.+\d+)[\s|.|-|]+`)
	regSp, _ := regexp.Compile(`Supports\s+Chrome\s+v(\d+-\d+)`)
	dr := reg.FindAll(info, -1)
	sp := regSp.FindAll(info, -1)
	for i, s := range sp {
		vers := strings.Split(string(s), "-")
		small, bigger := vers[0], vers[1]
		small = delReg(small)
		vr := delReg(string(dr[i]))
		sm, _ := strconv.ParseInt(small, 10, 64)
		bg, _ := strconv.ParseInt(bigger, 10, 64)
		now, _ := strconv.ParseInt(string(number), 10, 64)
		if now >= sm && now <= bg {
			file_vr = vr
			log.Println("找到浏览器对应驱动版本号: ", vr)
			return file_vr
		}
	}
	return file_vr
}

func delReg(s string) string {
	smList := strings.Split(s, "v")
	str := smList[len(smList)-1]
	return strings.Replace(str, " ", "", -1)
}

func getChromeVer() string {
	var ver string
	switch config.Environ.Sys {
	case "windows":
		ver = getWinChromeVer()
	default:
		ver = getUnixChromeVer()
	}
	log.Println("成功获取到本机Chrome版本: ", ver)
	return ver
}

func getWinChromeVer() string {
	cmd := exec.Command(fmt.Sprintf("%s/getWinReg.exe", config.Environ.Root), "--version")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		log.Fatal("获取Windows Chrome版本失败!请检查Chrome是否安装 Error: ", err)
	}
	ver := string(cmdOutput.Bytes())
	return strings.Replace(ver, "\n", "", -1)
}

func getUnixChromeVer() string {
	cmd := exec.Command(config.ChromeApp, "--version")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		log.Fatal("获取Mac Chrome版本失败!请检查Chrome是否安装 Error: ", err)
	}
	ver := string(cmdOutput.Bytes())
	verList := strings.Split(ver, ".")
	ver = verList[0]
	return ver
}
