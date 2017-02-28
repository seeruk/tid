# tid.fish - tid completions for fish shell.
# 
# To install the completions:
# $ mkdir -p ~/.config/fish/completions
# $ cp tid.fish ~/.config/fish/completions
# 
# Aims to support commands, options, and entry hashes

function __fish_tid_no_subcommand --description 'Test if tid is yet to be given a subcommand'
    for i in (commandline -opc)
        if contains -- $i add edit remove report resume start status stop
            return 1
        end
    end
    return 0
end

# List entries from the last 6 months
function __fish_tid_entries
    command tid report --start=(tiddate --months=-6) --end=(tiddate) \
        --no-summary --format="{{.Entry.Hash}}"\t"{{.Entry.Note}}"
end

# subcommands
# add
complete -c tid -f -n '__fish_tid_no_subcommand' -a add -d 'Add a timesheet entry.'

# edit
complete -c tid -f -n '__fish_tid_no_subcommand' -a edit -d 'Edit a timesheet entry.'
complete -c tid -n '__fish_seen_subcommand_from edit' -r -s d -l 'duration=DUR' -f -d 'A new duration to set on the entry. Mutually exclusive with offset.'
complete -c tid -n '__fish_seen_subcommand_from edit' -r -s n -l 'note=STR' -f -d 'A new note to set on the entry.'
complete -c tid -n '__fish_seen_subcommand_from edit' -r -s o -l 'offset=DUR' -f -d 'An offset to modify the duration by (can be negative). Mutually exclusive with duration.'
complete -c tid -f -n '__fish_seen_subcommand_from edit' -a '(__fish_tid_entries)'

# remove
complete -c tid -f -n '__fish_tid_no_subcommand' -a remove -d 'Remove a timesheet entry.'
complete -c tid -f -n '__fish_seen_subcommand_from remove' -a '(__fish_tid_entries)'

# report
complete -c tid -f -n '__fish_tid_no_subcommand' -a report -d 'Display a tabular timesheet report.'
complete -c tid -n '__fish_seen_subcommand_from report' -r -s e -l 'end=DATE' -f -d 'The end date of the report.'
complete -c tid -n '__fish_seen_subcommand_from report' -r -s f -l 'format=STR' -f -d 'Format string, uses table headers e.g. \'{{HASH}}\'.'
complete -c tid -n '__fish_seen_subcommand_from report' -l 'no-summary' -f -d 'Hide the summary?'
complete -c tid -n '__fish_seen_subcommand_from report' -r -s s -l 'start=DATE' -f -d 'The start date of the report.'

# resume
complete -c tid -f -n '__fish_tid_no_subcommand' -a resume -d 'Resume an existing timer.'
complete -c tid -f -n '__fish_seen_subcommand_from resume' -a '(__fish_tid_entries)'

# start
complete -c tid -f -n '__fish_tid_no_subcommand' -a start -d 'Start a new timer.'

# status
complete -c tid -f -n '__fish_tid_no_subcommand' -a status -d 'View the current status.'
complete -c tid -n '__fish_seen_subcommand_from status' -r -s f -l 'format=STR' -f -d 'Format string, uses table headers e.g. \'{{HASH}}\'.'
complete -c tid -f -n '__fish_seen_subcommand_from status' -a '(__fish_tid_entries)'

# stop
complete -c tid -f -n '__fish_tid_no_subcommand' -a stop -d 'Stop an existing timer.'
