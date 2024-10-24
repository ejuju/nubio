package httpmux

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

var subsitutions = map[rune]string{
	// Space characters.
	' ':  "-",
	'\n': "-",
	'\r': "-",
	'\t': "-",
	// Symbols and punctuation.
	'"':  "",
	'\'': "",
	'&':  "-",
	'@':  "-",
	// Letters with accents.
	'à': "a",
	'é': "e",
	'è': "e",
	'ê': "e",
	'ô': "o",
}

// Transforms the input string into a "slug" that can be used as a URI.
// This is a "best effort" algorithm, meaning that if no substitution is defined for the character,
// it will be included "as is" in the output.
//
// Example: "Alizée Doe" -> "alizee-doe".
func Slugify(raw string) (out string) {
	// Trim leading and trailing spaces, Normalize characters, transform to lowercase, and subsitute characters.
	raw = strings.TrimSpace(raw)
	raw = norm.NFKC.String(raw)
	raw = strings.ToLower(raw)
	for _, c := range raw {
		clean, ok := subsitutions[c]
		if !ok {
			clean = string(c)
		}
		out += string(clean)
	}
	return out
}
