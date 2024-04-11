#! /bin/bash
git add .
#read input
git commit -m "$1"
git push --set-upstream origin v2dev