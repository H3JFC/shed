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
