# softvowifi-netlink

[English](README.md) · **简体中文**

面向SoftVoWiFi维护的公开Linux netlink源码快照，保留
`github.com/iniwex5/netlink`模块和导入路径。它是独立维护仓，不是完整上游镜像，
也不属于GitHub上游fork network。

## 维护范围

- 保留iniwex5/netlink v1.3.3的netlink与`nl`生产源码。
- 修复XFRM通配selector的地址族解码：按内核返回的`msg.Sel.Family`处理，
  不让IPv6 `::/0`在`FAMILY_ALL`查询中变成IPv4 `0.0.0.0/0`。
- 保留聚焦的合成解码回归；没有完整上游测试或通用内核兼容保证。

## 使用

Linux内核接口受内核/命名空间/权限限制；修改网络通常需要`CAP_NET_ADMIN`。
`go.mod`声明Go 1.23，CI使用Linux/Go 1.26。请勿为了试用而修改宿主网络。

在已有Go模块中明确选择维护快照，导入路径保持不变：

```sh
go mod edit -replace=github.com/iniwex5/netlink=github.com/lilyzhaun/softvowifi-netlink@v0.0.0-20260915075719-be8893d91893
go get github.com/iniwex5/netlink@v1.3.3
go list -m -json github.com/iniwex5/netlink
```

以上固定可复现代码版本，不自动跟随main。英文README提供完整只读示例；
以`go list`的Replace字段确认实际消费来源，不把原模块版本误当未经修改的上游树。

## 测试与协作

```sh
go test -p 1 -race -shuffle=on -count=1 -timeout=120s ./...
go vet -p 1 -composites=false ./...
```

默认测试使用合成netlink消息，不修改宿主网络；CI不执行真实特权内核验收。
完整vet仍报告addr_linux.go既有未具名结构体字面量。CI关闭整个composites分析器，
其他默认分析继续执行；旧源码格式只报告、不重写，不冒称完整vet/格式门禁通过。
提交问题请使用[模板](https://github.com/lilyzhaun/softvowifi-netlink/issues/new/choose)，
流程见[贡献指南](CONTRIBUTING.md)，漏洞请按[安全策略](SECURITY.md)私密报告。
不要提交真实XFRM密钥、完整配置、凭据或原始抓包。

## 许可与来源

保留[Apache-2.0原文](LICENSE)、原版权及
[Go BSD-3-Clause文本](LICENSES/Go-BSD-3-Clause.txt)。
[PROVENANCE](PROVENANCE.md)中的复制来源/ipset义务未决项保持，
不能把根许可证解释为所有历史来源均已核定。[旧README](docs/README-imported.md)
仅作历史，不使用其中过期私有状态、安装路径或上游CI徽章。
