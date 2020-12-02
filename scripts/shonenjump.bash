export SHONENJUMP_SOURCED=1

# set error file location
if [[ "$(uname)" == "Darwin" ]]; then
    export SHONENJUMP_ERROR_PATH=~/Library/shonenjump/errors.log
elif [[ -n "${XDG_DATA_HOME}" ]]; then
    export SHONENJUMP_ERROR_PATH="${XDG_DATA_HOME}/shonenjump/errors.log"
else
    export SHONENJUMP_ERROR_PATH=~/.local/share/shonenjump/errors.log
fi

if [[ ! -d "$(dirname ${SHONENJUMP_ERROR_PATH})" ]]; then
    mkdir -p "$(dirname ${SHONENJUMP_ERROR_PATH})"
fi


# enable tab completion
_shonenjump() {
        local cur
        cur=${COMP_WORDS[*]:1}
        comps=$(shonenjump --complete $cur)
        while read i; do
            COMPREPLY=("${COMPREPLY[@]}" "${i}")
        done <<EOF
        $comps
EOF
}
complete -F _shonenjump j


# change pwd hook
shonenjump_add_to_database() {
    if [[ -f "${SHONENJUMP_ERROR_PATH}" ]]; then
        (shonenjump --add "$(pwd)" >/dev/null 2>>${SHONENJUMP_ERROR_PATH} &) &>/dev/null
    else
        (shonenjump --add "$(pwd)" >/dev/null &) &>/dev/null
    fi
}

case $PROMPT_COMMAND in
    *shonenjump*)
        ;;
    *)
        PROMPT_COMMAND="${PROMPT_COMMAND:+$(echo "${PROMPT_COMMAND}" | awk '{gsub(/; *$/,"")}1') ; }shonenjump_add_to_database"
        ;;
esac


# default shonenjump command
j() {
    if [[ ${1} == -* ]] && [[ ${1} != "--" ]]; then
        shonenjump ${@}
        return
    fi

    output="$(shonenjump ${@})"
    if [[ -d "${output}" ]]; then
				if [ -t 1 ]; then  # if stdout is a terminal, use colors
						echo -e "\\033[31m${output}\\033[0m"
				else
						echo -e "${output}"
				fi
        cd "${output}"
    else
        echo "shonenjump: directory '${@}' not found"
        echo "\n${output}\n"
        echo "Try \`shonenjump --help\` for more information."
        false
    fi
}


# jump to child directory (subdirectory of current path)
jc() {
    if [[ ${1} == -* ]] && [[ ${1} != "--" ]]; then
        shonenjump ${@}
        return
    else
        j $(pwd) ${@}
    fi
}


# open shonenjump results in file browser
jo() {
    if [[ ${1} == -* ]] && [[ ${1} != "--" ]]; then
        shonenjump ${@}
        return
    fi

    output="$(shonenjump ${@})"
    if [[ -d "${output}" ]]; then
        case ${OSTYPE} in
            linux*)
                xdg-open "${output}"
                ;;
            darwin*)
                open "${output}"
                ;;
            cygwin)
                cygstart "" $(cygpath -w -a ${output})
                ;;
            *)
                echo "Unknown operating system: ${OSTYPE}." 1>&2
                ;;
        esac
    else
        echo "shonenjump: directory '${@}' not found"
        echo "\n${output}\n"
        echo "Try \`shonenjump --help\` for more information."
        false
    fi
}


# open shonenjump results (child directory) in file browser
jco() {
    if [[ ${1} == -* ]] && [[ ${1} != "--" ]]; then
        shonenjump ${@}
        return
    else
        jo $(pwd) ${@}
    fi
}
