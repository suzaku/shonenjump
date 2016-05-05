# Shonenjump

shonenjump is a lightweight autojump clone written in Go.

# What does it do?

Quote from the description of [autojump](https://github.com/wting/autojump/):

> autojump is a faster way to navigate your filesystem. It works by maintaining a database of the directories you use the most from the command line.

> Directories must be visited first before they can be jumped to.

# Installation

1. Download the shonenjump binary for your platform, place it in a directory in your `$PATH`.
1. Download the setup [script](https://github.com/suzaku/shonenjump/blob/master/scripts/) for your shell and include it in your shell profile.

   For example, if you are using `zsh`, you can do the following:
    
   ```bash
   wget -O ~/.shonenjump.zsh https://github.com/suzaku/shonenjump/blob/master/scripts/shonenjump.zsh
   echo 'source $HOME/.shonenjump.zsh' >> ~/.zshrc
```

1. If you are using `zsh`, you'll need an extra step to setup tab completion.

   You need to place a script into the `zsh/site-functions` directory:
   ```bash
   cd $(brew --prefix)/share/zsh/site-functions/
   wget https://github.com/suzaku/shonenjump/blob/master/scripts/_j
   ```

# Usage

Usage is the same as [autojump](https://github.com/wting/autojump/#usage)
