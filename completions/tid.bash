# tid.bash - tid completions for bash shell.
# 
# To install the completions:
# @todo

# Get tid args, throw away any options, throw away tid command name.
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

_tid() {
    local cur prev

    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    RAW_PATH="$(_tid_get_args)"

    # At this point, RAW_PATH should only contain the things that aren't the current word we're on,
    # or the command name itself, or any options. It should just be full arguments that are left.

    echo ""
    echo "prev: $prev"
    echo "cur: $cur"
    echo "COMP_CWORD: $COMP_CWORD"
    echo "COMP_LINE: $COMP_LINE"
    echo "COMP_POINT: $COMP_POINT"
    echo "COMP_WORDBREAKS: $COMP_WORDBREAKS"
    echo "COMP_WORDS: $COMP_WORDS"
    echo "RAW_PATH: $RAW_PATH"

    # case "${prev}" in
    #
    # esac
}

complete -F _tid tid
