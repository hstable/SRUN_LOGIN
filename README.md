# SRUN_LOGIN

深澜网络准入认证系统一键登录

使用方法：

环境变量

`SRUN_UNAME`: 用户名

`SRUN_PASSWD`: 密码

`SRUN_HOST`: 域名，默认 10.248.98.2

`SRUN_KEEP_ALIVE`: 设为 1 或 true 则每隔 180 秒发送一次请求以保持在线，可以与 nohup 一起使用

### 二进制文件

1. 从 [releases](https://github.com/hstable/SRUN_LOGIN/releases) 下载二进制文件，例如 [SRUN_LOGIN_linux_amd64](https://github.com/hstable/SRUN_LOGIN/releases/latest/download/SRUN_LOGIN_linux_amd64)。

2. 赋予可执行权限：

   ```shell script
   chmod +x ./SRUN_LOGIN
   ```

3. 使用示例：

   ```shell script
   ./SRUN_LOGIN yourusername yourpassword
   ```

   更多参数:
   ```shell script
   SRUN_KEEP_ALIVE=1 SRUN_HOST=10.248.98.2 ./SRUN_LOGIN yourusername yourpassword
   ```

### Docker

```shell script
docker run -d --restart=always \
           -e SRUN_UNAME=yourusername \
           -e SRUN_PASSWD=yourpassword \
           -e SRUN_KEEP_ALIVE=1 \
           --name srun_login \
           crazcell/srun_login
```

