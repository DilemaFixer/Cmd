package main

import (
	"fmt"

	ctx "github.com/DilemaFixer/Cmd/context"
	p "github.com/DilemaFixer/Cmd/parser"
	rtr "github.com/DilemaFixer/Cmd/router"
)

// Example 3: Deployment Tool with Option Groups
// This example demonstrates exclusive and inclusive option groups
// Usage examples:
//   myapp deploy --environment=prod --docker --image=myapp:latest --registry=docker.io
//   myapp deploy --environment=staging --kubernetes --namespace=staging --replicas=3
//   myapp deploy --environment=dev --resources --memory=512 --cpu=2 --monitoring --metrics --alerts=slack
//   myapp backup --local --path=/backup --compress --encryption=aes256
//   myapp backup --s3 --bucket=my-backup --region=us-east-1 --access-key=AKIA...

func main() {
	// Example with exclusive deployment groups
	input := "deploy --environment=prod --docker --image=myapp:latest --registry=docker.io --resources --memory=1024 --cpu=4"

	parsedInput, err := p.ParseInput(input)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	context := ctx.NewContext(parsedInput)
	iterator := rtr.NewRoutingIterator(context)
	router := rtr.NewRouter()

	// Define deployment command with exclusive platform groups
	router.Endpoint("deploy").
		RequiredString("environment"). // Required: target environment

		// Exclusive deployment platform groups (only one can be used)
		ExclusiveGroup("docker", "--docker").
		RequiredString("image").    // Docker image name
		StringOption("registry").   // Docker registry
		StringOption("tag").        // Image tag
		StringOption("dockerfile"). // Path to Dockerfile
		BoolOption("build").        // Build image before deploy
		BoolOption("push").         // Push to registry
		EndGroup().
		ExclusiveGroup("kubernetes", "--kubernetes").
		RequiredString("namespace"). // K8s namespace
		IntOption("replicas").       // Number of replicas
		StringOption("config").      // Kubernetes config file
		StringOption("context").     // Kubectl context
		BoolOption("wait").          // Wait for deployment to complete
		EndGroup().
		ExclusiveGroup("serverless", "--serverless").
		RequiredString("provider"). // Cloud provider (aws, gcp, azure)
		StringOption("region").     // Target region
		StringOption("runtime").    // Runtime environment
		IntOption("memory").        // Memory allocation in MB
		IntOption("timeout").       // Timeout in seconds
		EndGroup().

		// Inclusive resource configuration group (optional)
		Group("resources", "--resources").
		IntOption("memory").     // Memory limit in MB
		IntOption("cpu").        // CPU limit (cores)
		StringOption("storage"). // Storage limit
		EndGroup().

		// Inclusive monitoring group (optional)
		Group("monitoring", "--monitoring").
		BoolOption("metrics").     // Enable metrics collection
		BoolOption("logging").     // Enable centralized logging
		StringOption("alerts").    // Alert destination (slack, email, etc.)
		StringOption("dashboard"). // Dashboard URL
		EndGroup().
		SetGroupsCanBeIgnored(false). // At least one platform group is required
		Handler(deployHandler).
		Register()

	// Define backup command with exclusive storage groups
	router.Endpoint("backup").
		BoolOption("verbose").  // Verbose output
		StringOption("output"). // Output filename

		// Exclusive storage destination groups
		ExclusiveGroup("local", "--local").
		RequiredString("path").     // Local backup path
		BoolOption("compress").     // Compress backup
		StringOption("encryption"). // Encryption method
		EndGroup().
		ExclusiveGroup("remote", "--remote").
		RequiredString("host"). // Remote host
		RequiredString("user"). // SSH user
		StringOption("key").    // SSH private key path
		IntOption("port").      // SSH port
		StringOption("path").   // Remote path
		EndGroup().
		ExclusiveGroup("cloud", "--s3").
		RequiredString("bucket").   // S3 bucket name
		RequiredString("region").   // AWS region
		StringOption("access-key"). // AWS access key
		StringOption("secret-key"). // AWS secret key
		StringOption("prefix").     // Object key prefix
		EndGroup().
		SetGroupsCanBeIgnored(false). // One storage group is required
		Handler(backupHandler).
		Register()

	// Execute the command
	router.Route(*context, iterator)
}

func deployHandler(ctx ctx.Context) error {
	environment, _ := ctx.GetValueAsString("environment")
	fmt.Printf("🚀 Deploying to %s environment\n", environment)

	// Check which exclusive platform group is being used
	if ctx.IsFlagExist("docker") {
		return handleDockerDeploy(ctx)
	} else if ctx.IsFlagExist("kubernetes") {
		return handleKubernetesDeploy(ctx)
	} else if ctx.IsFlagExist("serverless") {
		return handleServerlessDeploy(ctx)
	}

	return fmt.Errorf("no deployment platform specified")
}

func handleDockerDeploy(ctx ctx.Context) error {
	fmt.Println("🐳 Docker Deployment")

	image, _ := ctx.GetValueAsString("image")
	registry := ctx.GetValueOrDefault("registry", "docker.io")
	tag := ctx.GetValueOrDefault("tag", "latest")

	build := ctx.GetValueAsBool("build")
	push := ctx.GetValueAsBool("push")

	fmt.Printf("📦 Image: %s/%s:%s\n", registry, image, tag)

	if build {
		fmt.Println("🔨 Building Docker image...")
		dockerfile := ctx.GetValueOrDefault("dockerfile", "Dockerfile")
		fmt.Printf("  Using Dockerfile: %s\n", dockerfile)
	}

	if push {
		fmt.Println("⬆️  Pushing to registry...")
		fmt.Printf("  Registry: %s\n", registry)
	}

	// Handle optional resource configuration
	handleResourceConfiguration(ctx)

	// Handle optional monitoring configuration
	handleMonitoringConfiguration(ctx)

	fmt.Println("✅ Docker deployment completed!")
	return nil
}

func handleKubernetesDeploy(ctx ctx.Context) error {
	fmt.Println("☸️  Kubernetes Deployment")

	namespace, _ := ctx.GetValueAsString("namespace")
	replicas := ctx.GetValueOrDefault("replicas", "3")
	context_name := ctx.GetValueOrDefault("context", "default")

	wait := ctx.GetValueAsBool("wait")

	fmt.Printf("🏷️  Namespace: %s\n", namespace)
	fmt.Printf("🔢 Replicas: %s\n", replicas)
	fmt.Printf("🎯 Context: %s\n", context_name)

	if wait {
		fmt.Println("⏳ Waiting for deployment to complete...")
	}

	// Handle optional resource configuration
	handleResourceConfiguration(ctx)

	// Handle optional monitoring configuration
	handleMonitoringConfiguration(ctx)

	fmt.Println("✅ Kubernetes deployment completed!")
	return nil
}

func handleServerlessDeploy(ctx ctx.Context) error {
	fmt.Println("⚡ Serverless Deployment")

	provider, _ := ctx.GetValueAsString("provider")
	region := ctx.GetValueOrDefault("region", "us-east-1")
	runtime := ctx.GetValueOrDefault("runtime", "nodejs18.x")
	memory := ctx.GetValueOrDefault("memory", "128")
	timeout := ctx.GetValueOrDefault("timeout", "30")

	fmt.Printf("☁️  Provider: %s\n", provider)
	fmt.Printf("🌍 Region: %s\n", region)
	fmt.Printf("🔧 Runtime: %s\n", runtime)
	fmt.Printf("💾 Memory: %s MB\n", memory)
	fmt.Printf("⏰ Timeout: %s seconds\n", timeout)

	// Handle optional monitoring configuration
	handleMonitoringConfiguration(ctx)

	fmt.Println("✅ Serverless deployment completed!")
	return nil
}

func handleResourceConfiguration(ctx ctx.Context) {
	if ctx.IsFlagExist("resources") {
		fmt.Println("📊 Resource Configuration:")

		if ctx.IsFlagHaveValue("memory") {
			memory, _ := ctx.GetValueAsString("memory")
			fmt.Printf("  💾 Memory: %s MB\n", memory)
		}

		if ctx.IsFlagHaveValue("cpu") {
			cpu, _ := ctx.GetValueAsString("cpu")
			fmt.Printf("  🖥️  CPU: %s cores\n", cpu)
		}

		if ctx.IsFlagHaveValue("storage") {
			storage, _ := ctx.GetValueAsString("storage")
			fmt.Printf("  💿 Storage: %s\n", storage)
		}
	}
}

func handleMonitoringConfiguration(ctx ctx.Context) {
	if ctx.IsFlagExist("monitoring") {
		fmt.Println("📈 Monitoring Configuration:")

		metrics := ctx.GetValueAsBool("metrics")
		logging := ctx.GetValueAsBool("logging")

		fmt.Printf("  📊 Metrics: %v\n", metrics)
		fmt.Printf("  📝 Logging: %v\n", logging)

		if ctx.IsFlagHaveValue("alerts") {
			alerts, _ := ctx.GetValueAsString("alerts")
			fmt.Printf("  🚨 Alerts: %s\n", alerts)
		}

		if ctx.IsFlagHaveValue("dashboard") {
			dashboard, _ := ctx.GetValueAsString("dashboard")
			fmt.Printf("  📊 Dashboard: %s\n", dashboard)
		}
	}
}

func backupHandler(ctx ctx.Context) error {
	fmt.Println("💾 Creating Application Backup")

	verbose := ctx.GetValueAsBool("verbose")
	output := ctx.GetValueOrDefault("output", "backup.tar.gz")

	fmt.Printf("📄 Output file: %s\n", output)

	// Check which exclusive storage group is being used
	if ctx.IsFlagExist("local") {
		return handleLocalBackup(ctx, verbose)
	} else if ctx.IsFlagExist("remote") {
		return handleRemoteBackup(ctx, verbose)
	} else if ctx.IsFlagExist("s3") {
		return handleS3Backup(ctx, verbose)
	}

	return fmt.Errorf("no backup destination specified")
}

func handleLocalBackup(ctx ctx.Context, verbose bool) error {
	fmt.Println("🏠 Local Backup")

	path, _ := ctx.GetValueAsString("path")
	compress := ctx.GetValueAsBool("compress")
	encryption := ctx.GetValueOrDefault("encryption", "none")

	fmt.Printf("📁 Path: %s\n", path)
	fmt.Printf("🗜️  Compress: %v\n", compress)
	fmt.Printf("🔐 Encryption: %s\n", encryption)

	if verbose {
		fmt.Println("📋 Backup process:")
		fmt.Println("  → Creating archive...")
		fmt.Println("  → Adding application files...")
		fmt.Println("  → Adding database dump...")
		if compress {
			fmt.Println("  → Compressing archive...")
		}
		if encryption != "none" {
			fmt.Println("  → Encrypting backup...")
		}
	}

	fmt.Println("✅ Local backup completed!")
	return nil
}

func handleRemoteBackup(ctx ctx.Context, verbose bool) error {
	fmt.Println("🌐 Remote Backup")

	host, _ := ctx.GetValueAsString("host")
	user, _ := ctx.GetValueAsString("user")
	key := ctx.GetValueOrDefault("key", "~/.ssh/id_rsa")
	port := ctx.GetValueOrDefault("port", "22")
	remotePath := ctx.GetValueOrDefault("path", "/backup")

	fmt.Printf("🏠 Host: %s:%s\n", host, port)
	fmt.Printf("👤 User: %s\n", user)
	fmt.Printf("🔑 Key: %s\n", key)
	fmt.Printf("📁 Remote path: %s\n", remotePath)

	if verbose {
		fmt.Println("📋 Remote backup process:")
		fmt.Println("  → Establishing SSH connection...")
		fmt.Println("  → Creating local backup...")
		fmt.Println("  → Transferring to remote server...")
		fmt.Println("  → Verifying transfer...")
	}

	fmt.Println("✅ Remote backup completed!")
	return nil
}

func handleS3Backup(ctx ctx.Context, verbose bool) error {
	fmt.Println("☁️  AWS S3 Backup")

	bucket, _ := ctx.GetValueAsString("bucket")
	region, _ := ctx.GetValueAsString("region")
	accessKey := ctx.GetValueOrDefault("access-key", "from-environment")
	prefix := ctx.GetValueOrDefault("prefix", "backups/")

	fmt.Printf("🪣 Bucket: %s\n", bucket)
	fmt.Printf("🌍 Region: %s\n", region)
	fmt.Printf("🔑 Access Key: %s\n", accessKey)
	fmt.Printf("📂 Prefix: %s\n", prefix)

	if verbose {
		fmt.Println("📋 S3 backup process:")
		fmt.Println("  → Creating local backup...")
		fmt.Println("  → Configuring AWS credentials...")
		fmt.Println("  → Uploading to S3...")
		fmt.Println("  → Setting object metadata...")
		fmt.Println("  → Verifying upload...")
	}

	fmt.Println("✅ S3 backup completed!")
	return nil
}
