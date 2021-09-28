
# 编写目的

FCM官方文档中描述比较零散，我自己在调研接入方式时也耗费了一些时间，所以我想编写一个文档，对接入FCM前后端的改造点做一个完整的描述，也算是对FCM调研过程做一个总结，希望能给有相关需求的人一些参考，文档会分为FCM服务申请，WEB端SDK接入，服务器端消息发送三个部分来描述

# 参考资料
[FCM官方文档](https://firebase.google.com/docs/cloud-messaging)

## 一、FCM服务申请

谷歌推送以前叫GCM，在2019年4月11日迁移到了firebase cloud message（FCM），我们要使用FCM，需要在firebase上创建一个项目，所以需要有一个gmail邮箱

### 1.1 添加项目

访问[https://console.firebase.google.com/](https://console.firebase.google.com/)，登陆后可以看到以下界面：

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/create.png)


然后点击Add project，创建一个firebase项目

### 1.2 获取应用凭证

创建号项目之后，点击项目→进入控制台 → 点击设置

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/setting.png)

然后在常规信息栏中，查看到web应用凭证，这里也提供了npm和cdn的集成方式示例

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/config.png)

应用凭证是用来给到前端集成SDK与FCM创建长链接，这里我们这里的应用凭证是：

``` javascript
{
    apiKey: "AIzaSyAUVTaipaExXUTGGc7e-A3gUiA3Q8i7O8Y",
    authDomain: "shipment-portal-1f1b3.firebaseapp.com",
    projectId: "shipment-portal-1f1b3",
    storageBucket: "shipment-portal-1f1b3.appspot.com",
    messagingSenderId: "333922912996",
    appId: "1:333922912996:web:18467f5642e6fba00efaf1",
    measurementId: "G-5HNY007WZW"
}
```

然后在回到云消息传递信息栏，查看到服务器密钥

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/serverkey.png)

服务器密钥用来向FCM推送消息，然后由FCM通过长链接推送给web客户端

服务器密钥用来可以添加多个，这里我们使用的密钥是：

```
AAAAjaT6Kls:APA91bFki15CULQgVFGKn1N3cR-PlKhcU5wMEClTl5NtP9YQylHSBsJctlZzQ0sdp6TnOsp5jyzzqW4tzSVDneKCmP8XQ8ZgA8w6wk6AW28aAdZn6eyz0U-OA7TuJ2WlkBbS10KcxQHe
```


## 二、前端项目集成firebase sdk

前端可以通过npm或cdn的方式将fcm集成进自己的web项目，这里我用cdn的方式做示例，npm的操作方式可以参考：https://firebase.google.com/docs/web/setup

### 2.1 创建index.html

在界面上我们添加了一个accountId的输入框，点击登录之后，会将accountId和当前web客户端的token发送给服务器，用来模拟用户登陆的场景，界面如下：

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/front.png)

在onload中我们初始化firebase，并通过requestPermission方法请求用户允许发送通知，然后注册了一个messaging.onMessage监听，当FCM通过长链接推送消息到web客户端时，并且web应用位于前台时，会执行onMessage方法，请求用户允许通知界面如下：

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/allow.png)

index.html的所有代码如下：
```
<html xmlns="http://www.w3.org/1999/html">
<body>
<script src="https://www.gstatic.com/firebasejs/8.10.0/firebase-app.js"></script>
<script src="https://www.gstatic.com/firebasejs/8.10.0/firebase-messaging.js"></script>
<h1>FCM</h1>
account: <input id="accountId"/>
<input type="button" value="login" onclick="login()"/>
<div id="info"></div>
<div id="msg" style="font-size:64px;color:red">0</div>
</body>
<script>
    var accessToken = "123123"

    // 定义配置文件
    var firebaseConfig = {
        apiKey: "AIzaSyBqBrEAaRBXhqtOaHaknfNV4_27Go2P3zE",
        authDomain: "weqwe-d71c5.firebaseapp.com",
        projectId: "weqwe-d71c5",
        storageBucket: "weqwe-d71c5.appspot.com",
        messagingSenderId: "608358247003",
        appId: "1:608358247003:web:433e10a50f5f51a5db2d7e",
        measurementId: "G-4QGGNLHSWW"
    };
    function login() {
        var accountId = document.getElementById("accountId").value
        // 注册token到后台
        document.getElementById("info").innerHTML = "token: " + accessToken
        resp = fetch("http://localhost:8090/accessToken", {
            method: "POST",
            mode:"cors",
            headers: {
                "Content-type": 'application/json;charset=utf-8'
            },
            body: JSON.stringify({"accountId": accountId, "token": accessToken}),
        })
        console.log(resp);
    }
    window.onload  = function () {
        // 生成token
        firebase.initializeApp(firebaseConfig);
        var messaging = firebase.messaging();
        messaging.requestPermission()
            .then(function () {
                return messaging.getToken();
            })
            .then(function (token) {
                accessToken = token
            })
            .catch(function (err) {
                console.log('Unable to get permission to notify.', err);
            });

        // 注册监听事件
        messaging.onMessage(function (payload) {
            // 更新列表消息数量
            console.log('Message received. ', payload);
            var msg = document.getElementById("msg").innerHTML
            document.getElementById("msg").innerHTML = Number(msg) + 1
        });
    }

</script>
</html>
```
### 2.2 创建firebase-messaging-sw.js

当用户界面处于后台时，是无法触发onMessage方法，firebasebase基于service work提供了onBackgroundMessage方法，我们在index.html同目录下创建一个firebase-messaging-sw.js，然后注册一个onBackgroundMessage方法

代码如下：
```
importScripts('https://www.gstatic.com/firebasejs/8.10.0/firebase-app.js');
importScripts('https://www.gstatic.com/firebasejs/8.10.0/firebase-messaging.js');

// Initialize the Firebase app in the service worker by passing in the
// messagingSenderId.
firebase.initializeApp({
    apiKey: "AIzaSyAUVTaipaExXUTGGc7e-A3gUiA3Q8i7O8Y",
    authDomain: "shipment-portal-1f1b3.firebaseapp.com",
    projectId: "shipment-portal-1f1b3",
    storageBucket: "shipment-portal-1f1b3.appspot.com",
    messagingSenderId: "333922912996",
    appId: "1:333922912996:web:18467f5642e6fba00efaf1",
    measurementId: "G-5HNY007WZW"
});

// Retrieve an instance of Firebase Messaging so that it can handle background
// messages.
const messaging = firebase.messaging();

messaging.onBackgroundMessage((payload) => {
    console.log('Received background message ', payload);
    // console.log(self)
    var title = payload.data.title
    var options = {
        body: payload.data.body,
        icon: payload.data.url
    }
    self.registration.showNotification(title,
        options);
});
```
我们在onBackgroundMessage中通过showNotification弹出一个弹框，当web应用处于后台时（用户未停留在web界面），有FCM传递过来的新消息，会弹出一个弹框，效果如下：

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/backnotice.png)

## 三、后端项目发送通知消息
### 3.1 交互流程图

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/fcmflow.png)

### 3.2 创建Admin项目

创建一个Admin项目模拟我们的后台，主要是提供了查询在线用户列表，登陆，发送消息三个基本功能，用来模拟和web端的交互，admin的界面如下：

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/admin.png)

* 查询在线用户列表：用户初始化SDK成功之后，我们需要将用户的accountId和当前的token保存起来；
* 登陆：在用户web界面，点击login之后会将accountId和token传递给后台进行保存；
* 发送消息：输入需要发送的accountId和title body，url之后，点击send，消息内容达到后台之后，会通过POST请求发给FCM；

## 四、总结
### 4.1 项目代码：

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/project.png)

### 4.2 后端启动

后端是基于gin构建，直接启动main函数即可启动

### 4.3 前端启动

静态资源的部署有很多种方式，一般我们可以利用nginx去部署，这里我们利用http-server这种更加轻量的方式启动一个http-server服务器，http-server是一个轻量级的基于nodejs的http服务器，可以使任意一个目录成为服务器的目录，完全抛开后台的沉重工程，直接运行想要的js代码。

#### 4.2.1 安装node

直接下载安装即可，官网地址： https://nodejs.org

安装成功之后在命令行输入命令$ node -v以及$ npm -v检查版本，确认是否安装成功。

#### 4.2.2 安装http-server

打开终端输入：
```npm install http-server -g```

安装成功之后进入目标文件夹，这里就是我们项目的 front/src

然后输入: ```http-server```

这样就会在本地8080端口启动一个http服务器（如果8080被占用则会在8081端口启动，以此类推）,启动界面如下：

![avatar](https://raw.githubusercontent.com/caixunshi/go-fcm-example/main/doc/img/http-server.png)

用户web界面为：[http://localhost:8080/index.html](http://localhost:8080/index.html)

后端管理界面为：[http://localhost:8080/admin.html](http://localhost:8080/admin.html)

这样我们就可以在两个浏览器分别测试消息发送啦

PS: 用户web界面谷歌浏览器service work不支持非https，所以不能测试用户离开web界面的弹框提醒场景，只能测试用户停留在web界面接受消息的场景，需要看效果的可以用firefox

### 4.3 思考

* 用户初始化SDK并与FCM建立长链接是没有传入业务参数的，生成出来的token是跟我们的业务数据（如用户）无关的，所以我们需要客户端主动上报token 和 username，并在我们自己的后台维护这个关系，但这也衍生出来几个问题，需要根据具体的业务场景去制定解决方案。

* 用户在多个端（比如多个浏览器）登陆时，一个用户会有多个token，如果想要所有端都通知的话，就需要保留所有的token，否则就只有最后登陆的端能接收到消，用户量很大的情况下，保存这个关系需要成本；
  用户下线之后，服务器的token和user的关联关系如何删除？（直接关掉web界面或直接kill 掉 app）；