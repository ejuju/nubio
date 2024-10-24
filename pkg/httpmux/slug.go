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
func Slugify(v string) (out string) {
	// Trim leading and trailing spaces, Normalize characters, transform to lowercase, and subsitute characters.
	v = strings.TrimSpace(v)
	v = norm.NFKC.String(v)
	v = strings.ToLower(v)
	for _, c := range v {
		clean, ok := subsitutions[c]
		if !ok {
			clean = string(c)
		}
		out += string(clean)
	}
	return out
}
