package split

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	logFile          *os.File
	home             string
	keyValRe, credRe *regexp.Regexp
	configPath       string
	config           Config
	currentCredFile  string
)

func init() {
	keyValRe = regexp.MustCompile(`([A-Za-z][A-Za-z\d\-]*)\s*=\s*(\S+)`)
	credRe = regexp.MustCompile(`(\S+?)://(\S+?):(\S+?)@(\S+?)/(\S+)`)
}

func Split() error {
	var (
		err error
	)

	if err = prepareLog(); err != nil {
		return err
	}
	defer closeLog()

	if err = applyConfigFlags(); err != nil {
		log.Println("Error applying flag:", err)
		return err
	}

	//TODO: need to apply config; there is log related param in config

	if err = commands(); err != nil {
		log.Println(err)
		return err
	}

	return nil
	//return errors.New("debug!") //TODO: fix me! debug only
}

func commands() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working dir: %v", err)
	}
	log.Println("current directory:", cwd)
	log.Println("command:", flag.Args())
	cmd := flag.Arg(0)
	switch cmd {
	case "init":
	case "add":
	case "rm":
	case "list":
	case "store":
	case "get":
		return cmdGet()
	case "erase":
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return nil
}

func applyConfigFlags() error {
	var (
		err                error
		pathRegistered     bool
		path               string
		parts              []string
		cwd                string
		pathSep, pathCheck string
	)

	flag.StringVar(&configPath, "f", "", "config path")
	flag.Parse()

	if config, err = readConfig(configPath); err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}

	if cwd, err = os.Getwd(); err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	//pathSep = fmt.Sprintf("%v", os.PathSeparator)
	pathSep = string(os.PathSeparator)
	log.Println("separator:", pathSep)
	pathRegistered = false
	for _, path = range config.CredFiles {
		parts = strings.Split(path, pathSep)
		if len(parts) == 1 {
			pathCheck = pathSep + parts[0]
		} else if len(parts) > 1 {
			pathCheck = strings.Join(parts[:len(parts)-1], pathSep)
			if !strings.HasPrefix(pathCheck, pathSep) {
				pathCheck = pathSep + pathCheck
			}
		} else {
			continue
		}

		log.Println("check path:", path, pathCheck, cwd)
		if strings.HasPrefix(cwd, pathCheck) {
			pathRegistered = true
			break
		}
	}

	if !pathRegistered {
		return errors.New("path unregistered")
	}
	currentCredFile = path

	return nil
}

func prepareLog() error {
	var (
		err error
		ok  bool
	)

	home, ok = os.LookupEnv("HOME")
	if !ok {
		return errors.New("cannot find HOME environment variable")
	}

	logFile, err = os.OpenFile(home+"/.git-credential-split.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		// cannot open/create log file, bail out
		return err
	}

	log.SetOutput(logFile)
	return nil
}

func closeLog() {
	logFile.Close()
	log.SetOutput(os.Stderr)
}
