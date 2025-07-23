package lib

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

func CleanText(raw string) string {
	// 1. 줄바꿈 → 공백으로 변환
	cleaned := strings.ReplaceAll(raw, "\r", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", " ")

	// 2. 깨진 문자 제거 (� == \uFFFD)
	cleaned = strings.Map(func(r rune) rune {
		if r == '\uFFFD' || !utf8.ValidRune(r) {
			return -1
		}
		return r
	}, cleaned)

	// 3. 연속된 공백 정리 (2개 이상 → 1개)
	re := regexp.MustCompile(`\s+`)
	cleaned = re.ReplaceAllString(cleaned, " ")

	// 4. 양쪽 공백 제거
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}
