# Shonenjump ![build](https://github.com/suzaku/shonenjump/workflows/build/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/suzaku/shonenjump)](https://goreportcard.com/report/github.com/suzaku/shonenjump) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/suzaku/shonenjump)

shonenjump is a lightweight autojump clone written in Go.


<a href="https://www.buymeacoffee.com/suzaku" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>


# What does it do?

Quote from the description of [autojump](https://github.com/wting/autojump/):

> autojump is a faster way to navigate your filesystem. It works by maintaining a database of the directories you use the most from the command line.

> Directories must be visited first before they can be jumped to.

# How to use it?

Once you have `cd` into a directory, `shonenjump` will save it in a list.
The next time you can use the `j` shortcut to visit it.

For example, suppose that you have `cd` into a directory called `/usr/local/Very-Long-Dir-Name/Sub-Dir/target` after
`shonenjump` is enabled. You can then use `j long` or `j target` or `j vldn` to visit it.

Sometimes the first matched directory is not what you want, you can type `j <your key word>` and
then type Tab to trigger auto completion and see the options.

# Installation

## macOS

`brew install suzaku/homebrew-shonenjump/shonenjump`

## Linux

1. [Download](https://github.com/suzaku/shonenjump/releases) the shonenjump binary for your platform, place it in a directory in your `$PATH`.
1. [Download](https://github.com/suzaku/shonenjump/blob/master/scripts/) the setup script for your shell and include it in your shell profile.

   For example, if you are using `zsh`, you can do the following:

   ```bash
   wget -O ~/.shonenjump.zsh https://raw.githubusercontent.com/suzaku/shonenjump/master/scripts/shonenjump.zsh
   echo '. $HOME/.shonenjump.zsh' >> ~/.zshrc
    ```

1. If you are using `zsh`, you'll need an extra step to setup tab completion.

   You need to place a script into the `zsh/site-functions` directory:
   ```bash
   cd <Your Zsh Site-functions Dir>
   wget https://raw.githubusercontent.com/suzaku/shonenjump/master/scripts/_j
   ```
