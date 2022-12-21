FROM registry.xiaoyou.host/library/alpine:3.17.0
#WORKDIR /app
# 拷贝我们编译好的文件
COPY app /
# 声明暴露的端口
EXPOSE 8080
# 启动时直接运行app
ENTRYPOINT chmod +x /app && /app
