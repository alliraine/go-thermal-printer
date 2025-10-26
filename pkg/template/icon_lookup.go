package template

import (
	"fmt"
	"strings"
	"unicode"
)

var iconLookup map[string]string

func init() {
	iconLookup = make(map[string]string, len(iconData))
	for name := range iconWidgets {
		if _, ok := iconData[name]; !ok {
			panic(fmt.Sprintf("missing raster data for icon %s", name))
		}
		key := sanitizeIconKey(name)
		if existing, ok := iconLookup[key]; ok {
			if existing != name {
				panic(fmt.Sprintf("duplicate icon key %q for %q and %q", key, existing, name))
			}
			continue
		}
		iconLookup[key] = name
	}
}

func sanitizeIconKey(name string) string {
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

func lookupIconData(name string) ([]byte, string, error) {
	if entry, ok := iconData[name]; ok {
		return entry, name, nil
	}
	key := sanitizeIconKey(name)
	if canonical, ok := iconLookup[key]; ok {
		return iconData[canonical], canonical, nil
	}
	return nil, "", fmt.Errorf("icon %q not found", name)
}
