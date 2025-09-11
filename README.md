# Cmd

A powerful and type-safe Go library for building command-line interfaces with hierarchical commands, option groups, and automatic argument parsing and validation.

## Features

- **Hierarchical Commands**: Build complex nested command structures like `docker container run`
- **Type-Safe Options**: Automatic parsing and validation for strings, integers, floats, and booleans
- **Option Groups**: Create mutually exclusive or inclusive option groups
- **Fluent API**: Clean, readable command definitions with method chaining
- **Custom Error Handling**: Flexible error handling with custom error handlers
- **Context-Based**: Rich context object for accessing parsed arguments and flags

## Installation

```bash
go get github.com/DilemaFixer/Cmd
```

## Quick Start

```go
package main

import (
    "fmt"
    ctx "github.com/DilemaFixer/Cmd/context"
    p "github.com/DilemaFixer/Cmd/parser"
    rtr "github.com/DilemaFixer/Cmd/router"
)

func main() {
    // Parse command line input
    input := "server start --host=localhost --port=8080 --debug"
    parsedInput, err := p.ParseInput(input)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Create context and router
    context := ctx.NewContext(parsedInput)
    iterator := rtr.NewRoutingIterator(context)
    router := rtr.NewRouter()

    // Define command with options
    router.Endpoint("server").
        Endpoint("start").
        Description("Start HTTP server").
        StringOption("host").
        IntOption("port").
        BoolOption("debug").
        Handler(startHandler).
        Register()

    // Execute command
    router.Route(*context, iterator)
}

func startHandler(ctx ctx.Context) error {
    host := ctx.GetValueOrDefault("host", "localhost")
    port := ctx.GetValueOrDefault("port", "3000")
    debug := ctx.GetValueAsBool("debug")
    
    fmt.Printf("Starting server on %s:%s (debug: %v)\n", host, port, debug)
    return nil
}
```

## Parsing Input

### Parse String Input
```go
parsedInput, err := p.ParseInput("command subcommand --flag=value --bool-flag")
```

### Parse OS Arguments
```go
parsedInput, err := p.ParseOSArgs()
```

### Parse Custom Arguments
```go
args := []string{"command", "subcommand", "--flag=value"}
parsedInput, err := p.ParseArgs(args)
```

## Building Commands

### Simple Commands

```go
router := rtr.NewRouter()

// Single endpoint
router.Endpoint("version").
    Description("Show version information").
    BoolOption("short").
    Handler(versionHandler).
    Register()
```

### Nested Commands

```go
router.NewCmd("docker").
    NewSub("container").
        Endpoint("run").
            RequiredString("image").
            StringOption("name").
            IntOption("port").
            BoolOption("detach").
            Handler(containerRunHandler).
            Build().
        
        Endpoint("stop").
            RequiredString("container").
            BoolOption("force").
            Handler(containerStopHandler).
            Build().
        Build().
    Register()
```

## Option Types

The library supports various option types with automatic validation:

```go
router.Endpoint("config").
    // String options
    StringOption("name").        // Optional: --name=value
    RequiredString("config").    // Required: --config=value
    
    // Integer options  
    IntOption("count").          // Optional: --count=42
    RequiredInt("port").         // Required: --port=8080
    
    // Float options
    FloatOption("ratio").        // Optional: --ratio=3.14
    RequiredFloat("threshold").  // Required: --threshold=0.95
    
    // Boolean options (flags)
    BoolOption("enable").        // Optional: --enable
    RequiredBool("force").       // Required: --force
    
    Handler(configHandler).
    Register()
```

## Option Groups

### Exclusive Groups
Only one group can be active at a time:

```go
router.Endpoint("backup").
    // Mutually exclusive storage options
    ExclusiveGroup("local", "--local").
        RequiredString("path").
        BoolOption("compress").
        StringOption("encryption").
    EndGroup().
    
    ExclusiveGroup("remote", "--remote").
        RequiredString("host").
        RequiredString("user").
        StringOption("key").
    EndGroup().
    
    ExclusiveGroup("cloud", "--s3").
        RequiredString("bucket").
        RequiredString("region").
    EndGroup().
    
    SetGroupsCanBeIgnored(false). // At least one group required
    Handler(backupHandler).
    Register()
```

**Usage:**
```bash
# Valid - only one group active
myapp backup --local --path=/backup --compress

# Invalid - multiple exclusive groups
myapp backup --local --path=/backup --remote --host=server.com
```

### Inclusive Groups
Groups that can be used together:

```go
router.Endpoint("deploy").
    RequiredString("environment").
    
    // Optional resource configuration
    Group("resources", "--resources").
        IntOption("memory").
        IntOption("cpu").
        StringOption("storage").
    EndGroup().
    
    // Optional monitoring setup
    Group("monitoring", "--monitoring").
        BoolOption("metrics").
        BoolOption("logging").
        StringOption("alerts").
    EndGroup().
    
    Handler(deployHandler).
    Register()
```

**Usage:**
```bash
# Can use both groups together
myapp deploy --environment=prod --resources --memory=512 --cpu=2 --monitoring --metrics --logging
```

## Context API

Access parsed options and commands in your handlers:

### Flag Checking
```go
func handler(ctx ctx.Context) error {
    // Check if flag was provided
    if ctx.IsFlagExist("port") {
        // Flag exists
    }
    
    // Check if flag has a value
    if ctx.IsFlagHaveValue("host") {
        // Flag has non-empty value
    }
}
```

### Type-Safe Value Retrieval
```go
func handler(ctx ctx.Context) error {
    // String values
    host, err := ctx.GetValueAsString("host")
    if err != nil {
        return err
    }
    
    // Integer values
    port, err := ctx.GetValueAsInt("port")
    count32, err := ctx.GetValueAsInt32("count")
    count64, err := ctx.GetValueAsInt64("bigcount")
    
    // Float values
    ratio32, err := ctx.GetValueAsFloat32("ratio")
    ratio64, err := ctx.GetValueAsFloat64("precision")
    
    // Boolean values (returns false if flag doesn't exist)
    debug := ctx.GetValueAsBool("debug")
    
    // Default values
    timeout := ctx.GetValueOrDefault("timeout", "30s")
    
    return nil
}
```

### Command Information
```go
func handler(ctx ctx.Context) error {
    // Get main command
    cmd := ctx.GetCommand()
    
    // Check command
    if ctx.IsCommandEqual("serve") {
        // Handle serve command
    }
    
    // Check subcommands
    if ctx.IsSubcommandExist("start") {
        // Subcommand exists
    }
    
    // Get all subcommands
    subs := ctx.GetSubcommandsAsArr()
    
    // Get all flags
    flagsMap := ctx.GetFlagsAsMap()
    flagsArray := ctx.GetFlagsAsArr()
    keys := ctx.GetFlagsKeysAsArr()
    values := ctx.GetFlagsValuesAsArr()
    
    return nil
}
```

## Error Handling

### Custom Error Handler
```go
router := rtr.NewRouter()

router.CustomErrorHandler(func(err error, ctx ctx.Context) {
    fmt.Printf("Command failed: %v\n", err)
    fmt.Printf("Command: %s\n", ctx.GetCommand())
    fmt.Printf("Flags: %v\n", ctx.GetFlagsAsMap())
    os.Exit(1)
})
```

### Common Error Types
- **Missing required options**: `"Required 'config' flag not exist"`
- **Invalid option values**: `"Option port with type Int have error"`
- **Type validation**: `"Option debug with type Bool have value"`
- **Conflicting groups**: `"group requires solitude"`
- **Unknown commands**: `"Point with name 'unknown' not found"`

## Method Chaining

The fluent API allows clean command definitions:

```go
router.NewCmd("service").              // Returns *CmdWrapper
    Description("Service management").
    
    NewSub("database").                // Returns *CmdWrapper  
        Endpoint("migrate").           // Returns *EndPointWrapper
            Description("Run migrations").
            RequiredString("direction").
            BoolOption("dry-run").
            
            Group("connection", "--db").   // Returns *EndPointGroupWrapper
                RequiredString("host").
                IntOption("port").
            EndGroup().                    // Returns *EndPointWrapper
            
            Handler(migrateHandler).
            Build().                       // Returns parent *CmdWrapper
        Build().                           // Returns grandparent *CmdWrapper  
    Register()                             // Registers with router
```

## Data Structures

### ParsedInput
```go
type ParsedInput struct {
    Command     string        // Main command name
    Subcommands []string      // List of subcommands
    InputFlags  []InputFlag   // Parsed flags with values
}
```

### InputFlag  
```go
type InputFlag struct {
    Name  string    // Flag name (without --)
    Value string    // Flag value (empty for boolean flags)
}
```

## Examples

See the `example/` directory for complete working examples:

- **HTTPServerManager**: Basic server management with simple options
- **DatabaseMigrationTool**: Complex nested commands for database operations  
- **DeploymentTool**: Advanced option groups for deployment scenarios
- **FileManager**: Type validation and custom error handling
- **SimpleTestFunctionCall**: Minimal working example

## Error Handling Best Practices

1. **Always check required parameters** in your handlers
2. **Validate business logic** beyond type checking  
3. **Use custom error handlers** for better user experience
4. **Provide helpful error messages** with context

```go
func handler(ctx ctx.Context) error {
    // Check required parameters
    if !ctx.IsFlagExist("config") {
        return fmt.Errorf("config file is required")
    }
    
    // Validate values
    port, err := ctx.GetValueAsInt("port")
    if err != nil {
        return fmt.Errorf("invalid port: %v", err)
    }
    if port < 1 || port > 65535 {
        return fmt.Errorf("port must be between 1-65535, got: %d", port)
    }
    
    return nil
}
```

## License

This project is licensed under the MIT License.
