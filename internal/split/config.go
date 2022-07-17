package split

import (
	"bufio"
	"io"
	"log"
	"os"
)

type Config struct {
	Log       bool
	LogPath   *string
	CredFiles []string
}

func readConfig(configPath string) (config Config, err error) {
	var (
		configFile *os.File
		line       string
		output     interface{}
		ok         bool
		params     []string
	)

	log.Println("reading config:", configPath)
	if len(configPath) == 0 {
		if config, err = readConfig(home + "/.config/git-credential-split"); err != nil {
			if config, err = readConfig(home + "/.config/.git-credentials"); err != nil {
				return readConfig(home + "/.git-credentials")
			} else {
				return
			}
		} else {
			return
		}
	}

	if configFile, err = os.Open(configPath); err != nil {
		log.Println("error opening file", configPath, ":", err)
		return
	}
	defer configFile.Close()
	log.Println("open file", configPath, "success!")

	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > 0 {
			if output, ok = isParam(line); ok {
				if params, ok = output.([]string); ok && len(params) == 2 {
					_ = params
				}
			} else {
				config.CredFiles = append(config.CredFiles, line)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		if err == io.EOF {
			err = nil
		}
	}
	return
}

func isParam(text string) (interface{}, bool) {
	if keyValRe.MatchString(text) {

	}

	return text, false
}
