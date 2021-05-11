# Go SDK 使用指南

## 功能说明

wcs-go-sdk 是较为原始的封装，lib 部分没有引入 JSON 库，如果操作返回的结果是 JSON 字符串，您需要自己选择一个 JSON 库，比如 golang 自带的 json，并按照文档去正确解读。

## 版本说明

### 1.0.0.3


## 安装说明

1. 使用go module安装：`go get -u github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib`
2. 使用源码安装：`复制 /src/lib 到您的代码目录下`


## 初始化

在使用 SDK 之前，您需要获得一对有效的 AccessKey 和 SecretKey 签名授权。

可以通过如下方法获得：

1. 开通网宿云存储账号
2. 登录网宿 SI 平台，在安全管理-秘钥管理查看 AccessKey 和 SecretKey
3. 登录网宿 SI 平台，在安全管理-域名管理查看上传域名（UploadHost）和管理域名(ManageHost)。

 获取上面配置之后，调用如下代码进行初始化：

 ```
 auth := utility.NewAuth("<AccessKey>", "<SecretKey>")

 //config := core.NewDefaultConfig()

 config := core.NewConfig(false, "<UploadHost>", "<ManageHost>")

 ```

4. 更多范例请参考 src/examples
