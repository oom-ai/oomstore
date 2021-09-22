# sdk 包的 api 文档

### type OnlineFetaturesResponses 
```go
type OnlineFetaturesResponses map[string]FeatureValue
```

#### func NewOnlineFeatureResponses
```go
func NewOnlineFeatureResponses(rawResponse GetOnlineFeaturesResponse) OnlineFetaturesResponses
```
NewOnlineFeatureResponses 包装 grpc 接口返回的 GetOnlineFeaturesResponse，可以方便的获取具体类型的特征值。

#### func FeatureValue
```go
func (o OnlineFetaturesResponses) FeatureValue(featureName string) FeatureValue
```
获取指定 FeatureName 的 FeatureValue 类型， Featurevalue 提供了方便的 api 用于获取具体类型的值。


### type FeatureValue
```golang
type Reader struct {
 // 内涵隐藏或非导出字段
}
```

#### func (v featureValue) String
```go
func (v featureValue) String() (string, error)
```
返回特征的 string 类型值，注意只有值类型是 string 的特征允许使用该方法。

#### func (v featureValue) StringArray
```go
func (v featureValue) StringArray() ([]string, error) 
```
返回特征的 []string 类型值，注意只有值类型是 []string 的特征允许使用该方法。

#### func (v featureValue) Int64
```go
func (v featureValue) Int64() (int64, error) 
```
返回特征的 int64 类型值，注意只有值类型是 int64 的特征允许使用该方法。

#### func (v featureValue) Int64Array
```go
func (v featureValue) Int64Array() ([]int64, error) 
```
返回特征的 []int64 类型值，注意只有值类型是 []int64 的特征允许使用该方法。

#### func (v featureValue) Double
```go
func (v featureValue) Double() (float64, error) 
```
返回特征的 float64 类型值，注意只有值类型是 float64 的特征允许使用该方法。

#### func (v featureValue) DoubleArray
```go
func (v featureValue) DoubleArray() ([]float64, error) 
```
返回特征的 []float64 类型值，注意只有值类型是 []float64 的特征允许使用该方法。

#### func (v featureValue) Bool
```go
func (v featureValue) Bool() (bool, error) 
```
返回特征的 bool 类型值，注意只有值类型是 bool 的特征允许使用该方法。

#### func (v featureValue) BoolArray
```go
func (v featureValue) BoolArray() ([]bool, error) 
```
返回特征的 []bool 类型值，注意只有值类型是 []bool 的特征允许使用该方法。
