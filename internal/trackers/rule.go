package trackers

import (
	"net/url"
	"regexp"
	"strings"
)

type Rule struct {
	// The pattern for detecting the URL.
	UrlPattern string `json:"urlPattern"`

	// Whether the URL is always blocked by ClearURL regardless of content
	Blocked bool `json:"completeProvider"`

	// Each query rule may or may not be a regexp.
	// If it is not a regexp, it should be compared as a string.
	QueryRules []string `json:"rules,omitempty"`

	// Raw rules are applied directly to the raw URL.
	RawURLRules []string `json:"rawRules,omitempty"`

	// If any of the exceptions match, no replacement needs to be done.
	// Checked before applying any rules.
	Exceptions []string `json:"exceptions,omitempty"`

	// Redirections; should be checked before everything except Exceptions.
	// The first capture group (get with [regexp.Regexp.FindStringSubmatchIndex]) is the new URL.
	Redirections []string `json:"redirections,omitempty"`
}

type compiledRule struct {
	// The pattern for detecting the URL.
	UrlPattern *regexp.Regexp

	// Whether the URL is always blocked by ClearURL regardless of content.
	blocked bool

	// Regexp query rules.
	queryRules []*regexp.Regexp

	// String query rules, must be checked case-insensitive.
	queryStringRules []string

	// Raw rules are applied directly to the raw URL.
	rawRules []*regexp.Regexp

	// If any of the exceptions match, no replacement needs to be done.
	// Checked before applying any rules.
	exceptions []*regexp.Regexp

	// Redirections; should be checked before everything except Exceptions.
	// The first capture group (get with [regexp.Regexp.FindStringSubmatchIndex]) is the new URL.
	redirections []*regexp.Regexp
}

// Note: This function uses [regexp.MustCompile], and may panic if invalid regexps were passed.
// If [rule.Blocked == true], only [compiledRule.urlPattern] is compiled and [compiledRule.Blocked] is set.
func (r *Rule) Compile() *compiledRule {
	// Turn the patterns case insensitive
	mustCompile := func(s string) *regexp.Regexp { return regexp.MustCompile("(?i)" + s) }

	c := new(compiledRule)
	c.UrlPattern = mustCompile(r.UrlPattern)

	c.blocked = r.Blocked
	if c.blocked {
		return c
	}

	for _, tracker := range r.QueryRules {
		if rgx, err := regexp.Compile("(?i)" + tracker); err == nil {
			c.queryRules = append(c.queryRules, rgx)
		} else {
			c.queryStringRules = append(c.queryStringRules, tracker)
		}
	}

	for _, rawRule := range r.RawURLRules {
		c.rawRules = append(c.rawRules, mustCompile(rawRule))
	}

	for _, exc := range r.Exceptions {
		c.exceptions = append(c.exceptions, mustCompile(exc))
	}

	for _, redir := range r.Redirections {
		c.redirections = append(c.redirections, mustCompile(redir))
	}

	return c
}

// Assumes uri complies with [rule.UrlPattern]/
// Removes trackers from the [url.URL] if any are present.
// If the returned [url.URL] is nil, then the [url.URL] is fully blocked by ClearURL.
func (r *compiledRule) RemoveTrackersFrom(link *url.URL) *url.URL {
	if r.blocked {
		return nil
	}

	rawURL := link.String()
	for _, exc := range r.exceptions {
		if exc.MatchString(rawURL) {
			return link
		}
	}

	removed := *link
	r.removeQueryTrackers(&removed)
	removed = *r.removeRawTrackers(&removed)

	return &removed
}

func (r *compiledRule) removeRawTrackers(link *url.URL) *url.URL {
	rawURL := (*link).String()
	for _, rawRule := range r.rawRules {
		rawURL = rawRule.ReplaceAllString(rawURL, "")
	}
	// parsing errors cannot occur
	applied, _ := url.Parse(rawURL)
	return applied
}

func (r *compiledRule) removeQueryTrackers(link *url.URL) {
	queries := link.Query()

QueryLoop:
	for query := range queries {
		for _, regexRule := range r.queryRules {
			if regexRule.MatchString(query) {
				queries.Del(query)
				continue QueryLoop
			}
		}

		for _, strRule := range r.queryStringRules {
			if strings.EqualFold(strRule, query) {
				queries.Del(query)
				continue QueryLoop
			}
		}
	}

	link.RawQuery = queries.Encode()
}

// func (r *compiledRule) redirect(link *url.URL) {}
