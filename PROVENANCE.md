# netlink 本地副本来源与补丁

## 来源

- 上游仓库：`github.com/iniwex5/netlink`
- 版本：`v1.3.3`
- 原始提交：`433ed87f0847478420ddaab9fdfbbd0b912c786f`（`refs/tags/v1.3.3`）
- Go 模块归档：`/root/go/pkg/mod/cache/download/github.com/iniwex5/netlink/@v/v1.3.3.zip`
- 归档 SHA256：`0cd455b70a1e224419a1ced38cd882d167fb081b3bc228ad1efd03767a340c57`
- 模块校验和（go.sum h1）：`h1:Dk6b7Wao+v4leIqv9xnUSLCxpA7wzm664+adY+KD8Z0=`
- 许可证：Apache License 2.0（见同目录 `LICENSE`），保留上游版权与声明，未修改
- 最初导入使用本机 Go 模块缓存；2026-09-15 另只读核实[公开代理版本信息](https://proxy.golang.org/github.com/iniwex5/netlink/@v/v1.3.3.info) 的 Origin.Hash/Ref 与上述提交/tag 一致，重新读取[公开 zip](https://proxy.golang.org/github.com/iniwex5/netlink/@v/v1.3.3.zip) 的 SHA256 与上述归档一致，[公开校验和记录](https://sum.golang.org/lookup/github.com/iniwex5/netlink@v1.3.3) 返回上述 h1。此为固定版本的公开代理证据，不证明整个 Git 历史、当前 HEAD 或是否已修复。

## 导入内容

只导入实际编译所需源码，未做整体 vendor 或镜像：

- 根包非测试 Go 源：68 个
- `nl/` 非测试 Go 源：25 个
- `LICENSE`、`README.md`、`go.mod`、`go.sum`

未导入：上游全部 `*_test.go`、`cmd/`、`.github/`、`Makefile`、`CHANGELOG.md`、`testdata/`、`.gitignore`。未导入内容不影响本仓库编译。

## 本地行为补丁（两行，另有文件内修改声明）

文件 `xfrm_policy_linux.go`，`parseXfrmPolicy` 内：

```diff
-	policy.Dst = msg.Sel.Daddr.ToIPNet(msg.Sel.PrefixlenD, uint16(family))
-	policy.Src = msg.Sel.Saddr.ToIPNet(msg.Sel.PrefixlenS, uint16(family))
+	policy.Dst = msg.Sel.Daddr.ToIPNet(msg.Sel.PrefixlenD, msg.Sel.Family)
+	policy.Src = msg.Sel.Saddr.ToIPNet(msg.Sel.PrefixlenS, msg.Sel.Family)
```

原因：`XfrmPolicyGet` 以 `FAMILY_ALL` 调用 `parseXfrmPolicy`，而 `nl.XfrmAddress.ToIPNet` 仅在地址与前缀长度均为 0 时使用传入 family。旧代码传函数参数 `family`（即 `FAMILY_ALL`），导致全零 IPv6 selector 被解码成 `0.0.0.0/0`。改为使用内核实际回报的 `msg.Sel.Family`，保留原 `family` 参数仅作族过滤，未改为从调用方期望推断族，也未放宽指纹或隔离判定。

新增本地回归测试 `xfrm_policy_decoder_family_test.go`，用 `nl.XfrmUserpolicyInfo.Serialize` 构造原始字节，覆盖 IPv6/IPv4 通配符、非通配地址与族过滤。

## 自有维护迁移（Issue 35，远端切换未完成）

拟定新私有目标为 `lilyzhaun/softvowifi-netlink`，由 `lilyzhaun` 创建/确认和维护，`newdamm` 提供工具与迁移 PR。远端状态未确认，本任务未创建或推送；不以不可见响应断言仓库不存在。

初始维护快照冻结于 SoftVoWiFi `53c74b87a536eb684f0a84a474269a4cf596a1e2` 的本目录；它是上述保留子集加本地补丁，不是完整上游树。原有来源/缓存证据继续保留，新增公开代理核实仅限上节列明范围，不等于完整远端历史验证。当前模块路径 `github.com/iniwex5/netlink` 和 Apache-2.0 `LICENSE` 原文不变。

主仓库 `engine/third_party/manifest.json` 为机器事实源，导出 `.meta.json` 另记实际 commit/tree。主仓库 `docs/第三方Go依赖维护与导出.md` 规定 netlink 先行、实际版本固定、下游 swu-go 及主模块 replace 的后续切换。2026-09-15 源码就绪切片仅补充文件内修改声明、来源证据及随附许可，不改变上述两行行为，不免除发行义务，也不代表迁移完成。

## 2026-09-15 随附许可与仍待核定来源

- 固定 `v1.3.3.zip` 的完整成员清单中，LICENSE/NOTICE/COPYING/AUTHORS/COPYRIGHT 类文件仅有根 `LICENSE`，没有 NOTICE。因此本切片不编造 NOTICE；这不证明未核查的完整上游 Git 分发没有 NOTICE。
- `nl/nl_linux.go` 的 `nlmAlignOf` 保留原归属：`Taken from syscall/netlink_linux.go by The Go Authors under BSD-style license.` 本地 Go 源 `syscall/netlink_linux.go` 具有 `Copyright 2011 The Go Authors. All rights reserved.` 和 BSD-style 文件头；其函数与本库实现一致，仅本库对 `NLMSG_ALIGNTO` 使用 `syscall.` 限定。
- 补充全文 `LICENSES/Go-BSD-3-Clause.txt` 逐字取自本地完整 Go 安装的 `/usr/local/go/LICENSE`。该安装以 `GOROOT=/usr/local/go GOTOOLCHAIN=local /usr/local/go/bin/go env GOROOT GOVERSION` 核实为 `/usr/local/go`、`go1.22.1`，不是当前 engine 工具链或当年 netlink 导入版本。默认系统 Go 的 `/usr/lib/go-1.24/LICENSE` 缺失，未以其冒充全文来源。
- 补充许可及本地原文 SHA256：`2d36597f7117c38b006835ae7f537487207d8ec407aa9d9980794b2030cbc067`；核对源文件 `/usr/local/go/src/syscall/netlink_linux.go` SHA256：`569949cb048d294b11e02078525d4ac2be2b6d8520575fe4c9c3234ef418242b`。这是已验证 BSD 三条文本及函数来源家族，**原始导入时的 Go 版本/年份仍未知**；不能宣称找到了当年原件或整个组件已获法律放行。
- `link_linux.go` 第1548–1565行及 `link_tuntap_linux.go` 的对应注释明确写明 `The below comment taken from`，保留来自 [mistsys/tuntap 固定提交](https://github.com/mistsys/tuntap/commit/161418c25003bbee77d085a34af64d189df62bea) 的说明文字及链接，属于**保留的第三方复制材料**，不是“仅参考、无复制”。只读核实该提交的 [tun_linux.go](https://github.com/mistsys/tuntap/blob/161418c25003bbee77d085a34af64d189df62bea/tun_linux.go) 含对应文字；文件直接以 `package tuntap` 开头，无文件头版权声明。本轮不据此声称已穷尽其他代码的复制范围。
- 该 mistsys 提交的 [LICENSE](https://github.com/mistsys/tuntap/blob/161418c25003bbee77d085a34af64d189df62bea/LICENSE) 为 Apache-2.0，SHA256 为 `cb5e8e7e5f4a3988e1063c142c60dc2df75605f4c46515e776e3aca6df976e14`；固定提交递归树 API 返回9项且 `truncated=false`，没有 NOTICE。其第1–9节与本目录现有 Apache LICENSE 在归一化空白后相同；附录为通用占位模板，而本目录尾部保留 Vishvananda Ishaya / Docker 实际版权，故不声称两文件逐字相同。沿用现有 Apache 许可全文副本及原来源注释，不重复添加同条款许可或编造 NOTICE，不修改根许可/SPDX；这不是整个组件的法律放行。
- ipset 的可核实继承来源为 iniwex5 `v1.3.3` 公开代理文件快照：根 `ipset_linux.go` 与 `nl/ipset_linux.go` 均逐字相同。GitHub API 将 vishvananda 的引入提交解析为 [1e3d26b1b393006411752fc0196d65df9fc6fa15](https://github.com/vishvananda/netlink/commit/1e3d26b1b393006411752fc0196d65df9fc6fa15)（`Add support for ipset`）；该固定提交的 [LICENSE](https://github.com/vishvananda/netlink/blob/1e3d26b1b393006411752fc0196d65df9fc6fa15/LICENSE) SHA256 为 `04cf148229e562c643203d1ee8c43aae2022e5130e95ceb8ba8357a27c1a9abe`，与本目录 Apache LICENSE 相同。这不是当前文件与引入时文件完全一致或完整演化历史已核实的声明。
- 本次实际源码核对纠正“ipset 没有显式 C 源指针”的外部摘要：根文件第389行明确链接 [Olipro/ipset 固定 C 源](https://github.com/Olipro/ipset/blob/9f145b49100104d6570fe5c31a5236816ebb4f8f/kernel/net/netfilter/ipset/ip_set_hash_ipport.c#L30)，第401、425、451行还有 `Taken from` 及 C 文件路径；`nl/ipset_linux.go` 未发现对应显式链接。上述归属全部保留，不能以某文件缺指针或仓库根 Apache 许可推导“非 C 派生”或“没有 GPL 义务”。C 来源的具体复制范围及适用义务仍为 **unknown / 待核定**，本轮不擅改许可分类。
- 根 Apache `LICENSE`（含 Vishvananda Ishaya / Docker 版权）原字节不变。补充 BSD 文本不替换根许可，不改变 manifest 原许可哈希。
