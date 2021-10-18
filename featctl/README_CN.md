# featctl

用于管理特征仓库的 CLI 工具。

## Usage

**初始化特征仓库**

```
$ featctl init
```

**导入特征组**

下面展示如何将 `device.csv` 导入特征组 `device`。

首先，准备样例数据。

```sh
#!/usr/bin/env bash
set -euo pipefail

mkdir -p /tmp/featctl && cd /tmp/featctl

cat <<-EOF > device.csv
entity_key,brand,model,price
a9f0d6af575bb7e427fde2dcc81adbed,小米,MIX3,3999
134d9facd06ff355bf53846c0407d4f4,华为,P40,5299
0c66da7c680c4c44f33cb34881f1b104,苹果,IPHONE11,4999
EOF
```

随后，执行导入。

```sh
featctl import \
    --group device \
    --input-file device.csv \
    --separator "," \
    --description "test import"
```

**查看特征详情**
```
$ featctl describe feature price
Name:           price
Group:          device_info
Entity:         device
Category:       batch
ValueType:      int(11)
Description:    设备价格
Revision:       1634486400
DataTable:      batch_device_info_1634486400
CreateTime:     2021-09-10T15:20:43Z
ModifyTime:     2021-09-13T18:58:34Z
```

**创建特征配置**
```sh
featctl register feature \
    --name price \
    --group device \
    --category batch \
    --revision 20210909 \
    --revisions-limit 3 \
    --description "设备价格"
```

**注册特征实体**
```sh
featctl register entity device --length 32 --description "设备信息"
```

**修改特征配置**
```sh
# 启用特征并将其版本指定为 20210909
featctl update feature \
    --name price \
    --group device \
    --revision 20210909 \
    --status "enabled"
```

**查询特征值**

```
$ featctl query -g device -n brand,price -k a9f0d6af575bb7e427fde2dcc81adbed,134d9facd06ff355bf53846c0407d4f4
entity_key,brand,price
a9f0d6af575bb7e427fde2dcc81adbed,小米,3999
134d9facd06ff355bf53846c0407d4f4,华为,5299
```

**导出特征组**

```sh
featctl export --group device
```

**列举特征配置**

```sh
$ featctl list feature --group device
Name,Group,Revision,Status,Category,ValueType,Description,RevisionsLimit,CreateTime,ModifyTime
price,device,20210909,disabled,batch,int(11),设备价格,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
city,device,20210908,disabled,batch,int(11),城市,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z

$ featctl list feature
Name,Group,Revision,Status,Category,ValueType,Description,RevisionsLimit,CreateTime,ModifyTime
price,device,20210909,disabled,batch,int(11),设备价格,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
city,device,20210908,disabled,batch,int(11),城市,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
age,user,20210908,disabled,batch,int(11),年龄,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
```

**列举特征组的历史版本**
```sh
$ featctl list revision --group device
Group,Revision,Source,Description,CreateTime,ModifyTime
device,20210909,device_20210909,定时导入,2021-09-09T15:20:43Z,2021-09-09T15:20:43Z
device,20210908,device_20210908,手动触发,2021-09-08T15:20:43Z,2021-09-08T15:20:43Z
device,20210907,device_20210907,定时导入,2021-09-07T15:20:43Z,2021-09-07T15:20:43Z
```

## Config

featctl 默认读取 `$XDG_CONFIG_HOME/featctl/config.yaml` 作为配置文件（可通过 `--config` 手动指定）：

```yaml
host: 127.0.0.1
port: 4000
user: <user>
pass: <password>
```

如果没有提供配置文件，也可以在执行命令时手动指定以上参数。

## Development

Clone 项目之后先安装 [`pre-commit`](https://pre-commit.com/):

```sh
pip install pre-commit
pre-commit install
```

本地编译

```
make build
```

构建镜像

```
make image
```

推送镜像

```
make image-push
```

## TODO

- [x] `featctl init`
- [x] `featctl query`
- [x] `featctl export`
- [x] `featctl import`
- [x] `featctl list feature`
- [x] `featctl list revision`
- [x] `featctl register feature`
- [x] `featctl describe feature`
- [x] `featctl update   feature`
