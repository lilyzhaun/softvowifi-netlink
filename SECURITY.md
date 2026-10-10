# Security policy

## Scope and supported revisions

Security maintenance targets current `main` and explicitly reviewed fixes. This
repository has no published LTS support window or guarantee for old snapshots.
A consuming application that pins an older commit must separately review and
adopt fixes; documentation changes do not update its dependency automatically.

This is a low-level network library, not a security certification of a device,
kernel or complete VoWiFi product. Source and licensing limitations are recorded
in [PROVENANCE.md](PROVENANCE.md).

## Report privately

Use GitHub's **Security → Report a vulnerability**:

https://github.com/lilyzhaun/softvowifi-netlink/security/advisories/new

Do not disclose an unpatched vulnerability in a public issue or PR. If the private
reporting entry is temporarily unavailable, use GitHub's support/reporting tools;
do not replace it with a public dump of sensitive material.

Provide an affected commit, impact, relevant platform/privilege assumptions and a
minimal synthetic reproduction. Share only the information necessary to reproduce.
No response-time SLA is promised; disclosure timing should be coordinated with
the maintainer before publishing exploit details.

## Keep sensitive material out of reports

Never submit passwords, tokens, private keys, XFRM/AKA key material, real SIM
identifiers, complete device configuration, message bodies or unsanitized packet
captures. Use documentation-range IP addresses and synthetic values. Kernel
integration tests should be isolated and must clean up only their owned state.
