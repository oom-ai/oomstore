# featctl

用于管理特征仓库的 CLI 工具。

## Usage

**初始化特征仓库**

```
$ featctl init
```

**注册特征实体**
```sh
featctl register entity device --length 32 --description "设备信息"
```

**注册特征组**
```sh
featctl register group device_baseinfo --entity device --description "设备基础信息"
```

**列举特征组**
```sh
featctl list group --entity=device
```

**注册特征**
```sh
featctl register batch-feature model --group device --value-type "varchar(30)" --description 'phone model'
```

**导入批特征数据**

下面展示如何将 `device.csv` 导入特征组 `device`。

首先，准备样例数据。

```sh
#!/usr/bin/env bash
set -euo pipefail

mkdir -p /tmp/featctl && cd /tmp/featctl

cat <<-EOF > device.csv
device,brand,model,price
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

**查看实体详情**
```
$ featctl describe entity device
Name:           device
Value length:   32
Description:    registered device
CreateTime:     2021-09-10T15:20:43Z
ModifyTime:     2021-09-13T18:58:34Z
```

**列举实体**
```
$ featctl list entity
Name,Length,Description,CreateTime,ModifyTime
device,32,"registered device",2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
user,48,"registered user",2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
```

**查看特征组详情**
```
$ featctl describe group device_info
Name:           device_info
Entity:         device
Description:    device basic info
Revision:       1634486400
DataTable:      batch_device_info_1634486400
CreateTime:     2021-09-10T15:20:43Z
ModifyTime:     2021-09-13T18:58:34Z
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

**修改实体配置**
```sh
featctl update entity device \
    --description "registered device"
```

**修改特征配置**
```sh
featctl update feature price\
    --description "phone price"
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
