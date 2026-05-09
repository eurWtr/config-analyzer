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

	silent := flag.Bool("s", false, "Не выходить с ошибкой при наличии проблем")
	silentLong := flag.Bool("silent", false, "Не выходить с ошибкой при наличии проблем")
	stdin := flag.Bool("stdin", false, "Прочитать конфигурацию из stdin")
	recursive := flag.Bool("r", false, "Рекурсивный анализ директории")
	httpAddr := flag.String("http", "", "Запустить HTTP-сервер (например, :8080)")
	grpcAddr := flag.String("grpc", "", "Запустить gRPC-сервер (например, :9090)")
	outputFmt := flag.String("output", "text", "Формат вывода: text или json")
	timeout := flag.Duration("timeout", 10*time.Second, "Таймаут анализа")

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
		slog.Info("Запуск HTTP-сервера", "port", *httpAddr)
		srv := server.NewHTTPServer(*httpAddr, a)
		if err := srv.Start(); err != nil {
			slog.Error("Ошибка HTTP-сервера", "error", err)
			os.Exit(1)
		}
		return
	}

	if *grpcAddr != "" {
		slog.Info("Запуск gRPC-сервера", "port", *httpAddr)
		srv := server.NewGRPCServer(*grpcAddr, a)
		if err := srv.Start(); err != nil {
			slog.Error("Ошибка gRPC-сервера", "error", err)
			os.Exit(1)
		}
		return
	}
}

// recursionScan runs recursive scanning
func recursionScan(a *analyzer.Analyzer, isSilent bool) {
	if flag.NArg() < 1 {
		slog.Error("Использование: analyzer -r <директория>")
		os.Exit(1)
	}
	dir := flag.Arg(0)
	results, err := scanner.ScanDirectory(dir, a)
	if err != nil {
		slog.Error("Ошибка сканирования", "error", err)
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
		fmt.Println("Проблем не обнаружено.")
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
			fmt.Fprintln(os.Stderr, "Использование: analyzer [флаги] <путь_к_файлу>")
			flag.PrintDefaults()
			os.Exit(1)
		}
		filePath = flag.Arg(0)
		file, err := os.Open(filePath) // Open stream
		if err != nil {
			slog.Error("Ошибка открытия файла", "error", err)
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
		slog.Error("Ошибка анализа", "error", err)
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
		fmt.Println("Проблем не обнаружено.")
		return
	}

	if filePath != "" {
		fmt.Printf("=== Анализ: %s ===\n", filePath)
	}
	printIssues(result.Issues)
	fmt.Printf("\nВсего проблем: %d\n", len(result.Issues))

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
