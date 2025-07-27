# Cmd
Simpl console command parser

Create a simple command with basic options:

```go
router := NewRouter()

router.Endpoint("server").
    Description("Start HTTP server").
    StringOption("host").
    IntOption("port").
    BoolOption("debug").
    Handler(func(ctx Context) error {
        host, _ := ctx.GetValueAsString("host")
        port, _ := ctx.GetValueAsInt("port")
        debug, _ := ctx.GetValueAsBool("debug")
        
        fmt.Printf("Starting server on %s:%d (debug: %v)\n", host, port, debug)
        return nil
    }).
    Register()

// Usage: myapp server --host=localhost --port=8080 --debug
```

### Hierarchical Commands

Build complex nested command structures:

```go
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

// Usage examples:
// myapp docker container run --image=nginx --name=web --port=80 --detach
// myapp docker container stop --container=web --force
// myapp docker image build --dockerfile=Dockerfile --tag=myapp:latest
// myapp docker image pull --image=alpine --tag=3.18
```

### Option Groups

Create mutually exclusive or related option groups:

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

// Valid usage:
// myapp backup --verbose --local --path=/backup --compress
// myapp backup --remote --host=backup.com --user=admin --key=id_rsa
// myapp backup --s3 --bucket=my-backup --region=us-east-1

// Invalid usage (will cause error):
// myapp backup --local --path=/backup --remote --host=backup.com
```

### Advanced Example with Multiple Groups

```go
router.NewCmd("deploy").
    Description("Deploy application to various environments").
    
    NewSub("app").
        Endpoint("start").
            Description("Deploy and start application").
            
            // Basic options
            RequiredString("environment").
            BoolOption("force").
            IntOption("timeout").
            
            // Deployment method (exclusive)
            ExclusiveGroup("method", "--docker").
                RequiredString("image").
                StringOption("tag").
                IntOption("replicas").
            EndGroup().
            
            ExclusiveGroup("method", "--kubernetes").
                RequiredString("namespace").
                StringOption("config").
                BoolOption("dry-run").
            EndGroup().
            
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
            
            SetGroupsCanBeIgnored(true).
            Handler(deployHandler).
            Build().
        Build().
    Register()

// Usage:
// myapp deploy app start --environment=prod --docker --image=myapp:v1.2.3 --replicas=3 --resources --memory=512 --cpu=2
```

## Option Types

The router supports various option types with automatic parsing:

```go
router.Endpoint("config").
    StringOption("name").        // --name=value
    RequiredString("config").    // --config=value (required)
    IntOption("count").          // --count=42
    RequiredInt("port").         // --port=8080 (required)
    FloatOption("ratio").        // --ratio=3.14
    RequiredFloat("threshold").  // --threshold=0.95 (required)
    BoolOption("enable").        // --enable (flag)
    Handler(configHandler).
    Register()
```

## Context Usage

Access parsed options in your handlers:

```go
func serverHandler(ctx Context) error {
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

## Error Handling

The router provides detailed error information for invalid commands:

```go
// Invalid command: myapp deploy --docker --kubernetes
// Error: "conflicting groups: cannot use --kubernetes when --docker is active"

// Invalid option: myapp server --port=abc
// Error: "invalid value for --port: expected integer, got 'abc'"

// Missing required option: myapp backup --local
// Error: "missing required option --path for group 'local'"
```
