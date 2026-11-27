// Package redact hides credentials from what the CLI prints for debugging,
// such as the API responses and commands echoed by --verbose.
package redact

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"
)

// Placeholder replaces every redacted value
const Placeholder = "[REDACTED]"

// minTokenLength is the shortest run of characters taken as a token
const minTokenLength = 32

var (
	// secretKeyWords mark a key, flag or variable name as holding a
	// credential when the name contains any of them
	secretKeyWords = []string{
		"password", "passwd", "passphrase", "secret", "token", "credential",
		"authorization", "apikey", "api_key", "privatekey", "private_key",
		"accesskey", "access_key",
	}

	// secretKeyParts mark a name as holding a credential when they are
	// one of its parts, e.g. DB_PASS, SENTRY_DSN or docker.login
	secretKeyParts = map[string]bool{
		"key": true, "pass": true, "auth": true, "dsn": true, "login": true,
	}

	// secretFlags take a credential as their value
	secretFlags = map[string]bool{
		"token": true, "password": true, "pass": true, "secret": true,
		"client-secret": true, "api-key": true, "apikey": true,
		"api-token": true, "access-token": true, "auth-token": true,
	}

	// loginFlags take a credential as their value in a login command
	loginFlags = map[string]bool{"u": true, "username": true, "p": true}

	bearerRe      = regexp.MustCompile(`(?i)\b(bearer\s+)[A-Za-z0-9._~+/=\-]+`)
	urlUserinfoRe = regexp.MustCompile(`([A-Za-z][A-Za-z0-9+.\-]*://[^\s:/@]+:)[^\s/@]+@`)
	keySepRe      = regexp.MustCompile(`([A-Za-z0-9_.\-]+)("?[ \t]*[=:][ \t]*"?)`)
	valueRe       = regexp.MustCompile(`^[^\s"'&,;]+`)
	tokenRunRe    = regexp.MustCompile(`[A-Za-z0-9_+/\-]{32,}={0,2}`)
)

// IsSensitiveKey tells whether a key, flag or variable name looks like it
// holds a credential
func IsSensitiveKey(name string) bool {
	lower := strings.ToLower(name)

	for _, word := range secretKeyWords {
		if strings.Contains(lower, word) {
			return true
		}
	}

	if strings.HasSuffix(lower, "key") {
		return true
	}

	for _, part := range strings.FieldsFunc(lower, isNotAlphanumeric) {
		if secretKeyParts[part] {
			return true
		}
	}

	return false
}

// JSON redacts a JSON document: values under credential-like keys, and any
// token-looking value. When raw is not valid JSON it is redacted as Text.
func JSON(raw []byte) string {
	var (
		doc any
		dec = json.NewDecoder(bytes.NewReader(raw))
	)

	dec.UseNumber()
	if err := dec.Decode(&doc); err != nil {
		return Text(string(raw))
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		// trailing data; not a single JSON document
		return Text(string(raw))
	}

	var (
		out bytes.Buffer
		enc = json.NewEncoder(&out)
	)

	enc.SetEscapeHTML(false)
	if err := enc.Encode(redactValue(doc)); err != nil {
		return Text(string(raw))
	}

	return strings.TrimSuffix(out.String(), "\n")
}

// Text redacts credentials found in free text: bearer tokens, passwords in
// URLs, values of credential-like key=value pairs and token-looking values
func Text(s string) string {
	s = bearerRe.ReplaceAllString(s, "${1}"+Placeholder)
	s = urlUserinfoRe.ReplaceAllString(s, "${1}"+Placeholder+"@")

	s = redactKeyValues(s)

	return tokenRunRe.ReplaceAllStringFunc(s, func(run string) string {
		if looksLikeToken(run) {
			return Placeholder
		}
		return run
	})
}

// redactKeyValues redacts the values of credential-like pairs such as
// password=..., token: ... or "secret":"..."
func redactKeyValues(s string) string {
	var (
		redacted strings.Builder
		last     int
	)

	for _, m := range keySepRe.FindAllStringSubmatchIndex(s, -1) {
		if m[0] < last || !IsSensitiveKey(s[m[2]:m[3]]) {
			// either inside a value already redacted or not a credential
			continue
		}

		start, end := valueBounds(s, m[1], s[m[4]:m[5]])
		if start == end {
			// blank value; keeping it shows the credential is missing
			continue
		}

		redacted.WriteString(s[last:start])
		redacted.WriteString(Placeholder)
		last = end
	}

	redacted.WriteString(s[last:])

	return redacted.String()
}

// valueBounds finds the value of a credential-like pair that starts at
// from, right after the given separator. A quoted value runs through its
// closing quote, skipping escaped characters, and keeps its quotes; any
// other value runs up to the next space or delimiter.
func valueBounds(s string, from int, sep string) (start, end int) {
	var quote byte

	start = from
	if strings.HasSuffix(sep, `"`) {
		// the separator took the opening quote, as in "secret":"..."
		quote = '"'
	} else if from < len(s) && (s[from] == '"' || s[from] == '\'') {
		quote = s[from]
		start++
	}

	if quote != 0 {
		if n := closingQuote(s[start:], quote); n >= 0 {
			return start, start + n
		}
	}

	return start, start + len(valueRe.FindString(s[start:]))
}

// closingQuote returns the index of the first unescaped quote in s, or -1
func closingQuote(s string, quote byte) int {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case quote:
			return i
		}
	}

	return -1
}

// Args redacts the arguments of a command line: values of credential flags
// (--token, --password, and -u/-p of a login command), credential-like
// KEY=VALUE pairs and token-looking values
func Args(args []string) []string {
	var (
		redacted    = make([]string, len(args))
		isLogin     bool
		redactValue bool
	)

	for _, arg := range args {
		if arg == "login" {
			isLogin = true
			break
		}
	}

	for i, arg := range args {
		if redactValue {
			redacted[i] = Placeholder
			redactValue = false
			continue
		}

		if strings.HasPrefix(arg, "-") {
			flag, value, hasValue := strings.Cut(arg, "=")

			if isSecretFlag(flag, isLogin) {
				if hasValue {
					redacted[i] = flag + "=" + Placeholder
				} else {
					redacted[i] = arg
					redactValue = true
				}
				continue
			}

			if hasValue {
				redacted[i] = flag + "=" + redactArg(value)
				continue
			}
		}

		redacted[i] = redactArg(arg)
	}

	return redacted
}

// redactArg redacts a single argument: the value of a credential-like
// KEY=VALUE pair, or any token-looking value
func redactArg(arg string) string {
	if key, _, isPair := strings.Cut(arg, "="); isPair && !strings.HasPrefix(arg, "-") && IsSensitiveKey(key) {
		return key + "=" + Placeholder
	}

	return Text(arg)
}

// isSecretFlag tells whether the given flag (with its dashes) takes a credential as value
func isSecretFlag(flag string, isLogin bool) bool {
	name := strings.ToLower(strings.TrimLeft(flag, "-"))

	return secretFlags[name] || (isLogin && loginFlags[name])
}

func redactValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for key, item := range v {
			if IsSensitiveKey(key) && !isBlank(item) {
				v[key] = Placeholder
			} else {
				v[key] = redactValue(item)
			}
		}
	case []any:
		for i := range v {
			v[i] = redactValue(v[i])
		}
	case string:
		return Text(v)
	}

	return value
}

// isBlank tells whether a value holds nothing worth hiding; keeping those
// as they are shows when a credential is missing
func isBlank(value any) bool {
	switch v := value.(type) {
	case nil, bool:
		return true
	case string:
		return v == ""
	}

	return false
}

// looksLikeToken tells whether a run of characters looks like a generated
// credential: at least minTokenLength characters with letters and digits,
// either in mixed case or with no separators at all - which leaves out
// file paths and names-like-this-1234
func looksLikeToken(run string) bool {
	core := strings.TrimRight(run, "=")

	if len(core) < minTokenLength || isPath(core) {
		return false
	}

	var upper, lower, digit, separator bool
	for _, r := range core {
		switch {
		case r >= 'a' && r <= 'z':
			lower = true
		case r >= 'A' && r <= 'Z':
			upper = true
		case r >= '0' && r <= '9':
			digit = true
		default:
			separator = true
		}
	}

	if !digit || (!upper && !lower) {
		return false
	}

	return (upper && lower) || !separator
}

// isPath tells whether a run looks like an absolute file path rather than
// a base64 token: it starts with a slash and has no long segment
func isPath(run string) bool {
	if !strings.HasPrefix(run, "/") {
		return false
	}

	for _, segment := range strings.Split(run, "/") {
		if len(segment) >= minTokenLength {
			return false
		}
	}

	return true
}

func isNotAlphanumeric(r rune) bool {
	return (r < 'a' || r > 'z') && (r < '0' || r > '9')
}
