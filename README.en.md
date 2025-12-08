# Go SDK User Guide

## 语言 / Language
- [简体中文](README.md)
- [English](README.en.md)

## Feature Description
wcs-go-sdk is a relatively raw encapsulation. The lib part does not include a JSON library. If the operation returns a JSON string, you need to choose a JSON library yourself, such as the json library that comes with golang, and correctly interpret it according to the documentation.

## Version Description
### 1.0.0.3
## Installation Instructions
1. Install using go module: `go get -u github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib`
2. Install using source code: `Copy /src/lib to your code directory`

## Initialization
Call the following code for initialization:
- AccessKey&SecretKey: The key to access the API, which can be obtained in the console security settings.
- UploadHost: Upload domain name, which can be obtained in the console space overview.
- ManageHost: Management domain name, which can be obtained in the console space overview.

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
    // Configure authentication AK and SK
    auth := utility.NewAuth("<AccessKey>", "<SecretKey>")
    
    // Configure whether to use https, upload domain name and management domain name (Note: No need to add http:// prefix to the domain name)
    config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")
    
    // Perform operations
}
 ```

## Usage
### Calculate Upload Token
```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
)
func main() {
    // Assemble put_policy according to actual needs
    put_policy := `{"scope":"aaa","deadline":"1893427200000"}`
    
    auth := utility.NewAuth("aaasdfasf", "bbbsdfdsfdsafsdf")
    fmt.Printf(auth.CreateUploadToken(put_policy))
}
```

### Simple Upload
Upload the complete file in a single request. If the file size exceeds 2G, you must use chunked upload.

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
    
    // Perform operations
    su := core.NewSimpleUpload(auth, config, nil)
    
    // Set token validity period
    deadline := time.Now().Add(time.Second*3600).Unix() * 1000
    put_policy := `{"scope": "bucketName","deadline": "` + strconv.FormatInt(deadline, 10) + `"}`
    response, err := su.UploadFile(`C:\Windows\WindowsShell.Manifest`, put_policy, "WindowsShell.txt", nil)
    if nil != err {
        fmt.Println("UploadFile() failed:", err)
        return
    }
    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =")
        fmt.Println(string(body))
    }
}
```

For more examples, please refer to: src/examples/test_simple_upload/

### Chunked Upload
Split a file into a series of data chunks of a specific size, upload these chunks to the server separately, and then merge these chunks into a resource on the server after all uploads are completed.

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

    // Set token validity period
    deadline := time.Now().Add(time.Second*3600).Unix() * 1000
    put_policy := `{"scope": "` + bucket + `","deadline": "` + strconv.FormatInt(deadline, 10) + `"}`

    var response *http.Response
    var err error
    // Upload by specifying block size, the block size unit is MB, and it must be an integer multiple of 4
    response, err = su.UploadFileWithBlockSize("localfilename", put_policy, key, put_extra, block_size)

    // Upload by specifying concurrency
    response, err = su.UploadFileConcurrent(localfilename, put_policy, key, put_extra, pool_size)

    if nil != err {
        fmt.Println("UploadFile() failed:", err)
        return
    }
    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =")
        fmt.Println(string(body))
    }
}
```

For more examples, please refer to: src/examples/test_wcslib/test_wcslib.go

### Delete Resource
Delete a specified resource file.

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
        fmt.Sprintf("Delete() failed: %s", err)
        return
    }

    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =")
        fmt.Println(string(body))
    }
}
```
For more examples, please refer to: src/examples/test_wcslib/test_wcslib.go

### Get File Information
Get the information description of a file, including file name, file size, ETag of the file, file upload time, file expiration time, file storage type, etc.

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
        fmt.Println("Failed, StatusCode =")
        fmt.Println(string(body))
    }
}
```
For more examples, please refer to: src/examples/test_wcslib/test_wcslib.go


### List Buckets
Get the bucket list.

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
        fmt.Println("Failed, StatusCode =")
        fmt.Println(string(body))
    }
}
```
For more examples, please refer to: src/examples/test_wcslib/test_wcslib.go


### List Resources
List resources in a specified bucket. If there are many files, you need to list them in batches.

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

    // Specify the number of files returned in a single time, maximum 1000
    limit := 1000

    // Specify the list sorting method: 0 means priority to list files in the directory; 1 means priority to list folders in the directory.
    mode := 0

    // List files with the specified prefix
    prefix := "prefix"

    // The starting position of the listing, the position mark returned by the last listing can be used as the starting point information for this listing 
    marker := "marker"
    // List resources, note: the last two parameters are deprecated, use "" to fill
    response, err := bm.List("bucketName", limit, prefix, mode, marker, "", "")
    if nil != err {
        fmt.Println("List() failed:", err)
        return
    }
    body, _ := ioutil.ReadAll(response.Body)
    if http.StatusOK == response.StatusCode {
        fmt.Println(string(body))
    } else {
        fmt.Println("Failed, StatusCode =")
        fmt.Println(string(body))
    }
}
```
For more examples, please refer to: src/examples/test_wcslib/test_wcslib.go

### Fetch Resource
Pull resources from the specified URL to cloud storage.

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
        fmt.Println("Failed, StatusCode =")
        fmt.Println(string(body))
    }
}
```
For more examples, please refer to: src/examples/test_wcslib/test_wcslib.go

### Audio and Video Processing
Perform audio and video encoding, format conversion and other operations.

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
        fmt.Println("Failed, StatusCode =")
        fmt.Println(string(body))
    }
}
```
For more examples, please refer to: src/examples/test_wcslib/test_wcslib.go

### Image Detection
Intelligent detection function for specified image resources (URL).

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
        fmt.Println("Failed, StatusCode =")
        fmt.Println(string(body))
    }
}
```

### Calculate Local File Etag
```
package main

import (
    "fmt"
    "github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
)
func main() {
    etag := utility.ComputeEtag([]byte(data))
    fmt.Println("ETag =")
}
```
For more examples, please refer to: src/examples/test_wcslib/test_wcslib.go