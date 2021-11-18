# apply 命令文档

apply 的效果是：接受一个资源（entity，group，feature）文件，使用 apply 将数据库中的资源设置成资源文件声明的样子（通过创建和更新）。但对于没有声明的字段，不做修改。

下面是几种资源文件的例子：

```yaml
kind: Entity
name: user
length: 8
description: entity description
batch-features:
- group: device
  description: a description
  features:
  - name: model
    db-type-value: varchar(16)
    description: 'device model info'
  - name: price
    db-type-value: int
    description: 'device price'
- group: user
  description: a description
  features:
  - name: age
    db-type-value: int
    description: 'user age'
  - name: gender
    db-type-value: int
    description: 'user gender'
---
kind: Entity
name: user
length: 8
description: 'User ID'
---
kind: FeatureGroup
name: account
entity-name: user
category: batch
description: 'user account info'
---
kind: FeatureGroup
name: user-device
entity-name: device
category: batch
description: 'phone info'
features:
- name: model
  db-type-value: varchar(16)
  description: 'device model info'
- name: price
  db-type-value: int
  description: 'device price'
- name: radio
  db-type-value: int
  description: 'radio info'
---
kind: Feature
name: model
group-name: device
category: batch
description: 'device info'
---
```

问题：
* 每种资源文件需要对应到一个 struct，这个 struct 长什么样子？
* 如何将资源文件正确解析到正确的 struct？
* 拿到 struct 后，如何和数据库中的数据对比，找到需要新增和更新的 action，正确拼装参数，调用相应 api？
