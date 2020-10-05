// launch this application
package main

import ( // {{{
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
) // }}}

// TODO 202 Error, Means Sign error.

const ( // {{{
	// id and key
	ID  = "3cd47726b488d5ad"
	KEY = "WiAQEJoYTzUcF9We6HnwkFAy7NByYfMt"

	// flags
	FLAG_HELP      = "-h"
	FLAG_HELP_LONG = "--help"
	FLAG_FROM      = "-f"
	FLAG_FROM_LONG = "--from"
	FLAG_TO        = "-t"
	FLAG_TO_LONG   = "--to"

	// language code
	LANG_AUTO = "auto"
	LANG_CN   = "zh-CHS"
	LANG_EN   = "en"
	LANG_JA   = "ja"
	LANG_KO   = "ko"
	LANG_FR   = "fr"
	LANG_RU   = "ru"
	LANG_DE   = "de"
) // }}}

// include form and to language
type Language struct { // {{{
	From string
	To   string
} // }}}

// launch translate
func main() { // {{{
	lang, word := ParseArgs(os.Args)
	fmt.Println("Translate:", word)
	fmt.Println("From", lang.From, "To", lang.To)
	translate(lang, word)
} // }}}

// translate
func translate(lang Language, word string) { // {{{
	URL := "https://openapi.youdao.com/api"
	Url, err := url.Parse(URL)
	if err != nil {
		fmt.Println("URL Params Parse Error")
		os.Exit(1)
	}
	salt := strconv.FormatInt(time.Now().UnixNano(), 10)
	curtime := strconv.FormatInt(time.Now().Unix(), 10)
	params := Url.Query()
	params.Set("q", url.QueryEscape(word))
	params.Set("from", "auto")
	params.Set("to", "auto")
	params.Set("appKey", ID)
	params.Set("salt", salt)
	params.Set("curtime", curtime)
	params.Set("sign", Sign(url.QueryEscape(word), salt, curtime))
	params.Set("signType", "v3")
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Request Services Error")
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read Response Body Error")
		os.Exit(1)
	}

	fmt.Println(string(body))
} // }}}

// sign
func Sign(input string, salt string, curtime string) string { // {{{
	Input := []rune(input)
	Lenth := len(Input)
	if Lenth > 20 {
		Pre := string(Input[:10])
		Suf := string(Input[Lenth-10:])
		input = Pre + strconv.Itoa(Lenth) + Suf
	}
	signString := ID + input + salt + curtime + KEY
	s := sha256.New().Sum([]byte(signString))
	return bytes2string(s)
} // }}}

// bytes convert to string
func bytes2string(bytes []byte) string { // {{{
	var builder strings.Builder
	for _, b := range bytes {
		h := int(b) & 0xFF
		hexString := strconv.FormatInt(int64(h), 16)
		if len(hexString) == 1 {
			builder.WriteString("0")
			builder.WriteString(hexString)
		}
		builder.WriteString(hexString)
	}
	return builder.String()
} // }}}

// parse input args
func ParseArgs(args []string) (lang Language, word string) { // {{{
	language := Language{
		From: LANG_AUTO,
		To:   LANG_AUTO,
	}
	for i, arg := range args {
		// args[0] always is binary file.
		if i == 0 {
			continue
		}

		if arg == FLAG_HELP || arg == FLAG_HELP_LONG {
			Usage()
			os.Exit(0)
		}

		if strings.HasPrefix(arg, FLAG_FROM) || strings.HasPrefix(arg, FLAG_FROM_LONG) {
			lang := strings.Split(arg, "=")[1]
			if IsLang(lang) {
				language.From = lang
			}
		}

		if strings.HasPrefix(arg, FLAG_TO) || strings.HasPrefix(arg, FLAG_TO_LONG) {
			lang := strings.Split(arg, "=")[1]
			if IsLang(lang) {
				language.To = lang
			}
		}
	}
	// last item in args is translate word
	return language, args[len(args)-1]
} // }}}

// if lang in ["cn", "en", "ja", "ko", "fr", "ru", "de"], return true
func IsLang(lang string) bool { // {{{
	langs := []string{"cn", "en", "ja", "ko", "fr", "ru", "de"}
	for _, l := range langs {
		if l == lang {
			return true
		}
	}
	return false
} // }}}

// print usage to stdin
func Usage() { // {{{
	fmt.Println("tr [flag] [word]")
	fmt.Println("")
	fmt.Println("flag:")
	fmt.Println(" -h, --help    show this help")
	fmt.Println(" -f=<LANG>, --form=<LANG>  set form language")
	fmt.Println(" -t=<LANG>,   --to=<LANG>  set to language")
	fmt.Println("")
	fmt.Println("<LANG> can be:")
	fmt.Println(" 中文Chinese  - cn")
	fmt.Println(" 英文English  - en")
	fmt.Println(" 日文Japanese - ja")
	fmt.Println(" 韩文Koren    - ko")
	fmt.Println(" 法文French   - fr")
	fmt.Println(" 俄文Russian  - ru")
	fmt.Println(" 德文German   - de")
	fmt.Println("")
	fmt.Println("example:")
	fmt.Println(" $ tr hello # form and to default to auto")
	fmt.Println(" $ tr -t=fa world")
	fmt.Println("")
	fmt.Println("This project in https://github.com/hbk01/tr")
} // }}}
