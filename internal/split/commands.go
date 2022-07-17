package split

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type credentialItem struct {
	proto    string
	username string
	password string
	host     string
	path     string
}

func cmdGet() error {
	var (
		err      error
		param    string
		matches  [][]string
		paramMap map[string]string
		credFile *os.File
		scan     *bufio.Scanner
		credList []credentialItem
		credItem credentialItem
		ok       bool
	)

	if credFile, err = os.Open(currentCredFile); err != nil {
		return fmt.Errorf("error opening credential file %s: %v", currentCredFile, err)
	}
	defer credFile.Close()

	scan = bufio.NewScanner(credFile)
	for scan.Scan() {
		param = scan.Text()
		if !credRe.MatchString(param) {
			return fmt.Errorf("got invalid match in %s: %s", currentCredFile, param)
		}
		log.Println("match for param:", param)

		matches = credRe.FindAllStringSubmatch(param, -1)
		log.Println("sub matches:", len(matches))
		if len(matches) > 0 {
			credItem = credentialItem{
				proto:    matches[0][1],
				username: matches[0][2],
				password: matches[0][3],
				host:     matches[0][4],
				path:     matches[0][5],
			}
			log.Printf("insert to the list: %v\n", credItem)
			credList = append(credList, credItem)
		}
	}
	log.Println("number of credential read:", len(credList))

	if err = scan.Err(); err != nil {
		return err
	}

	paramMap = make(map[string]string)

	scan = bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		param = scan.Text()
		log.Println("got input:", param)

		if !keyValRe.MatchString(param) {
			return fmt.Errorf("got invalid parameter from git: %s", param)
		}

		matches = keyValRe.FindAllStringSubmatch(param, -1)
		if len(matches) > 0 && len(matches[0]) == 3 {
			paramMap[strings.ToLower(matches[0][1])] = matches[0][2]
		}
	}

	if err = scan.Err(); err != nil {
		return err
	}

	log.Printf("input: %+v\n", paramMap)
	ok = false
	for _, credItem = range credList {
		log.Printf("item: %+v\n", credItem)
		if credItem.proto == paramMap["protocol"] && credItem.host == paramMap["host"] &&
			pathMatch(credItem.path, paramMap["path"]) {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("username=%s", credItem.username))
			fmt.Fprintln(os.Stdout, fmt.Sprintf("password=%s", credItem.password))
			log.Println("got credential match:", credItem)
			ok = true
			break
		}
	}

	if !ok {
		return errors.New("no credential match")
	}

	return nil
}

func pathMatch(path1, path2 string) bool {
	if path1 == path2 {
		return true
	}

	if strings.HasSuffix(path1, "*") && strings.HasPrefix(path2, strings.TrimSuffix(path1, "*")) {
		return true
	}

	if strings.HasSuffix(path2, "*") && strings.HasPrefix(path1, strings.TrimSuffix(path2, "*")) {
		return true
	}

	return false
}
