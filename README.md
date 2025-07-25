# Cmd
Simpl console command parser

TODO: add group item settings like if --local is set you can't use any another group or group is undependent 
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
