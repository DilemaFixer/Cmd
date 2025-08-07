# Cmd

TODO: In context struct add func MustExist(targetFlags []string) error , return err if flag is exist 

A simple and powerful Go library for building command-line interfaces with hierarchical commands, option groups, and type-safe argument parsing.

## Quick Start

```bash
go get github.com/DilemaFixer/Cmd
```

## Features

- **Hierarchical Commands**: Build complex nested command structures like `docker container run`
- **Type-Safe Options**: Automatic parsing and validation for strings, integers, floats, and booleans
- **Option Groups**: Create mutually exclusive or related option groups
- **Fluent API**: Clean, readable command definitions with method chaining
- **Custom Error Handling**: Flexible error handling with custom handlers
- **Context-Based**: Rich context object for accessing parsed arguments

## Basic Usage

Create a simple command with basic options:

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
    input := "server --host=localhost --port=8080 --debug"
    parsedInput, err := p.ParseInput(input)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Create context and router
    context := ctx.NewContext(parsedInput)
    iterator := rtr.NewRoutingIterator(context)
    router := rtr.NewRouter()

    // Define command
    router.Endpoint("server").
        Description("Start HTTP server").
        StringOption("host").
        IntOption("port").
        BoolOption("debug").
        Handler(serverHandler).
        Register()

    // Route and execute
    router.Route(*context, iterator)
}

func serverHandler(ctx ctx.Context) error {
    host, _ := ctx.GetValueAsString("host")
    port, _ := ctx.GetValueAsInt("port")
    debug, _ := ctx.GetValueAsBool("debug")
    
    fmt.Printf("Starting server on %s:%d (debug: %v)\n", host, port, debug)
    return nil
}
```

**Usage**: `myapp server --host=localhost --port=8080 --debug`

## Hierarchical Commands

Build complex nested command structures:

```go
router := rtr.NewRouter()

router.NewCmd("docker").
    NewSub("container").
        Endpoint("run").
            Description("Run a new container").
            RequiredString("image").
            StringOption("name").
            IntOption("port").
            BoolOption("detach").
            Handler(containerRunHandler).
            Build().
        
        Endpoint("stop").
            Description("Stop running container").
            RequiredString("container").
            BoolOption("force").
            Handler(containerStopHandler).
            Build().
        Build().
    
    NewSub("image").
        Endpoint("build").
            Description("Build image from Dockerfile").
            RequiredString("dockerfile").
            StringOption("tag").
            BoolOption("no-cache").
            Handler(imageBuildHandler).
            Build().
        
        Endpoint("pull").
            Description("Pull image from registry").
            RequiredString("image").
            StringOption("tag").
            Handler(imagePullHandler).
            Build().
        Build().
    Register()
```

**Usage examples**:
```bash
myapp docker container run --image=nginx --name=web --port=80 --detach
myapp docker container stop --container=web --force
myapp docker image build --dockerfile=Dockerfile --tag=myapp:latest
myapp docker image pull --image=alpine --tag=3.18
```

## Option Groups

### Exclusive Groups

Create mutually exclusive option groups where only one can be used:

```go
router.Endpoint("backup").
    Description("Create application backup").
    
    // Common options
    BoolOption("verbose").
    StringOption("output").
    
    // Exclusive storage groups
    ExclusiveGroup("local", "--local").
        RequiredString("path").
        BoolOption("compress").
        StringOption("encryption").
    EndGroup().
    
    ExclusiveGroup("remote", "--remote").
        RequiredString("host").
        RequiredString("user").
        StringOption("key").
        IntOption("port").
    EndGroup().
    
    ExclusiveGroup("cloud", "--s3").
        RequiredString("bucket").
        RequiredString("region").
        StringOption("access-key").
    EndGroup().
    
    SetGroupsCanBeIgnored(false). // One group is required
    Handler(backupHandler).
    Register()
```

**Valid usage**:
```bash
myapp backup --verbose --local --path=/backup --compress
myapp backup --remote --host=backup.com --user=admin --key=id_rsa
myapp backup --s3 --bucket=my-backup --region=us-east-1
```

**Invalid usage** (will cause error):
```bash
myapp backup --local --path=/backup --remote --host=backup.com
```

### Inclusive Groups

Create related option groups that can be used together:

```go
router.Endpoint("deploy").
    Description("Deploy application").
    RequiredString("environment").
    
    // Resource limits (inclusive group)
    Group("resources", "--resources").
        IntOption("memory").
        IntOption("cpu").
        StringOption("storage").
    EndGroup().
    
    // Monitoring options (inclusive group)
    Group("monitoring", "--monitoring").
        BoolOption("metrics").
        BoolOption("logging").
        StringOption("alerts").
    EndGroup().
    
    Handler(deployHandler).
    Register()
```

**Usage**:
```bash
myapp deploy --environment=prod --resources --memory=512 --cpu=2 --monitoring --metrics --logging
```

## Option Types

The library supports various option types with automatic parsing and validation:

```go
router.Endpoint("config").
    StringOption("name").        // --name=value
    RequiredString("config").    // --config=value (required)
    IntOption("count").          // --count=42
    RequiredInt("port").         // --port=8080 (required)
    FloatOption("ratio").        // --ratio=3.14
    RequiredFloat("threshold").  // --threshold=0.95 (required)
    BoolOption("enable").        // --enable (flag, no value)
    Handler(configHandler).
    Register()
```

## Context Usage

Access parsed options in your handlers with type safety:

```go
func serverHandler(ctx ctx.Context) error {
    // Check if option exists
    if !ctx.IsFlagExist("port") {
        return fmt.Errorf("port is required")
    }
    
    // Check if option has value
    if !ctx.IsFlagHaveValue("host") {
        return fmt.Errorf("host cannot be empty")
    }
    
    // Get typed values with error handling
    port, err := ctx.GetValueAsInt("port")
    if err != nil {
        return fmt.Errorf("invalid port: %v", err)
    }
    
    host, err := ctx.GetValueAsString("host")
    if err != nil {
        return fmt.Errorf("invalid host: %v", err)
    }
    
    debug, _ := ctx.GetValueAsBool("debug") // defaults to false if not set
    
    // Use default values
    timeout := ctx.GetValueOrDefault("timeout", "30s")
    
    fmt.Printf("Server: %s:%d (debug: %v, timeout: %s)\n", 
        host, port, debug, timeout)
    
    return nil
}
```

### Context Methods

**Flag Checking**:
- `IsFlagExist(name string) bool` - Check if flag was provided
- `IsFlagHaveValue(name string) bool` - Check if flag has a value

**Type-Safe Getters**:
- `GetValueAsString(name string) (string, error)`
- `GetValueAsInt(name string) (int, error)`
- `GetValueAsInt32(name string) (int32, error)`
- `GetValueAsInt64(name string) (int64, error)`
- `GetValueAsFloat32(name string) (float32, error)`
- `GetValueAsFloat64(name string) (float64, error)`
- `GetValueAsBool(name string) (bool, error)`

**Utility Methods**:
- `GetValueOrDefault(name, defaultValue string) string`
- `GetCommand() string`
- `IsCommandEqual(target string) bool`
- `IsSubcommandExist(target string) bool`
- `GetSubcommandsAsArr() []string`
- `GetFlagsAsMap() map[string]string`

## Advanced Example

Complete example with multiple nested commands and groups:

```go
package main

import (
    "fmt"
    ctx "github.com/DilemaFixer/Cmd/context"
    p "github.com/DilemaFixer/Cmd/parser"
    rtr "github.com/DilemaFixer/Cmd/router"
)

func main() {
    input := "service database migrate --direction=up --dry-run --db --host=localhost --port=5432"
    parsedInput, _ := p.ParseInput(input)
    context := ctx.NewContext(parsedInput)
    iterator := rtr.NewRoutingIterator(context)
    router := rtr.NewRouter()

    router.NewCmd("service").
        Description("Service management").
        
        NewSub("database").
            Endpoint("migrate").
                Description("Run database migrations").
                RequiredString("direction").
                BoolOption("dry-run").
                
                Group("connection", "--db").
                    RequiredString("host").
                    IntOption("port").
                    StringOption("database").
                EndGroup().
                
                Handler(migrateHandler).
                Build().
            Build().
        Register()

    router.Route(*context, iterator)
}

func migrateHandler(ctx ctx.Context) error {
    direction, _ := ctx.GetValueAsString("direction")
    dryRun, _ := ctx.GetValueAsBool("dry-run")
    host, _ := ctx.GetValueAsString("host")
    port, _ := ctx.GetValueAsInt("port")
    
    fmt.Printf("Migration: %s (dry-run: %v)\n", direction, dryRun)
    fmt.Printf("Database: %s:%d\n", host, port)
    
    return nil
}
```

## Error Handling

The library provides detailed error information for invalid commands and supports custom error handlers:

```go
router := rtr.NewRouter()

// Custom error handler
router.CustomErrorHandler(func(err error, ctx ctx.Context) {
    fmt.Printf("Command failed: %v\n", err)
    fmt.Printf("Command: %s\n", ctx.GetCommand())
    fmt.Printf("Flags: %v\n", ctx.GetFlagsAsMap())
    os.Exit(1)
})
```

**Common error types**:
- Missing required options: `"missing required option --path"`
- Invalid option values: `"invalid value for --port: expected integer, got 'abc'"`
- Conflicting groups: `"conflicting groups: cannot use --kubernetes when --docker is active"`
- Unknown commands: `"Point with name 'unknown' not found"`

## Method Chaining

The fluent interface allows for clean, readable command definitions:

```go
// Each method returns the appropriate wrapper for continued chaining
router.NewCmd("service").                    // Returns *CmdWrapper
    Description("Service management").
    
    NewSub("database").                      // Returns *CmdWrapper
        Endpoint("migrate").                 // Returns *EndPointWrapper
            Description("Run migrations").
            RequiredString("direction").
            BoolOption("dry-run").
            
            Group("connection", "--db").     // Returns *EndPointGroupWrapper
                RequiredString("host").
                IntOption("port").
                StringOption("database").
            EndGroup().                      // Returns *EndPointWrapper
            
            Handler(migrateHandler).
            Build().                         // Returns *CmdWrapper (parent)
        Build().                             // Returns *CmdWrapper (grandparent)
    Register()                               // Registers with router
```
