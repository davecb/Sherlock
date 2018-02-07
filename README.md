# Sherlock

In agile companies, one still needs a
"Sherlock Holmes"[1] style of log reader:  
```
Read the log from a failing program,
but only show us the new lines, the 
lines that weren't there until we 
had the problem. 
```

However, the original implementation, 
"antilog" was a batch program run by cron. 
It's more common now to want immediate 
notification. 

Sherlock is therefore a streams program, one that 
can tail multiple logs, remove all the uninteresting 
content, and notify operations
when a new, unexpected message shows up.

It can be used as a background process, a
retrospective log analyzer, a tool to explore logs
interactively or even an old-fashioned cron job,
all using the same basic algorithm. 

[1. "It is an old maxim of mine that when you have 
excluded the impossible, whatever remains, 
however improbable, must be the truth." 
Sherlock Holmes, _Adventure of the Beryl Coronet._

By this logic, 
anything seen before the error ocurred is an
"impossible" cause, and should be excluded.]
