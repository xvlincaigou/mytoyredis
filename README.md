# mytoyreids

这只是一个玩具，并且大部分代码并不是原创的。
参见[https://www.build-redis-from-scratch.dev/en/introduction](https://www.build-redis-from-scratch.dev/en/introduction)

## 使用方法（我是用的是Ubuntu 20.04）

1. 安装Go语言环境：

   ```shell
   sudo apt install golang-go
   ```

2. 安装redis：

   ```shell
   sudo apt install redis-tools
   ```

3. 克隆代码库：

   ```shell
   git clone https://github.com/xvlincaigou/mytoyredis.git
   ```

4. 开启该服务：

   ```shell
   cd mytoyredis
   go run resp.go handler.go aof.go main.go
   ```

5. 开启redis-cli：

   ```shell
   redis-cli
   ```

6. 使用redis-cli与server交互：

   ```shell
   PING

   PONG
   
   SET xl mxy
   
   OK
   
   SET mxy xl 
   
   OK
   
   GET xl
   
   mxy
   
   GET mxy
   
   xl
   ```

祝你玩得愉快！