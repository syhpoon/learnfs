// Based on https://github.com/chrismcguire/gobberish/blob/master/gobberish.go
package lib

import (
	"math/rand"
	"unicode/utf8"
)

var chars = []rune(`abcdefghijklmnopqrstuvwxyz` +
	`ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789` +
	`АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ` +
	`абвгдеёжзийклмнопрстуфхцчшщъыьэюя` +
	`-_` +
	`的一是了我不人在他有这个上们来到时大地为` +
	`子中你说生国年着就那和要她出也得里后自以` +
	`会家可下而过天去能对小多然于心学么之都好` +
	`看起发当没成只如事把还用第样道想作种开美` +
	`总从无情己面最女但现前些所同日手又行意动` +
	`方期它头经长儿回位分爱老因很给名法间斯知` +
	`世什两次使身者被高已`)

// Generate a random utf-8 string with length of <= `length`.
func RandomBytes(length int, rnd *rand.Rand) string {
	var s []rune
	totalLength := 0

	for {
		r := chars[rnd.Intn(len(chars))]
		l := utf8.RuneLen(r)

		if totalLength+l > length {
			break
		}

		s = append(s, r)
		totalLength += l
	}

	return string(s)
}
