# Grafalog 导入日志小工具
* 支持日志推送到Mysql
* 支持自定义解析日志格式
* 实现接口可支持自定义推送数据源 , 比如Zabbix



# Start

* 默认输出到终端

```
package main

import (
    "os"

    "github.com/relunctance/grafalog"
)

func main() {
    f, err := os.Open("test.logs")
    if err != nil {
        panic("open test.logs is faild")
    }
    defer f.Close()
    g := grafalog.New(f)
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


