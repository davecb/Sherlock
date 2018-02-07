# How to develop Sherlock filters

Get yesterdays logs, or logs from a day when you didn't have the current problem, into a file. 

- name it &lt;something>.log 
- take another copy and call it &lt;something>.template

The template file will be used by Sherlock to create a filter that removes
matching lines from today's logs, leaving only lines that haven't been seen before

- vi &lt;something>.template
- remove the dates from the beginning of the line. 
  
If you're in an editor with commands like vi, type 

	:s/^...// 

and add dots until you're just past the end of the dates. That will remove the beginning of every line.

Repeat until the date and any uninteresting or repeated junk is gone.
  
- now sort and unique the lines in the file. In vi, 

	:sort u

If you're in a hurry, skip right to "Filtering out yesterday's bugs", and try this template file out against testerday's collection of logs.

You will now have lots of lines that are similar, like
	
	starting pass 2
	starting pass 3
	starting pass 4
  
- keep "starting pass" from the first line and remove everything else. We're looking for expressions that will match all the starting page lines, and take them out.

This will be your first rule: lines containing "starting pass" are uninteresting, so they'll be filtered out

- remove all the other  "starting pass" lines from the .template file

You will be left with a line that says "starting pass" and a bunch of other lines, such as

	ending pass 9
	ending pass 10

- Just as before, keep "ending pass" as a rule and delete the repeat lines from the file.

Keep repeating the last few steps until you are at the bottom of the file.
Yes, some of the rules will say to remove things like "ERROR, pregnant moose", but if pregnant moose haven't caused problems before today, then they won't today.

You're allowed to come back to this file afterwards and investigate the errors and warnings, but that's not what we're doing today. We're looking for new errors, *not* old ones.

You should have a file full of rules, each describing something to remove from the logs. Congratulations, it's time to try them out.

# Filtering out yesterday's bugs
Run

	sherlock --compile <something>.template
	
and then

	sherlock <something> 
		
The First command creates a filter, the second takes the &lt;something>.filter file and run it against a &lt;something>.log file in real time.

If it suceeded, it will return nothing, because it's removed eveything "uninteresting". If it returns log lines, you will need to add lies to the template file that will match them., and run sherlock --compile &lt;something>.template again until it returns nothing.

# Filtering the good stuff out of today's file
- Collect all of today's files and call them &lt;something>.log.
run 

	sherlock &lt;something>


Whatever it prints out is peculiar to today. If the messages you see is "mose has exploded", then you know a moose exploded today and not yesterday, and may identify the problem.


