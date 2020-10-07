autoinc
[![Go](https://github.com/issue9/autoinc/workflows/Go/badge.svg)](https://github.com/issue9/autoinc/actions?query=workflow%3AGo)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/issue9/autoinc/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/autoinc)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/autoinc)](https://pkg.go.dev/github.com/issue9/autoinc)
======

autoinc 提供了一个简单的 ID 自增功能

```go
ai := autoinc.New(0, 1, 1)
for i:=0; i<10; i++ {
    fmt.Println(ai.MustID())
}
```

安装
----

```shell
go get github.com/issue9/autoinc
```

版权
----

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
