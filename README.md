[![release](https://img.shields.io/github/release/relunctance/grafalog?style=flat-square)](https://github.com/relunctance/grafalog/releases)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/relunctance/grafalog?tab=doc)


# Grafalog 导入日志小工具
* 支持日志推送到Mysql
* 支持自定义解析日志格式
* 实现接口可支持自定义推送数据源 , 比如Zabbix


# Install 
```
go get -u -v github.com/relunctance/grafalog
```

# Start

* 默认输出到终端

```
package main

import (
    "os"

    "github.com/relunctance/grafalog"
)

func main() {
    g := grafalog.New("./test.logs")
    err = g.Run() // default output os.Stdout
    if err != nil {
        panic(err)
    }
}
```

# Example 示例

* [日志导入Mysql示例](https://github.com/relunctance/grafalog/blob/master/example/mysql/main.go)
* [日志输出到终端示例](https://github.com/relunctance/grafalog/blob/master/example/default/main.go)

# Contribute

* Please feel free to make suggestions, create issues, fork the repository and send pull requests!


