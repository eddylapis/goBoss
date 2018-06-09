// 下面是编译的源代码, 因为mac下无法编译通过registry库
package driver  // 去掉这行

//package main
//
//import (
//	"fmt"
//	"log"
//	"github.com/golang/sys/windows/registry"
//	"goBoss/config"
//	"strings"
//)
//
//func main() {
//	k, err := registry.OpenKey(registry.CURRENT_USER, config.ChromeReg, registry.ALL_ACCESS)
//	if err != nil {
//		log.Fatal("获取Windows Chrome版本失败!请检查Chrome是否安装 Error: ", err)
//	}
//	s, _, err := k.GetStringValue("version")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer k.Close()
//	verList := strings.Split(s, ".")
//	ver := verList[0]
//	fmt.Println(ver)
//}