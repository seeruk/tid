#!/usr/bin/env bash

# tid.bash - tid completions for bash shell.
# 
# Bash completion installation seems to be quite platform dependant, you might just have to Google
# this one guys.
#
# Supports all commands, subcommands, aliases, options, and fills in arguments for entry hashes,
# timesheet dates, and workspace names.

# _tid_get_args takes the current command line, throws away any options, throws away the tid command
# name, and then echoes the resulting words on new lines for consumption with readarray.
_tid_get_args() {
    # Create array from COMP_LINE
    IFS=", " read -r -a RAW_PATH <<< "$COMP_LINE"

    # Remove `tid` part
    unset 'RAW_PATH[0]'

    # Get index of last item, so we can remove it if necessary.
    end=${#RAW_PATH[@]}

    # Remove current word if we're in the middle of typing something
    if [ "$cur" != "" ] && [ "${RAW_PATH[$end]}" == "$cur" ]; then
        unset "RAW_PATH[$end]"
    fi

    # Remove options and echo result
    for item in "${RAW_PATH[@]}"; do
        if [[ "$item" != -* ]]; then
            echo "$item"
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
                    echo "entry" ;;
                "report"|"rep")
                    echo "report" ;;
                "resume"|"res")
                    echo "resume" ;;
                "start")
                    echo "start" ;;
                "status"|"st")
                    echo "status" ;;
                "stop")
                    echo "stop" ;;
                "timesheet"|"t")
                    echo "timesheet" ;;
                "workspace"|"w")
                    echo "workspace" ;;
                *)
                    echo "$path" ;;
            esac
        ;;
        # Commands with subcommands
        "entry")
            case "$next" in
                "create"|"c")
                    echo "$path create" ;;
                "delete"|"d")
                    echo "$path delete" ;;
                "list"|"ls")
                    echo "$path list" ;;
                "update"|"u")
                    echo "$path update" ;;
                *)
                    echo "$path" ;;
            esac
        ;;
        "timesheet")
            case "$next" in
                "delete"|"d")
                    echo "$path delete" ;;
                "list"|"ls")
                    echo "$path list" ;;
                *)
                    echo "$path" ;;
            esac
        ;;
        "workspace")
            case "$next" in
                "create"|"c")
                    echo "$path create" ;;
                "delete"|"d")
                    echo "$path delete" ;;
                "list"|"ls")
                    echo "$path list" ;;
                "switch"|"s")
                    echo "$path switch" ;;
                *)
                    echo "$path" ;;
            esac
        ;;
        # Fall back to current path
        *)
            echo "$path" ;;
    esac
}

_tid_get_path() {
    local path

    path=""

    while read line; do
        path="$(_tid_paths "$path" "$line")"
    done <<< "$(_tid_get_args)"

    echo "$path"
}

_tid_entries() {
    echo $(command tid entry list --start=$(command tiddate --months=-6) --end=$(command tiddate) \
        --format="{{.ShortHash}}")
}

_tid_timesheets() {
    echo $(command tid timesheet list --start=$(command tiddate --years=-1) --end=$(command tiddate) \
        --format="{{.Key}}")
}

_tid_workspaces() {
    echo $(command tid w ls | awk '{ print $1 }' | sort)
}

# _tid provides completions for tid.
_tid() {
    local cur opts

    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=""

    case "$(_tid_get_path)" in
        # No command
        "")
            case "$cur" in
                -*)
                    opts="--help" ;;
                *)
                    opts="entry report resume start status stop timesheet workspace" ;;
            esac
        ;;
        # Commands
        "entry")
            opts="create delete list update" ;;
        "report")
            opts="--date -d --end -e --format -f --no-summary --start -s" ;;
        "resume")
            opts="$(_tid_entries)" ;;
        "status")
            case "$cur" in
                -*)
                    opts="--format -f" ;;
                *)
                    opts="$(_tid_entries)" ;;
            esac
        ;;
        "timesheet")
            opts="delete list" ;;
        "workspace")
            opts="create delete list switch" ;;
        # Sub-commands
        "entry create")
            opts="--date -d" ;;
        "entry delete")
            opts="$(_tid_entries)" ;;
        "entry list")
            opts="--date -d --end -e --format -f --start -s" ;;
        "entry update")
            case "$cur" in
                -*)
                    opts="--duration -d --note -n --offset -o" ;;
                *)
                    opts="$(_tid_entries)" ;;
            esac
        ;;
        "timesheet delete")
            opts="$(_tid_timesheets)" ;;
        "timesheet list")
            opts="--end -e --format -f --start -s" ;;
        "workspace delete")
            opts="$(_tid_workspaces)" ;;
        "workspace switch")
            opts="$(_tid_workspaces)" ;;
    esac

    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )

    return 0
}

complete -F _tid tid
