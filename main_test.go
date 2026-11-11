package main

import (
	"bytes"
	"testing"
)

func TestExample(t *testing.T) {
	code, out, errOut := runFixture(t, "testdata/example.json", "2026-11-10T12:00:00Z")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	want := "next_high 2026-11-10T15:00:00Z 4.2\nnext_low 2026-11-10T21:30:00Z 0.4\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestOnlyLow(t *testing.T) {
	code, out, errOut := runFixture(t, "testdata/only_low.json", "2026-11-10T12:00:00Z")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	want := "next_low 2026-11-10T21:30:00Z 0.4\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestNothingAfter(t *testing.T) {
	code, out, errOut := runFixture(t, "testdata/nothing_after.json", "2026-11-10T12:00:00Z")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	want := "no upcoming tide\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestLowBeforeHighPrintsHighFirst(t *testing.T) {
	code, out, errOut := runFixture(t, "testdata/low_before_high.json", "2026-11-10T12:00:00Z")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	want := "next_high 2026-11-10T18:00:00Z 4.2\nnext_low 2026-11-10T13:00:00Z 0.4\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestMissingFlags(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run(nil, &out, &errOut); code != 2 {
		t.Fatalf("no args: exit %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q", out.String())
	}

	out.Reset()
	errOut.Reset()
	if code := run([]string{"--file", "testdata/example.json"}, &out, &errOut); code != 2 {
		t.Fatalf("missing --now: exit %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q", out.String())
	}

	out.Reset()
	errOut.Reset()
	if code := run([]string{"--now", "2026-11-10T12:00:00Z"}, &out, &errOut); code != 2 {
		t.Fatalf("missing --file: exit %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q", out.String())
	}
}

func TestBadJSON(t *testing.T) {
	code, out, _ := runFixture(t, "testdata/bad.json", "2026-11-10T12:00:00Z")
	if code != 2 {
		t.Fatalf("exit %d", code)
	}
	if out != "" {
		t.Fatalf("stdout %q", out)
	}
}

func TestBadKind(t *testing.T) {
	code, out, _ := runFixture(t, "testdata/bad_kind.json", "2026-11-10T12:00:00Z")
	if code != 2 {
		t.Fatalf("exit %d", code)
	}
	if out != "" {
		t.Fatalf("stdout %q", out)
	}
}

func TestBadNow(t *testing.T) {
	code, out, _ := runFixture(t, "testdata/example.json", "not-a-time")
	if code != 2 {
		t.Fatalf("exit %d", code)
	}
	if out != "" {
		t.Fatalf("stdout %q", out)
	}
}

func runFixture(t *testing.T, file, now string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run([]string{"--file", file, "--now", now}, &out, &errOut)
	return code, out.String(), errOut.String()
}
