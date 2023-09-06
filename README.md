# alinLab-webhook
一个webhook在线测试的网站.集成在个人网站里面  https://www.weridolin.cn/#/lab/index 里面


# 功能
提供一个独立webhook的地址,供项目中需要用到webhook时可以测试.前端实时展示webhook的调用结果.
![逻辑](./docs/%E9%80%BB%E8%BE%91%E5%9B%BE.png)

# 技术栈
- 数据库 mysql
- 后端框架 go-zero 方便后期统一部署
- orm: gorm
- 通讯: http socketio

# 注
考虑到ws/socket io 不同节点下的转发的问题，引入rabbitmq作为转发的中间件