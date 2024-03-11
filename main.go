package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"tehdas/tehdas"
)

const name = ".tehdas/"

func init() {
	var err error
	cfg, err = os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	cfg, err = filepath.Abs(filepath.Join(cfg, name))
	Must(err, "error getting absolute path for cfg")
	cfgFile, err = filepath.Abs(filepath.Join(cfg, ".conf"))
	Must(err, "error getting absolute path for cfg file")

	cwd, err = os.Getwd()
	Must(err, "error getting current working directory")
	if len(os.Args) > 1 {
		arg := os.Args[1]
		Must(err, "error getting absolute path for arg")
		source := filepath.Join(cwd, arg)
		fileinfo, err := os.Stat(source)
		printSLn("WD", cwd, "ARG", arg, "SOURCE", source)
		switch {
		case os.IsNotExist(err):
			log.Fatalln("file does not exist:", arg)
		case fileinfo.IsDir():
			entry = "."
		default:
			entry = source
		}
		return
	}

	// if no args, use current directory
	entry = "."
}

// printSLn -- print structured line
func printSLn(pairs ...string) {
	if len(pairs)%2 != 0 {
		log.Fatalln("input must be even length, got", len(pairs), "instead")
	}

	maxL := 0
	for i := 0; i < len(pairs); i += 2 {
		if len(pairs[i]) > maxL {
			maxL = len(pairs[i])
		}
	}

	for i := 0; i < len(pairs); i += 2 {
		k := pairs[i]
		v := pairs[i+1]
		fmt.Printf("%-*s:   %s\n", maxL, k, v)
	}
}

var (
	cwd     string
	cfg     string
	cfgFile string
	entry   string
	// buildName string
)

func configureTehdas() string {
	fmt.Println("where do you want binaries to be compiled to:")
	var tar string
	_, err := fmt.Scanln(&tar)
	// check that err isnt' unexpected new line
	if err.Error() != "unexpected newline" {
		Must(err, "error scanning input")
	}
	if tar == "default" || tar == "" {
		tar = cfg
	}

	return tar
}

func main() {
	Must(EnsureExists(cfg, true), "error ensuring dir", cfg, "exists")
	Must(EnsureExists(cfgFile, false), "error ensuring file", cfgFile, "exists")

	printSLn("CFG", cfg, "CFGFILE", cfgFile, "ENTRY", entry)
	// open in rw mode
	f, err := os.OpenFile(cfgFile, os.O_RDWR, 0644)
	Must(err, "error opening file", cfgFile)
	defer f.Close()

	if err := tehdas.Inst.Decode(f); tehdas.IsEmptyErr(err) {
		tar := configureTehdas()
		tehdas.Inst.Add("target", tar)
		Must(tehdas.Inst.Encode(f), "error encoding file", cfgFile)
	} else {
		Must(err, "error decoding file", cfgFile)
	}

	target := tehdas.Inst.MustGet("target")
	Must(EnsureExists(target, true), "error ensuring dir", target, "exists")

	fmt.Println("Building project from entry:", entry, "to target:", target)
	buildFn := inferBuildCommand()
	Must(buildFn(entry, target), "error building", entry)
}

func inferBuildCommand() func(string, string) error { return Go }

func Go(path, tar string) error {
	// if current folder does not have go.mod, fail
	if _, err := os.Stat(filepath.Join(cwd, "go.mod")); os.IsNotExist(err) {
		return fmt.Errorf("no go.mod found in current directory; is not a compileable go project")
	}
	args := []string{"go", "build", "-o", tar, path}
	fmt.Println("Running command:", strings.Join(args, " "))

	cmd := exec.Command("go", "build", "-o", tar, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Must(err error, args ...string) {
	msg := strings.Join(args, " ")
	if err != nil {
		if len(msg) > 0 {
			fmt.Println(msg)
		}
		panic(err)
	}
}

func EnsureExists(tar string, dir bool) error {
	fileInfo, err := os.Stat(tar)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalln("os.Stat caused unexpected error", err)
	}
	if dir {
		if err == nil {
			if fileInfo.IsDir() {
				return nil
			}
			return fmt.Errorf("path exists but is not a directory: %s", tar)
		}

		if os.IsNotExist(err) {
			err = os.MkdirAll(tar, 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory: %s", err)
			}
			return nil
		}

		return fmt.Errorf("unknown error: %s", err)
	}

	if os.IsNotExist(err) {
		f, err := os.Create(tar)
		if err != nil {
			return fmt.Errorf("failed to create directory: %s", err)
		}
		f.Close()
		return nil
	}

	return nil
}
