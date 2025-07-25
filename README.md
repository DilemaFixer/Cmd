# Cmd
Simpl console command parser


```
command: "backup",
subcommands: ["one", "two"]
option_groups: {
  storage: {
    required: true,           
    exclusive: true,          
    groups: {
      local: {
        trigger: "--local",  
        options: ["--path", "--compress", "--encrypt"]
      },
      remote: {
        trigger: "--remote", 
        options: ["--host", "--user", "--key", "--port"]
      },
      cloud: {
        trigger: "--s3",
        options: ["--bucket", "--region", "--access-key"]
      }
    }
  }
}
```
