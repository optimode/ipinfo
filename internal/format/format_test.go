package format

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/optimode/ipinfo/internal/api"
)

func sampleResponse() *api.Response {
	return &api.Response{
		Query:       "8.8.8.8",
		Status:      "success",
		CountryCode: "US",
		RegionName:  "California",
		City:        "Mountain View",
		ISP:         "Google LLC",
		Proxy:       false,
		Hosting:     true,
		Mobile:      false,
	}
}

func TestPrinter_Table(t *testing.T) {
	var buf bytes.Buffer
	p := New(FormatTable, &buf)
	p.PrintHeader()
	p.Print(sampleResponse())

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + separator + row), got %d:\n%s", len(lines), out)
	}
	if !strings.Contains(lines[0], "| IP |") {
		t.Errorf("header missing: %s", lines[0])
	}
	if !strings.Contains(lines[1], "|----|") {
		t.Errorf("separator missing: %s", lines[1])
	}
	if !strings.Contains(lines[2], "8.8.8.8") {
		t.Errorf("data row missing IP: %s", lines[2])
	}
	if !strings.Contains(lines[2], "Google LLC") {
		t.Errorf("data row missing ISP: %s", lines[2])
	}
}

func TestPrinter_Summary(t *testing.T) {
	var buf bytes.Buffer
	p := New(FormatSummary, &buf)
	p.PrintHeader()
	p.Print(sampleResponse())

	out := buf.String()
	if !strings.Contains(out, "8.8.8.8\tUS\tCalifornia\tMountain View\tGoogle LLC") {
		t.Errorf("unexpected summary output: %s", out)
	}
	if !strings.Contains(out, "hosting=true") {
		t.Errorf("missing hosting flag: %s", out)
	}
}

func TestPrinter_JSON(t *testing.T) {
	var buf bytes.Buffer
	p := New(FormatJSON, &buf)
	p.PrintHeader()
	p.Print(sampleResponse())

	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, buf.String())
	}
	if result["query"] != "8.8.8.8" {
		t.Errorf("query = %v, want 8.8.8.8", result["query"])
	}
}

func TestPrinter_CSV(t *testing.T) {
	var buf bytes.Buffer
	p := New(FormatCSV, &buf)
	p.PrintHeader()
	p.Print(sampleResponse())

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (header + row), got %d:\n%s", len(lines), buf.String())
	}
	if lines[0] != "ip,country,region,city,isp,proxy,hosting,mobile" {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
	if !strings.Contains(lines[1], "8.8.8.8") {
		t.Errorf("data row missing IP: %s", lines[1])
	}
}

func TestPrinter_PrintError_Table(t *testing.T) {
	var buf bytes.Buffer
	p := New(FormatTable, &buf)
	p.PrintError("1.2.3.4", "some error")

	out := buf.String()
	if !strings.Contains(out, "1.2.3.4") {
		t.Errorf("error row missing IP: %s", out)
	}
	if !strings.Contains(out, "ERROR") {
		t.Errorf("error row missing ERROR: %s", out)
	}
}

func TestPrinter_PrintError_JSON(t *testing.T) {
	var buf bytes.Buffer
	p := New(FormatJSON, &buf)
	p.PrintError("1.2.3.4", "some error")

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["status"] != "error" {
		t.Errorf("status = %q, want %q", result["status"], "error")
	}
	if result["query"] != "1.2.3.4" {
		t.Errorf("query = %q, want %q", result["query"], "1.2.3.4")
	}
}

func TestPrinter_HeaderOnlyOnce(t *testing.T) {
	var buf bytes.Buffer
	p := New(FormatTable, &buf)
	p.PrintHeader()
	p.PrintHeader()
	p.PrintHeader()

	out := buf.String()
	count := strings.Count(out, "| IP |")
	if count != 1 {
		t.Errorf("header printed %d times, want 1", count)
	}
}

func TestStr(t *testing.T) {
	if str("") != "-" {
		t.Error("empty string should return -")
	}
	if str("   ") != "-" {
		t.Error("whitespace-only should return -")
	}
	if str("hello") != "hello" {
		t.Error("non-empty should return as-is")
	}
}
