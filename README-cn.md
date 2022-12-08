# AWS Operator
基础设施 AWS 运营商。 它管理和使用 AWS API 资源。 它提供简单的资源统计，同时保持资源调用的高性能。

[English](README.md)

## 目标和状态
- 该运算符的主要目标是：
```text
    1. 轻松的获取 aws 服务账单
    2. 对 devops 友好（监控、可用性、可扩展性和备份）
    3. 创建及管理 AWS Service 资源
```

## 架构设计
- default URL path
```text
/check
    health status check
/-/reload
    reload config file
/metrics
/
    default page
/api/v1/
    BasicAuth page
```

## 使用它
- build package
```
# 执行 go build, 并制作 images
bash cmd/linux_build.sh v0.4
```

- kubernetes deploy
```
kubectl create namespace go
kubectl -n go create configmap go-default-service-configmap --from-file=configs/config.yaml
kubectl apply -f build/go_default_service-deploy.yaml
```
- native deploy (systemd)
```text
version=v0.4
basedir=/usr/local/infra-aws-operator
mkdir -p ${basedir}

tar -zxf infra-aws-operator-linux-amd64-v0.4.tar.gz -C ${basedir}
chmod +x ${basedir}/infra-aws-operator
cp /usr/local/infra-aws-operator/configs/infra-aws-operator.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now infra-aws-operator
systemctl status infra-aws-operator
```

### demo
- initial page

![初始页面演示](./doc/img/init_demo.png)
