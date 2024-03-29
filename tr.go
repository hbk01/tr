package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO 202 Error, Means Sign error.

var ErrorCodeMap = map[string]string{
	"101":   "缺少必填的参数",
	"102":   "不支持的语言类型",
	"103":   "翻译文本过长",
	"104":   "不支持的API类型",
	"105":   "不支持的签名类型",
	"106":   "不支持的响应类型",
	"107":   "不支持的传输加密类型",
	"108":   "应用ID无效，注册账号，登录后台创建应用和实例并完成绑定，可获得应用ID和应用密钥等信息",
	"109":   "batchLog格式不正确",
	"110":   "无相关服务的有效实例",
	"111":   "开发者账号无效",
	"112":   "请求服务无效",
	"113":   "q不能为空",
	"114":   "不支持的图片传输方式",
	"201":   "解密失败，可能为DES,BASE64,URLDecode的错误",
	"202":   "签名检验失败",
	"203":   "访问IP地址不在可访问IP列表",
	"205":   "请求的接口与应用的平台类型不一致，如有疑问请参考入门指南",
	"206":   "因为时间戳无效导致签名校验失败",
	"207":   "重放请求",
	"301":   "辞典查询失败",
	"302":   "翻译查询失败",
	"303":   "服务端的其它异常",
	"304":   "会话闲置太久超时",
	"401":   "账户已经欠费停",
	"402":   "offlinesdk不可用",
	"411":   "访问频率受限,请稍后访问",
	"412":   "长请求过于频繁，请稍后访问",
	"1001":  "无效的OCR类型",
	"1002":  "不支持的OCR image类型",
	"1003":  "不支持的OCR Language类型",
	"1004":  "识别图片过大",
	"1201":  "图片base64解密失败",
	"1301":  "OCR段落识别失败",
	"1411":  "访问频率受限",
	"1412":  "超过最大识别字节数",
	"2003":  "不支持的语言识别Language类型",
	"2004":  "合成字符过长",
	"2005":  "不支持的音频文件类型",
	"2006":  "不支持的发音类型",
	"2201":  "解密失败",
	"2301":  "服务的异常",
	"2411":  "访问频率受限,请稍后访问",
	"2412":  "超过最大请求字符数",
	"3001":  "不支持的语音格式",
	"3002":  "不支持的语音采样率",
	"3003":  "不支持的语音声道",
	"3004":  "不支持的语音上传类型",
	"3005":  "不支持的语言类型",
	"3006":  "不支持的识别类型",
	"3007":  "识别音频文件过大",
	"3008":  "识别音频时长过长",
	"3009":  "不支持的音频文件类型",
	"3010":  "不支持的发音类型",
	"3201":  "解密失败",
	"3301":  "语音识别失败",
	"3302":  "语音翻译失败",
	"3303":  "服务的异常",
	"3411":  "访问频率受限,请稍后访问",
	"3412":  "超过最大请求字符数",
	"4001":  "不支持的语音识别格式",
	"4002":  "不支持的语音识别采样率",
	"4003":  "不支持的语音识别声道",
	"4004":  "不支持的语音上传类型",
	"4005":  "不支持的语言类型",
	"4006":  "识别音频文件过大",
	"4007":  "识别音频时长过长",
	"4201":  "解密失败",
	"4301":  "语音识别失败",
	"4303":  "服务的异常",
	"4411":  "访问频率受限,请稍后访问",
	"4412":  "超过最大请求时长",
	"5001":  "无效的OCR类型",
	"5002":  "不支持的OCR image类型",
	"5003":  "不支持的语言类型",
	"5004":  "识别图片过大",
	"5005":  "不支持的图片类型",
	"5006":  "文件为空",
	"5201":  "解密错误，图片base64解密失败",
	"5301":  "OCR段落识别失败",
	"5411":  "访问频率受限",
	"5412":  "超过最大识别流量",
	"9001":  "不支持的语音格式",
	"9002":  "不支持的语音采样率",
	"9003":  "不支持的语音声道",
	"9004":  "不支持的语音上传类型",
	"9005":  "不支持的语音识别 Language类型",
	"9301":  "ASR识别失败",
	"9303":  "服务器内部错误",
	"9411":  "访问频率受限（超过最大调用次数）",
	"9412":  "超过最大处理语音长度",
	"10001": "无效的OCR类型",
	"10002": "不支持的OCR image类型",
	"10004": "识别图片过大",
	"10201": "图片base64解密失败",
	"10301": "OCR段落识别失败",
	"10411": "访问频率受限",
	"10412": "超过最大识别流量",
	"11001": "不支持的语音识别格式",
	"11002": "不支持的语音识别采样率",
	"11003": "不支持的语音识别声道",
	"11004": "不支持的语音上传类型",
	"11005": "不支持的语言类型",
	"11006": "识别音频文件过大",
	"11007": "识别音频时长过长，最大支持30s",
	"11201": "解密失败",
	"11301": "语音识别失败",
	"11303": "服务的异常",
	"11411": "访问频率受限,请稍后访问",
	"11412": "超过最大请求时长",
	"12001": "图片尺寸过大",
	"12002": "图片base64解密失败",
	"12003": "引擎服务器返回错误",
	"12004": "图片为空",
	"12005": "不支持的识别图片类型",
	"12006": "图片无匹配结果",
	"13001": "不支持的角度类型",
	"13002": "不支持的文件类型",
	"13003": "表格识别图片过大",
	"13004": "文件为空",
	"13301": "表格识别失败",
	"15001": "需要图片",
	"15002": "图片过大（1M）",
	"15003": "服务调用失败",
	"17001": "需要图片",
	"17002": "图片过大（1M）",
	"17003": "识别类型未找到",
	"17004": "不支持的识别类型",
	"17005": "服务调用失败",
}

const (
	// id and key
	ID  = "3cd47726b488d5ad"
	KEY = "WiAQEJoYTzUcF9We6HnwkFAy7NByYfMt"
	API = "https://openapi.youdao.com/api"

	// language code
	LANG_AUTO = "auto"
	LANG_CN   = "zh-CHS"
	LANG_EN   = "en"
	LANG_JA   = "ja"
	LANG_KO   = "ko"
	LANG_FR   = "fr"
	LANG_RU   = "ru"
	LANG_DE   = "de"
)

type Language struct {
	From string
	To   string
}

// Json Root Tag
type Response struct {
	ErrorCode   string   `json:"errorCode"`   // 错误码，一定存在，无错误时为0
	Translation []string `json:"translation"` // 翻译，查询正确时（错误码=0）一定存在
	Web         []Web    `json:"web"`         // 网络释义（不一定有）
	Basic       Basic    `json:"basic"`       // 基本词典（查词时才有）
}

type Web struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

type Basic struct {
	Phonetic    string   `json:"phonetic"`
	Us_phonetic string   `json:"us-phonetic"`
	Uk_phonetic string   `json:"uk-phonetic"`
	Explains    []string `json:"explains"`
}

// translate
func translate(lang Language, word string) {
	Url, err := url.Parse(API)
	if err != nil {
		fmt.Println("URL Params Parse Error")
		os.Exit(1)
	}
	salt := strconv.FormatInt(time.Now().UnixNano(), 10)
	curtime := strconv.FormatInt(time.Now().Unix(), 10)
	params := Url.Query()
	params.Set("q", word)
	params.Set("from", lang.From)
	params.Set("to", lang.To)
	params.Set("appKey", ID)
	params.Set("salt", salt)
	params.Set("curtime", curtime)
	params.Set("sign", Sign(word, salt, curtime))
	params.Set("signType", "v3")
	params.Set("ext", "mp3")
	params.Set("voice", "0")
	params.Set("strict", "false")
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Request Services Error")
		os.Exit(1)
	}
	defer resp.Body.Close()
	parseJson(json.NewDecoder(resp.Body))
}

// Parse json and print the result.
func parseJson(dec *json.Decoder) {
	var response Response

	for {
		if err := dec.Decode(&response); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
	}

	if *Raw {
		b, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			fmt.Println("Can not marshal the json.")
			printResult(response)
		} else {
			fmt.Printf("%v\n", string(b))
		}
	} else if response.ErrorCode != "0" {
		printHumanError(response.ErrorCode)
	} else {
		printResult(response)
	}
}

// print error code use human message
func printHumanError(code string) {
	for k, v := range ErrorCodeMap {
		if k == code {
			fmt.Println(k, "->", v)
		}
	}
	fmt.Println("ErrorCode:", code)
}

func printResult(resp Response) { // {{{

	if *Clean {
		printArray(resp.Translation)
		os.Exit(0)
	}

	if resp.Translation != nil {
		fmt.Print("\nTranslation: ")
		printArray(resp.Translation)
	}
	fmt.Println("")
	if resp.Basic.Phonetic != "" {
		fmt.Println("Phonetic:", resp.Basic.Phonetic)
	}

	if resp.Basic.Us_phonetic != "" {
		fmt.Println("US-Phonetic:", resp.Basic.Us_phonetic)
	}

	if resp.Basic.Uk_phonetic != "" {
		fmt.Println("UK-Phonetic:", resp.Basic.Uk_phonetic)
	}

	if len(resp.Basic.Explains) > 0 {
		fmt.Println("")
		fmt.Print("Explains: ")
		printArrayFull(resp.Basic.Explains, "", "\n\t")
	}

	if len(resp.Web) > 0 {
		fmt.Println("")
		fmt.Println("Web: ")
		for _, web := range resp.Web {
			fmt.Print("\n" + web.Key + ": ")
			printArrayFull(web.Value, "", "\n\t\t")
		}
	}
} // }}}

func printArray(v []string) { // {{{
	printArrayFull(v, ", ", "")
}

func printArrayFull(v []string, split string, head string) {
	var builder strings.Builder
	for _, i := range v {
		builder.WriteString(head)
		builder.WriteString(i)
		builder.WriteString(split)
	}
	str := builder.String()
	fmt.Println(str[:len(str)-len(split)])
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
	s := sha256.Sum256([]byte(signString))
	return hex.EncodeToString(s[:])
} // }}}

// cast:     "cn", "en", "ja", "ko", "fr", "ru", "de"
//
// to  : "zh-CHS", "en", "ja", "ko", "fr", "ru", "de"
//
// if the lang is not in the array, return auto.
func convertLanguage(lang string) string { // {{{
	if lang == "cn" {
		return LANG_CN
	} else if in(lang, []string{"en", "ja", "ko", "fr", "ru", "de"}) {
		return lang
	} else {
		return LANG_AUTO
	}
} // }}}

// if lang in ["cn", "en", "ja", "ko", "fr", "ru", "de"], return true
func IsLang(lang string) bool { // {{{
	langs := []string{"cn", "en", "ja", "ko", "fr", "ru", "de"}
	return in(lang, langs)
} // }}}

// if item in array, return true, else false.
func in(item string, array []string) bool { // {{{
	for _, i := range array {
		if i == item {
			return true
		}
	}
	return false
} // }}}
