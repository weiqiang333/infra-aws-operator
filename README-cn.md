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

- 更多版本功能信息请查看: [version info](./doc/version)

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
bash cmd/linux_build.sh v0.5
```

- kubernetes deploy
```
kubectl create namespace go
kubectl -n go create configmap infra-aws-operator-configmap --from-file=configs/config.yaml
kubectl apply -f build/go_default_service-deploy.yaml
```
- native deploy (systemd)
```text
version=v0.5
basedir=/usr/local/infra-aws-operator
mkdir -p ${basedir}

wget https://github.com/weiqiang333/infra-aws-operator/releases/download/${version}/infra-aws-operator-linux-amd64-${version}.tar.gz
tar -zxf infra-aws-operator-linux-amd64-${version}.tar.gz -C ${basedir}
chmod +x ${basedir}/infra-aws-operator
cp /usr/local/infra-aws-operator/configs/infra-aws-operator.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now infra-aws-operator
systemctl status infra-aws-operator
```

### demo
- initial page

![初始页面演示](./doc/img/init_demo.png)

- /metrics
```
# HELP aws_cost_billcost_service_by_day 按天计费服务
# TYPE aws_cost_billcost_service_by_day gauge
aws_cost_billcost_service_by_day{date_end="2023-03-16",date_start="2023-03-15",service_name="AWS Cost Explorer"} 0.02
aws_cost_billcost_service_by_day{date_end="2023-03-16",date_start="2023-03-15",service_name="Amazon Elastic Compute Cloud - Compute"} 1.5568
aws_cost_billcost_service_by_day{date_end="2023-03-16",date_start="2023-03-15",service_name="Amazon Lightsail"} 1.22976
...
# HELP aws_cost_billcost_service_by_month 按月计费服务
# TYPE aws_cost_billcost_service_by_month gauge
aws_cost_billcost_service_by_month{date_end="2022-12-01",date_start="2022-11-01",service_name="AWS Cost Explorer"} 0.28
aws_cost_billcost_service_by_month{date_end="2022-12-01",date_start="2022-11-01",service_name="AWS Support (Developer)"} 20.190525
aws_cost_billcost_service_by_month{date_end="2022-12-01",date_start="2022-11-01",service_name="Amazon Elastic Compute Cloud - Compute"} 42.0774433528
aws_cost_billcost_service_by_month{date_end="2022-12-01",date_start="2022-11-01",service_name="Amazon Lightsail"} 14.4703713538
aws_cost_billcost_service_by_month{date_end="2022-12-01",date_start="2022-11-01",service_name="EC2 - Other"} 83.2234017271
...
# HELP aws_lightsail_gb_month_network_in 光帆实例当月已使用传入流量 GB
# TYPE aws_lightsail_gb_month_network_in gauge
aws_lightsail_gb_month_network_in{name="aus-sa-1",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 11.258531582541764
aws_lightsail_gb_month_network_in{name="aus-sa-2",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 17.0006951848045
aws_lightsail_gb_month_network_in{name="aus-sa-3",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 46.21729747578502
...
# HELP aws_lightsail_gb_month_network_out 光帆实例当月已使用传出流量 GB
# TYPE aws_lightsail_gb_month_network_out gauge
aws_lightsail_gb_month_network_out{name="aus-sa-1",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 11.034040275029838
aws_lightsail_gb_month_network_out{name="aus-sa-2",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 16.252381362952292
aws_lightsail_gb_month_network_out{name="aus-sa-3",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 44.11481489613652
...
# HELP aws_lightsail_gb_per_month_allocated_transfer 光帆实例流量包免费流量 GB
# TYPE aws_lightsail_gb_per_month_allocated_transfer gauge
aws_lightsail_gb_per_month_allocated_transfer{name="aus-sa-1",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 1024
aws_lightsail_gb_per_month_allocated_transfer{name="aus-sa-2",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 1024
aws_lightsail_gb_per_month_allocated_transfer{name="aus-sa-3",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 1024
...
# HELP aws_lightsail_gb_remain_month_network 光帆实例当月流量套餐剩余 sum GB
# TYPE aws_lightsail_gb_remain_month_network gauge
aws_lightsail_gb_remain_month_network{name="aus-sa-1",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 1001.7074281424284
aws_lightsail_gb_remain_month_network{name="aus-sa-2",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 990.7469234522432
aws_lightsail_gb_remain_month_network{name="aus-sa-3",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 933.6678876280785
...
# HELP aws_lightsail_gb_use_month_network 光帆实例当月流量使用 sum GB
# TYPE aws_lightsail_gb_use_month_network gauge
aws_lightsail_gb_use_month_network{name="aus-sa-1",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 22.292571857571602
aws_lightsail_gb_use_month_network{name="aus-sa-2",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 33.25307654775679
aws_lightsail_gb_use_month_network{name="aus-sa-3",public_ip_address="3.3.3.3",regions="ap-southeast-2"} 90.33211237192154
...
# HELP aws_lightsail_instances 光帆实例信息状态
# TYPE aws_lightsail_instances gauge
aws_lightsail_instances{blueprint_name="Debian",cpu_count="1",created_at="2022-11-29 06:29:03.553 +0000 UTC",gb_per_month_allocated_transfer="2048.000000",name="sgp-2",private_ip_address="3.3.3.3",public_ip_address="3.3.3.3",ram_size_in_gb="1.000000",regions="ap-southeast-1"} 1
aws_lightsail_instances{blueprint_name="Debian",cpu_count="1",created_at="2022-11-29 06:37:45.32 +0000 UTC",gb_per_month_allocated_transfer="2048.000000",name="kr-s-1",private_ip_address="3.3.3.3",public_ip_address="3.3.3.3",ram_size_in_gb="1.000000",regions="ap-northeast-2"} 1
aws_lightsail_instances{blueprint_name="Debian",cpu_count="1",created_at="2022-11-29 06:41:04.499 +0000 UTC",gb_per_month_allocated_transfer="2048.000000",name="jp-tk-31",private_ip_address="3.3.3.3",public_ip_address="3.3.3.3",ram_size_in_gb="1.000000",regions="ap-northeast-1"} 1
...
```