package main

import (
	"config-analyzer/internal/server"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"config-analyzer/internal/analyzer"
	"config-analyzer/internal/models"
	"config-analyzer/internal/scanner"
)

func main() {

	silent := flag.Bool("s", false, "Do not exit with error code when issues are found")
	silentLong := flag.Bool("silent", false, "Do not exit with error code when issues are found")
	stdin := flag.Bool("stdin", false, "Read configuration from stdin")
	recursive := flag.Bool("r", false, "Recursively analyze a directory")
	httpAddr := flag.String("http", "", "Start HTTP server (e.g. :8080)")
	grpcAddr := flag.String("grpc", "", "Start gRPC server (e.g. :9090)")
	outputFmt := flag.String("output", "text", "Output format: text or json")
	timeout := flag.Duration("timeout", 10*time.Second, "Analysis timeout")

	flag.Parse()

	var logHandler slog.Handler

	if *httpAddr != "" || *grpcAddr != "" {
		logHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		logHandler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})
	}

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	isSilent := *silent || *silentLong

	a := analyzer.New()

	start(httpAddr, grpcAddr, a)

	if *recursive {
		recursionScan(a, isSilent)
	}

	readConfig(stdin, a, isSilent, timeout, outputFmt)

}

// start launches the server depending on selected mode
func start(httpAddr *string, grpcAddr *string, a *analyzer.Analyzer) {
	if *httpAddr != "" {
		slog.Info("Starting HTTP server", "port", *httpAddr)
		srv := server.NewHTTPServer(*httpAddr, a)
		if err := srv.Start(); err != nil {
			slog.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
		return
	}

	if *grpcAddr != "" {
		slog.Info("Starting gRPC server", "port", *httpAddr)
		srv := server.NewGRPCServer(*grpcAddr, a)
		if err := srv.Start(); err != nil {
			slog.Error("gRPC server error", "error", err)
			os.Exit(1)
		}
		return
	}
}

// recursionScan runs recursive scanning
func recursionScan(a *analyzer.Analyzer, isSilent bool) {
	if flag.NArg() < 1 {
		slog.Error("Usage: analyzer -r <directory>")
		os.Exit(1)
	}
	dir := flag.Arg(0)
	results, err := scanner.ScanDirectory(dir, a)
	if err != nil {
		slog.Error("Scanning error", "error", err)
		os.Exit(1)
	}

	hasIssues := false
	for _, result := range results {
		if result.HasIssues() {
			hasIssues = true
			fmt.Printf("\n=== %s ===\n", result.FilePath)
			printIssues(result.Issues)
		}
	}

	if !hasIssues {
		fmt.Println("No issues found.")
	} else if !isSilent {
		os.Exit(1)
	}
	return
}

// readConfig reads the configuration
func readConfig(stdin *bool, a *analyzer.Analyzer, isSilent bool, timeout *time.Duration, outputFmt *string) {

	var reader io.Reader
	var filePath string

	if *stdin {
		reader = os.Stdin
	} else {
		if flag.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "Usage: analyzer [flags] <path_to_file>")
			flag.PrintDefaults()
			os.Exit(1)
		}
		filePath = flag.Arg(0)
		file, err := os.Open(filePath) // Open stream
		if err != nil {
			slog.Error("File open error", "error", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	result, err := a.Analyze(ctx, models.AnalysisRequest{
		Reader:   reader,
		FilePath: filePath,
	})
	if err != nil {
		slog.Error("Analysis error", "error", err)
		os.Exit(1)
	}

	if *outputFmt == "json" {
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
		if result.HasIssues() && !isSilent {
			os.Exit(1)
		}
		return
	}

	if !result.HasIssues() {
		fmt.Println("No issues found.")
		return
	}

	if filePath != "" {
		fmt.Printf("=== Analysis: %s ===\n", filePath)
	}
	printIssues(result.Issues)
	fmt.Printf("\nTotal issues: %d\n", len(result.Issues))

	if !isSilent {
		os.Exit(1)
	}
}

// printIssues prints the list of found issues
func printIssues(issues []models.Issue) {
	for _, issue := range issues {
		fmt.Println(issue.String())
	}
}
