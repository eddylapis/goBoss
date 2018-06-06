package main

import (
	"fmt"
	cf "goBoss/config"
	"goBoss/page"
	"log"
	"os"

	"github.com/fedesog/webdriver"
)

func main() {
	setLog()
	chromeDriver := webdriver.NewChromeDriver("/Users/wuranxu/Downloads/chromedriver")
	lg := &page.Login{Driver: chromeDriver}

	lg.Start()
	lg.OpenBrowser()
	lg.Login()
	page.TearDown(lg)
	msg := &page.Message{Driver: chromeDriver, Session: lg.Session}
	msg.Listen()
	defer page.TearDown(lg)
}

func setLog() {
	//set logfile Stdout
	logFile, logErr := os.OpenFile(fmt.Sprintf("%s/boss.log", cf.Environ.Root), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", logFile, "cServer start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
