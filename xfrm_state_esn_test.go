//go:build linux

package netlink

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/iniwex5/netlink/nl"
)

// These response bytes do not use writeReplayEsn or a replay decoder helper.
func replayWords(words ...uint32) []byte {
	raw := make([]byte, len(words)*4)
	for i, word := range words {
		binary.NativeEndian.PutUint32(raw[i*4:], word)
	}
	return raw
}

func stateReplayMessage(esn bool, window uint8, payloads ...[]byte) []byte {
	info := &nl.XfrmUsersaInfo{}
	info.Family, info.Sel.Family = nl.FAMILY_V4, nl.FAMILY_V4
	info.Id.Proto, info.Id.Spi = 50, nl.Swap32(0x10203040)
	info.ReplayWindow = window
	if esn {
		info.Flags = nl.XFRM_STATE_ESN
	}
	raw := append([]byte(nil), info.Serialize()...)
	for _, payload := range payloads {
		raw = append(raw, nl.NewRtAttr(nl.XFRMA_REPLAY_ESN_VAL, payload).Serialize()...)
	}
	return raw
}

// Project public fields without requiring a new library declaration to compile
// against the old decoder. The RED must reach parseXfrmState, not fail to build.
type replayESNProjection struct {
	OSeq, Seq, OSeqHi, SeqHi, ReplayWindow uint32
	BitMap                                 []uint32
}

func stateReplayProjection(t *testing.T, state *XfrmState) *replayESNProjection {
	t.Helper()
	raw, err := json.Marshal(state)
	if err != nil {
		t.Fatal("cannot project parsed replay state")
	}
	var projection struct{ ReplayESN *replayESNProjection }
	if json.Unmarshal(raw, &projection) != nil {
		t.Fatal("cannot read parsed replay projection")
	}
	return projection.ReplayESN
}

func TestParseXfrmStateReadsESNAndExtendedReplay(t *testing.T) {
	for _, tc := range []struct {
		name string
		esn  bool
		want replayESNProjection
	}{
		{"esn64", true, replayESNProjection{0xfffffffe, 7, 1, 3, 64, []uint32{1, 0x80000000}}},
		{"extended-window-without-esn", false, replayESNProjection{9, 7, 0, 0, 128, []uint32{1, 2, 3, 4}}},
		{"esn-outbound-zero-window", true, replayESNProjection{9, 7, 2, 3, 0, nil}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			words := []uint32{uint32(len(tc.want.BitMap)), tc.want.OSeq, tc.want.Seq, tc.want.OSeqHi, tc.want.SeqHi, tc.want.ReplayWindow}
			words = append(words, tc.want.BitMap...)
			wire := stateReplayMessage(tc.esn, 0, replayWords(words...))
			original := bytes.Clone(wire)
			state, err := parseXfrmState(wire, FAMILY_ALL)
			if err != nil {
				t.Fatalf("valid replay response rejected: %v", err)
			}
			if state.ESN != tc.esn || state.ReplayWindow != int(tc.want.ReplayWindow) {
				t.Fatal("actual decoder lost the ESN flag or extended replay window")
			}
			got := stateReplayProjection(t, state)
			if got == nil || got.OSeq != tc.want.OSeq || got.Seq != tc.want.Seq || got.OSeqHi != tc.want.OSeqHi || got.SeqHi != tc.want.SeqHi || got.ReplayWindow != tc.want.ReplayWindow || !slices.Equal(got.BitMap, tc.want.BitMap) {
				t.Fatal("actual decoder lost replay sequence words or bitmap")
			}
			if !bytes.Equal(wire, original) {
				t.Fatal("response decoding mutated its input")
			}
			if len(tc.want.BitMap) != 0 {
				wire[len(wire)-1] ^= 0xff
				if !slices.Equal(stateReplayProjection(t, state).BitMap, tc.want.BitMap) {
					t.Fatal("replay bitmap aliases the received message")
				}
			}
		})
	}
}

func TestParseXfrmStateKeepsLegacyReplay(t *testing.T) {
	wire := stateReplayMessage(false, 32)
	wire = append(wire, nl.NewRtAttr(nl.XFRMA_REPLAY_VAL, replayWords(9, 7, 5)).Serialize()...)
	state, err := parseXfrmState(wire, FAMILY_ALL)
	if err != nil || state.ESN || state.ReplayWindow != 32 || state.Replay == nil || *state.Replay != (XfrmReplayState{OSeq: 9, Seq: 7, BitMap: 5}) {
		t.Fatal("ordinary non-ESN replay compatibility changed")
	}
	if stateReplayProjection(t, state) != nil {
		t.Fatal("legacy replay was invented as extended replay")
	}
}

func TestParseXfrmStateRejectsMalformedESNReplay(t *testing.T) {
	valid := replayWords(2, 9, 7, 1, 3, 64, 1, 0)
	cases := []struct {
		name     string
		payloads [][]byte
	}{
		{"missing-bitmap", [][]byte{valid[:24]}},
		{"extra-bitmap-word", [][]byte{append(bytes.Clone(valid), 0, 0, 0, 0)}},
		{"partial-word", [][]byte{valid[:29]}},
		{"unbounded-bitmap-length", [][]byte{replayWords(0xffffffff, 0, 0, 0, 0, 0)}},
		{"window-outside-bitmap", [][]byte{replayWords(1, 0, 0, 0, 0, 33, 0)}},
		{"window-without-bitmap", [][]byte{replayWords(0, 0, 0, 0, 0, 1)}},
		{"duplicate-context", [][]byte{valid, valid}},
	}
	for size := 0; size < 24; size++ {
		cases = append(cases, struct {
			name     string
			payloads [][]byte
		}{fmt.Sprintf("short-header-%02d", size), [][]byte{valid[:size]}})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parseXfrmState(stateReplayMessage(true, 0, tc.payloads...), FAMILY_ALL); err == nil {
				t.Fatal("malformed extended replay response was accepted")
			}
		})
	}
}

func TestParseXfrmStateKeepsFamilyFilterWithESN(t *testing.T) {
	wire := stateReplayMessage(true, 0, replayWords(1, 9, 7, 1, 3, 32, 1))
	if _, err := parseXfrmState(wire, FAMILY_V6); err != familyError {
		t.Fatal("ESN decoding changed the address-family filter")
	}
}
