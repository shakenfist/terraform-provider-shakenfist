#!/bin/bash

echo "--- Determine verison number ---"
PREVIOUS=`git tag | egrep "^v" | sort -n | tail -1 | sed 's/^v//'`

echo
echo -n "What is the version number (previous was $PREVIOUS)? "
read VERSION

echo "Now we need to bump go.mod."
read DUMMY
vim go.mod
go mod download github.com/shakenfist/client-go

echo "--- Setup ---"
echo "Do you want to apply a git tag for this release (yes to tag)?"
read TAG
set -x

if [ "%$TAG%" == "%yes%" ]
then
  git tag -s "v$VERSION" -m "Release v$VERSION"
  git push origin "v$VERSION"
fi
