# change pwd hook
shonenjump_chpwd() {
    if [[ -f "${SHONENJUMP_ERROR_PATH}" ]]; then
        shonenjump --add "$(pwd)" >/dev/null 2>>${SHONENJUMP_ERROR_PATH} &!
    else
        shonenjump --add "$(pwd)" >/dev/null &!
    fi
}

typeset -gaU chpwd_functions
chpwd_functions+=shonenjump_chpwd


# default shonenjump command
j() {
    if [[ ${1} == -* ]] && [[ ${1} != "--" ]]; then
        shonenjump ${@}
        return
    fi

    setopt localoptions noautonamedirs
    local output="$(shonenjump ${@})"
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

    setopt localoptions noautonamedirs
    local output="$(shonenjump ${@})"
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
                echo "Unknown operating system: ${OSTYPE}" 1>&2
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
