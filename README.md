# Translate For CLI Version

This project is [Android Translate's](https://github.com/hbk01/Translate) CLI version.

# Install

## Dependency

- git
- golang
- make

## Compile

```shell
git clone https://github.com/hbk01/tr
cd tr
make clean install
```

# Usage

**Warning: It may not update, use `tr -h` to get usage.**

```
USAGE:
  tr [FLAG] [word]

FLAG:
  -h, --help    show this help
  -c, --clean   only show Translation
  -f=<LANG>, --form=<LANG>  set form language
  -t=<LANG>,   --to=<LANG>  set to language

LANG:
  cn - 中文Chinese
  en - 英文English
  ja - 日文Japanese
  ko - 韩文Koren
  fr - 法文French
  ru - 俄文Russian
  de - 德文German

EXAMPLE:
  $ tr hello
  $ tr -f=en -t=fr world
```

