# oomctl

用于管理特征仓库的 CLI 工具。

## Usage

**初始化特征仓库**

```
$ oomctl init
```

**注册特征实体**
```sh
oomctl register entity device --length 32 --description "设备信息"
```

**注册特征组**
```sh
oomctl register group device_baseinfo --entity device --description "设备基础信息"
```

**列举特征组**
```sh
oomctl list group --entity=device
```

**注册特征**
```sh
oomctl register batch-feature model --group device --db-value-type "varchar(30)" --description 'phone model'
```

**导入批特征数据**

下面展示如何将 `device.csv` 导入特征组 `device`。

首先，准备样例数据。

```sh
#!/usr/bin/env bash
set -euo pipefail

mkdir -p /tmp/oomctl && cd /tmp/oomctl

cat <<-EOF > device.csv
device,brand,model,price
a9f0d6af575bb7e427fde2dcc81adbed,小米,MIX3,3999
134d9facd06ff355bf53846c0407d4f4,华为,P40,5299
0c66da7c680c4c44f33cb34881f1b104,苹果,IPHONE11,4999
EOF
```

随后，执行导入。

```sh
oomctl import \
    --group device \
    --input-file device.csv \
    --separator "," \
    --description "test import"
```

**查看实体详情**
```
$ oomctl describe entity device
Name:           device
Value length:   32
Description:    registered device
CreateTime:     2021-09-10T15:20:43Z
ModifyTime:     2021-09-13T18:58:34Z
```

**列举实体**
```
$ oomctl list entity
Name,Length,Description,CreateTime,ModifyTime
device,32,"registered device",2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
user,48,"registered user",2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
```

**查看特征组详情**
```
$ oomctl describe group device_info
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
$ oomctl describe feature price
Name:           price
Group:          device_info
Entity:         device
Category:       batch
DBValueType:    int(11)
ValueType:      int32
Description:    设备价格
Revision:       1634486400
DataTable:      batch_device_info_1634486400
CreateTime:     2021-09-10T15:20:43Z
ModifyTime:     2021-09-13T18:58:34Z
```

**创建特征配置**
```sh
oomctl register feature \
    --name price \
    --group device \
    --category batch \
    --revision 20210909 \
    --revisions-limit 3 \
    --description "设备价格"
```

**修改实体配置**
```sh
oomctl update entity device \
    --description "registered device"
```

**修改特征组配置**
```sh
oomctl update group device_info\
    --description "device basic info"
```

**修改特征配置**
```sh
oomctl update feature price\
    --description "phone price"
```

**查询特征值**

```
$ oomctl query -g device -n brand,price -k a9f0d6af575bb7e427fde2dcc81adbed,134d9facd06ff355bf53846c0407d4f4
entity_key,brand,price
a9f0d6af575bb7e427fde2dcc81adbed,小米,3999
134d9facd06ff355bf53846c0407d4f4,华为,5299
```

**导出特征组**

```sh
oomctl export --group device
```

**列举特征配置**

```sh
$ oomctl list feature --group device
Name,Group,Revision,Status,Category,DBValueType,ValueType,Description,RevisionsLimit,CreateTime,ModifyTime
price,device,20210909,disabled,batch,int(11),int32,设备价格,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
city,device,20210908,disabled,batch,int(11),int32,城市,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z

$ oomctl list feature
Name,Group,Revision,Status,Category,DBValueType,ValueType,Description,RevisionsLimit,CreateTime,ModifyTime
price,device,20210909,disabled,batch,int(11),int32,设备价格,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
city,device,20210908,disabled,batch,int(11),int32,城市,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
age,user,20210908,disabled,batch,int(11),int32,年龄,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
```

**列举特征组的历史版本**
```sh
$ oomctl list revision --group device
Group,Revision,Source,Description,CreateTime,ModifyTime
device,20210909,device_20210909,定时导入,2021-09-09T15:20:43Z,2021-09-09T15:20:43Z
device,20210908,device_20210908,手动触发,2021-09-08T15:20:43Z,2021-09-08T15:20:43Z
device,20210907,device_20210907,定时导入,2021-09-07T15:20:43Z,2021-09-07T15:20:43Z
```

## Config

oomctl 默认读取 `$XDG_CONFIG_HOME/oomctl/config.yaml` 作为配置文件（可通过环境变量`OOMCLI_CONFIG` 或参数 `--config` 指定）：

```yaml
online-store:
  backend: postgres
  postgres:
    host: 127.0.0.1
    port: 5432
    user: user
    password: password
    database: onlinestore

offline-store:
  backend: postgres
  postgres:
    host: 127.0.0.1
    port: 5432
    user: user
    password: password
    database: offlinestore

metadata-store:
  backend: postgres
  postgres:
    host: 127.0.0.1
    port: 5432
    user: user
    password: password
    database: metadatastore
```

## Development

Clone 项目之后先安装 [`pre-commit`](https://pre-commit.com/):

```sh
pip install pre-commit
pre-commit install
```

编译

```
make build
```

测试

```
make test
make integration-test
```
