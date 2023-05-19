package marc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewLeader(t *testing.T) {
	t.Parallel()

	leaderBytes := []byte("01848nam a2200385 i 4500")

	want := Leader{
		raw:           leaderBytes,
		dataOffset:    385,
		Status:        byte('n'),
		Type:          byte('a'),
		BibLevel:      byte('m'),
		Control:       byte(' '),
		EncodingLevel: byte(' '),
		Form:          byte('i'),
		Multipart:     byte(' '),
	}

	got, _ := NewLeader(leaderBytes)

	opt := cmp.AllowUnexported(Leader{})

	if !cmp.Equal(want, got, opt) {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestNewLeader_ErrorsOnBadOffset(t *testing.T) {
	t.Parallel()

	leader_bytes := []byte("ZZZZZnamZa22ZZZZZzZZ4500")
	_, err := NewLeader(leader_bytes)
	if err == nil {
		t.Error("want error for invalid input")
	}
}

func TestNewLeader_ErrorsOnShortLeader(t *testing.T) {
	t.Parallel()

	leader_bytes := []byte("01848nam a2200385 i 450")
	_, err := NewLeader(leader_bytes)
	if err == nil {
		t.Error("want error for invalid input")
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	want := "=LDR  01848nam a2200385 i 4500"

	leader_bytes := []byte("01848nam a2200385 i 4500")
	l, _ := NewLeader(leader_bytes)
	got := l.String()

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestRaw(t *testing.T) {
	t.Parallel()

	want := "01848nam a2200385 i 4500"

	leader_bytes := []byte("01848nam a2200385 i 4500")
	l, _ := NewLeader(leader_bytes)
	got := l.Raw()

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
