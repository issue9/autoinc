autoinc [![Build Status](https://travis-ci.org/issue9/autoinc.svg?branch=master)](https://travis-ci.org/issue9/autoinc)
======

autoinc提供了一个简单的ID自增功能。
```go
ai := autoinc.New(0, 1, 1)
for i:=0; i<10; i++ {
    fmt.Println(ai.ID())
}
```

### 安装

```shell
go get github.com/issue9/autoinc
```


### 文档

[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/issue9/autoinc)
[![GoDoc](https://godoc.org/github.com/issue9/autoinc?status.svg)](https://godoc.org/github.com/issue9/autoinc)


### 版权

本项目采用[MIT](http://opensource.org/licenses/MIT)开源授权许可证，完整的授权说明可在[LICENSE](LICENSE)文件中找到。
