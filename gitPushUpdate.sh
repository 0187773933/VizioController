#!/bin/bash
function is_int() { return $(test "$@" -eq "$@" > /dev/null 2>&1); }
ssh-add -D
git init
git config --global --unset user.name
git config --global --unset user.email
git config user.name "0187773933"
git config user.email "collincerbus@student.olympic.edu"
ssh-add -k /Users/morpheous/.ssh/githubWinStitch

LastCommit=$(git log -1 --pretty="%B" | xargs)
# https://stackoverflow.com/a/3626205
if $(is_int "${LastCommit}");
    then
    NextCommitNumber=$((LastCommit+1))
else
   echo "Not an integer Resetting"
   NextCommitNumber=1
fi
git add .
git tag -l | xargs git tag -d
if [ -n "$1" ]; then
   git commit -m "$1"
   git tag v1.0.$1
else
   git commit -m "$NextCommitNumber"
   git tag v1.0.$NextCommitNumber
fi
git remote add origin git@github.com:0187773933/VizioController.git
git push origin --tags
git push origin master