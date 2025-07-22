package common

import "regexp"

var (
	customEmojiRegex  = regexp.MustCompile(`^<a?:\w+:\d+>$`)
	// This regex matches most common Unicode emoji characters.
	unicodeEmojiRegex = regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{2600}-\x{26FF}]`)
)

func IsValidEmoji(emoji string) bool {
	return customEmojiRegex.MatchString(emoji) || unicodeEmojiRegex.MatchString(emoji)
}