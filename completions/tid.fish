# tid.fish - tid completions for fish shell.
# 
# To install the completions:
# $ mkdir -p ~/.config/fish/completions
# $ cp tid.fish ~/.config/fish/completions

function __fish_tid_get_args --description 'Get tid args, throw away options. Echo on new lines'
    if test (count (commandline -opc)) -lt 2
        return 0
    end

    for i in (commandline -opc)[2..-1]
        # Skip options (they start with "-")
        if test (string sub -s1 -l1 -- $i) = "-"
            continue
        end

        echo "$i"
    end
end

function __fish_tid_paths
    set -l path $argv[1]
    set -l next $argv[2]

    switch $path
        # Root and it's subcommands
        case ""
            switch $next
                case "entry" "e"
                    echo "entry"
                case "report" "rep"
                    echo "report"
                case "resume" "res"
                    echo "resume"
                case "start"
                    echo "start"
                case "status" "st"
                    echo "status"
                case "stop"
                    echo "stop"
                case "timesheet" "t"
                    echo "timesheet"
                case "workspace" "w"
                    echo "workspace"
                case "*"
                    echo "$path"
            end
        # Commands with subcommands
        case "entry"
            switch $next
                case "create" "c"
                    echo "$path create"
                case "delete" "d"
                    echo "$path delete"
                case "list" "ls"
                    echo "$path list"
                case "update" "u"
                    echo "$path update"
                case "*"
                    echo "$path"
            end
        case "timesheet"
            switch $next
                case "delete" "d"
                    echo "$path delete"
                case "list" "ls"
                    echo "$path list"
                case "*"
                    echo "$path"
            end
        case "workspace"
            switch $next
                case "create" "c"
                    echo "$path create"
                case "delete" "d"
                    echo "$path delete"
                case "list" "ls"
                    echo "$path list"
                case "switch" "s"
                    echo "$path switch"
                case "*"
                    echo "$path"
            end
        case "*"
            echo "$path"
    end
end

function __fish_tid_is_at_path --description 'Test if the current tid commandline is at the given path or aliast'
    set -l args (__fish_tid_get_args)
    set -l idx 1
    set -l path ""

    # Build path in loop, we should be building full names. Make function that takes path, and next
    # path item, then works has a switch for all of the possible items below it? Maybe it returns
    # the new path?
    for a in $args
        set path (__fish_tid_paths "$path" "$a")
    end

    if test "$path" = "$argv"
        return 0
    else
        return 1
    end
end

# List entries from the last 6 months
function __fish_tid_entries
    command tid entry list --start=(tiddate --months=-6) --end=(tiddate) \
        --format="{{.ShortHash}}"\t"{{.Note}}"
end

function __fish_tid_timesheets
    command tid timesheet list --start=(tiddate --years=-1) --end=(tiddate) \
        --format="{{.Key}}"
end

function __fish_tid_workspaces
    command tid w ls | awk '{ print $1 }' | sort
end

# No command:
complete -c tid -n '__fish_tid_is_at_path ""' -a 'entry' -f -d 'Manage timesheet entries.'
complete -c tid -n '__fish_tid_is_at_path ""' -a 'report' -f -d 'Display a timesheet report.'
complete -c tid -n '__fish_tid_is_at_path ""' -a 'resume' -f -d 'Resume an existing timer.'
complete -c tid -n '__fish_tid_is_at_path ""' -a 'start' -f -d 'Start a new timer.'
complete -c tid -n '__fish_tid_is_at_path ""' -a 'status' -f -d 'View the current status.'
complete -c tid -n '__fish_tid_is_at_path ""' -a 'stop' -f -d 'Stop an existing timer.'
complete -c tid -n '__fish_tid_is_at_path ""' -a 'timesheet' -f -d 'Manage timesheets.'
complete -c tid -n '__fish_tid_is_at_path ""' -a 'workspace' -f -d 'Manage workspaces.'

# Commands
# entry, e
complete -c tid -n '__fish_tid_is_at_path "entry"' -a 'create' -f -d 'Create a new timesheet entry.'
complete -c tid -n '__fish_tid_is_at_path "entry"' -a 'delete' -f -d 'Delete a timesheet entry.'
complete -c tid -n '__fish_tid_is_at_path "entry"' -a 'list' -f -d 'List timesheet entries.'
complete -c tid -n '__fish_tid_is_at_path "entry"' -a 'update' -f -d 'Update a timesheet entry.'

# report
complete -c tid -n '__fish_tid_is_at_path "report"' -r -s e -l 'date=DATE' -f -d 'The exact date of a timesheet to show a report for.'
complete -c tid -n '__fish_tid_is_at_path "report"' -r -s e -l 'end=DATE' -f -d 'The end date of the report. (Default: today)'
complete -c tid -n '__fish_tid_is_at_path "report"' -r -s f -l 'format=STR' -f -d 'Output formatting string. Uses Go templates.'
complete -c tid -n '__fish_tid_is_at_path "report"' -l 'no-summary' -f -d 'Hide the summary?'
complete -c tid -n '__fish_tid_is_at_path "report"' -r -s s -l 'start=DATE' -f -d 'The start date of the report. (Default: today)'

# resume
complete -c tid -f -n '__fish_tid_is_at_path "resume"' -a '(__fish_tid_entries)'

# status
complete -c tid -n '__fish_tid_is_at_path "status"' -r -s f -l 'format=STR' -f -d 'Format string, uses Go templates.'
complete -c tid -n '__fish_tid_is_at_path "status"' -f -a '(__fish_tid_entries)'

# timesheet, t
complete -c tid -n '__fish_tid_is_at_path "timesheet"' -a 'delete' -f -d 'Delete a timesheet, and it\'s entries.'
complete -c tid -n '__fish_tid_is_at_path "timesheet"' -a 'list' -f -d 'List timesheets.'

# workspace, w
complete -c tid -n '__fish_tid_is_at_path "workspace"' -a 'create' -f -d 'Create a new workspace.'
complete -c tid -n '__fish_tid_is_at_path "workspace"' -a 'delete' -f -d 'Delete a workspace.'
complete -c tid -n '__fish_tid_is_at_path "workspace"' -a 'list' -f -d 'List available workspaces.'
complete -c tid -n '__fish_tid_is_at_path "workspace"' -a 'switch' -f -d 'Switch to another workspace.'

# Sub-commands
# entry create
complete -c tid -n '__fish_tid_is_at_path "entry create"' -r -s d -l 'date=DATE' -f -d 'When did you start working?'

# entry delete
complete -c tid -n '__fish_tid_is_at_path "entry delete"' -f -a '(__fish_tid_entries)'

# entry list
complete -c tid -n '__fish_tid_is_at_path "entry list"' -r -s e -l 'date=DATE' -f -d 'The exact date of a timesheet to show a listing for.'
complete -c tid -n '__fish_tid_is_at_path "entry list"' -r -s e -l 'end=DATE' -f -d 'The end date of the listing. (Default: today)'
complete -c tid -n '__fish_tid_is_at_path "entry list"' -r -s f -l 'format=STR' -f -d 'Output formatting string. Uses Go templates.'
complete -c tid -n '__fish_tid_is_at_path "entry list"' -r -s s -l 'start=DATE' -f -d 'The start date of the listing. (Default: today)'

# entry update
complete -c tid -n '__fish_tid_is_at_path "entry update"' -r -s d -l 'duration=DUR' -f -d 'A new duration to set on the entry. Mutually exclusive with offset.'
complete -c tid -n '__fish_tid_is_at_path "entry update"' -r -s n -l 'note=STR' -f -d 'A new note to set on the entry.'
complete -c tid -n '__fish_tid_is_at_path "entry update"' -r -s o -l 'offset=DUR' -f -d 'An offset to modify the duration by (can be negative). Mutually exclusive with duration.'
complete -c tid -n '__fish_tid_is_at_path "entry update"' -f -a '(__fish_tid_entries)'

# timesheet delete
complete -c tid -n '__fish_tid_is_at_path "timesheet delete"' -f -a '(__fish_tid_timesheets)'

# timesheet list
complete -c tid -n '__fish_tid_is_at_path "timesheet list"' -r -s e -l 'end=DATE' -f -d 'The end date of the listing. (Default: today)'
complete -c tid -n '__fish_tid_is_at_path "timesheet list"' -r -s f -l 'format=STR' -f -d 'Output formatting string. Uses Go templates.'
complete -c tid -n '__fish_tid_is_at_path "timesheet list"' -r -s s -l 'start=DATE' -f -d 'The start date of the listing. (Default: last monday)'

# workspace delete
complete -c tid -n '__fish_tid_is_at_path "workspace delete"' -f -a '(__fish_tid_workspaces)'

# workspace switch
complete -c tid -n '__fish_tid_is_at_path "workspace switch"' -f -a '(__fish_tid_workspaces)'
