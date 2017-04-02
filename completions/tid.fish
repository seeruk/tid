# tid.fish - tid completions for fish shell.
# 
# To install the completions:
# $ mkdir -p ~/.config/fish/completions
# $ cp tid.fish ~/.config/fish/completions
# 
# Aims to support commands, options, and entry hashes

function __fish_tid_no_command --description 'Test if tid is yet to be given a subcommand'
    for i in (commandline -opc)
        if contains -- $i 'entry' 'e' 'report' 'rep' 'resume' 'res' 'start' 'status' 'st' 'stop' 'timesheet' 't' 'workspace' 'w'
            return 1
        end
    end
    return 0
end

function __fish_tid_has_command --description 'Test if tid has been given the given command(s)'
    set -l targs (commandline -opc)

    if test (count $targs) -lt 2
        return 1
    end

    if contains -- $targs[2] $argv
        return 0
    end

    return 1
end

function __fish_tid_has_command_only --description 'Test if tid has been given only the given command(s)'
    set -l targs (commandline -opc)

    if test (count $targs) -ne 2
        return 1
    end

    if contains -- $targs[2] $argv
        return 0
    end

    return 1
end

function __fish_tid_has_subcommand --description 'Test if tid has been given the given subcommands(s)'
    set -l targs (commandline -opc)

    if test (count $targs) -lt 3
        return 1
    end

    if contains -- $targs[3] $argv
        return 0
    end

    return 1
end

function __fish_tid_has_command_and_subcommand --description 'Test if tid has been given the given command and subcommand'
    if test (count $argv) -lt 2
        return 1
    end

    # Commands
    if not __fish_tid_has_command (string split "|" $argv[1])
        return 1
    end

    # Sub-commands
    if not __fish_tid_has_subcommand (string split "|" $argv[2])
        return 1
    end

    return 0
end

function __fish_tid_has_command_and_no_subcommand --description 'Test if tid has been given the given command and not the given subcommand'
    if test (count $argv) -lt 2
        return 1
    end

    # Commands
    if not __fish_tid_has_command (string split "|" $argv[1])
        return 1
    end

    # Sub-commands
    if __fish_tid_has_subcommand (string split "|" $argv[2])
        return 1
    end

    return 0
end

# List entries from the last 6 months
function __fish_tid_entries
    command tid entry list --start=(tiddate --months=-6) --end=(tiddate) \
        --format="{{.ShortHash}}"\t"{{.Note}}"
end

# commands
complete -c tid -n '__fish_tid_no_command' -a 'entry' -f -d 'Manage timesheet entries.'
complete -c tid -n '__fish_tid_no_command' -a 'report' -f -d 'Display a timesheet report.'
complete -c tid -n '__fish_tid_no_command' -a 'resume' -f -d 'Resume an existing timer.'
complete -c tid -n '__fish_tid_no_command' -a 'start' -f -d 'Start a new timer.'
complete -c tid -n '__fish_tid_no_command' -a 'status' -f -d 'View the current status.'
complete -c tid -n '__fish_tid_no_command' -a 'stop' -f -d 'Stop an existing timer.'
complete -c tid -n '__fish_tid_no_command' -a 'timesheet' -f -d 'Manage timesheets.'
complete -c tid -n '__fish_tid_no_command' -a 'workspace' -f -d 'Manage workspaces.'

# sub-commands
# entry, e
complete -c tid -n '__fish_tid_has_command_only entry e' -a 'create' -f -d 'Create a new timesheet entry.'
complete -c tid -n '__fish_tid_has_command_only entry e' -a 'delete' -f -d 'Delete a timesheet entry.'
complete -c tid -n '__fish_tid_has_command_only entry e' -a 'list' -f -d 'List timesheet entries.'
complete -c tid -n '__fish_tid_has_command_only entry e' -a 'update' -f -d 'Update a timesheet entry.'

# timesheet, t
complete -c tid -n '__fish_tid_has_command_only timesheet t' -a 'delete' -f -d 'Delete a timesheet, and it\'s entries.'
complete -c tid -n '__fish_tid_has_command_only timesheet t' -a 'list' -f -d 'List timesheets.'

# workspace, w
complete -c tid -n '__fish_tid_has_command_only workspace w' -a 'create' -f -d 'Create a new workspace.'
complete -c tid -n '__fish_tid_has_command_only workspace w' -a 'delete' -f -d 'Delete a workspace.'
complete -c tid -n '__fish_tid_has_command_only workspace w' -a 'list' -f -d 'List available workspaces.'
complete -c tid -n '__fish_tid_has_command_only workspace w' -a 'switch' -f -d 'Switch to another workspace.'


# sub-sub-commands
# entry create
complete -c tid -n '__fish_tid_has_command_and_subcommand "entry|e" "create|c"' -r -s d -l 'date=DATE' -f -d 'When did you start working?'

# entry delete
complete -c tid -n '__fish_tid_has_command_and_subcommand "entry|e" "delete|d"' -a '(__fish_tid_entries)'

# entry update
complete -c tid -n '__fish_tid_has_command_and_subcommand "entry|e" "update|u"' -r -s d -l 'duration=DUR' -f -d 'A new duration to set on the entry. Mutually exclusive with offset.'
complete -c tid -n '__fish_tid_has_command_and_subcommand "entry|e" "update|u"' -r -s n -l 'note=STR' -f -d 'A new note to set on the entry.'
complete -c tid -n '__fish_tid_has_command_and_subcommand "entry|e" "update|u"' -r -s o -l 'offset=DUR' -f -d 'An offset to modify the duration by (can be negative). Mutually exclusive with duration.'
complete -c tid -n '__fish_tid_has_command_and_subcommand "entry|e" "update|u"' -a '(__fish_tid_entries)'

# report
complete -c tid -n '__fish_tid_has_command_only report' -r -s e -l 'date=DATE' -f -d 'The exact date of a timesheet to show a report for.'
complete -c tid -n '__fish_tid_has_command_only report' -r -s e -l 'end=DATE' -f -d 'The end date of the report.'
complete -c tid -n '__fish_tid_has_command_only report' -r -s f -l 'format=STR' -f -d 'Format string, uses Go templates.'
complete -c tid -n '__fish_tid_has_command_only report' -l 'no-summary' -f -d 'Hide the summary?'
complete -c tid -n '__fish_tid_has_command_only report' -r -s s -l 'start=DATE' -f -d 'The start date of the report.'

# resume
complete -c tid -f -n '__fish_tid_has_command_only resume' -a '(__fish_tid_entries)'

# status
complete -c tid -n '__fish_tid_has_command_only status' -r -s f -l 'format=STR' -f -d 'Format string, uses Go templates.'
complete -c tid -f -n '__fish_tid_has_command_only status' -a '(__fish_tid_entries)'


# The above is just getting out of hand... We're looking for something that can take in a
# fully-named expected command "path", and return either true or false to say whether we're at that
# point or not. This will need to be an internal option built into eidolon/console:
#
# $ tid e u --duration=DUR --ecint-is-at-path="entry update"
# $ echo $status
# 0
#
# $ tid e u --duration=DUR --ecint-is-at-path="entry list"
# $ echo $status
# 1
#
# To do this, we'll need to process the --ecint-is-at-args as if it is real input, and get to the
# point where we have the path. Once we're at that point, we can confirm that the path matches.
#
# This is only possible by exploiting the knowledge that the application has about it's current
# state when it's running. In other words, getting reliable completions without using something like
# this may be next to impossible (or so pointlessly complex it's not worth trying).
