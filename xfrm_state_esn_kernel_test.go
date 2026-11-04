//go:build linux

package netlink

import (
	"context"
	"net"
	"os"
	"os/exec"
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestXfrmStateESNKernelReadback(t *testing.T) {
	if os.Getenv("VOWIFI_NETLINK_KERNEL_TEST") != "1" {
		t.Skip("requires an explicitly isolated network namespace")
	}
	parent := os.Getenv("VOWIFI_NETLINK_PARENT_NS")
	current, err := os.Readlink("/proc/self/ns/net")
	if err != nil || parent == "" || current == parent {
		t.Fatal("kernel test must not run in the parent network namespace")
	}
	ip := func(args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := exec.CommandContext(ctx, "ip", args...).Run(); err != nil {
			t.Fatalf("isolated iproute2 setup failed: %v", err)
		}
	}
	if states, err := XfrmStateList(FAMILY_ALL); err != nil || len(states) != 0 {
		t.Fatal("isolated namespace must start without states")
	}
	if policies, err := XfrmPolicyList(FAMILY_ALL); err != nil || len(policies) != 0 {
		t.Fatal("isolated namespace must start without policies")
	}
	query := &XfrmState{Src: net.ParseIP("192.0.2.1"), Dst: net.ParseIP("192.0.2.2"), Proto: XFRM_PROTO_ESP, Spi: 0x10203040}
	// Independent iproute2 construction seeds nonzero high words; the library's
	// Add path intentionally does not restore a previously observed replay state.
	ip("xfrm", "state", "add", "src", "192.0.2.1", "dst", "192.0.2.2", "proto", "esp", "spi", "0x10203040",
		"reqid", "77", "mode", "tunnel", "flag", "esn", "replay-window", "64",
		"aead", "rfc4106(gcm(aes))", "0x000102030405060708090a0b0c0d0e0fa0a1a2a3", "128",
		"replay-oseq", "4294967294", "replay-oseq-hi", "1", "replay-seq", "7", "replay-seq-hi", "3")
	t.Cleanup(func() {
		if err := XfrmStateDel(query); err != nil {
			t.Errorf("exact synthetic SA cleanup failed: %v", err)
		}
		if states, err := XfrmStateList(FAMILY_ALL); err != nil || len(states) != 0 {
			t.Error("synthetic SA remains after cleanup")
		}
		if policies, err := XfrmPolicyList(FAMILY_ALL); err != nil || len(policies) != 0 {
			t.Error("unexpected policy remains after cleanup")
		}
	})
	state, err := XfrmStateGet(query)
	if err != nil {
		t.Fatalf("real ESN state GET failed: %v", err)
	}
	replay := state.ReplayESN
	if !state.ESN || state.ReplayWindow != 64 || replay == nil || replay.OSeq != 0xfffffffe || replay.OSeqHi != 1 || replay.Seq != 7 || replay.SeqHi != 3 || replay.ReplayWindow != 64 || len(replay.BitMap) != 2 {
		t.Fatal("real GET did not retain ESN window, high words, and bitmap")
	}
	want := *replay
	want.BitMap = slices.Clone(replay.BitMap)
	state.Limits.TimeSoft, state.Limits.TimeHard = 60, 120
	// Readback fields are not requests to rewind a live SA's sequence state.
	state.ReplayESN = &XfrmReplayStateESN{}
	if err := XfrmStateUpdate(state); err != nil {
		t.Fatalf("same-SA lifetime update failed: %v", err)
	}
	after, err := XfrmStateGet(query)
	if err != nil || !after.ESN || after.ReplayWindow != 64 || after.ReplayESN == nil || !reflect.DeepEqual(*after.ReplayESN, want) || after.Limits.TimeSoft != 60 || after.Limits.TimeHard != 120 {
		t.Fatal("same-SA update lost replay state or requested lifetime")
	}
}
