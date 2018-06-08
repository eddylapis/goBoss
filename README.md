# goBoss

  这是基于go语言编写的一款boss直聘机器人软件(牛人版)。附上[Python版本](https://github.com/wuranxu/Boss),
  无需配置Go环境, 我会提供windows和macos的可执行程序。
  
# 闪光点

  - 自动回复boss消息
  
    回复消息有3种类型。可自行修改, 传入关键字即可(忽略大小写如b站)。消息同一个人只会回复一次。

    - 大厂
    
    - 普通
    
    - 黑名单
      
  - 自动发送简历
  
    大厂专属, 先声明这里的大厂指的是心仪的公司, 而本人比较心仪这种公司, 所以改不了口了。

    当自动回复以后, 大厂的回复中包含"简历"的子字符串, 则会自动发送您的附件简历。
    
  - 自动刷新消息
  
    随时已读, 给人随时随地无时无刻不在的感觉。

# 效果图

  - 自动回复(这里我特意注册了招聘者的号)
  
  ![image.png](https://upload-images.jianshu.io/upload_images/6053915-a571a172db5f84b4.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

  ![image.png](https://upload-images.jianshu.io/upload_images/6053915-53b65f6096ece8ae.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

  ![image.png](https://upload-images.jianshu.io/upload_images/6053915-d4ee051d3a068c83.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
  
  map中key为boss名字, 可能会有重名情况。但是目前我只遍历前5条数据, 暂时还能用。value为发送消息/简历的状态, 如果key未找到说明没有回复过这个人, value为false代表简历未发送但是消息已发送, true代表消息和简历都已经发送。
  
  
# 快速开始

### 下载
- ```git clone https://github.com/wuranxu/goBoss.git```

- 下载zip文件并解压

### 修改json配置文件

百度API文字识别(每日500次免费)，进入[官网](http://ai.baidu.com/tech/ocr/general)申请并配置。

- app_id

- api_key

- secret_key

用户密码配置

- user(boss直聘手机号)

- password(boss直聘登录密码)

其他配置

下面是我本人的配置, 注意, star_reply字段里的第一个%s代表对方姓名, 第二个%s代表对方公司名。如果去掉的话会报错(设计如此, 后续可修改), 黑名单我就不放出来了哈。O(∩_∩)O~

```Javascript
{  

"star_company": [
    "百度", "阿里", "口碑", "天猫", "盒马", "UC", "淘宝", "蚂蚁", "支付宝", "今日头条", "字节跳动", "腾讯", "滴滴", "bili", "美团", 
    "点评", "饿了么", "京东", "喜马拉雅", "盛大", "拼多多", "链家", "58", "沪江", "bili", "哔哩", "二三四五", "2345", "猫眼", 
    "陆金所", "小红书", "七牛", "musical", "虎扑", "小度", "唯品会", "苏宁", "平安", "携程", "有赞", "哈罗", "运满满", "蔚来",
    "巨人", "游族", "易果", "爱奇艺", "美味不用等", "号店", "360", "拍拍贷", "b站", "网易"
  ],
  "star_reply": "%s您好, 十分荣幸能受到大厂: %s的亲睐, 这是程序自动下发的消息, 如果您需要我的简历, 请在回复中带上\"简历\"字样。项目地址:https://github.com/wuranxu/goBoss",
  "black_reply": "您好, 暂时没有兴趣, 抱歉~",
  "common_reply": "您好, 这是一条由直聘机器人自动发送的消息, 请等待我本人查看..."

}

```

- 下载chromedriver

  **chrome浏览器是使用本软件的前提, 我们需要chromedriver驱动**
  
  - 前往[淘宝镜像](https://npm.taobao.org/mirrors/chromedriver)

  - 选择与你浏览器版本对应的驱动
  
  如我的版本是66, 确实很6哈哈哈哈...
  
  ![image.png](https://upload-images.jianshu.io/upload_images/6053915-912f39ffc4192811.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

  ![image.png](https://upload-images.jianshu.io/upload_images/6053915-253b1f19e4f8a76c.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
  
  - 找到适配版本和操作系统的浏览器驱动
  
    ![image.png](https://upload-images.jianshu.io/upload_images/6053915-ea88308400da03a8.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

    可以看到, 2.38和2.39都满足我的需求, 现在我选2.39.
    
    ![image.png](https://upload-images.jianshu.io/upload_images/6053915-94dcd77d5ac95fcf.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
    
    mac/windows选择自己的驱动
    
  - 将下载好的驱动放入driver目录里
  

  - 运行

  之后就可以双击main.exe(windows)或者main挂起你的聊天机器人了。
  
  注意: **windows下要用管理员身份开启main.exe, 而且最好杀毒软件信任。**
  

# todolist
  还有很多不完善, 没做好的。之后填坑, 首当其冲就是解决用户需要手动安装浏览器驱动的问题。
  
- 发简历后邮件通知
  
- 低薪过滤
- 工作地点筛选
- chromedriver自动下载
- 对方连续发送表情时会接收不到新消息的bug(因为表情不是文本, 在web页面属于icon)
- 去除time.Sleep这种丑陋的等待元素方式


