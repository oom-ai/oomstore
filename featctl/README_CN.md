# featctl

用于管理特征仓库的 CLI 工具。

## Usage

**初始化特征仓库**

```
$ featctl init
```

**导入特征组**

下面展示如何将 `device.csv` 导入特征组 `device`。

首先，准备样例数据和 schema 模板。

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

cat <<-EOF > schema-template.sql
CREATE TABLE {{TABLE_NAME}} (
    entity_key VARCHAR(32) COMMENT '设备ID' PRIMARY KEY,
    brand      VARCHAR(16) COMMENT '设备厂商',
    model      VARCHAR(32) COMMENT '设备型号',
    price      INT         COMMENT '设备价格'
);
EOF
```

随后，执行导入。

```sh
featctl import \
    --group device \
    --revision 20210909 \
    --schema-template schema-template.sql \
    --input-file device.csv \
    --has-header \
    --separator "," \
    --description "test import"
```

**查看特征详情**
```
$ featctl describe --group device --name price
Name:           price
Group:          device
Revision:       20210909
Status:         disabled
Category:       batch
ValueType:      int(11)
Description:    设备价格
RevisionsLimit: 3
CreateTime:     2021-09-10T15:20:43Z
ModifyTime:     2021-09-13T18:58:34Z
```

**创建特征配置**
```sh
featctl create feature \
    --name price \
    --group device \
    --category batch \
    --revision 20210909 \
    --revisions-limit 3 \
    --description "设备价格"
```

**修改特征配置**
```sh
# 启用特征并将其版本指定为 20210909
featctl set \
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

将特征组 `device` 的全部特征下载到 `/tmp/featctl/device-exported.csv`。

```sh
featctl export --group device --output-file device-exported.csv
```

**列举特征配置**
```sh
$ featctl list feature --group device
Name,Group,Revision,Status,Category,ValueType,Description,RevisionsLimit,CreateTime,ModifyTime
price,device,20210909,disabled,batch,int(11),设备价格,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
city,device,20210908,disabled,batch,int(11),设备价格,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
age,device,20210908,disabled,batch,int(11),设备价格,3,2021-09-10T15:20:43Z,2021-09-13T18:58:34Z
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

## Dependency

对于 Linux 用户，下载并压缩 [tidb-toolkit-v4.0.13](https://download.pingcap.org/tidb-toolkit-v4.0.13-linux-amd64.tar.gz)，将其中的 `bin` 目录添加到环境变量 `PATH` 中：

```sh
export PATH="$PATH:/path/to/tidb-toolkit/bin"
```

对于 MacOS 用户，手动编译 dumpling 和 tidb-lightning 后，放在 ~/softwares/tidb-toolkit 下，随后将它们添加到环境变量 `PATH` 中，同上。注意，需要使用 `chmod +x` 赋予执行权限。

生产环境建议采用容器的方式运行，docker 镜像已经打包好了依赖：

```sh
docker run --rm aiinfra/featctl:latest sh -c 'featctl ...'
```

如果需要和宿主机交换文件，可以将宿主机目录挂载到容器的 `/work` 目录。

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
- [x] `featctl export`
- [x] `featctl import`
- [x] `featctl create feature`
- [x] `featctl set`
- [x] `featctl describe`
- [x] `featctl query`
- [x] `featctl list features`
- [ ] `featctl list revisions`
