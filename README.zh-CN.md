# go-ext

[English](README.md) · [简体中文](README.zh-CN.md)

[![CI](https://github.com/qizhaoTan/go-ext/actions/workflows/ci.yml/badge.svg)](https://github.com/qizhaoTan/go-ext/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/qizhaoTan/go-ext.svg)](https://pkg.go.dev/github.com/qizhaoTan/go-ext)
[![Go Report Card](https://goreportcard.com/badge/github.com/qizhaoTan/go-ext)](https://goreportcard.com/report/github.com/qizhaoTan/go-ext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Go 泛型数据结构与函数式工具库——补齐标准库没有覆盖的部分。

> **状态：早期开发中。** 仓库骨架和设计已确定，但下列包尚未实现。当前无法导入使用，
> API 仍可能调整。如果想知道首个版本何时发布，可以 watch 本仓库。

## 规划中的包

| 包 | 职责 |
| --- | --- |
| `slices` | 切片上的即时求值函数式工具——`Map`、`Filter`、`Reduce`、`GroupBy`、`Chunk` 等 |
| `iters` | 基于 `iter.Seq` 的**惰性**管道——组合过程不产生中间分配，可处理无限序列 |
| `sets` | 泛型哈希集合，含完整集合代数 |
| `maps` | `Keys`、`Values`、`Filter`、`Invert`、`Merge` 等 |
| `ptr` | 值与指针互转，用于可选字段 |
| `constraints` | 数值类型约束 |

## 设计目标

以下是实现过程中要遵守的约束。

**零依赖。** 永不引入标准库以外的依赖。一个工具库如果拖进一整棵依赖树，代价比它
替代的代码更高。

**纯函数。** 不修改入参。所有返回的集合都是新分配的，且不与入参共享底层数组，
因此结果可以放心交给调用方，不会有别名带来的意外。

**`nil` 即空。** `nil` 切片和 map 遍历零次，因此它会作为空集合自然地流过每个函数
——不需要防御性判断，API 里也不设特例。

**不重复标准库。** 排序、查找、相等性判断已经在标准库的 `slices` 和 `maps` 里了，
本库只填补空缺。这也是 `constraints` 不会定义 `Ordered` 的原因——`cmp.Ordered`
已经存在。

**用自由函数，不用链式调用。** Go 的方法无法引入新的类型参数，因此 `.Map().Filter()`
这样的链式 API 只能在构造时锁定元素类型，永远无法做类型变换。标准库出于同样的原因
做了同样的选择。

**调用方的 bug 就 panic。** 非正数的分块大小没有正确答案。返回一个随意的结果只会
把调用处的 bug 藏起来，panic 不会。

## 即时求值还是惰性求值

两个集合包会提供相同的操作，但求值策略不同，在两者之间做选择是本库要求使用者做的
主要决策。

`slices` 是**即时求值**：每一步都分配新切片并完整执行。这是合适的默认选择——实现
简单、小数据量下更快，而且中间结果就是普通切片，可以直接在调试器里查看。

`iters` 是**惰性求值**：在序列被消费之前什么都不会执行，全程不构造中间切片，且只
计算真正被需要的元素。正因如此，作用在无限序列上的管道才能正常终止。

经验法则：数据源很大、计算代价高、或者无界时用 `iters`，其余情况用 `slices`。

## 环境要求

Go 1.24 及以上。`iters` 包构建在 `iter.Seq` 之上，这也是版本下限无法再降低的原因。

## 开发

```bash
git clone https://github.com/qizhaoTan/go-ext.git
cd go-ext
make check   # 格式化 + vet + 测试
```

| 目标 | 作用 |
| --- | --- |
| `make test` | `go test -race ./...` |
| `make cover` | 生成覆盖率报告 `coverage.html` |
| `make lint` | `golangci-lint run` |
| `make fmt` | `gofmt -w .` |
| `make check` | 格式化 + vet + 测试，每次提交前执行 |

CI 会在 Linux、macOS、Windows 三个平台上，用 Go 1.24 和当前 stable 两个版本，
开启竞态检测运行测试。

## 参与贡献

请阅读 [CONTRIBUTING.md](CONTRIBUTING.md)。在库还处于搭建阶段时，提交大型 PR 之前
先开一个 issue 讨论，可以节省你的时间。

## 许可证

[MIT](LICENSE) © qizhaoTan
