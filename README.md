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

First, install the Rapture binary.

(TODO: binary instructions here, for building from source or download.)

Then configure your shell to load Rapture at start:

    $ echo 'eval "$( rapture setup )"' >> ~/.bash_profile

*NOTE:* You may need to add this line to `~/.bashrc` instead on some systems.

(TODO: Specify setup for other shells.)

Finally, open a new terminal window to verify that Rapture is automatically loaded:

    $ rapture check
    OK: Rapture is set up correctly


## Configuration

No configuration is required to start using Rapture, but Rapture will store configuration in `config.json`, `aliases.json`, and `accounts.json` in the `~/.rapture` directory.


## Environment Variables

Rapture sets the `RAPTURE_ROLE` environment variable with the role ARN or alias of the currently-assumed role.


## Caveats

Rapture assumes the use of `AWS_*` environment variables for determining the root identity.

Rapture does _not_ manage your secrets for you. I recommend [Vaulted](https://github.com/miquella/vaulted) for managing storing AWS access keys (and other secrets) securely in an easily manageable format and for loading them into your environment.

## License

This software is public domain. No rights are reserved. See LICENSE for more information.
