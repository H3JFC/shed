module h3jfc/shed

go 1.25.1

require (
	github.com/golang-migrate/migrate/v4 v4.19.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/spf13/cobra v1.10.1
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
)

replace github.com/mattn/go-sqlite3 => github.com/jgiannuzzi/go-sqlite3 v1.14.17-0.20230327162135-f208443ec79d
