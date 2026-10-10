> **Historical README, archived 2026-10-10.** Earlier private status, invitations,
> upstream badges and installation commands below are historical, not current
> instructions. See the [current README](../README.md). Only relative documentation
> links were adjusted when archiving; the source notes and original examples remain.

> **Private maintenance repository (Issue 35, checked 2026-09-15)**: [lilyzhaun/softvowifi-netlink](https://github.com/lilyzhaun/softvowifi-netlink) exists and is PRIVATE. `lilyzhaun` owns the repository and final integration; `newdamm` contributes reviewed maintenance PRs. The write invitation is **PENDINGINVITATION** at this check, not confirmed access; only `newdamm` can accept it. No admin permission was granted.
> Module path remains `github.com/iniwex5/netlink`. SoftVoWiFi's verified migration candidate consumes pinned commit `be8893d918930ed1f6c790c41e51c13c20b20cf9` (`v0.0.0-20260915075719-be8893d91893`), not the latest maintenance `main` HEAD. This documentation-only update advances `main` without changing that pin, production code, module files or licenses; the main-repository migration still awaits final acceptance.
> The initial snapshot was imported from SoftVoWiFi commit `8ca6114de2c915473165bf0698a09ced0f9108e3`, path `engine/third_party/netlink`. [PROVENANCE.md](../PROVENANCE.md), the original README below (including historical commands), [LICENSE](../LICENSE) and [Go BSD license](../LICENSES/Go-BSD-3-Clause.txt) are preserved. Unknown provenance and licensing limits remain unresolved; PRIVATE is not legal or public-release clearance. This is a private snapshot import, not a GitHub fork-network claim.
> Canonical policy: `engine/third_party/manifest.json` in SoftVoWiFi; handoff: `docs/第三方Go依赖维护与导出.md`. Initial-import notes, upstream badges and examples below are historical, not current integration instructions or device validation. Future code changes require review and a separately verified pin update.

# netlink - netlink library for go #

![Build Status](https://github.com/vishvananda/netlink/actions/workflows/main.yml/badge.svg) [![GoDoc](https://godoc.org/github.com/vishvananda/netlink?status.svg)](https://godoc.org/github.com/vishvananda/netlink)

The netlink package provides a simple netlink library for go. Netlink
is the interface a user-space program in linux uses to communicate with
the kernel. It can be used to add and remove interfaces, set ip addresses
and routes, and configure ipsec. Netlink communication requires elevated
privileges, so in most cases this code needs to be run as root. Since
low-level netlink messages are inscrutable at best, the library attempts
to provide an api that is loosely modeled on the CLI provided by iproute2.
Actions like `ip link add` will be accomplished via a similarly named
function like AddLink(). This library began its life as a fork of the
netlink functionality in
[docker/libcontainer](https://github.com/docker/libcontainer) but was
heavily rewritten to improve testability, performance, and to add new
functionality like ipsec xfrm handling.

## Local Build and Test ##

You can use go get command:

    go get github.com/vishvananda/netlink

Testing dependencies:

    go get github.com/vishvananda/netns

Testing (requires root):

    sudo -E go test github.com/vishvananda/netlink

## Examples ##

Add a new bridge and add eth1 into it:

```go
package main

import (
    "fmt"
    "github.com/vishvananda/netlink"
)

func main() {
    la := netlink.NewLinkAttrs()
    la.Name = "foo"
    mybridge := &netlink.Bridge{LinkAttrs: la}
    err := netlink.LinkAdd(mybridge)
    if err != nil  {
        fmt.Printf("could not add %s: %v\n", la.Name, err)
    }
    eth1, _ := netlink.LinkByName("eth1")
    netlink.LinkSetMaster(eth1, mybridge)
}

```
Note `NewLinkAttrs` constructor, it sets default values in structure. For now
it sets only `TxQLen` to `-1`, so kernel will set default by itself. If you're
using simple initialization(`LinkAttrs{Name: "foo"}`) `TxQLen` will be set to
`0` unless you specify it like `LinkAttrs{Name: "foo", TxQLen: 1000}`.

Add a new ip address to loopback:

```go
package main

import (
    "github.com/vishvananda/netlink"
)

func main() {
    lo, _ := netlink.LinkByName("lo")
    addr, _ := netlink.ParseAddr("169.254.169.254/32")
    netlink.AddrAdd(lo, addr)
}

```

## Future Work ##

Many pieces of netlink are not yet fully supported in the high-level
interface. Aspects of virtually all of the high-level objects don't exist.
Many of the underlying primitives are there, so its a matter of putting
the right fields into the high-level objects and making sure that they
are serialized and deserialized correctly in the Add and List methods.

There are also a few pieces of low level netlink functionality that still
need to be implemented. Routing rules are not in place and some of the
more advanced link types. Hopefully there is decent structure and testing
in place to make these fairly straightforward to add.
