package main

import (
	"fmt"
	"os"

	ctx "github.com/DilemaFixer/Cmd/context"
	p "github.com/DilemaFixer/Cmd/parser"
	rtr "github.com/DilemaFixer/Cmd/router"
)

// Example 4: File Manager with Type Validation and Error Handling
// This example demonstrates proper type validation, required fields, and custom error handling
// Usage examples:
//   myapp file copy --source=/home/user/file.txt --destination=/backup/ --buffer-size=8192 --verify
//   myapp file compress --input=/data/logs --output=/compressed/logs.tar.gz --level=9 --threads=4
//   myapp file sync --source=/local/data --target=/remote/backup --batch-size=100 --retry=3
//   myapp file search --pattern="*.log" --directory=/var/log --max-depth=3 --case-sensitive

func main() {
	// Example with type validation and required parameters
	input := "file copy --source=/home/user/document.pdf --destination=/backup/ --buffer-size=8192 --threads=2 --verify"

	parsedInput, err := p.ParseInput(input)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	context := ctx.NewContext(parsedInput)
	iterator := rtr.NewRoutingIterator(context)
	router := rtr.NewRouter()

	// Set up custom error handler for better error messages
	router.CustomErrorHandler(customErrorHandler)

	// Define file management commands with strict type validation
	router.NewCmd("file").
		Endpoint("copy").
		RequiredString("source").      // Required source path
		RequiredString("destination"). // Required destination path
		IntOption("buffer-size").      // Buffer size in bytes (default: 4096)
		IntOption("threads").          // Number of parallel threads
		BoolOption("verify").          // Verify copy integrity
		BoolOption("preserve-attrs").  // Preserve file attributes
		BoolOption("overwrite").       // Overwrite existing files
		FloatOption("throttle").       // Throttle speed in MB/s
		Handler(fileCopyHandler).
		Build().
		Endpoint("compress").
		Description("Compress files and directories").
		RequiredString("input").  // Required input path
		RequiredString("output"). // Required output file
		RequiredInt("level").     // Required compression level (1-9)
		IntOption("threads").     // Number of threads for compression
		BoolOption("recursive").  // Compress directories recursively
		StringOption("format").   // Compression format (zip, tar.gz, 7z)
		FloatOption("max-size").  // Maximum file size to include (MB)
		Handler(fileCompressHandler).
		Build().
		Endpoint("sync").
		Description("Synchronize directories").
		RequiredString("source"). // Required source directory
		RequiredString("target"). // Required target directory
		IntOption("batch-size").  // Files to process in each batch
		IntOption("retry").       // Number of retry attempts
		FloatOption("timeout").   // Timeout for operations in seconds
		BoolOption("delete").     // Delete files not in source
		BoolOption("dry-run").    // Preview changes without applying
		StringOption("exclude").  // Pattern for files to exclude
		Handler(fileSyncHandler).
		Build().
		Endpoint("search").
		Description("Search for files with advanced criteria").
		RequiredString("pattern").    // Required search pattern
		RequiredString("directory").  // Required search directory
		IntOption("max-depth").       // Maximum search depth
		IntOption("min-size").        // Minimum file size in bytes
		IntOption("max-size").        // Maximum file size in bytes
		BoolOption("case-sensitive"). // Case sensitive search
		BoolOption("include-hidden"). // Include hidden files
		StringOption("type").         // File type filter (file, dir, link)
		Handler(fileSearchHandler).
		Build().
		Endpoint("monitor").
		Description("Monitor directory for changes").
		RequiredString("directory").   // Required directory to monitor
		IntOption("interval").         // Check interval in seconds
		FloatOption("size-threshold"). // Size change threshold in MB
		BoolOption("recursive").       // Monitor subdirectories
		StringOption("output").        // Output log file
		Handler(fileMonitorHandler).
		Build().
		Register()

	// Execute the command
	router.Route(*context, iterator)
}

// Custom error handler that provides detailed error information
func customErrorHandler(err error, ctx ctx.Context) {
	fmt.Printf("❌ Command failed: %s\n", ctx.GetCommand())

	// Show the command path that failed
	subcommands := ctx.GetSubcommandsAsArr()
	if len(subcommands) > 0 {
		fmt.Printf("📍 Subcommand: %s\n", subcommands[len(subcommands)-1])
	}

	// Show provided flags for debugging
	fmt.Printf("🏷️  Provided flags: %v\n", ctx.GetFlagsKeysAsArr())

	// Display the error with better formatting
	fmt.Printf("🚫 Error: %v\n", err)

	// Provide helpful hints based on error type
	errorMsg := err.Error()
	if contains(errorMsg, "Required") && contains(errorMsg, "not exist") {
		fmt.Println("💡 Hint: Make sure all required flags are provided with --flag-name=value")
	} else if contains(errorMsg, "invalid value") {
		fmt.Println("💡 Hint: Check that numeric values are valid integers or floats")
	} else if contains(errorMsg, "Bool have value") {
		fmt.Println("💡 Hint: Boolean flags should not have values, use --flag instead of --flag=value")
	}

	fmt.Println("\n📚 Use --help for command usage information")
	os.Exit(1)
}

// Helper function to check if string contains substring
func contains(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func fileCopyHandler(ctx ctx.Context) error {
	fmt.Println("📁 File Copy Operation")

	// Get required parameters
	source, err := ctx.GetValueAsString("source")
	if err != nil {
		return fmt.Errorf("source path is required: %v", err)
	}

	destination, err := ctx.GetValueAsString("destination")
	if err != nil {
		return fmt.Errorf("destination path is required: %v", err)
	}

	fmt.Printf("📂 Source: %s\n", source)
	fmt.Printf("📁 Destination: %s\n", destination)

	// Handle optional integer parameters with validation
	bufferSize := 4096 // default
	if ctx.IsFlagExist("buffer-size") {
		if size, err := ctx.GetValueAsInt("buffer-size"); err != nil {
			return fmt.Errorf("invalid buffer-size: %v", err)
		} else if size <= 0 {
			return fmt.Errorf("buffer-size must be positive, got: %d", size)
		} else {
			bufferSize = size
		}
	}

	threads := 1 // default
	if ctx.IsFlagExist("threads") {
		if t, err := ctx.GetValueAsInt("threads"); err != nil {
			return fmt.Errorf("invalid threads count: %v", err)
		} else if t < 1 || t > 32 {
			return fmt.Errorf("threads must be between 1 and 32, got: %d", t)
		} else {
			threads = t
		}
	}

	// Handle optional float parameter
	var throttleSpeed float64
	if ctx.IsFlagExist("throttle") {
		if speed, err := ctx.GetValueAsFloat64("throttle"); err != nil {
			return fmt.Errorf("invalid throttle speed: %v", err)
		} else if speed <= 0 {
			return fmt.Errorf("throttle speed must be positive, got: %.2f", speed)
		} else {
			throttleSpeed = speed
		}
	}

	// Handle boolean flags
	verify := ctx.GetValueAsBool("verify")
	preserveAttrs := ctx.GetValueAsBool("preserve-attrs")
	overwrite := ctx.GetValueAsBool("overwrite")

	fmt.Printf("⚙️  Buffer Size: %d bytes\n", bufferSize)
	fmt.Printf("🧵 Threads: %d\n", threads)

	if throttleSpeed > 0 {
		fmt.Printf("🐌 Throttle: %.2f MB/s\n", throttleSpeed)
	}

	fmt.Printf("✅ Verify: %v\n", verify)
	fmt.Printf("📋 Preserve Attributes: %v\n", preserveAttrs)
	fmt.Printf("📝 Overwrite: %v\n", overwrite)

	// Simulate copy operation
	fmt.Println("\n🔄 Copying files...")
	fmt.Println("  → Analyzing source...")
	fmt.Println("  → Preparing destination...")
	fmt.Println("  → Copying data...")

	if verify {
		fmt.Println("  → Verifying integrity...")
	}

	fmt.Println("✅ Copy operation completed successfully!")
	return nil
}

func fileCompressHandler(ctx ctx.Context) error {
	fmt.Println("🗜️ File Compression")

	// Get and validate required parameters
	input, err := ctx.GetValueAsString("input")
	if err != nil {
		return fmt.Errorf("input path is required: %v", err)
	}

	output, err := ctx.GetValueAsString("output")
	if err != nil {
		return fmt.Errorf("output path is required: %v", err)
	}

	level, err := ctx.GetValueAsInt("level")
	if err != nil {
		return fmt.Errorf("compression level is required: %v", err)
	}

	// Validate compression level range
	if level < 1 || level > 9 {
		return fmt.Errorf("compression level must be between 1-9, got: %d", level)
	}

	fmt.Printf("📂 Input: %s\n", input)
	fmt.Printf("📄 Output: %s\n", output)
	fmt.Printf("📊 Compression Level: %d\n", level)

	// Handle optional parameters
	threads := 1
	if ctx.IsFlagExist("threads") {
		if t, err := ctx.GetValueAsInt("threads"); err != nil {
			return fmt.Errorf("invalid threads: %v", err)
		} else {
			threads = t
		}
	}

	format := ctx.GetValueOrDefault("format", "tar.gz")
	recursive := ctx.GetValueAsBool("recursive")

	var maxSize float64
	if ctx.IsFlagExist("max-size") {
		if size, err := ctx.GetValueAsFloat64("max-size"); err != nil {
			return fmt.Errorf("invalid max-size: %v", err)
		} else if size <= 0 {
			return fmt.Errorf("max-size must be positive, got: %.2f", size)
		} else {
			maxSize = size
		}
	}

	fmt.Printf("🧵 Threads: %d\n", threads)
	fmt.Printf("📦 Format: %s\n", format)
	fmt.Printf("♻️  Recursive: %v\n", recursive)

	if maxSize > 0 {
		fmt.Printf("📏 Max Size: %.2f MB\n", maxSize)
	}

	fmt.Println("\n🔄 Compressing files...")
	fmt.Printf("  → Scanning %s...\n", input)
	fmt.Println("  → Building file list...")
	fmt.Printf("  → Creating %s archive...\n", format)
	fmt.Printf("  → Compressing at level %d...\n", level)

	fmt.Println("✅ Compression completed successfully!")
	return nil
}

func fileSyncHandler(ctx ctx.Context) error {
	fmt.Println("🔄 Directory Synchronization")

	// Get required paths
	source, err := ctx.GetValueAsString("source")
	if err != nil {
		return fmt.Errorf("source directory is required: %v", err)
	}

	target, err := ctx.GetValueAsString("target")
	if err != nil {
		return fmt.Errorf("target directory is required: %v", err)
	}

	fmt.Printf("📂 Source: %s\n", source)
	fmt.Printf("🎯 Target: %s\n", target)

	// Handle optional parameters with validation
	batchSize := 50 // default
	if ctx.IsFlagExist("batch-size") {
		if size, err := ctx.GetValueAsInt("batch-size"); err != nil {
			return fmt.Errorf("invalid batch-size: %v", err)
		} else if size < 1 || size > 1000 {
			return fmt.Errorf("batch-size must be between 1-1000, got: %d", size)
		} else {
			batchSize = size
		}
	}

	retries := 0
	if ctx.IsFlagExist("retry") {
		if r, err := ctx.GetValueAsInt("retry"); err != nil {
			return fmt.Errorf("invalid retry count: %v", err)
		} else if r < 0 || r > 10 {
			return fmt.Errorf("retry count must be between 0-10, got: %d", r)
		} else {
			retries = r
		}
	}

	var timeout float64 = 30.0 // default 30 seconds
	if ctx.IsFlagExist("timeout") {
		if t, err := ctx.GetValueAsFloat64("timeout"); err != nil {
			return fmt.Errorf("invalid timeout: %v", err)
		} else if t <= 0 {
			return fmt.Errorf("timeout must be positive, got: %.2f", t)
		} else {
			timeout = t
		}
	}

	delete := ctx.GetValueAsBool("delete")
	dryRun := ctx.GetValueAsBool("dry-run")
	exclude := ctx.GetValueOrDefault("exclude", "")

	fmt.Printf("📦 Batch Size: %d\n", batchSize)
	fmt.Printf("🔄 Retries: %d\n", retries)
	fmt.Printf("⏰ Timeout: %.1f seconds\n", timeout)
	fmt.Printf("🗑️  Delete Extra: %v\n", delete)

	if exclude != "" {
		fmt.Printf("🚫 Exclude Pattern: %s\n", exclude)
	}

	if dryRun {
		fmt.Println("🔍 DRY RUN MODE - No changes will be applied")
	}

	fmt.Println("\n🔄 Synchronizing...")
	fmt.Println("  → Scanning source directory...")
	fmt.Println("  → Scanning target directory...")
	fmt.Println("  → Computing differences...")
	fmt.Println("  → Processing in batches...")

	if !dryRun {
		fmt.Println("✅ Synchronization completed successfully!")
	} else {
		fmt.Println("ℹ️  Dry run completed - use without --dry-run to apply changes")
	}

	return nil
}

func fileSearchHandler(ctx ctx.Context) error {
	fmt.Println("🔍 File Search")

	// Get required parameters
	pattern, err := ctx.GetValueAsString("pattern")
	if err != nil {
		return fmt.Errorf("search pattern is required: %v", err)
	}

	directory, err := ctx.GetValueAsString("directory")
	if err != nil {
		return fmt.Errorf("search directory is required: %v", err)
	}

	fmt.Printf("🎯 Pattern: %s\n", pattern)
	fmt.Printf("📁 Directory: %s\n", directory)

	// Handle optional parameters
	maxDepth := -1 // unlimited
	if ctx.IsFlagExist("max-depth") {
		if depth, err := ctx.GetValueAsInt("max-depth"); err != nil {
			return fmt.Errorf("invalid max-depth: %v", err)
		} else if depth < 0 {
			return fmt.Errorf("max-depth must be non-negative, got: %d", depth)
		} else {
			maxDepth = depth
		}
	}

	var minSize, maxSize int
	if ctx.IsFlagExist("min-size") {
		if size, err := ctx.GetValueAsInt("min-size"); err != nil {
			return fmt.Errorf("invalid min-size: %v", err)
		} else if size < 0 {
			return fmt.Errorf("min-size must be non-negative, got: %d", size)
		} else {
			minSize = size
		}
	}

	if ctx.IsFlagExist("max-size") {
		if size, err := ctx.GetValueAsInt("max-size"); err != nil {
			return fmt.Errorf("invalid max-size: %v", err)
		} else if size < 0 {
			return fmt.Errorf("max-size must be non-negative, got: %d", size)
		} else {
			maxSize = size
		}
	}

	// Validate size range
	if minSize > 0 && maxSize > 0 && minSize > maxSize {
		return fmt.Errorf("min-size (%d) cannot be greater than max-size (%d)", minSize, maxSize)
	}

	caseSensitive := ctx.GetValueAsBool("case-sensitive")
	includeHidden := ctx.GetValueAsBool("include-hidden")
	fileType := ctx.GetValueOrDefault("type", "all")

	// Validate file type
	validTypes := []string{"file", "dir", "link", "all"}
	isValidType := false
	for _, valid := range validTypes {
		if fileType == valid {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid file type: %s (must be one of: file, dir, link, all)", fileType)
	}

	if maxDepth >= 0 {
		fmt.Printf("📏 Max Depth: %d\n", maxDepth)
	} else {
		fmt.Println("📏 Max Depth: unlimited")
	}

	if minSize > 0 {
		fmt.Printf("📐 Min Size: %d bytes\n", minSize)
	}
	if maxSize > 0 {
		fmt.Printf("📏 Max Size: %d bytes\n", maxSize)
	}

	fmt.Printf("🔤 Case Sensitive: %v\n", caseSensitive)
	fmt.Printf("👻 Include Hidden: %v\n", includeHidden)
	fmt.Printf("📋 Type Filter: %s\n", fileType)

	fmt.Println("\n🔍 Searching...")
	fmt.Println("  → Scanning directories...")
	fmt.Println("  → Applying filters...")
	fmt.Println("  → Matching patterns...")

	// Simulate found results
	fmt.Println("\n📋 Results:")
	fmt.Println("  → /var/log/application.log (1.2MB)")
	fmt.Println("  → /var/log/error.log (256KB)")
	fmt.Println("  → /var/log/debug.log (5.1MB)")

	fmt.Println("✅ Search completed: 3 files found")
	return nil
}

func fileMonitorHandler(ctx ctx.Context) error {
	fmt.Println("👀 File Monitor")

	// Get required directory
	directory, err := ctx.GetValueAsString("directory")
	if err != nil {
		return fmt.Errorf("directory to monitor is required: %v", err)
	}

	fmt.Printf("📁 Monitoring: %s\n", directory)

	// Handle optional parameters
	interval := 5 // default 5 seconds
	if ctx.IsFlagExist("interval") {
		if i, err := ctx.GetValueAsInt("interval"); err != nil {
			return fmt.Errorf("invalid interval: %v", err)
		} else if i < 1 || i > 3600 {
			return fmt.Errorf("interval must be between 1-3600 seconds, got: %d", i)
		} else {
			interval = i
		}
	}

	var sizeThreshold float64 = 1.0 // default 1MB
	if ctx.IsFlagExist("size-threshold") {
		if threshold, err := ctx.GetValueAsFloat64("size-threshold"); err != nil {
			return fmt.Errorf("invalid size-threshold: %v", err)
		} else if threshold <= 0 {
			return fmt.Errorf("size-threshold must be positive, got: %.2f", threshold)
		} else {
			sizeThreshold = threshold
		}
	}

	recursive := ctx.GetValueAsBool("recursive")
	outputFile := ctx.GetValueOrDefault("output", "monitor.log")

	fmt.Printf("⏰ Check Interval: %d seconds\n", interval)
	fmt.Printf("📊 Size Threshold: %.2f MB\n", sizeThreshold)
	fmt.Printf("♻️  Recursive: %v\n", recursive)
	fmt.Printf("📄 Output Log: %s\n", outputFile)

	fmt.Println("\n👀 Starting monitor...")
	fmt.Println("  → Setting up file watchers...")
	fmt.Println("  → Establishing baseline...")
	fmt.Println("  → Monitor active (Press Ctrl+C to stop)")

	fmt.Println("✅ Monitor started successfully!")
	return nil
}
