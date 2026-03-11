package format

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/optimode/ipinfo/internal/api"
)

// Format constants.
const (
	FormatTable   = "table"
	FormatSummary = "summary"
	FormatJSON    = "json"
	FormatCSV     = "csv"
)

// Printer writes formatted output for IP responses.
type Printer struct {
	format    string
	out       io.Writer
	mu        sync.Mutex
	csvWriter *csv.Writer
	header    bool
}

// New creates a new Printer for the given format and writer.
func New(format string, out io.Writer) *Printer {
	p := &Printer{
		format: format,
		out:    out,
	}
	if format == FormatCSV {
		p.csvWriter = csv.NewWriter(out)
	}
	return p
}

// PrintHeader prints the header row for table and csv formats.
// Safe to call multiple times – only prints once.
func (p *Printer) PrintHeader() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.header {
		return
	}
	p.header = true
	switch p.format {
	case FormatTable:
		fmt.Fprintln(p.out, "| IP | Country | Region | City | ISP | Proxy | Hosting | Mobile |")
		fmt.Fprintln(p.out, "|----|---------|--------|------|-----|-------|---------|--------|")
	case FormatCSV:
		p.csvWriter.Write([]string{"ip", "country", "region", "city", "isp", "proxy", "hosting", "mobile"}) //nolint
		p.csvWriter.Flush()
	}
}

// Print outputs a single response in the configured format.
func (p *Printer) Print(r *api.Response) {
	p.mu.Lock()
	defer p.mu.Unlock()
	switch p.format {
	case FormatTable:
		fmt.Fprintf(p.out, "| %s | %s | %s | %s | %s | %v | %v | %v |\n",
			r.Query,
			str(r.CountryCode),
			str(r.RegionName),
			str(r.City),
			str(r.ISP),
			r.Proxy,
			r.Hosting,
			r.Mobile,
		)
	case FormatSummary:
		fmt.Fprintf(p.out, "%s\t%s\t%s\t%s\t%s\tproxy=%v\thosting=%v\tmobile=%v\n",
			r.Query,
			str(r.CountryCode),
			str(r.RegionName),
			str(r.City),
			str(r.ISP),
			r.Proxy,
			r.Hosting,
			r.Mobile,
		)
	case FormatJSON:
		b, _ := json.Marshal(r)
		fmt.Fprintln(p.out, string(b))
	case FormatCSV:
		p.csvWriter.Write([]string{ //nolint
			r.Query,
			str(r.CountryCode),
			str(r.RegionName),
			str(r.City),
			str(r.ISP),
			fmt.Sprintf("%v", r.Proxy),
			fmt.Sprintf("%v", r.Hosting),
			fmt.Sprintf("%v", r.Mobile),
		})
		p.csvWriter.Flush()
	}
}

// PrintError outputs an error row in the configured format.
func (p *Printer) PrintError(ip, message string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	switch p.format {
	case FormatTable:
		fmt.Fprintf(p.out, "| %s | ERROR | %s | - | - | - | - | - |\n", ip, message)
	case FormatSummary:
		fmt.Fprintf(p.out, "%s\tERROR\t%s\n", ip, message)
	case FormatJSON:
		b, _ := json.Marshal(map[string]string{"query": ip, "status": "error", "message": message})
		fmt.Fprintln(p.out, string(b))
	case FormatCSV:
		p.csvWriter.Write([]string{ip, "ERROR", message, "", "", "", "", ""}) //nolint
		p.csvWriter.Flush()
	}
}

func str(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
