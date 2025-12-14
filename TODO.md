# Features

## Brackets Tech Debt

- [x] Combine ParseCommand & ParseParameters & ParseSecrets(tbd) to return (\*Bracket, error)
  - [x] Add Secret type (Maybe combine) due to similarities to Parameter / ValuedParameter
    ```go
        type Secret struct {
            Key         string
            Description string
        }
        type ValuedSecret struct {
            Key         string
            Value       string
        }
    ```
  - [x] Create Brackets type
    ```go
        type Brackets struct {
            Command     string
            Parameters  Parameters
            Secrets     Secrets
        }
    ```

## Secrets

- [x] Interpolate {{!secret|description}} as a valid secret Parameters
  - [x] Brackets.ParseSecret
- [x] "shed add secret --description|-d <KEY> <VALUE>"
  - [x] store.AddSecret
- [x] "shed list secrets"
  - [x] store.ListSecrets
- [x] "shed edit secret --description|-d <KEY> <VALUE>"
  - [x] store.UpdateSecret
- [x] "shed rm secret <KEY>"
  - [x] store.RemoveSecret

## Next

- [x] Modify shed add|edit to check for new secrets
  - Warn user if secret does not exist yet
  - Prompt user to create secret with secret command

## Run Command

- [ ] "shed run <COMMAND_NAME>"
  - [ ] Execute Library to split command and then execute
- [x] update shed command describe <COMMAND_NAME> to list secrets
