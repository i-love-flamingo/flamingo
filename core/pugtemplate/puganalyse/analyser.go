package puganalyse

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type (
	AtomicDesignAnalyser struct {
		baseDir        string
		HasError       bool
		AllowedPugLibs []string
		CheckCount     int
	}

	JsDependencyAnalyser struct {
		baseDir             string
		HasError            bool
		AllowedJsLibFolders []string
		CheckCount          int
	}
)

func NewAtomicDesignAnalyser(baseDir string) AtomicDesignAnalyser {
	return AtomicDesignAnalyser{
		baseDir:        baseDir,
		AllowedPugLibs: []string{"/", "/shared"},
	}
}

func (a *AtomicDesignAnalyser) CheckPugImports() {

	a.checkPugsInDir(filepath.Join(a.baseDir, "atom"), false, nil)
	a.checkPugsInDir(filepath.Join(a.baseDir, "molecule"), false, []string{"atom"})
	a.checkPugsInDir(filepath.Join(a.baseDir, "organism"), false, []string{"atom", "molecule"})
	a.checkPugsInDir(filepath.Join(a.baseDir, "template"), true, []string{"atom", "molecule", "organism", "template"})
	a.checkPugsInDir(filepath.Join(a.baseDir, "page"), true, []string{"atom", "molecule", "organism", "template"})
}

func (a *AtomicDesignAnalyser) checkPugsInDir(dir string, checkExtends bool, allowed []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(fmt.Sprintf("Warning - checking %v Failed: (%v)", dir, err))
		return
	}

	allowedPathsInAllLibs := []string{}
	for _, lib := range a.AllowedPugLibs {
		for _, allowedPath := range allowed {
			allowedPathsInAllLibs = append(allowedPathsInAllLibs, filepath.Join(lib, allowedPath))
		}
	}

	for _, f := range files {
		filePath := filepath.Join(dir, f.Name())
		if f.IsDir() {
			a.checkPugsInDir(filePath, checkExtends, allowed)
		}
		if strings.HasSuffix(f.Name(), ".pug") {
			a.CheckCount++
			b, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}

			re := regexp.MustCompile(`include (.*)`)
			matches := re.FindAllStringSubmatch(string(b), -1)
			for _, match := range matches {
				if !inSlicePrefix(match[1], allowedPathsInAllLibs) {
					relFilePath, _ := filepath.Rel(a.baseDir, filePath)
					fmt.Println(fmt.Sprintf("🐛 ERROR  File %v Contains unallowed include to %v", relFilePath, match[1]))
					fmt.Println(fmt.Sprintf("       Allowed is %#v", allowedPathsInAllLibs))

					a.HasError = true
				}
			}

			re = regexp.MustCompile(`extends (.*)`)
			matches = re.FindAllStringSubmatch(string(b), -1)
			for _, match := range matches {
				relFilePath, _ := filepath.Rel(a.baseDir, filePath)
				if !checkExtends {
					fmt.Println(fmt.Sprintf("🐛 ERROR  File %v Contains unallowed extend to %v (extends not allowed for this element)", relFilePath, match[1]))
					a.HasError = true
					continue
				}
				if !inSlicePrefix(match[1], allowedPathsInAllLibs) {
					fmt.Println(fmt.Sprintf("🐛 ERROR  File %v Contains unallowed extend to %v", relFilePath, match[1]))
					fmt.Println(fmt.Sprintf("       Allowed is %#v", allowedPathsInAllLibs))
					a.HasError = true
				}
			}
		}
	}
}

func NewJsDependencyAnalyser(baseDir string) JsDependencyAnalyser {
	return JsDependencyAnalyser{
		baseDir: baseDir,
	}
}

func (a *JsDependencyAnalyser) Check() {
	a.checkJsDep(a.baseDir)
}

func (a *JsDependencyAnalyser) checkJsDep(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range files {
		filePath := filepath.Join(dir, f.Name())
		relDir, _ := filepath.Rel(a.baseDir, dir)
		parts := strings.Split(relDir, "/")
		relativeDeep := len(parts)
		relFilePath, _ := filepath.Rel(a.baseDir, filePath)

		if f.IsDir() {
			a.checkJsDep(filePath)
		}
		if strings.HasSuffix(f.Name(), ".js") {
			a.CheckCount++
			b, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}

			re := regexp.MustCompile(`import (.*) from (.*)`)
			matches := re.FindAllStringSubmatch(string(b), -1)

			for _, match := range matches {
				importPath := match[2]
				importPathRelativeBackCount := strings.Count(importPath, "../")
				if importPathRelativeBackCount > relativeDeep {
					fmt.Println(fmt.Sprintf("🐛  ERROR File %v (%v) Contains import outside basePath! %v (%v)", relFilePath, relativeDeep, importPath, importPathRelativeBackCount))
				}
			}
		}
	}
}

func inSlicePrefix(checkString string, allowedPrefixes []string) bool {
	for _, allowedPrefix := range allowedPrefixes {
		if strings.HasPrefix(checkString, allowedPrefix) {
			return true
		}
	}
	return false
}
