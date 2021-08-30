// launch this application
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// TODO 202 Error, Means Sign error.

var (
	Clean = flag.Bool("clean", false, "Only show 'translate' item")
	// 设置翻译语言，值必须有一个 - 符号，-符号前面的为翻译前的语言，-符号后面的为翻译后的语言
	Lang = flag.String("lang", "auto/auto", "language to translate, format of '[from]/[to]', ")
	From = flag.String("from", "auto", "set 'from' language to translate, alias of '-lang=from/'")
	To   = flag.String("to", "auto", "set 'to' language to translate, alias of '-lang=/to'")
	Arch = flag.Bool("arch", false, "print the device arch type.")
	Raw  = flag.Bool("raw", false, "print the responses(json data) from server.")

	Help = flag.Bool("help", false, "print the help")
)

// launch translate
func main() {
	flag.Parse()

	if *Help {
		flag.Usage()

		fmt.Println("Languages: ")
		fmt.Println("  cn: 中文(Chinese)")
		fmt.Println("  en: 英文(English)")
		fmt.Println("  ja: 日文(Japanese)")
		fmt.Println("  ko: 韩文(Koren)")
		fmt.Println("  fr: 法文(French)")
		fmt.Println("  ru: 俄文(Russian)")
		fmt.Println("  de: 德文(German)")
		os.Exit(1)
	}

	if *Arch {
		fmt.Println(runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// 将最后一个参数当作翻译文本
	word := flag.Args()[len(flag.Args())-1]

	// 处理 language
	var lang Language
	if *Lang != "auto/auto" {
		tempLang := strings.Split(*Lang, "/")
		f, t := tempLang[0], tempLang[1]
		if f != "" {
			lang.From = f
		}

		if t != "" {
			lang.To = t
		}
	} else {
		lang.From = *From
		lang.To = *To
	}

	if !*Clean {
		// when Clean is false, print this!
		fmt.Printf("Language: %s\tText: %s\n", lang.From+"/"+lang.To, word)
	}
	translate(lang, word)
}
