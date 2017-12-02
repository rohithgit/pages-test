#!/bin/sh

#Find the Process ID for syncapp running instance

echo "stopping scheduler-pagetest"
ps -ef | grep scheduler-pagetest | grep -v grep | grep 4031 | awk '{print $2}' | xargs kill -9 

