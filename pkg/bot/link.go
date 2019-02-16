package bot

import "strings"

const (
	magnetPrefix = "magnet:?"
	eD2kPrefix   = "ed2k://"
	ftpPrefix    = "ftp://"
)

func isMagnetLink(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), magnetPrefix)
}

func isED2kLink(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), eD2kPrefix)
}
func isFTPLink(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), ftpPrefix)
}
