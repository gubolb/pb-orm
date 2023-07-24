package orm

import "github.com/fatih/structtag"

func extractOrmNameFromTag(rawTag string) string {
	tags, err := structtag.Parse(rawTag)
	if err != nil || tags == nil {
		return ""
	}

	ormTag, err := tags.Get("orm")
	if err != nil || ormTag == nil {
		return ""
	}

	return ormTag.Name
}

func isOmitable(rawTag string) bool {
	tags, err := structtag.Parse(rawTag)
	if err != nil || tags == nil {
		return true
	}

	ormTag, err := tags.Get("orm")
	if err != nil || ormTag == nil {
		return true
	}

	return ormTag.HasOption("omitempty")
}
