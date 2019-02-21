package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func run(cmd string, args ...string) error {
	return exec.Command(cmd, args...).Run()
}

func getreplacements() []string {
	in, err := ioutil.ReadFile("go.mod.local")
	if err != nil {
		return nil
	}

	return strings.Split(strings.TrimSpace(string(in)), "\n")
}

func m() error {
	switch os.Args[1] {
	case "gomod-local":
		if err := run("git", "config", "filter.gomod.clean", "go run flamingo.me/flamingo/v3/build/dev gomod-clean"); err != nil {
			return err
		}
		if err := run("git", "config", "filter.gomod.smudge", "go run flamingo.me/flamingo/v3/build/dev gomod-smudge"); err != nil {
			return err
		}
		if err := run("git", "config", "filter.gomod.required", "true"); err != nil {
			return err
		}

	case "gomod-unlocal":
		if err := run("git", "config", "filter.gomod.clean", ""); err != nil {
			return err
		}
		if err := run("git", "config", "filter.gomod.smudge", ""); err != nil {
			return err
		}
		if err := run("git", "config", "filter.gomod.required", "false"); err != nil {
			return err
		}

	case "gomod-smudge":
		r := getreplacements()
		if len(r) == 0 {
			return nil
		}

		in, _ := ioutil.ReadAll(os.Stdin)

		for _, r := range r {
			r := strings.Split(r, "=")
			in = append(in, []byte("\nreplace "+r[0]+" => "+r[1])...)
		}

		fmt.Fprint(os.Stdout, string(in))

	case "gomod-clean":
		in, _ := ioutil.ReadAll(os.Stdin)
		r := strings.Split(string(in), "\n")
		for i, rr := range r {
			if strings.Contains(rr, "../") {
				r[i] = ""
			}
		}

		fmt.Fprint(os.Stdout, strings.Replace(strings.Join(r, "\n"), "\n\n", "\n", -1))

		in, err := ioutil.ReadFile("go.mod")
		if err != nil {
			return err
		}
		if bytes.Contains(in, []byte("=> ../")) {
			if _, err := fmt.Fprintln(os.Stderr, "can not commit local go.mod replacements!"); err != nil {
				return err
			}
			os.Exit(1)
		}
	}

	return nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("flamingo build/dev arguments:")
		fmt.Println("gomod-local: install git filter")
		fmt.Println("gomod-clean: check go.mod")
		return
	}

	if err := m(); err != nil {
		log.Fatal(err)
	}
}
