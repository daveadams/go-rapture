# rapture

Rapture is a shell-integrated CLI tool for assuming AWS IAM roles easily and
quickly.

## Usage Example

    $ rapture whoami
    arn:aws:iam::999988887777:user/janesmith

    $ rapture alias set admin arn:aws:iam::000011110000:role/admin-power
    alias 'admin' was set to 'arn:aws:iam::000011110000:role/admin-power'

    $ rapture alias ls
    arn:aws:iam::000011110000:role/admin-power admin
    arn:aws:iam::302830283028:role/marketing-access marketing

    $ rapture assume admin
    Assumed role 'arn:aws:iam::000011110000:role/admin-power'

    $ rapture whoami
    arn:aws:sts::000011110000:assumed-role/admin-power/rapture-janesmith

    $ rapture resume
    Resumed base credentials

    $ rapture whoami
    arn:aws:iam::999988887777:user/janesmith


## Prerequisites

* A machine that can run golang binaries
* Bash, zsh, or fish


## Installation

First, install the latest Rapture binary for your platform [from Github](https://github.com/daveadams/go-rapture/releases)
and copy it to a directory in your PATH. Or if you prefer you can build it yourself (Go 1.12 or higher required):

    $ go install github.com/daveadams/go-rapture/cmd/rapture@latest

To configure your shell to load Rapture at startup, follow the instructions
below for your specific shell.

### Bash/Zsh Installation

In your shell startup file (usually `~/.bash_profile` or `~/.bashrc`; for Zsh
usually `~/.zshrc`), add the following line:

    eval "$( command rapture shell-init )"

### Fish Installation

In your Fish startup file (`~/.config/fish/config.fish`), add the following line:

    eval ( command rapture shell-init )

### Verifying the setup

Finally, open a new terminal window to verify that Rapture is automatically loaded:

    $ rapture check
    OK: Rapture is set up correctly


## Upgrading from Bash Rapture

Simply remove the command `source ~/.rapture/rapture.sh` from your shell startup
script, and replace it with the new command mentioned above. The old Rapture
configuration files will continue to work in the same way.


## Configuration

No configuration is required to start using Rapture, but Rapture will store
configuration in `config.json`, `aliases.json`, and `accounts.json` in the
`~/.rapture` directory.


## Environment Variables

Rapture exports the `RAPTURE_ROLE` environment variable with the user-supplied
identifier of the currently-assumed role, either the role alias, or the ARN.

Rapture also exports the `RAPTURE_ASSUMED_ROLE_ARN` environment variable to
the full ARN of the currently assumed role.

Both of these environment variables are unset when base credentials are loaded.


## Vaulted integration

Rapture does _not_ manage your secrets for you. But Rapture is integrated with
[Vaulted](https://github.com/miquella/vaulted) for providing secure storage of
AWS access keys (and other secrets) in an easily manageable format and for
loading them into your environment.

If you have your AWS credentials configured in a vault named `default`, then
you can simply run:

    $ rapture init

This will run `vaulted` on your behalf to load the credentials from the `default`
Vault into your current environment. Or you can specify a different vault name:

    $ rapture init awsvault

If you set a value for `default_vault` in `~/.rapture/config.json`, Rapture
will use that name instead of `default` as the default vault to decrypt.


## License

This software is public domain. No rights are reserved. See LICENSE for more information.
