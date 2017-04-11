# tid.bash - tid completions for bash shell.
# 
# To install the completions:
# @todo

_tid() {
    local cur prev opts

    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="entry report resume start status stop timesheet workspace"

    if [[ ${cur} == -* || ${COMP_CWORD} -eq 1 ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi

    # case "${prev}" in
    #
    # esac
}

complete -F _tid tid
