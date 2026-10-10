# Contributing

English and 简体中文 reports are welcome. This repository maintains a retained
source snapshot; it does not automatically track or contain every upstream change.

## Before opening an issue

- Check the [README](README.md), [source notes](PROVENANCE.md) and existing issues.
- Include the exact maintained commit/module replacement, Go version, OS,
  architecture and relevant kernel/namespace conditions.
- Prefer a minimal synthetic decoder case. Do not change the host network to
  collect a reproduction unless you understand and can clean up your own state.
- Remove keys, tokens, full configuration, subscriber identifiers and raw captures.
  Vulnerabilities belong in the [private security channel](SECURITY.md).

## Pull requests

1. Branch from current `main`; keep each PR focused.
2. For a behavior fix, first add a regression that demonstrates the defect. Do not
   delete or weaken existing tests to make a change pass.
3. Run `gofmt` on changed Go files, then:

   ```sh
   go test -p 1 -race -shuffle=on -count=1 -timeout=120s ./...
   go vet -p 1 -composites=false ./...
   git diff --check
   ```

4. State exactly which tests ran, which needed privileges and which did not run.
   Synthetic tests do not prove live XFRM or device compatibility.
5. Preserve the original module path, license texts, copyright and attribution.
   Dependency updates or changes to copied material need explicit review of
   provenance and redistribution obligations.
6. Use descriptive Conventional Commits and your existing legitimate Git identity.
   Do not add generated tool banners, fabricated authors or attribution trailers.

The PR template requests the reviewed head/base, compatibility effects and evidence.
The composites analyzer has a documented legacy exception in the README; it is
disabled in CI, not silently treated as a passing full-vet check. Run full vet
when working on related source and report new warnings rather than hiding them.
Do not rewrite shared history or force-push over other contributors' work. A
documentation-only maintenance commit does not require consumers to upgrade pins.

## Privileged testing

Default CI only covers the retained synthetic tests. Test kernel mutations in a
disposable network namespace or VM, document kernel/capabilities, and clean up
owned interfaces, routes, rules and XFRM state. Never flush another owner's state.

## Conduct

Follow the [code of conduct](CODE_OF_CONDUCT.md). Reviews should explain technical
tradeoffs and evidence, not make assumptions about a contributor's identity.
