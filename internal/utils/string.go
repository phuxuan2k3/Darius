package utils

func Shorten(s string) string {
	if len(s) <= 60 {
		return s
	}
	return s[:30] + "..." + s[len(s)-30:]
}
