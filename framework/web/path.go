package web

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

type (
	// Path is a matchable routing path
	Path struct {
		path          string // `/path/to/:something/$id<[0-9]+>/*foo
		parts         []part // {'path/to', ':something', '$id<[0-9]+>', '*foo'}
		params        []string
		trailingSlash bool
		normalize     map[string]struct{}
	}

	// Match is the result of matching a path
	Match struct {
		Values map[string]string
	}

	part interface {
		match(path string) (matched bool, key, value string, length int)
		read(segment string) (leftover, paramname string, err error)
		render(values map[string]string, normalize map[string]struct{}) (string, []string, error)
	}

	partFixed struct {
		part   string
		length int
	}

	partParam struct {
		name, suffix string
	}

	partRegex struct {
		name  string
		regex *regexp.Regexp
	}

	partWildcard struct {
		name string
	}
)

func (p *partFixed) read(path string) (string, string, error) {
	pos := strings.IndexAny(path, ":*$")

	if pos < 0 {
		p.length = len(path)
		if strings.HasSuffix(path, "/") {
			p.length--
		}
		p.part = path[:p.length]
		return "", "", nil
	}

	p.part = path[:pos-1]
	p.length = pos - 1

	return path[pos-1:], "", nil
}

func (p *partFixed) match(path string) (matched bool, key, value string, length int) {
	if len(path) >= p.length && path[:p.length] == p.part {
		return true, "", "", p.length
	}
	return false, "", "", 0
}

func (p *partFixed) render(values map[string]string, normalize map[string]struct{}) (string, []string, error) {
	return p.part, []string{}, nil
}

func (p *partParam) read(path string) (string, string, error) {
	parts := strings.SplitN(path, "/", 2)

	if len(parts) > 0 {
		p.name = parts[0][1:]
	}

	if len(parts) > 1 {
		return `/` + parts[1], p.name, nil
	}

	p.name = path[1:]

	parts = strings.SplitN(p.name, ".", 2)
	if len(parts) > 1 {
		p.suffix = `.` + parts[1]
	}
	p.name = parts[0]

	return "", p.name, nil
}

func (p *partParam) match(path string) (matched bool, key, value string, length int) {
	parts := strings.SplitN(path, "/", 2)

	if len(parts) < 1 {
		return false, "", "", 0
	}

	if !strings.HasSuffix(parts[0], p.suffix) {
		return false, "", "", 0
	}

	val, _ := url.QueryUnescape(parts[0][:len(parts[0])-len(p.suffix)])
	return true, p.name, val, len(parts[0])
}

func (p *partParam) render(values map[string]string, normalize map[string]struct{}) (string, []string, error) {
	if value, ok := values[p.name]; ok {
		if _, ok := normalize[p.name]; ok {
			value = URLTitle(value)
		}
		return url.QueryEscape(value) + p.suffix, []string{p.name}, nil
	}
	return "", []string{}, errors.New("param " + p.name + " not found")
}

var partRegexMatch = regexp.MustCompile(`([^<]*)<([^>]+)>(.*)`)

func (p *partRegex) read(path string) (string, string, error) {
	var matches = partRegexMatch.FindStringSubmatch(path[1:])

	if matches == nil {
		return "", "", errors.New("no regex found")
	}

	p.name, p.regex = matches[1], regexp.MustCompile(`^`+matches[2])

	return matches[3], p.name, nil
}

func (p *partRegex) match(path string) (matched bool, key, value string, length int) {
	var matches = p.regex.FindStringSubmatch(path)
	if matches == nil {
		return false, "", "", 0
	}
	return true, p.name, matches[0], len(matches[0])
}

func (p *partRegex) render(values map[string]string, normalize map[string]struct{}) (string, []string, error) {
	if value, ok := values[p.name]; ok {
		if _, ok := normalize[p.name]; ok {
			value = URLTitle(value)
		}
		if p.regex.FindStringSubmatch(value) != nil {
			return value, []string{p.name}, nil
		}
		return "", []string{}, errors.New("param " + p.name + " in wrong format")
	}
	return "", []string{}, errors.New("param " + p.name + " not found")
}

func (p *partWildcard) read(path string) (string, string, error) {
	parts := strings.SplitN(path, "/", 2)

	if len(parts) > 0 {
		p.name = parts[0][1:]
	}

	if len(parts) > 1 {
		return `/` + parts[1], p.name, nil
	}

	p.name = path[1:]
	return "", p.name, nil
}

func (p *partWildcard) match(path string) (matched bool, key, value string, length int) {
	return true, p.name, path, len(path)
}

func (p *partWildcard) render(values map[string]string, normalize map[string]struct{}) (string, []string, error) {
	if value, ok := values[p.name]; ok {
		if _, ok := normalize[p.name]; ok {
			value = URLTitle(value)
		}
		return value, []string{p.name}, nil
	}
	return "", []string{}, nil
}

// NewPath returns a new path
func NewPath(path string) (*Path, error) {
	var newPath = &Path{
		path:          path,
		trailingSlash: strings.HasSuffix(path, "/"),
	}

	var current part
	var param string
	var err error

	for len(path) > 1 {
		if path[0] != '/' {
			return nil, fmt.Errorf("path %q corrupted", path)
		}
		path = path[1:]

		switch path[0] {
		case ':':
			current = new(partParam)
			path, param, err = current.read(path)
			if err != nil {
				return nil, err
			}
			newPath.parts = append(newPath.parts, current)

		case '$':
			current = new(partRegex)
			path, param, err = current.read(path)
			if err != nil {
				return nil, err
			}
			newPath.parts = append(newPath.parts, current)

		case '*':
			current = new(partWildcard)
			path, param, err = current.read(path)
			if err != nil {
				return nil, err
			}
			newPath.parts = append(newPath.parts, current)

		default:
			current = new(partFixed)
			path, param, err = current.read(path)
			if err != nil {
				return nil, err
			}
			newPath.parts = append(newPath.parts, current)
		}

		if param != "" {
			newPath.params = append(newPath.params, param)
		}
	}

	sort.Strings(newPath.params)

	return newPath, nil
}

// Match matches a given path
func (p *Path) Match(path string) *Match {
	var match = &Match{
		Values: make(map[string]string),
	}

	for _, part := range p.parts {
		if len(path) < 1 {
			return nil
		}

		if path[0] != '/' {
			return nil
		}
		// prefix /
		path = path[1:]

		matched, key, value, length := part.match(path)

		// log.Printf("%#v == %v (%d) %s", part, matched, length, value)

		if !matched {
			return nil
		}

		if key != "" {
			match.Values[key] = value
		}
		path = path[length:]
	}

	if len(path) > 0 && path != "/" {
		return nil
	}

	return match
}

// Render a path for a given list of values
func (p *Path) Render(values map[string]string, usedValues map[string]struct{}) (string, error) {
	var path string

	for _, part := range p.parts {
		val, used, err := part.render(values, p.normalize)
		if err != nil {
			return "", err
		}

		path += `/` + val
		for _, u := range used {
			usedValues[u] = struct{}{}
		}
	}

	if len(path) == 0 {
		path = "/"
	} else if p.trailingSlash && !strings.HasSuffix(path, "/") {
		path += "/"
	}

	query := url.Values{}
	queryUsed := false
	for k, v := range values {
		if _, ok := usedValues[k]; !ok {
			queryUsed = true
			query.Set(k, v)
		}
	}

	if queryUsed {
		return path + `?` + query.Encode(), nil
	}
	return path, nil
}
