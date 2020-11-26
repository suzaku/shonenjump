set -gx AUTOJUMP_SOURCED 1

# Set ostype, if not set
if not set -q OSTYPE
    set -gx OSTYPE (bash -c 'echo ${OSTYPE}')
end

# set error file location
if test (uname) = "Darwin"
    set -gx AUTOJUMP_ERROR_PATH ~/Library/shonenjump/errors.log
else if test -d "$XDG_DATA_HOME"
    set -gx AUTOJUMP_ERROR_PATH $XDG_DATA_HOME/shonenjump/errors.log
else
    set -gx AUTOJUMP_ERROR_PATH ~/.local/share/shonenjump/errors.log
end

if test ! -d (dirname $AUTOJUMP_ERROR_PATH)
    mkdir -p (dirname $AUTOJUMP_ERROR_PATH)
end


# change pwd hook
function __aj_add --on-variable PWD
    status --is-command-substitution; and return
    shonenjump --add (pwd) >/dev/null 2>>$AUTOJUMP_ERROR_PATH &
end


# misc helper functions
function __aj_err
    # TODO(ting|#247): set error file location
    echo -e $argv 1>&2; false
end

# default shonenjump command
function j
    switch "$argv"
        case '-*' '--*'
            shonenjump $argv
        case '*'
            set -l output (shonenjump $argv)
            # Check for . and attempt a regular cd
            if test "$output" = "."
                cd $argv
            else
                if test -d "$output"
                    set_color red
                    echo $output
                    set_color normal
                    cd $output
                else
                    __aj_err "shonenjump: directory '"$argv"' not found"
                    __aj_err "\n$output\n"
                    __aj_err "Try `shonenjump --help` for more information."
                end
            end
    end
end


# jump to child directory (subdirectory of current path)
function jc
    switch "$argv"
        case '-*'
            j $argv
        case '*'
            j (pwd) $argv
    end
end


# open shonenjump results in file browser
function jo
    set -l output (shonenjump $argv)
    if test -d "$output"
        switch $OSTYPE
            case 'linux*'
                xdg-open (shonenjump $argv)
            case 'darwin*'
                open (shonenjump $argv)
            case cygwin
                cygstart "" (cygpath -w -a (pwd))
            case '*'
                __aj_err "Unknown operating system: \"$OSTYPE\""
        end
    else
        __aj_err "shonenjump: directory '"$argv"' not found"
        __aj_err "\n$output\n"
        __aj_err "Try `shonenjump --help` for more information."
    end
end


# open shonenjump results (child directory) in file browser
function jco
    switch "$argv"
        case '-*'
            j $argv
        case '*'
            jo (pwd) $argv
    end
end
