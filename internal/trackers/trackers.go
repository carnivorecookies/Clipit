package trackers

import (
	_ "embed"
	"encoding/json"
	"net/url"
)

//go:embed data.minify.json
var jsonRules []byte

var rules = make(map[string]*compiledRule)

// Initializes `trackers`.
// Can be called multiple times to update trackers list.
func Init() {
	var rawJson map[string]map[string]Rule
	_ = json.Unmarshal(jsonRules, &rawJson)

	for k, r := range rawJson["providers"] {
		rules[k] = r.Compile()
	}
}

// [Init] must be called before [RemoveFrom].
// uri.Scheme must be http or https.
func RemoveFrom(link *url.URL) *url.URL {
	if link.Scheme != "http" && link.Scheme != "https" {
		return link
	}

	removed := *link

	for _, r := range rules {
		if r.UrlPattern.MatchString(removed.String()) {
			removed = *r.RemoveTrackersFrom(&removed)
		}
	}

	return &removed
}
