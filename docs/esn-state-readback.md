# ESN state readback

This change is developed on `fix/esn-state-readback-20261105` from
`49ba356baca7c7e3770e6ff3f65caea20ee9b468`, maintained by `lilyzhaun`.

The scope is Linux XFRM state response decoding: the ESN flag, replay
window, sequence words, and replay bitmap must come from the response.
Ordinary non-ESN state decoding and request construction remain unchanged.
Reading dynamic replay state does not authorize restoring or rewriting it.

The existing decoder has a direct regression: a response containing the
ESN flag and a 64-packet extended replay window is returned as non-ESN
with a zero replay window. The non-ESN control retains its original window.

Tests use independently assembled Linux UAPI messages. This work uses the
retained implementation and public interface semantics, without importing
another project's implementation or configuration. The module path,
Apache-2.0 license, dependencies, and existing provenance remain unchanged.

`XfrmState.ReplayESN` exposes the complete returned replay context. Its
sequence words and bitmap are readback-only: Add and Update do not use
them to restore counters. The `ESN` flag still comes from the actual
state flag, not from the presence of a large replay window.

Malformed, truncated, duplicate, or inconsistent replay attributes return
an error. The returned bitmap owns its storage independently of the input
message. Legacy replay responses keep their existing representation.

The optional `TestXfrmStateESNKernelReadback` requires iproute2, XFRM
privileges, `VOWIFI_NETLINK_KERNEL_TEST=1`, and
`VOWIFI_NETLINK_PARENT_NS` set to the parent process's network namespace
identifier. Run it only in a fresh `unshare --net` child, not in the host
namespace. It seeds synthetic sequence words with iproute2, reads them
through the real library GET path, updates only lifetime, verifies replay
preservation, and deletes the exact synthetic SA. Ordinary CI skips this
explicitly privileged test and still runs the message-level regressions.

Message-level regression and isolated kernel readback cover this decoder's
scope. They do not establish negotiation, an ESP data plane, Android support,
or a state-restoration API. Kernel behavior remains platform-specific; passing
these tests is not full upstream or cross-kernel compatibility certification.
