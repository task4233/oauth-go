package common

import (
	"fmt"
	"net/url"
	"strings"
)

const blacklistInQuery = "&="

func AreTwoUnorderedSlicesSame[E comparable](s []E, t []E) bool {
	if len(s) != len(t) {
		return false
	}

	mpS := map[E]int{}
	for _, e := range s {
		mpS[e]++
	}
	for _, e := range t {
		mpS[e]--
	}
	for _, v := range mpS {
		if v != 0 {
			return false
		}
	}

	return true
}

func ConstructURLWithQueries(uri string, queryParameters map[string]string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("failed url.Parse: %w", err)
	}

	q := u.Query()
	for k, v := range queryParameters {
		if strings.ContainsAny(k, blacklistInQuery) || strings.ContainsAny(v, blacklistInQuery) {
			return "", fmt.Errorf("invalid query parameter: %s=%s", k, v)
		}
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
