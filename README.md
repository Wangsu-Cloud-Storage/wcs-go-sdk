# Go SDK 使用指南

## 功能说明
wcs-go-sdk 是较为原始的封装，lib 部分没有引入 JSON 库，如果操作返回的结果是 JSON 字符串，您需要自己选择一个 JSON 库，比如 golang 自带的 json，并按照文档去正确解读。

## 版本说明
### 1.0.0.3
## 安装说明
1. 使用go module安装：`go get -u github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib`
2. 使用源码安装：`复制 /src/lib 到您的代码目录下`

## 初始化
调用如下代码进行初始化：
- AccessKey&SecretKey：访问API的密钥，可在控制台安全设置中获取。
- UploadHost：上传域名，可在控制台空间概览获取。
- ManageHost：管理域名，可在控制台空间概览获取。

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")
    
    // 执行操作
}
 ```

## 使用
### 简单上传
单次请求直接上传完整的文件，若文件大小超过2G，必须使用分片上传。

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")
    
    // 执行操作
    su := core.NewSimpleUpload(auth, config, nil)
    
    // 设置token有效期
    deadline := time.Now().Add(time.Second*3600).Unix() * 1000
    put_policy := "{\"scope\": \"bucketName\",\"deadline\": \"" + strconv.FormatInt(deadline, 10) + "\"}"
    response, err := su.UploadFile(`C:\Windows\WindowsShell.Manifest`, put_policy, "WindowsShell.txt", nil)
    if nil != err {
        fmt.Println("UploadFile() failed:", err)
        return
    }
    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =", response.StatusCode)
        fmt.Println(string(body))
    }
}
```

更多实例参考：src/examples/test_simple_upload/

### 分片上传
将一个文件切割为一系列特定大小的数据片，将这些数据片分别上传到服务端，全部上传完后再在服务端将这些数据片合并成为一个资源。

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    bucket := "bucketName"
    key := "keyName"

    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")

    su := core.NewSliceUpload(auth, config, nil)
    put_extra := core.NewPutExtra(days_remain_to_delete)

    // 设置token有效期
    deadline := time.Now().Add(time.Second*3600).Unix() * 1000
    put_policy := "{\"scope\": \"" + bucket + "\",\"deadline\": \"" + strconv.FormatInt(deadline, 10) + "\"}"

    // 指定块大小的方式上传，块最小为4M，并且要为4M的整数倍
    response, err = su.UploadFileWithBlockSize("localfilename", put_policy, key, put_extra, block_size)

    // 指定并发数的方式上传
    response, err = su.UploadFileConcurrent(localfilename, put_policy, key, put_extra, pool_size)

    if nil != err {
        fmt.Println("UploadFile() failed:", err)
        return
    }
    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =", response.StatusCode)
        fmt.Println(string(body))
    }
}
```

更多实例参考：src/examples/test_wcslib/test_wcslib.go

### 删除资源
删除一个指定资源文件。

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")

    bm := core.NewBucketManager(auth, config, nil)
    response, err := bm.Delete("bucket", "key")
    if nil != err {
        Exit(-3, fmt.Sprintf("Delete() failed: %s", err))
        return
    }

    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =", response.StatusCode)
        fmt.Println(string(body))
    }
}
```
更多实例参考：src/examples/test_wcslib/test_wcslib.go

### 获取文件信息
获取一个文件的信息描述，包括文件名，文件大小，文件的ETag、文件上传时间、文件过期时间、文件存储类型等信息。

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")

    bm := core.NewBucketManager(auth, config, nil)
    response, err := bm.Stat(bucket, key)
    if nil != err {
        Exit(-3, fmt.Sprintf("Delete() failed: %s", err))
        return
    }

    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =", response.StatusCode)
        fmt.Println(string(body))
    }
}
```
更多实例参考：src/examples/test_wcslib/test_wcslib.go


### 列举空间
获取空间列表。

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")

    bm := core.NewBucketManager(auth, config, nil)
    response, err := bm.ListBucket()
    if nil != err {
        Exit(-3, fmt.Sprintf("Delete() failed: %s", err))
        return
    }

    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =", response.StatusCode)
        fmt.Println(string(body))
    }
}
```
更多实例参考：src/examples/test_wcslib/test_wcslib.go


### 列举资源
列举指定空间内的资源，如文件较多需要分批列举

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")
    
    bm := core.NewBucketManager(auth, config, nil)

    // 指定单次返回的文件数量，最大1000
    limit := 1000

    // 指定列表排序方式：0代表优先列出目录下的文件；1代表优先列出目录下的文件夹。
    mode := 0

    // 列举指定前缀的文件
    prefix := "prefix"

    // 列举的起始位置，可用上次列举返回的位置标记，作为本次列举的起点信息 
    marker ：= "marker"
    // 列举资源
    response, err := bm.List("bucketName", limit, prefix, mode, marker)
    if nil != err {
        fmt.Println("List() failed:", err)
        return
    }
    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =", response.StatusCode)
        fmt.Println(string(body))
    }
}
```
更多实例参考：src/examples/test_wcslib/test_wcslib.go

### 拉取资源
将指定URL的资源拉取到云存储。

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")

    fm := core.NewFileManager(auth, config, nil)
    response, err := fm.Fetch(fetch_url, bucket, key, prefix, md5, decompression, notify_url, force, separate)
    if nil != err {
        Exit(-3, fmt.Sprintf("Fetch() failed: %s", err))
        return
    }

    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =", response.StatusCode)
        fmt.Println(string(body))
    }
}
```
更多实例参考：src/examples/test_wcslib/test_wcslib.go

### 音视频处理
执行音视频编码和格式转换等操作。

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")

    body := "bucket=aW1hZ2Vz&key=bGVodS5tcDQ==&fops=YXZ0aHVtYi9mbHYvcy80ODB4Mzg0fHNhdmVhcy9hVzFoWjJWek9tZHFhQzVtYkhZPQ==&force=1&separate=1"
    response, err := core.FOps(auth, config, nil, body)
    if nil != err {
        fmt.Println("FOps() failed:", err)
        return
    }
    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =", response.StatusCode)
        fmt.Println(string(body))
    }
}
```
更多实例参考：src/examples/test_wcslib/test_wcslib.go

### 图片鉴定
对指定图片资源（URL)进行智能鉴定的功能。

```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)
func main() {
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")

    image_op := core.NewImageOp(auth, config, nil)
    response, err := image_op.ImageDetect("image", "type", "bucket")
    if nil != err {
        fmt.Println("FOps() failed:", err)
        return
    }
    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =", response.StatusCode)
        fmt.Println(string(body))
    }
}
```
更多实例参考：src/examples/test_wcslib/test_wcslib.go
