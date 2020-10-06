# Translate For CLI Version

This project is [Android Translate's](https://github.com/hbk01/Translate) CLI version.

# Install

## Dependency

- [ ] git
- [ ] golang
- [ ] make

## Compile

```shell
git clone https://github.com/hbk01/tr
cd tr
make clean install
```

# Usage

```
USAGE:
  tr [FLAG] [word]

FLAG:
  -h, --help    show this help
  -f=<LANG>, --form=<LANG>  set form language
  -t=<LANG>,   --to=<LANG>  set to language

<LANG> can be:
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

