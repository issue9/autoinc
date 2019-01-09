autoinc
[![Build Status](https://travis-ci.org/issue9/autoinc.svg?branch=master)](https://travis-ci.org/issue9/autoinc)
[![Build status](https://ci.appveyor.com/api/projects/status/lkit5143eg8f4wuc?svg=true)](https://ci.appveyor.com/project/caixw/autoinc)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/issue9/autoinc/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/autoinc)
======

autoinc 提供了一个简单的 ID 自增功能。
```go
ai := autoinc.New(0, 1, 1)
for i:=0; i<10; i++ {
    fmt.Println(ai.MustID())
}
```

### 安装

```shell
go get github.com/issue9/autoinc
```


### 文档

[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/issue9/autoinc)
[![GoDoc](https://godoc.org/github.com/issue9/autoinc?status.svg)](https://godoc.org/github.com/issue9/autoinc)


### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
