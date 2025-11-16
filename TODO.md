# TODO LIST

- [] Require name to be snake_case and restrictions on characters
- [] Add Description to commands
- [] Each Parameter gets a description, parameters must be snake_case too
- [] Each Parameter is a required parameter
- [] Interpolation Engine
- [] How do I capture optionally describing each parameter?

## Commands TODO

- [] "shed add command <NAME> --description|-d <CLI_COMMAND>"
  - [] "ERROR ON DUPLICATE NAME"
- [] "shed rm command <NAME>"
  - [] "ERROR ON NAME DOES NOT EXIST"
- [] "shed edit command <NAME> --description|-d <CLI_COMMAND>"
  - [] "ERROR ON NAME DOES NOT EXIST"
- [] "shed describe command <NAME>"
  - [] "ERROR ON NAME DOES NOT EXIST"
- [] "shed list commands"

## Examples

shed add command post_server curl -XPOST --data '{"name":"hello world","password": {{ $PASSWORD }},"email":"hector.friedman.cintron@gmail.com"}' -w "%{http_code}" http://localhost:8000/auth/register

shed add command post_server curl -XPOST --data '{"name":"hello world","password": {{ desc($PASSWORD, "foo bar's company password") }},"email":"hector.friedman.cintron@gmail.com"}' -w "%{http_code}" http://localhost:8000/auth/register

shed add command post_server curl -XPOST --data '{"name":"hello world","password": {{ with_default($PASSWORD, "default_password") }},"email":"hector.friedman.cintron@gmail.com"}' -w "%{http_code}" http://localhost:8000/auth/register

### Simple Add command

```zsh
# "add" adds a command name post_server with password and server as variables and secret_token as a token
shed add cmd post_server curl -XPOST --data '{"foo":"bar"}' -H Authentication:{{password}} {{server}} -H Token:{{!secret_token}}
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
