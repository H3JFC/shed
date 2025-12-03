# TODO LIST

- [x] Interpolation Engine
- [x] Require name to be snake_case and restrictions on characters
- [x] Name less than 32 char
- [x] Each Parameter gets a description
- [x] Each Parameter is considered required
- [x] How do I capture optionally describing each parameter?
- [x] Add Description to Commands table
- [x] SQLC in Github Actions

## Commands TODO

### CREATE

- [x] "shed add <NAME> --description|-d <COMMAND_DESCRIPTION> <CLI_COMMAND>"
- [x] "shed cp <COMMAND_SRC_NAME> <COMMAND_DEST_NAME> {jsonValueParams}" # in form {[<param>]:<value>}
- [x] "shed rm command <NAME>"
- [x] "shed list"
- [x] "shed describe <COMMAND_NAME>"
- [x] "shed edit <COMMAND_NAME> --description|-d <COMMAND_DESCRIPTION> --name|-n <NEW_COMMAND_NAME> <CLI_COMMAND> {jsonValueParams}"

### Review

- [x] shed add
- [x] shed cp
- [x] shed rm
- [x] shed list
- [x] shed edit
- [x] shed describe

## Examples

```zsh
# "add" adds a command name post_server with password and server as variables and secret_token as a token
# ignore implementation of secret_token for now
# unique by name
shed add post_server curl -XPOST --data '{"foo":"bar"}' -H Authentication:{{password|"password for service"}} {{server|"url for service"}} -H Token:{{!secret_token}}
```

```zsh
# "cp" copies a command into a new command with the new name... kwargs at the end interpolate any parameters only, secret names provided in this context throw an error
shed cp post_server post_server_2 {"server":"foobar.com"} # valid

shed cp post_server post_server_2 {"server":"foobar.com", "secret_token":"test-token"} # invalid & throws error

```

```zsh
# 'rm' removes the created command
shed rm post_server     # prompts user are you sure y/n?
shed rm post_server --f # force doesn't prompt
shed rm does_not_exist  # throws error
```

```zsh
# 'describe' shows the in the $EDITOR of choice or a default in read-only mode
shed describe post_server     # shows up in a buffered editor window
shed describe does_not_exist  # throws error
```

```zsh
# 'edit' shows the in the $EDITOR of choice or a default in edit mode, on close the updates are saved
shed edit post_server     # shows up in a buffered editor window
shed edit does_not_exist  # throws error
```

```zsh
# 'list' shows the lists of commands in name |
shed list
```

### How it works

#### Commands

- add
- rm
- cp
- describe
- edit

#### SQL

- Add "configuration" column JSONB_ARRAY with migration
  - https://sqlite.org/json1.html

- Add SQLc or Go Struct to handle config

  ```json
  [
    {
      "name": "param_name_1",
      "description": "param_desc_1"
    },
    {
      "name": "param_name_2",
      "description": "param_desc_2"
    }
  ]
  ```

#### Lower Level functionality

- parse {{parameter_name|optional description}} into go struct of

  ```Go
  type Parameter struct {
      Name string
      Description string
  }
  ```

- Update SQLC queries
- Open up Command in buffer with chosen $editor
  - $editor determination
  - temporary file that get read on close and removed once a save succeeds

- rm user prompting for confirm remove

#### SQLc

```
# AI Related response
3. Using JSONB in SQLite (Advanced)
Although SQLite supports jsonb(), it's not recommended for direct use in applications. Instead, use json() or json_array() and let SQLite handle storage as text.

If you do use JSONB:

-- Insert JSONB blob
INSERT INTO items (data) VALUES (jsonb('[1,2,3]'));

But ensure your Go scanner can handle []byte. Again, use sqlc overrides:

overrides:
  - column: "data"
    go_type: "json.RawMessage"

And scan into a []byte or json.RawMessage field.
```

## Future

- MCP Groups in MCP server so you don't host all of your tools
- Client Server Grouping
  - One Machine to host MCP tools
