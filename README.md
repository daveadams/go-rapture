# rapture

Rapture is a shell-integrated CLI tool for assuming AWS IAM roles easily and
quickly.

## NOTICE

The golang version of Rapture is still in the initial development stages. Be
careful! The docs below are incomplete and partially wrong!

## Usage Example

    $ rapture whoami
    arn:aws:iam::999988887777:user/janesmith

    $ rapture alias set admin arn:aws:iam::000011110000:role/admin-power
    rapture: alias 'admin' was set to 'arn:aws:iam::000011110000:role/admin-power'

    $ rapture alias ls
    admin      arn:aws:iam::000011110000:role/admin-power
    marketing  arn:aws:iam::302830283028:role/marketing-access

    $ rapture assume admin
    rapture: Assumed assumed-role admin-power in account 000011110000

    $ rapture whoami
    arn:aws:sts::000011110000:assumed-role/admin-power/rapture-janesmith

    $ rapture resume
    rapture: Resumed user janesmith in account 999988887777

    $ rapture whoami
    arn:aws:iam::999988887777:user/janesmith


## Prerequisites

* A machine that can run golang binaries
* Bash, zsh, or fish


## Installation

First, install the latest Rapture binary for your platform [from Github](https://github.com/daveadams/go-rapture/releases)
and copy it to a directory in your PATH. Or if you prefer you can build it yourself:

    $ go get -u github.com/daveadams/go-rapture/cmd/rapture

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

Simply remove the old `source ~/.rapture/rapture.sh` from your shell startup
script, and replace it with the new command. The old Rapture configuration files
will continue to work in the same way.


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


## Caveats

Rapture assumes the use of `AWS_*` environment variables for determining the root identity.

Rapture does _not_ manage your secrets for you. I recommend [Vaulted](https://github.com/miquella/vaulted) for managing storing AWS access keys (and other secrets) securely in an easily manageable format and for loading them into your environment.

## License

This software is public domain. No rights are reserved. See LICENSE for more information.
