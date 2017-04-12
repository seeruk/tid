#!/usr/bin/env bash

# tid.bash - tid completions for bash shell.
# 
# To install the completions:
# @todo
#
# Concerns in order:
# - command and sub-command completion
# - entry hash, timesheet date, and workspace name completion
# - option completion

# _tid_get_args takes the current command line, throws away any options, throws away the tid command
# name, and then echoes the resulting words on new lines for consumption with readarray.
_tid_get_args() {
    # Create array from COMP_LINE
    IFS=", " read -r -a RAW_PATH <<< "$COMP_LINE"

    # Remove `tid` part
    unset 'RAW_PATH[0]'

    # Remove current word if we're in the middle of typing something
    if [ "$cur" != "" ] && [ "${RAW_PATH[-1]}" == "$cur" ]; then
        unset 'RAW_PATH[-1]'
    fi

    # Remove options and echo result
    for i in "${!RAW_PATH[@]}"; do
        if [[ "${RAW_PATH[$i]}" != -* ]]; then
            echo "${RAW_PATH[$i]}"
        fi
    done
}

_tid_paths() {
    local path next

    path="$1"
    next="$2"

    case "$path" in
        # Root and it's subcommands
        "")
            case "$next" in
                "entry"|"e")
                    echo "entry"
                ;;
                "report"|"rep")
                    echo "report"
                ;;
                "resume"|"res")
                    echo "resume"
                ;;
                "start")
                    echo "start"
                ;;
                "status"|"st")
                    echo "status"
                ;;
                "stop")
                    echo "stop"
                ;;
                "timesheet"|"t")
                    echo "timesheet"
                ;;
                "workspace"|"w")
                    echo "workspace"
                ;;
                *)
                    echo "$path"
                ;;
            esac
        ;;
        # Commands with subcommands
        "entry")
            case "$next" in
                "create"|"c")
                    echo "$path create"
                ;;
                "delete"|"d")
                    echo "$path delete"
                ;;
                "list"|"ls")
                    echo "$path list"
                ;;
                "update"|"u")
                    echo "$path update"
                ;;
                *)
                    echo "$path"
                ;;
            esac
        ;;
        "timesheet")
            case "$next" in
                "delete"|"d")
                    echo "$path delete"
                ;;
                "list"|"ls")
                    echo "$path list"
                ;;
                *)
                    echo "$path"
                ;;
            esac
        ;;
        "workspace")
            case "$next" in
                "create"|"c")
                    echo "$path create"
                ;;
                "delete"|"d")
                    echo "$path delete"
                ;;
                "list"|"ls")
                    echo "$path list"
                ;;
                "switch"|"s")
                    echo "$path switch"
                ;;
                *)
                    echo "$path"
                ;;
            esac
        ;;
        *)
            echo "$path"
        ;;
    esac
}

_tid_is_at_path() {
    local args idx path

    readarray -t args <<< "$(_tid_get_args)"

    idx=1
    path=""

    for i in "${!args[@]}"; do
        path="$(_tid_paths "$path" "${args[$i]}")"
    done

    if [ "$path" == "$1" ]; then
        return 0
    else
        return 1
    fi
}

# _tid provides completions for tid.
_tid() {
    local cur prev

    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    readarray -t RAW_PATH <<< "$(_tid_get_args)"

    # At this point, RAW_PATH should only contain the things that aren't the current word we're on,
    # or the command name itself, or any options. It should just be full arguments that are left.

    if $(_tid_is_at_path "entry update"); then
        echo "AT PATH 'entry update'"
    else
        echo "NOT AT PATH 'entry update'"
    fi

    echo ""
    echo "prev: $prev"
    echo "cur: $cur"
    echo "COMP_CWORD: $COMP_CWORD"
    echo "COMP_LINE: $COMP_LINE"
    echo "COMP_POINT: $COMP_POINT"
    echo "COMP_WORDBREAKS: $COMP_WORDBREAKS"
    echo "COMP_WORDS: $COMP_WORDS"
    echo "RAW_PATH: ${RAW_PATH[@]}"

    # case "${prev}" in
    #
    # esac
}

complete -F _tid tid
