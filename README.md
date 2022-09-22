# 租辆酷车项目源码

## 项目介绍
实战共享出行-汽车分时租赁小程序。你将掌握敏捷开发，
领域驱动开发的理念，使用typescript进行小程序前端开
发。后端采用go微服务架构，使用k8s+docker在云端进行
部署。

## 项目亮点
* Google架构师亲授地道的Google设计理念
* 敏捷式开发工程化实践
* 领域驱动设计前沿理念
* Go+Typescript稀缺双语式教学
* 贴合企业级工程化部署：k8s+docker上云

## 如何运行后端
1. 根据您的实际情况修改各子文件夹下的main.go->flag.String语句key所对应的value
1. 使用server/shared/mongo/setup.js脚本初始化数据库
1. `cd/server`
1. `go run 子文件夹/main.go`

## 如何编译以及运行小程序
1. `cd wx/miniprogram`
1. `npm install`
1. 打开小程序开发工具
1. 确保在详情->本地设置中勾选`启用自定义处理命令`
1. 若请求地址非https协议，请在详情->本地设置中勾选`不校验合法域名`
1. 点击工具->构建npm
1. 点击编译
1. 汽车的二维码在项目根目录的qrcode子目录下（这边只生成了10辆车，若需添加，请自行查找找工具生成）

## 如何部署
本项目采用k8s+docker在云端进行部署。
1. deployment文件夹下是各服务部署所需要的配置文件，请根据实际情况修改
1. 在deployment/config下创建各服务需要的ConfigMap和Secret
1. 参考k8s和docker文档进行部署