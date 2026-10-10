# softvowifi-netlink

[![CI](https://github.com/lilyzhaun/softvowifi-netlink/actions/workflows/ci.yml/badge.svg)](https://github.com/lilyzhaun/softvowifi-netlink/actions/workflows/ci.yml)

**English** · [简体中文](README.zh-CN.md)

A public, maintained Linux netlink snapshot used by SoftVoWiFi. It exposes the
retained link, address, route, rule and XFRM APIs under the original module path
`github.com/iniwex5/netlink`.

This is an independent maintained snapshot, **not an official upstream release
or a complete mirror**. It is not part of GitHub's upstream fork network.

## What is maintained here?

- Production netlink and `nl` sources retained from `iniwex5/netlink v1.3.3`.
- An XFRM policy decoder fix that uses the kernel selector's family, rather than
  `FAMILY_ALL`, when decoding wildcard selectors. In particular, `::/0` must not
  turn into `0.0.0.0/0`.
- Focused, synthetic decoder regressions. The complete upstream test suite and
  examples were not imported; passing this repository's CI is not full upstream
  or kernel compatibility certification.

The code remains the authority for the retained API. Background, source pointers
and licensing qualifications are in [PROVENANCE.md](PROVENANCE.md).

## Requirements

- Linux for kernel-facing operations; usable APIs depend on the host kernel and
  network namespace.
- `go.mod` declares Go 1.23. CI uses Go 1.26 on Linux.
- Network mutations normally require `CAP_NET_ADMIN`. Some read-only operations
  can run without it, subject to namespace, seccomp and platform policy.

Do not run network-mutating code in your host namespace just to try an example.
Use a disposable, isolated environment for privileged integration tests.

## Use the maintained snapshot

The repository URL and module path deliberately differ. In an existing Go module:

```sh
go mod edit -replace=github.com/iniwex5/netlink=github.com/lilyzhaun/softvowifi-netlink@v0.0.0-20260915075719-be8893d91893
go get github.com/iniwex5/netlink@v1.3.3
go list -m -json github.com/iniwex5/netlink
```

The example pins a reproducible maintained code revision. It does not follow
`main` automatically. Keep importing `github.com/iniwex5/netlink`, and inspect
the `Replace` entry from `go list` to confirm what your build actually consumes.
Do not confuse the upstream `v1.3.3` requirement with an unmodified upstream tree.

## Read-only example

This lists interfaces in the process's current network namespace; it does not
create links, change addresses or install routes or IPsec state.

```go
package main

import (
	"fmt"
	"log"

	"github.com/iniwex5/netlink"
)

func main() {
	links, err := netlink.LinkList()
	if err != nil {
		log.Fatal(err)
	}
	for _, link := range links {
		fmt.Printf("%d: %s\n", link.Attrs().Index, link.Attrs().Name)
	}
}
```

## Develop and test

```sh
git clone https://github.com/lilyzhaun/softvowifi-netlink.git
cd softvowifi-netlink
go test -p 1 -race -shuffle=on -count=1 -timeout=120s ./...
go vet -p 1 -composites=false ./...
```

The retained decoder tests use constructed netlink messages, not host network
mutations. CI does not run as root or exercise real kernel XFRM interoperability.
Full `go vet` currently reports an existing unkeyed `nl.IfaCacheInfo` literal in
`addr_linux.go`. CI disables the **whole composites analyzer**, not just that
literal; the other default analyzers still run. Legacy formatting differences
are reported without rewriting retained source. These are explicit limitations,
not a claim that full vet or a full formatting check passed.
The [contribution guide](CONTRIBUTING.md) explains the review and test boundaries.

## Support and security

- [Bug reports and feature requests](https://github.com/lilyzhaun/softvowifi-netlink/issues/new/choose)
- [Security policy and private reporting](SECURITY.md)
- [Code of conduct](CODE_OF_CONDUCT.md)

Use synthetic reproductions. Never attach real XFRM keys, credentials, private
configuration or unsanitized packet captures to a public issue.

## License and provenance

The original [Apache-2.0 LICENSE](LICENSE) and its copyright notices are retained.
The Go-derived alignment helper has an additional
[BSD-3-Clause notice](LICENSES/Go-BSD-3-Clause.txt).

See [PROVENANCE.md](PROVENANCE.md), including its unresolved copied-source and
ipset obligation notes. The repository's root license is not a blanket statement
that every historical attribution question has been resolved. No original
license text or copyright has been replaced by this documentation refresh.

The [earlier README](docs/README-imported.md) is archived for source context;
its old installation commands and upstream badges are not current guidance.
