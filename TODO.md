# Features

## Secrets

- [ ] "shed add secret --description|-d <KEY> <VALUE>"
- [ ] "shed list secrets"
- [ ] "shed edit secret --description|-d <KEY> <VALUE>"
- [ ] "shed rm secret <KEY>"
- [ ] Interpolate {{!secret|description}} as a valid secret Parameters
  - Add Secret and ValuedSecrets types
  - Brackets.ParseSecret
- [ ] Modify shed add|edit to check for new secrets
  - Warn user if secret does not exist yet
  - Prompt user to create secret with secret command

## Brackets Tech Debt

- [ ] Combine ParseCommand & ParseParameters & ParseSecrets(tbd) to return (\*Bracket, error)
  - [ ] Add Secret type (Maybe combine) due to similarities to Parameter / ValuedParameter
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
  - [ ] Create Brackets type
    ```go
        type Brackets struct {
            Command     string
            Parameters  Parameters
            Secrets     Secrets
        }
    ```
