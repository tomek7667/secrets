package secrets_test

import (
	"testing"

	"github.com/tomek7667/secrets/internal/secrets"
)

func TestPatternMatchesAllWildCards(t *testing.T) {
	type scenario struct {
		Pattern        string
		KeysMatching   []string
		KeysUnmatching []string
	}
	scenarios := map[string]scenario{
		"surrounded": {
			Pattern: "*abc*",
			KeysMatching: []string{
				"abc",
				"xxxabc",
				"abcyyy",
				"xxxxabcyyyy",
			},
			KeysUnmatching: []string{
				"abzc",
				"xxabzxx",
				"aabbcc",
				"",
			},
		},
		"internal": {
			Pattern: "ab*cd",
			KeysMatching: []string{
				"abZZZcd",
				"abcd",
			},
			KeysUnmatching: []string{
				"xxabcd",
				"xxabxxcd",
				"xxabxxcdyy",
				"abxxcdyy",
			},
		},
		"path-like matching": {
			Pattern: "*/anything/inside/*.txt",
			KeysMatching: []string{
				"essa///anything/inside/abc.txt",
				"essa/d/a/anything/inside/abc.txt",
				"/anything/inside/def.txt",
				"/anything/inside/CASE.txt",
				"/anything/inside/abc.txt",
				"/anything/inside/abc.txt.txt",
				"/anything/inside/abc.txt/abc.txt",
			},
			KeysUnmatching: []string{
				"essa///not/inside/abc.txt",
				"essa/d/a/anything/inside/abc.zip",
				"/anything/inside/abc.txt/",
			},
		},
	}
	for name, scenario := range scenarios {
		t.Run(name, func(tt *testing.T) {
			for _, matchingKey := range scenario.KeysMatching {
				matches := secrets.PatternMatches(matchingKey, scenario.Pattern)
				if !matches {
					tt.Errorf("'%s' should match pattern '%s'", matchingKey, scenario.Pattern)
				}
			}
			for _, unmatchingKey := range scenario.KeysUnmatching {
				matches := secrets.PatternMatches(unmatchingKey, scenario.Pattern)
				if matches {
					tt.Errorf("'%s' should not match pattern '%s'", unmatchingKey, scenario.Pattern)
				}
			}
		})
	}
}
