package netlink

// Focused regression for the local two-line decoder patch. Root-free and
// socket-free: bytes come from the real nl.XfrmUserpolicyInfo.Serialize, the
// same struct the kernel fills in a GETPOLICY reply.

import (
	"net"
	"testing"

	"github.com/iniwex5/netlink/nl"
)

func rawUserpolicyInfo(t *testing.T, family uint16, prefixSource, prefixDestination uint8,
	source, destination [16]byte,
) []byte {
	t.Helper()
	info := &nl.XfrmUserpolicyInfo{}
	info.Sel.Family = family
	info.Sel.PrefixlenS, info.Sel.PrefixlenD = prefixSource, prefixDestination
	info.Sel.Saddr = nl.XfrmAddress(source)
	info.Sel.Daddr = nl.XfrmAddress(destination)
	info.Index = 1
	info.Dir = uint8(XFRM_DIR_OUT)
	return append([]byte(nil), info.Serialize()...)
}

func TestParseXfrmPolicyPreservesSelectorFamily(t *testing.T) {
	var empty [16]byte
	ipv6Host := [16]byte{}
	copy(ipv6Host[:], net.ParseIP("2001:db8::1").To16())
	ipv4Host := [16]byte{}
	copy(ipv4Host[0:4], net.ParseIP("192.0.2.1").To4())

	for _, test := range []struct {
		name             string
		raw              []byte
		wantSource       string
		wantDestination  string
		wantSourceIsIPv6 bool
	}{
		{
			name: "ipv6-wildcard", raw: rawUserpolicyInfo(t, nl.FAMILY_V6, 0, 0, empty, empty),
			wantSource: "::/0", wantDestination: "::/0", wantSourceIsIPv6: true,
		},
		{
			name: "ipv4-wildcard", raw: rawUserpolicyInfo(t, nl.FAMILY_V4, 0, 0, empty, empty),
			wantSource: "0.0.0.0/0", wantDestination: "0.0.0.0/0", wantSourceIsIPv6: false,
		},
		{
			name: "ipv6-host-prefix", raw: rawUserpolicyInfo(t, nl.FAMILY_V6, 128, 128, ipv6Host, ipv6Host),
			wantSource: "2001:db8::1/128", wantDestination: "2001:db8::1/128", wantSourceIsIPv6: true,
		},
		{
			name: "ipv4-host-prefix", raw: rawUserpolicyInfo(t, nl.FAMILY_V4, 32, 32, ipv4Host, ipv4Host),
			wantSource: "192.0.2.1/32", wantDestination: "192.0.2.1/32", wantSourceIsIPv6: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			policy, err := parseXfrmPolicy(test.raw, FAMILY_ALL)
			if err != nil {
				t.Fatalf("parseXfrmPolicy: %v", err)
			}
			if policy.Src.String() != test.wantSource || policy.Dst.String() != test.wantDestination {
				t.Fatalf("selector=%s -> %s want %s -> %s", policy.Src, policy.Dst, test.wantSource, test.wantDestination)
			}
			if (policy.Src.IP.To4() == nil) != test.wantSourceIsIPv6 {
				t.Fatalf("source family=%s isIPv6=%t", policy.Src, policy.Src.IP.To4() == nil)
			}
		})
	}
}

func TestParseXfrmPolicyKeepsFamilyFilter(t *testing.T) {
	var empty [16]byte
	ipv6Wildcard := rawUserpolicyInfo(t, nl.FAMILY_V6, 0, 0, empty, empty)
	if _, err := parseXfrmPolicy(ipv6Wildcard, FAMILY_V4); err != familyError {
		t.Fatalf("FAMILY_V4 filter on IPv6 message: %v", err)
	}
	policy, err := parseXfrmPolicy(ipv6Wildcard, FAMILY_V6)
	if err != nil {
		t.Fatalf("FAMILY_V6 filter on IPv6 message: %v", err)
	}
	if policy.Src.String() != "::/0" || policy.Dst.String() != "::/0" {
		t.Fatalf("family-scoped selector=%s -> %s want ::/0", policy.Src, policy.Dst)
	}
}
