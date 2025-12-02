# 前置依赖

#启动项目前，请确保本地已安装以下工具：
Go：版本 1.16+
Git：版本 2.30+
MySQL：版本 5.7+ 或 8.0+

# loginpage

前后端同仓的任务管理系统，后端基于Go后端开发，支持用户登录，注册，修改密码。

# 本地启动步骤

打开终端，执行以下命令克隆仓库：

# 仓库地址

git clone https://github.com/nanywy/LoginPage.git

# 进入项目根目录

cd loginpage

# 创建MySQL数据库（在MySQL内操作）

CREATE DATABASE login_information;

# 在项目根目录的 loginPage/ 文件夹下，新建 config.env 文件；

MYSQL_USER=root          #MySQL默认用户名（一般无需修改）

MYSQL_PASSWORD=你的MySQL密码 #安装MySQL时设置的密码

MYSQL_HOST=127.0.0.1     #本地数据库地址（默认无需修改）

MYSQL_PORT=3306          #MySQL默认端口（默认无需修改）

MYSQL_DB=login_information #数据库名（无需修改）

# 返回终端，安装Go依赖

go mod tidy

# 启动后端

go run loginPage/login.go

启动成功标识：终端输出 [GIN-debug] Listening and serving HTTP on :8080

# 验证后端

打开浏览器访问 http://localhost:8080/page/secondpage/secondpage/secondloginpage.html，页面显示 {"status":"ok"} 即正常
