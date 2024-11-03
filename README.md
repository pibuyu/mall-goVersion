mall项目-mall-portal模块Go语言重构版。完全适配mall-app-web-master前端。

How to run?

1.前端：
首先下载mall-app-web-master代码：https://github.com/macrozheng/mall-app-web
仅需修改一处：
appConfig.js文件中：
  export const API_BASE_URL = 'http://localhost:8085';
修改为
  export const API_BASE_URL = 'http://localhost:9090';
然后在HbuilderX->运行->运行到浏览器->chrome启动前端项目即可。

2.后端：
将config.example.ini复制一份命名为config.ini，并放在/config文件夹下，填写mysql，redis以及mongodb配置信息即可。
需要用到的mysql服务、redis服务和mongodb服务的配置过程见：
  https://www.macrozheng.com/mall/start/mall_deploy_windows.html
然后在根目录执行 go run main.go即可启动后端项目。
