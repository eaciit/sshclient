package sshclient

import (
	"fmt"

	"github.com/eaciit/errorlib"
	// "io"
	// "os"
	// "log"
	"strings"
)

const (
	LIST         = "ls -l %v"
	LIST_PARAM   = "ls -l %v | awk '{ print $1,\"||\",$2,\"||\",$3,\"||\",$4,\"||\",$5,\"||\",$6,\"||\",$7,\"||\",$8,\"||\",$9,$10,$11}'"
	SEARCH       = "find %v -name *%v* -ls | awk '{print $3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13}'"
	SEARCH_PARAM = "find %v -name *%v* -ls | awk '{ print $3,\"||\",$4,\"||\",$5,\"||\",$6,\"||\",$7,\"||\",$8,\"||\",$9,\"||\",$10,\"||\",$11,$12,$13}'"
	// SEARCH       = "find %v -name *%v* | xargs -r -l ls -l"
	// SEARCH_PARAM = "find %v -name *%v* | xargs -r -l ls -l | awk '{ print $1,\"||\",$2,\"||\",$3,\"||\",$4,\"||\",$5,\"||\",$6,\"||\",$7,\"||\",$8,\"||\",$9,$10,$11}'"
	// SEARCH           = "ls -l -R %s | grep %s"
	// SEARCH_PARAM     = "ls -l -R %s | grep %s | awk '{ print $1,\"||\",$2,\"||\",$3,\"||\",$4,\"||\",$5,\"||\",$6,\"||\",$7,\"||\",$8,\"||\",$9,$10,$11}'"
	MKDIR            = "mkdir -m %v '%v'"
	REMOVE           = "rm -f %v"
	REMOVE_RECURSIVE = "rm -R -f %v"
	RENAME           = "mv '%v' '%v'"
	MKFILE           = "echo \"%v\" > '%v'"
	CHMOD            = "chmod %v %v"
	CHOWN            = "chwon %v:%v %v"
	CHOWN_RECURSIVE  = "chwon -R %v:%v %v"
	CAT              = "cat %v"
)

func List(s SshSetting, path string, isParseble bool) (string, error) {
	if path == "" {
		return "", errorlib.Error("", "", "LIST", "Path is Undivined")
	}

	cmd := ""

	if isParseble {
		cmd = fmt.Sprintf(LIST_PARAM, path)
	} else {
		cmd = fmt.Sprintf(LIST, path)
	}

	res, e := s.GetOutputCommandSsh(cmd)

	if e != nil {
		e = errorlib.Error("", "", "LIST", e.Error())
	}

	startStr := strings.Index(res, "\n")

	if res != "" {
		return res[startStr:], e
	}

	return res, e
}

func Search(s SshSetting, path string, isParseble bool, search string) (string, error) {
	if path == "" {
		return "", errorlib.Error("", "", "SEARCH", "Path is Undivined")
	}

	if search == "" {
		return "", errorlib.Error("", "", "SEARCH", "Search param is undivined")
	}

	cmd := ""

	if isParseble {
		cmd = fmt.Sprintf(SEARCH_PARAM, path, search)
	} else {
		cmd = fmt.Sprintf(SEARCH, path, search)
	}

	res, e := s.GetOutputCommandSsh(cmd)

	if e != nil {
		e = errorlib.Error("", "", "SEARCH", e.Error())
	}

	return res, e
}

func MakeDir(s SshSetting, path string, permission string) error {
	if path == "" {
		return errorlib.Error("", "", "MAKEDIR", "Path is Undivined")
	}

	if permission == "" {
		permission = "755"
	}

	cmd := fmt.Sprintf(MKDIR, permission, path)
	_, e := s.GetOutputCommandSsh(cmd)

	if e != nil {
		e = errorlib.Error("", "", "MAKEDIR", e.Error())
	}

	return e
}

func Rename(s SshSetting, oldPath string, newPath string) error {
	if oldPath == "" {
		return errorlib.Error("", "", "RENAME", "Old Path is Undivined")
	}

	if newPath == "" {
		return errorlib.Error("", "", "RENAME", "New Path is Undivined")
	}

	cmd := fmt.Sprintf(RENAME, oldPath, newPath)
	_, e := s.GetOutputCommandSsh(cmd)

	if e != nil {
		e = errorlib.Error("", "", "RENAME", e.Error())
	}

	return e
}

func Remove(s SshSetting, recursive bool, paths ...string) map[string]error {
	var es map[string]error

	if len(paths) < 1 {
		es[""] = errorlib.Error("", "", "REMOVE", "Paths is Undivined")
		return es
	}

	for _, path := range paths {
		cmd := ""
		if recursive {
			cmd = fmt.Sprintf(REMOVE_RECURSIVE, path)
		} else {
			cmd = fmt.Sprintf(REMOVE, path)
		}

		_, e := s.GetOutputCommandSsh(cmd)

		if e != nil {
			if es == nil {
				es = map[string]error{}
			}
			es[path] = errorlib.Error("", "", "REMOVE", e.Error())
		}
	}
	return es
}

func MakeFile(s SshSetting, content string, path string, permission string, format bool) error {
	if path == "" {
		return errorlib.Error("", "", "MAKEFILE", "Path is Undivined")
	}

	cmd := fmt.Sprintf(MKFILE, content, path)
	_, e := s.GetOutputCommandSsh(cmd)

	if e != nil {
		return errorlib.Error("", "", "MAKEFILE", e.Error())
	}

	if permission == "" {
		permission = "755"
		format = false
	}

	e = Chmod(s, path, permission, format)

	if e != nil {
		return errorlib.Error("", "", "MAKEFILE", e.Error())
	}

	return e
}

func constructPermission(strPermission string) (result string, err error) {
	permission_map := map[string]string{
		"rwx": "7",
		"rw-": "6",
		"r-x": "5",
		"r--": "4",
		"-wx": "3",
		"-w-": "2",
		"--x": "1",
		"---": "0",
	}

	if len(strPermission) == 9 {
		owner := permission_map[strPermission[:3]]
		group := permission_map[strPermission[3:6]]
		other := permission_map[strPermission[6:]]

		result = owner + group + other
	} else {
		err = errorlib.Error("", "", "Construct Permission", "The permission is not valid")
	}

	return
}

func Chmod(s SshSetting, path string, permission string, format bool) (e error) {
	if format {
		permission, e = constructPermission(permission)
		if e != nil {
			return errorlib.Error("", "", "CHMOD", e.Error())
		}
	}

	if path == "" {
		return errorlib.Error("", "", "CHMOD", "Path is Undivined")
	}

	if permission == "" {
		return errorlib.Error("", "", "CHMOD", "Permission is Undivined")
	}

	cmd := fmt.Sprintf(CHMOD, permission, path)
	_, e = s.GetOutputCommandSsh(cmd)

	if e != nil {
		return errorlib.Error("", "", "CHMOD", e.Error())
	}

	return e
}

func Chown(s SshSetting, path string, user string, group string, recursive bool) error {
	if path == "" {
		return errorlib.Error("", "", "CHOWN", "Path is Undivined")
	}

	if user == "" {
		return errorlib.Error("", "", "CHOWN", "User is Undivined")
	}

	if group == "" {
		return errorlib.Error("", "", "CHOWN", "Group is Undivined")
	}

	cmd := ""

	if recursive {
		cmd = fmt.Sprintf(CHOWN_RECURSIVE, user, group, path)
	} else {
		cmd = fmt.Sprintf(CHOWN, user, group, path)
	}

	_, e := s.GetOutputCommandSsh(cmd)

	if e != nil {
		return errorlib.Error("", "", "CHOWN", e.Error())
	}

	return e
}

func Cat(s SshSetting, path string) (string, error) {
	if path == "" {
		return "", errorlib.Error("", "", "CAT", "Path is Undivined")
	}

	cmd := fmt.Sprintf(CAT, path)
	res, e := s.GetOutputCommandSsh(cmd)
	if e != nil {
		e = errorlib.Error("", "", "CAT", e.Error())
	}

	return res, e
}

/*func UploadFile(s SshSetting, content io.Reader, size int64, filename string, destination string) error {
	if destination == "" {
		return errorlib.Error("", "", "UPLOAD_FILE", "Destination is Undivined")
	}

	e := s.SshCopyByFile(content, size, os.FileMode.ModeAppend, filename, destination)

	if e != nil {
		e = errorlib.Error("", "", "UPLOAD_FILE", e.Error())
	}

	return e
}*/
