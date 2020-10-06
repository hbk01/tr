#!/bin/env bash

echo "check environment"
if [ ! -f release.note ]; then
    echo "Please edit 'release.note' file for this release note."
    exit 1
fi

if [ ! $1 ]; then
    echo "Please input tag name."
    exit 1
fi

echo "clean and build"
rm -rf ./bin
go build -o ./bin/tr

if [ $? -ne 0 ]; then
    echo "Build Failed."
    exit 1
fi

echo "check gh auth"
gh auth status

if [ $? -ne 0 ]; then
    echo "gh not login, starting login..."
    gh auth login
fi

echo "release as draft..."
gh release create $1 --draft --title "Release $1" --notes-file ./release.note ./bin/tr

if [ $? -eq 0 ]; then
    echo "release sussess, please view https://github.com/hbk01/tr/releases/tag/$1 to comfirm and publish!"
    echo "clean..."
    rm -rf ./bin
    if [ -f release.note ]; then
        rm release.note
    fi
fi

