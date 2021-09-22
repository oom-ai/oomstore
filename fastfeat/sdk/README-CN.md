# sdk 包使用方式

该包可以方便的获取特征的具体值，下面是使用方法：

```go
package FastFeatAdapter

import (
	"context"
	"fmt"
	"log"

	"gitlab.pri.ibanyu.com/server/fastfeat/pub.git/FastFeatAdapter"
	. "gitlab.pri.ibanyu.com/server/fastfeat/pub.git/idl/grpc/fastfeat"

	"gitlab.pri.ibanyu.com/server/fastfeat/pub.git/sdk"
)

func ExampleGetOnlineFeatures() {
	response, err := FastFeatAdapter.GetOnlineFeatures(context.Background(), &GetOnlineFeaturesRequest{
		EntityKey:    "user_id:3132393327",
		FeatureNames: []string{sdk.FeatureV1UserStreamKol, sdk.FeatureV1UserStreamIsNewUserBeforeTimestamp, sdk.FeatureV1UserStreamIsNewUserBeforeTimestamp},
	})
	if err != nil {
		panic(err)
	}

	userLastClick5PostId, err := response.FeatureValue(sdk.FeatureV1UserStreamClickedPostidWdw5TpeIntArr).StringArray()
	if err != nil {
		if err == sdk.FastFeatEmptyValue {
			fmt.Printf("feature %s not found\n", sdk.FeatureV1UserStreamClickedPostidWdw5TpeIntArr)
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("feature %s value: %v", sdk.FeatureV1UserStreamClickedPostidWdw5TpeIntArr, userLastClick5PostId)
	}

	kolValue, err := response.FeatureValue(sdk.FeatureV1UserStreamKol).Int64()
	if err != nil {
		if err == sdk.FastFeatEmptyValue {
			log.Printf("feature %s not found", sdk.FeatureV1UserStreamKol)
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("feature %s value: %d\n", sdk.FeatureV1UserStreamKol, kolValue)
	}


	isNewUserValue, err := response.FeatureValue(sdk.FeatureV1UserStreamIsNewUserBeforeTimestamp).Bool()
	if err != nil {
		if err == sdk.FastFeatEmptyValue {
			log.Printf("feature %s not found", sdk.FeatureV1UserStreamIsNewUserBeforeTimestamp)
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("feature %s value: %t\n", sdk.FeatureV1UserStreamIsNewUserBeforeTimestamp, isNewUserValue)
	}

	// Output:
	// feature v1:user:stream:clicked_postid_wdw_5_tpe_int_arr not found
	// feature v1:user:stream:kol not found
	// feature v1:user:stream:is_new_user_before_timestamp value: false
}
```


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
