# MCP-SERVER

## Features

### Tech Debt

- [x] Automatically increment tag for release job
- [x] CI badges
- [x] Shed Logo

### MCP

- [ ] Modify Logger so that there are no colors and it logs directly to stdout
- [ ] Daemonize: https://ieftimov.com/posts/four-steps-daemonize-your-golang-programs/
  - [ ] Daemon Config
  - [ ] Reload config on SIGHUP
  - [ ] Kill on SIGINT or SIGKILL
- [ ] create shed mcp command
  - [ ] shed mcp --port <PORT>
- [ ] on-startup / reload pulls all commands from database and serves them as tools
  - [ ] install https://github.com/mark3labs/mcp-go
- [ ] whenever a command is run via "shed", it touches a server.reload file
  - [ ] watch server.reload file, on change, reloadCommands

  ```go
  // Option 2: Touch a file + fsnotify (clean & fast)
  // How it works

  // Command updates SQLite

  // Command also touches a file (reload.signal)

  // Server watches the file and reloads immediately

  // Command side
  // touch /tmp/server.reload

  // Server side (Go)
  watcher, _ := fsnotify.NewWatcher()
  watcher.Add("/tmp/server.reload")

  go func() {
      for {
          select {
          case <-watcher.Events:
              reloadCommands()
          case err := <-watcher.Errors:
              log.Println("watch error:", err)
          }
      }
  }()
  ```
