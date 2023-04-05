package utils

import "regexp"

var urlRegex = regexp.MustCompile(`^((https|http|ftp|rtsp|mms)?:\/\/)[^\s]+`)

func IsUrl(url string) bool {
	return urlRegex.MatchString(url)
}
