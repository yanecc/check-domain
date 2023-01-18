# Build

``` shell
git clone git@github.com:18183883296/check-domain.git
go build check-domain
```

# Usage

## 初始化配置文件

``` shell
check-domain
```

初次运行会生成 `config.toml` 配置文件，内容如下：

``` toml
# api获取地址：https://user.whoisxmlapi.com/products
apiKey = "apiKey"

accurateMode = true

# 免费账户有500次Whois查询额度及100次/月域名可用性检测额度
useWhois = false
```

check-domain 程序调用 WhoisXMLApi 的接口来检查域名是否可注册，需要先访问 [WhoisXml Api](https://user.whoisxmlapi.com/) 注册账号，获取 `apiKey` 填入配置文件。默认配置 `accurateMode = true` 可以更准确地检查域名可用性，但会耗时稍长。

## 批量查询

### 无标签查询

``` shell
check-domain
```

直接运行check-domain，从同目录下 `domains.txt` 文件中逐行读取**完整域名**，检查域名是否可注册。

### 指定读入文件

``` shell
check-domain -path ./method/domains.txt
```

`-path` 标签用于指定包含待检查域名的文本文件路径。

### 指定域名后缀

``` shell
check-domain -suffix app
```

`-suffix` 标签用于为待检查域名设定统一后缀，从文件中逐行读取**域名前缀**。
