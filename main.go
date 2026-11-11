package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("tideglass", flag.ContinueOnError)
	fs.SetOutput(stderr)
	file := fs.String("file", "", "predictions JSON file")
	nowStr := fs.String("now", "", "current time, RFC3339 UTC")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *file == "" || *nowStr == "" || fs.NArg() != 0 {
		fmt.Fprintln(stderr, "usage: tideglass --file predictions.json --now 2026-11-10T12:00:00Z")
		return 2
	}

	now, err := parseUTC(*nowStr)
	if err != nil {
		fmt.Fprintf(stderr, "bad --now: %v\n", err)
		return 2
	}

	events, err := loadPredictions(*file)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return 2
	}

	highIdx, lowIdx := -1, -1
	for i, ev := range events {
		if !ev.when.After(now) {
			continue
		}
		if ev.kind == "H" && (highIdx < 0 || ev.when.Before(events[highIdx].when)) {
			highIdx = i
		}
		if ev.kind == "L" && (lowIdx < 0 || ev.when.Before(events[lowIdx].when)) {
			lowIdx = i
		}
	}

	if highIdx < 0 && lowIdx < 0 {
		fmt.Fprintln(stdout, "no upcoming tide")
		return 0
	}
	if highIdx >= 0 {
		ev := events[highIdx]
		fmt.Fprintf(stdout, "next_high %s %.1f\n", ev.when.UTC().Format(time.RFC3339), ev.height)
	}
	if lowIdx >= 0 {
		ev := events[lowIdx]
		fmt.Fprintf(stdout, "next_low %s %.1f\n", ev.when.UTC().Format(time.RFC3339), ev.height)
	}
	return 0
}

type event struct {
	when   time.Time
	height float64
	kind   string
}

func loadPredictions(path string) ([]event, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	trimmed := bytes.TrimSpace(b)
	if len(trimmed) == 0 || trimmed[0] != '[' {
		return nil, fmt.Errorf("bad JSON")
	}

	var raw []struct {
		Time     *string  `json:"time"`
		HeightFt *float64 `json:"height_ft"`
		Kind     *string  `json:"kind"`
	}
	if err := json.Unmarshal(trimmed, &raw); err != nil {
		return nil, fmt.Errorf("bad JSON")
	}

	events := make([]event, 0, len(raw))
	for _, row := range raw {
		if row.Time == nil || row.HeightFt == nil || row.Kind == nil {
			return nil, fmt.Errorf("bad JSON")
		}
		if *row.Kind != "H" && *row.Kind != "L" {
			return nil, fmt.Errorf("bad kind")
		}
		when, err := parseUTC(*row.Time)
		if err != nil {
			return nil, fmt.Errorf("bad time: %v", err)
		}
		events = append(events, event{when: when, height: *row.HeightFt, kind: *row.Kind})
	}
	return events, nil
}

func parseUTC(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	_, offset := t.Zone()
	if offset != 0 {
		return time.Time{}, fmt.Errorf("time must be UTC")
	}
	return t.UTC(), nil
}
