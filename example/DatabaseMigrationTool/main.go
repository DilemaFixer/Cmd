package main

import (
	"fmt"

	ctx "github.com/DilemaFixer/Cmd/context"
	p "github.com/DilemaFixer/Cmd/parser"
	rtr "github.com/DilemaFixer/Cmd/router"
)

// Example 2: Database Migration Tool with Nested Commands
// This example demonstrates hierarchical command structure like "db migrate up"
// Usage examples:
//   myapp db migrate up --steps=5 --dry-run
//   myapp db migrate down --steps=2
//   myapp db migrate status --detailed
//   myapp db schema dump --format=sql --output=backup.sql
//   myapp db connection test --host=localhost --port=5432 --database=myapp

func main() {
	// Simulate complex nested command
	input := "db migrate up --steps=3 --dry-run"

	parsedInput, err := p.ParseInput(input)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	context := ctx.NewContext(parsedInput)
	iterator := rtr.NewRoutingIterator(context)
	router := rtr.NewRouter()

	// Build complex nested command structure
	router.NewCmd("db").
		// Migration sub-commands
		NewSub("migrate").
		Endpoint("up").
		IntOption("steps").    // Number of migrations to run
		BoolOption("dry-run"). // Preview changes without applying
		BoolOption("verbose"). // Detailed output
		Handler(migrateUpHandler).
		Build().
		Endpoint("down").
		Description("Rollback migrations").
		RequiredInt("steps").  // Required: number of migrations to rollback
		BoolOption("dry-run"). // Preview changes without applying
		BoolOption("force").   // Force rollback even if risky
		Handler(migrateDownHandler).
		Build().
		Endpoint("status").
		Description("Show migration status").
		BoolOption("detailed"). // Show detailed migration info
		BoolOption("pending").  // Show only pending migrations
		Handler(migrateStatusHandler).
		Build().
		Endpoint("create").
		Description("Create new migration file").
		RequiredString("name"). // Migration name
		StringOption("type").   // Migration type (table, data, etc.)
		Handler(migrateCreateHandler).
		Build().
		Build(). // End migrate sub-command

		// Schema sub-commands
		NewSub("schema").
		Endpoint("dump").
		RequiredString("format"). // Output format (sql, json, yaml)
		StringOption("output").   // Output file path
		BoolOption("data").       // Include data in dump
		BoolOption("compress").   // Compress output
		Handler(schemaDumpHandler).
		Build().
		Endpoint("load").
		Description("Import database schema").
		RequiredString("file").      // Input file path
		BoolOption("drop-existing"). // Drop existing tables
		BoolOption("ignore-errors"). // Continue on errors
		Handler(schemaLoadHandler).
		Build().
		Build(). // End schema sub-command

		// Connection sub-commands
		NewSub("connection").
		Endpoint("test").
		Description("Test database connection").
		RequiredString("host").     // Database host
		RequiredInt("port").        // Database port
		RequiredString("database"). // Database name
		StringOption("user").       // Username
		StringOption("password").   // Password
		BoolOption("ssl").          // Use SSL connection
		IntOption("timeout").       // Connection timeout
		Handler(connectionTestHandler).
		Build().
		Endpoint("info").
		Description("Show database connection info").
		BoolOption("show-config"). // Show current configuration
		Handler(connectionInfoHandler).
		Build().
		Build(). // End connection sub-command

		Register()

	// Execute the command
	router.Route(*context, iterator)
}

func migrateUpHandler(ctx ctx.Context) error {
	fmt.Println("🔄 Running Database Migrations (UP)")

	steps := ctx.GetValueOrDefault("steps", "all")
	dryRun := ctx.GetValueAsBool("dry-run")
	verbose := ctx.GetValueAsBool("verbose")

	if dryRun {
		fmt.Println("🔍 DRY RUN MODE - No changes will be applied")
	}

	fmt.Printf("📊 Steps to run: %s\n", steps)

	if verbose {
		fmt.Println("📝 Detailed output enabled")
		fmt.Println("  → Found 3 pending migrations")
		fmt.Println("  → 001_create_users_table.sql")
		fmt.Println("  → 002_add_user_indexes.sql")
		fmt.Println("  → 003_create_orders_table.sql")
	}

	if !dryRun {
		fmt.Println("✅ Migrations applied successfully!")
	} else {
		fmt.Println("ℹ️  Dry run completed - use without --dry-run to apply changes")
	}

	return nil
}

func migrateDownHandler(ctx ctx.Context) error {
	fmt.Println("⬇️  Rolling Back Database Migrations")

	steps, err := ctx.GetValueAsInt("steps")
	if err != nil {
		return fmt.Errorf("invalid steps value: %v", err)
	}

	dryRun := ctx.GetValueAsBool("dry-run")
	force := ctx.GetValueAsBool("force")

	fmt.Printf("📊 Rolling back %d migrations\n", steps)

	if force {
		fmt.Println("⚠️  FORCE mode enabled - ignoring safety checks")
	}

	if dryRun {
		fmt.Println("🔍 DRY RUN MODE - Showing what would be rolled back")
	} else {
		fmt.Println("✅ Rollback completed successfully!")
	}

	return nil
}

func migrateStatusHandler(ctx ctx.Context) error {
	fmt.Println("📋 Migration Status")

	detailed := ctx.GetValueAsBool("detailed")
	pendingOnly := ctx.GetValueAsBool("pending")

	if pendingOnly {
		fmt.Println("📌 Pending Migrations:")
		fmt.Println("  → 003_create_orders_table.sql")
		fmt.Println("  → 004_add_foreign_keys.sql")
	} else {
		fmt.Println("Applied: 2, Pending: 2, Total: 4")

		if detailed {
			fmt.Println("\n✅ Applied Migrations:")
			fmt.Println("  → 001_create_users_table.sql (2024-01-15)")
			fmt.Println("  → 002_add_user_indexes.sql (2024-01-16)")
			fmt.Println("\n📌 Pending Migrations:")
			fmt.Println("  → 003_create_orders_table.sql")
			fmt.Println("  → 004_add_foreign_keys.sql")
		}
	}

	return nil
}

func migrateCreateHandler(ctx ctx.Context) error {
	name, _ := ctx.GetValueAsString("name")
	migrationType := ctx.GetValueOrDefault("type", "table")

	fmt.Printf("📝 Creating new migration: %s\n", name)
	fmt.Printf("🏷️  Type: %s\n", migrationType)

	filename := fmt.Sprintf("005_%s.sql", name)
	fmt.Printf("📄 Generated file: %s\n", filename)

	return nil
}

func schemaDumpHandler(ctx ctx.Context) error {
	format, _ := ctx.GetValueAsString("format")
	output := ctx.GetValueOrDefault("output", fmt.Sprintf("schema.%s", format))
	includeData := ctx.GetValueAsBool("data")
	compress := ctx.GetValueAsBool("compress")

	fmt.Printf("💾 Dumping database schema to: %s\n", output)
	fmt.Printf("📋 Format: %s\n", format)
	fmt.Printf("📊 Include data: %v\n", includeData)
	fmt.Printf("🗜️  Compress: %v\n", compress)

	fmt.Println("✅ Schema dump completed successfully!")
	return nil
}

func schemaLoadHandler(ctx ctx.Context) error {
	file, _ := ctx.GetValueAsString("file")
	dropExisting := ctx.GetValueAsBool("drop-existing")
	ignoreErrors := ctx.GetValueAsBool("ignore-errors")

	fmt.Printf("📥 Loading schema from: %s\n", file)

	if dropExisting {
		fmt.Println("⚠️  Will drop existing tables")
	}

	if ignoreErrors {
		fmt.Println("🤷 Will ignore errors and continue")
	}

	fmt.Println("✅ Schema loaded successfully!")
	return nil
}

func connectionTestHandler(ctx ctx.Context) error {
	host, _ := ctx.GetValueAsString("host")
	port, _ := ctx.GetValueAsInt("port")
	database, _ := ctx.GetValueAsString("database")
	user := ctx.GetValueOrDefault("user", "postgres")
	ssl := ctx.GetValueAsBool("ssl")
	timeout := ctx.GetValueOrDefault("timeout", "30")

	fmt.Println("🔌 Testing Database Connection...")
	fmt.Printf("🏠 Host: %s:%d\n", host, port)
	fmt.Printf("🗃️  Database: %s\n", database)
	fmt.Printf("👤 User: %s\n", user)
	fmt.Printf("🔒 SSL: %v\n", ssl)
	fmt.Printf("⏰ Timeout: %s seconds\n", timeout)

	fmt.Println("✅ Connection test successful!")
	fmt.Println("📊 Response time: 45ms")

	return nil
}

func connectionInfoHandler(ctx ctx.Context) error {
	showConfig := ctx.GetValueAsBool("show-config")

	fmt.Println("ℹ️  Database Connection Information")
	fmt.Println("Status: Connected")
	fmt.Println("Server Version: PostgreSQL 14.2")
	fmt.Println("Client Version: PostgreSQL 14.2")

	if showConfig {
		fmt.Println("\n⚙️  Current Configuration:")
		fmt.Println("  Host: localhost")
		fmt.Println("  Port: 5432")
		fmt.Println("  Database: myapp_dev")
		fmt.Println("  SSL Mode: prefer")
	}

	return nil
}
