package automate

import (
	"fmt"
	"strconv"
	"strings"
)

// qualifiers in ascending order of precedence, as used by Spring Initializr
var versionQualifiers = []string{"M", "RC", "BUILD-SNAPSHOT", "RELEASE"}

type version struct {
	numbers   [3]int
	qualifier int
	iteration int
}

// parseVersion parses versions like 4.1.2, 4.1.2.RELEASE, 4.2.0.RC1 or 4.2.0-M1
func parseVersion(raw string) (v version, err error) {
	raw = strings.TrimSpace(raw)
	parts := strings.SplitN(raw, ".", 4)
	if len(parts) == 3 {
		if patch, qualifier, found := strings.Cut(parts[2], "-"); found {
			parts = []string{parts[0], parts[1], patch, qualifier}
		}
	}
	if len(parts) < 3 {
		return v, fmt.Errorf("invalid version %q", raw)
	}

	for i := 0; i < 3; i++ {
		if v.numbers[i], err = strconv.Atoi(parts[i]); err != nil {
			return v, fmt.Errorf("invalid version %q", raw)
		}
	}

	v.qualifier = len(versionQualifiers) - 1
	if len(parts) == 4 {
		qualifier := strings.ToUpper(parts[3])
		for i := len(versionQualifiers) - 1; i >= 0; i-- {
			if strings.HasPrefix(qualifier, versionQualifiers[i]) {
				v.qualifier = i
				v.iteration, _ = strconv.Atoi(strings.TrimPrefix(qualifier, versionQualifiers[i]))
				break
			}
		}
		if qualifier == "SNAPSHOT" {
			v.qualifier = 2
		}
	}

	return
}

func (v version) compare(other version) int {
	for i := 0; i < 3; i++ {
		if v.numbers[i] != other.numbers[i] {
			return v.numbers[i] - other.numbers[i]
		}
	}
	if v.qualifier != other.qualifier {
		return v.qualifier - other.qualifier
	}
	return v.iteration - other.iteration
}

// versionInRange tells whether the version satisfies a Spring Initializr
// version range, such as "4.0.0.RELEASE" (minimum) or "[4.0.0.RELEASE,4.1.0.M1)"
func versionInRange(raw, versionRange string) (bool, error) {
	versionRange = strings.TrimSpace(versionRange)
	if versionRange == "" {
		return true, nil
	}

	v, err := parseVersion(raw)
	if err != nil {
		return false, err
	}

	if !strings.HasPrefix(versionRange, "[") && !strings.HasPrefix(versionRange, "(") {
		lower, err := parseVersion(versionRange)
		if err != nil {
			return false, err
		}
		return v.compare(lower) >= 0, nil
	}

	bounds := strings.Split(versionRange[1:len(versionRange)-1], ",")
	if len(bounds) != 2 {
		return false, fmt.Errorf("invalid version range %q", versionRange)
	}

	lower, err := parseVersion(bounds[0])
	if err != nil {
		return false, err
	}
	upper, err := parseVersion(bounds[1])
	if err != nil {
		return false, err
	}

	lowerCmp, upperCmp := v.compare(lower), v.compare(upper)

	lowerOk := lowerCmp > 0 || (lowerCmp == 0 && versionRange[0] == '[')
	upperOk := upperCmp < 0 || (upperCmp == 0 && versionRange[len(versionRange)-1] == ']')

	return lowerOk && upperOk, nil
}
