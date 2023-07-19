package re

const (
	SRE_FLAG_TEMPLATE   = 1 << iota // template mode (disable backtracking)
	SRE_FLAG_IGNORECASE             // i, case insensitive
	SRE_FLAG_LOCALE                 // honour system locale
	SRE_FLAG_MULTILINE              // m, treat target as multiline string
	SRE_FLAG_DOTALL                 // s, treat target as a single string
	SRE_FLAG_UNICODE                // use unicode "locale"
	SRE_FLAG_VERBOSE                // ignore whitespace and comments
	SRE_FLAG_DEBUG                  // debugging
	SRE_FLAG_ASCII                  // use ascii "locale"
)
