package router

import (
	"errors"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

type (
	// Path is a matchable routing path
	Path struct {
		path   string // `/path/to/:something/$id<[0-9]+>/*foo
		parts  []part // {'path/to', ':something', '$id<[0-9]+>', '*foo'}
		params []string
	}

	// Match is the result of matching a path
	Match struct {
		Values map[string]string
	}

	part interface {
		match(path string) (matched bool, key, value string, length int)
		read(segment string) (leftover, paramname string)
		render(values map[string]string) (string, error)
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

func (p *partFixed) read(path string) (string, string) {
	pos := strings.IndexAny(path, ":*$")

	if pos < 0 {
		p.part = path
		p.length = len(path)
		return "", ""
	}

	p.part = path[:pos-1]
	p.length = pos - 1

	return path[pos-1:], ""
}

func (p *partFixed) match(path string) (matched bool, key, value string, length int) {
	if len(path) >= p.length && path[:p.length] == p.part {
		return true, "", "", p.length
	}
	return false, "", "", 0
}

func (p *partFixed) render(values map[string]string) (string, error) {
	return p.part, nil
}

func (p *partParam) read(path string) (string, string) {
	parts := strings.SplitN(path, "/", 2)

	if len(parts) > 0 {
		p.name = parts[0][1:]
	}

	if len(parts) > 1 {
		return `/` + parts[1], p.name
	}

	p.name = path[1:]

	parts = strings.SplitN(p.name, ".", 2)
	if len(parts) > 1 {
		p.suffix = `.` + parts[1]
	}
	p.name = parts[0]

	return "", p.name
}

func (p *partParam) match(path string) (matched bool, key, value string, length int) {
	parts := strings.SplitN(path, "/", 2)

	if len(parts) < 1 {
		return false, "", "", 0
	}

	if !strings.HasSuffix(parts[0], p.suffix) {
		return false, "", "", 0
	}

	return true, p.name, parts[0][:len(parts[0])-len(p.suffix)], len(parts[0])
}

func (p *partParam) render(values map[string]string) (string, error) {
	if value, ok := values[p.name]; ok {
		return url.QueryEscape(value) + p.suffix, nil
	}
	return "", errors.New("param " + p.name + " not found")
}

var partRegexMatch = regexp.MustCompile(`([^<]*)<([^>]+)>(.*)`)

func (p *partRegex) read(path string) (string, string) {
	var matches = partRegexMatch.FindStringSubmatch(path[1:])

	if matches == nil {
		panic("ain't no regex? o_O")
	}

	p.name, p.regex = matches[1], regexp.MustCompile(`^`+matches[2])

	return matches[3], p.name
}

func (p *partRegex) match(path string) (matched bool, key, value string, length int) {
	var matches = p.regex.FindStringSubmatch(path)
	if matches == nil {
		return false, "", "", 0
	}
	return true, p.name, matches[0], len(matches[0])
}

func (p *partRegex) render(values map[string]string) (string, error) {
	if value, ok := values[p.name]; ok {
		if p.regex.FindStringSubmatch(value) != nil {
			return value, nil
		}
		return "", errors.New("param " + p.name + " in wrong format")
	}
	return "", errors.New("param " + p.name + " not found")
}

func (p *partWildcard) read(path string) (string, string) {
	parts := strings.SplitN(path, "/", 2)

	if len(parts) > 0 {
		p.name = parts[0][1:]
	}

	if len(parts) > 1 {
		return `/` + parts[1], p.name
	}

	p.name = path[1:]
	return "", p.name
}

func (p *partWildcard) match(path string) (matched bool, key, value string, length int) {
	return true, p.name, path, len(path)
}

func (p *partWildcard) render(values map[string]string) (string, error) {
	if value, ok := values[p.name]; ok {
		return value, nil
	}
	return "", nil
}

// NewPath returns a new path
func NewPath(path string) *Path {
	var newPath = &Path{
		path: path,
	}

	var current part
	var param string

	for len(path) > 1 {
		if path[0] != '/' {
			panic("Path " + path + " corrupted")
		}
		path = path[1:]

		switch path[0] {
		case ':':
			current = new(partParam)
			path, param = current.read(path)
			newPath.parts = append(newPath.parts, current)

		case '$':
			current = new(partRegex)
			path, param = current.read(path)
			newPath.parts = append(newPath.parts, current)

		case '*':
			current = new(partWildcard)
			path, param = current.read(path)
			newPath.parts = append(newPath.parts, current)

		default:
			current = new(partFixed)
			path, param = current.read(path)
			newPath.parts = append(newPath.parts, current)
		}

		if param != "" {
			newPath.params = append(newPath.params, param)
		}
	}

	sort.Strings(newPath.params)

	return newPath
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

		log.Printf("%#v == %v (%d) %s", part, matched, length, value)

		if !matched {
			return nil
		}

		if key != "" {
			match.Values[key] = value
		}
		path = path[length:]
	}

	//log.Printf("%s", path)

	if len(path) > 1 {
		return nil
	}

	if len(path) == 1 && path != "/" {
		return nil
	}

	return match
}

// Render a path for a given list of values
func (p *Path) Render(values map[string]string) (string, error) {
	var path string

	for _, part := range p.parts {
		val, err := part.render(values)
		if err != nil {
			return "", err
		}

		//log.Printf("%#v: %s", part, val)

		path += `/` + val
	}

	if len(path) == 0 {
		path = "/"
	}

	return path, nil
}
