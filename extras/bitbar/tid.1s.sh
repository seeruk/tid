#!/bin/bash

# To install:
# 	1. Install BitBar (https://github.com/matryer/bitbar)
#	2. From the BitBar menu, click Preferences > Open Plugin Folder
#	3. Copy this script into the BitBar plugin folder
#	4. Check "TID=..." below to ensure your tid binary can be found

# <bitbar.title>tid</bitbar.title>
# <bitbar.version>v0.1</bitbar.version>
# <bitbar.author>Callum Jones</bitbar.author>
# <bitbar.author.github>cj123</bitbar.author.github>
# <bitbar.desc>Show tid status information and stop and resume tid timers.</bitbar.desc>
# <bitbar.dependencies>tid</bitbar.dependencies>
# <bitbar.image>http://i.imgur.com/PPtl120.png</bitbar.image>

# you may need to update TID to point to where your tid binary
# is located if `which tid` provides no results
TID=$(which tid || echo "$HOME/Projects/go/bin/tid")

# keep 'note' at end of this string so we can be lazy and slice to the end of the array :P
# without having to deal with spaces
status=`$TID status -f="{{.ShortHash}} {{.Duration}} {{.IsRunning}} {{.Note}}"`
arr=($status)

hash=${arr[0]}
duration=${arr[1]}
isRunning=${arr[2]}
note=${arr[@]:3}

if [ "$isRunning" == "true" ]; then
	echo "$duration | color=green"
	echo "---"
	echo "$note ($hash)"
	echo "Stop | bash=$TID param1=stop param2=$hash terminal=false"
else
	echo "$duration | color=red"
	echo "---"
	echo "$note ($hash)"
	echo "Resume | bash=$TID param1=resume param2=$hash terminal=false"
fi
