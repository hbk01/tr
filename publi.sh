#!/bin/env bash

echo "- check environment"
if [ ! -f release.note ]; then
    echo "Please edit 'release.note' file for this release note."
    exit 1
fi

if [ ! $1 ]; then
    echo "Please input tag name."
    exit 1
fi

echo "- clean"
rm -rf ./bin

echo "- build binary files"
os=linux
echo "  - build $os"
echo "    - build amd64"
GOOS=$os GOARCH=amd64 go build -o ./bin/tr_$1-$os-amd64
echo "    - build 386"
GOOS=$os GOARCH=386 go build -o ./bin/tr_$1-$os-386
echo "    - build arm"
GOOS=$os GOARCH=arm go build -o ./bin/tr_$1-$os-arm

os=freebsd
echo "  - build $os"
echo "    - build amd64"
GOOS=$os GOARCH=amd64 go build -o ./bin/tr_$1-$os-amd64
echo "    - build 386"
GOOS=$os GOARCH=386 go build -o ./bin/tr_$1-$os-386
echo "    - build arm"
GOOS=$os GOARCH=arm go build -o ./bin/tr_$1-$os-arm

os=darwin
echo "  - build $os"
echo "    - build amd64"
GOOS=$os GOARCH=amd64 go build -o ./bin/tr_$1-macOS-amd64

os=windows
echo "  - build $os"
echo "    - build x64"
GOOS=$os GOARCH=amd64 go build -o ./bin/tr_$1-$os-x64.exe
echo "    - build x86"
GOOS=$os GOARCH=386 go build -o ./bin/tr_$1-$os-x86.exe

os=android
echo "  - build $os"
echo "    - build arm64"
GOOS=$os GOARCH=arm64 go build -o ./bin/tr_$1-$os-arm64

echo "- build finished."

if [ $? -ne 0 ]; then
    echo "Build Failed."
    exit 1
fi

echo "- check gh auth"
gh auth status 2>&1 >> /dev/null

if [ $? -ne 0 ]; then
    echo "gh not login, starting login..."
    gh auth login
fi

echo "files will be upload to assets: "
files=$(ls ./bin)
for file in $files
do
    echo "  $file"
done

echo "- create tag $1"
git tag $1 -m "Release $1"
echo "- pushing tags"
git push --tags

echo "- release as draft..."
url=$(gh release create $1 --draft --title "Release $1" --notes-file ./release.note ./bin/tr*)

if [ $? -eq 0 ]; then
    echo "release sussess, please view $url to comfirm and publish!"
    echo "- clean..."
    rm -rf ./bin
    if [ -f release.note ]; then
        rm release.note
    fi
fi

