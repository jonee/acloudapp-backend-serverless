
# thanks http://www.anyexample.com/linux_bsd/bash/recursively_delete_backup_files.xml

find ./ -name '*~' -exec rm '{}' \; -print -or -name ".*~" -exec rm {} \; -print

