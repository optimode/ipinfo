package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/optimode/ipinfo/internal/api"
	"github.com/optimode/ipinfo/internal/cache"
	"github.com/optimode/ipinfo/internal/format"
	"github.com/optimode/ipinfo/internal/input"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagFormat      string
	flagSummary     bool
	flagFile        string
	flagAPIKey      string
	flagConcurrency int
	flagNoCache     bool
)

var rootCmd = &cobra.Command{
	Use:   "ipinfo [IP...]",
	Short: "IP geolocation lookup via ip-api.com Pro",
	Long: `ipinfo queries the ip-api.com Pro API for IP address geolocation and metadata.

Examples:
  ipinfo 8.8.8.8 1.1.1.1
  ipinfo -s 8.8.8.8
  echo "8.8.8.8" | ipinfo -s
  cat ips.txt | ipinfo --format table
  grep 'login' /var/log/mail.log | awk '{print $7}' | sort -u | ipinfo -s
  ipinfo --file ips.txt --format csv`,
	RunE: run,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&flagFormat, "format", "f", "table",
		"Output format: table, summary, json, csv")
	rootCmd.PersistentFlags().BoolVarP(&flagSummary, "summary", "s", false,
		"Shorthand for --format summary")
	rootCmd.PersistentFlags().StringVar(&flagFile, "file", "",
		"Input file with one IP per line")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "",
		"ip-api.com Pro API key (overrides config/env)")
	rootCmd.PersistentFlags().IntVarP(&flagConcurrency, "concurrency", "c", 0,
		"Parallel requests (0 = use config, default 5)")
	rootCmd.PersistentFlags().BoolVar(&flagNoCache, "no-cache", false,
		"Disable /24 subnet cache")

	viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key"))         //nolint
	viper.BindPFlag("concurrency", rootCmd.PersistentFlags().Lookup("concurrency")) //nolint
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/ipinfo")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("IPINFO")

	viper.SetDefault("concurrency", 5)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		}
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Resolve API key
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		return fmt.Errorf("API key is required. Set via --api-key, IPINFO_API_KEY env, or /etc/ipinfo/config.yaml")
	}

	// Resolve concurrency: flag > config default
	concurrency := flagConcurrency
	if concurrency <= 0 {
		concurrency = viper.GetInt("concurrency")
	}
	if concurrency <= 0 {
		concurrency = 5
	}

	// -s is shorthand for --format summary
	if flagSummary {
		flagFormat = format.FormatSummary
	}

	// Validate format
	switch flagFormat {
	case format.FormatTable, format.FormatSummary, format.FormatJSON, format.FormatCSV:
	default:
		return fmt.Errorf("invalid format %q, must be one of: table, summary, json, csv", flagFormat)
	}

	// Collect IPs from all available sources
	var ips []string

	// CLI arguments
	ips = append(ips, args...)

	// --file flag
	if flagFile != "" {
		f, err := os.Open(flagFile)
		if err != nil {
			return fmt.Errorf("cannot open file: %w", err)
		}
		defer f.Close()
		ips = append(ips, input.FromReader(f)...)
	}

	// stdin/pipe (always read if available)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		ips = append(ips, input.FromReader(os.Stdin)...)
	}

	if len(ips) == 0 {
		return fmt.Errorf("no input: provide IPs as arguments, --file, or via stdin pipe")
	}

	if len(ips) == 0 {
		return fmt.Errorf("no IPs to process")
	}

	// Setup
	client := api.New(apiKey)
	printer := format.New(flagFormat, os.Stdout)
	var subnetCache *cache.SubnetCache
	if !flagNoCache {
		subnetCache = cache.New()
	}

	// Print header once before concurrent output
	printer.PrintHeader()

	// Worker pool
	type job struct{ ip string }
	jobs := make(chan job, len(ips))
	for _, ip := range ips {
		jobs <- job{ip}
	}
	close(jobs)

	var wg sync.WaitGroup
	hasError := false
	var errMu sync.Mutex

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				processIP(j.ip, client, printer, subnetCache, &hasError, &errMu)
			}
		}()
	}

	wg.Wait()

	if hasError {
		os.Exit(1)
	}
	return nil
}

func processIP(ip string, client *api.Client, printer *format.Printer, subnetCache *cache.SubnetCache, hasError *bool, errMu *sync.Mutex) {
	// Check cache
	if subnetCache != nil {
		key := cache.Key(ip)
		if cached := subnetCache.Get(ip, key); cached != nil {
			fmt.Fprintf(os.Stderr, "(cached: %s → %s.x/24)\n", ip, key)
			printer.Print(cached)
			return
		}
	}

	// API lookup
	resp, err := client.Lookup(ip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		printer.PrintError(ip, err.Error())
		errMu.Lock()
		*hasError = true
		errMu.Unlock()
		return
	}

	// Store in cache
	if subnetCache != nil {
		key := cache.Key(ip)
		subnetCache.Set(key, resp)
	}

	printer.Print(resp)
}
