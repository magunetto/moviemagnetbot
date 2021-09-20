package uri

// Schemes for different protocols.
const (
	SchemeMagnet  = "magnet"
	SchemeED2k    = "ed2k"
	SchemeFTP     = "ftp"
	SchemeUnknown = ""
)

// IsMagnet returns true if given string is a magnet link.
func IsMagnet(s string) bool {
	return GetScheme(s) == SchemeMagnet
}

// IsMagnet returns true if given string is an ED2k link.
func IsED2k(s string) bool {
	return GetScheme(s) == SchemeED2k
}

// IsMagnet returns true if given string is a FTP link.
func IsFTP(s string) bool {
	return GetScheme(s) == SchemeFTP
}

// GetScheme returns the scheme of the given link if it has one, otherwise empty string will be returned.
func GetScheme(s string) string {
	return NewURI(s).Scheme
}
