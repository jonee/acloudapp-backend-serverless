// Created by Jonee Ryan Ty
// Copyright ACloudApp

/**
 * Generic Utility functions
 */

package utilities

import (
	"log"
	"strconv"
	"strings"
	"time"

	crand "crypto/rand"
	"encoding/base64"

	"math/big"
	"math/rand"
	"net/smtp"
)

// thanks http://stackoverflow.com/questions/15334220/encode-decode-base64
func DecodeB64(sToDecode string) string {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(sToDecode)))
	l, _ := base64.StdEncoding.Decode(base64Text, []byte(sToDecode))
	// log.Printf("base64: %s\n", base64Text[:l])
	return string(base64Text[:l])
}

func EncodeB64(s string) string {
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(base64Text, []byte(s))
	return string(base64Text)
}

// via mailgun vendor
// randomString generates a string of given length, but random content.
// All content will be within the ASCII graphic character set.
// (Implementation from Even Shaw's contribution on
// http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go).
func GetRandomString(n int, prefix string) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return prefix + string(bytes)
}

func GetWeirdString() string {
	max := *big.NewInt(99999999999)
	cryptoRand, _ := crand.Int(crand.Reader, &max)
	myRand := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Int()) + strconv.FormatInt(rand.Int63(), 10) + strconv.FormatInt(cryptoRand.Int64(), 10)
	return myRand

	// GetWeirdString
}

func Implode(mapStrStr map[string]string, del string) string {
	ret := ""

	for _, v := range mapStrStr {
		ret = ret + del + v
	}
	if ret != "" {
		ret = ret[len(del):]
	}

	return ret
}

func InArrayStr(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1

	arrayLen := len(array)
	for i := 0; i < arrayLen; i++ {
		if val == array[i] {
			index = i
			exists = true
			i = arrayLen // break from loop
		}
	}

	return
}

func MaskString(s string, masker string) string {
	return strings.Repeat(masker, len(s))
}

func SanitizeFilename(s string) string {
	s = strings.Replace(s, " ", "", -1) // remove spaces
	// todo remove non number or letter
	return s
}

// thanks https://gist.github.com/jpillora/cb46d183eca0710d909a
/*
func SendEmailGmailSmtp(from string, pass string, to string, subject string, body string) error {
	return SendEmailSmtp(from, pass, to, subject, body, "smtp.gmail.com", 587)
}
*/

func SendEmailSmtp(from string, password string, to string, subject string, body string, smtpServer string, smtpPort int) error {
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(smtpServer+":"+strconv.Itoa(smtpPort), smtp.PlainAuth("", from, password, smtpServer), from, []string{to}, []byte(msg))

	return err
}

func SubStr(str string, start int, length int) string {
	if start < 0 {
		start = 0
	}
	if length == 0 {
		length = len(str)
	}
	strArr := strings.Split(str, "")
	i := 0
	str = ""
	for i = start; i < start+length; i++ {
		if i >= len(strArr) {
			break
		}
		str += strArr[i]
	}
	return str
}

// thanks http://stackoverflow.com/questions/36530251/golang-time-since-with-months-and-years
func TimeDiff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// thanks https://github.com/PavelVershinin/GoWeb/blob/master/web/strings.go
func UcFirst(str string) string {
	first := SubStr(str, 0, 1)
	last := SubStr(str, 1, 0)
	return strings.ToUpper(first) + strings.ToLower(last)
}

func UcWords(str string) string {
	strArr := strings.Split(str, " ")
	for i := 0; i < len(strArr); i++ {
		strArr[i] = UcFirst(strArr[i])
	}
	return strings.Join(strArr, " ")
}

func GetExtensionFromContentType(mime string) string {
	ext := ""
	mapExt := map[string]string{"image/png": "png", "image/jpeg": "jpeg"}

	ext, _ = mapExt[mime]

	if ext == "" {
		log.Println("unknown extension: ", mime)
	}

	return ext
}
