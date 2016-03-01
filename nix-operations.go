package sshclient

import (
	"fmt"
	"github.com/eaciit/errorlib"
	"io"
	// "os"
)

const (
	TYPE_FOLDER = "dir"
	TYPE_FILE   = "file"
)

const (
	LIST             = "ls -all %s"
	MKDIR            = "mkdir -m %s %s"
	REMOVE           = "rm -f %s"
	REMOVE_RECURSIVE = "rm -R -f %s"
	RENAME           = "mv %s %s"
	MKFILE           = "echo \"%s\" > %s"
	CHMOD            = "chmod %v %v"
	CHOWN            = "chwon %v:%v %v"
	CHOWN_RECURSIVE  = "chwon -R %v:%v %v"
)

func List(s SshSetting, path string) (string, error) {
	if path == "" {
		return "", errorlib.Error("", "", "LIST", "Path is Undivined")
	}

	cmd := fmt.Sprintf(LIST, path)

	res, e := s.GetOutputCommandSsh(cmd)

	if e != nil {
		e = errorlib.Error("", "", "LIST", e.Error())
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

func MakeFile(s SshSetting, content string, path string, permission string) error {
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
	}

	e = Chmod(s, path, permission)

	if e != nil {
		return errorlib.Error("", "", "MAKEFILE", e.Error())
	}

	return e
}

func Chmod(s SshSetting, path string, permission string) error {
	if path == "" {
		return errorlib.Error("", "", "CHMOD", "Path is Undivined")
	}

	if permission == "" {
		return errorlib.Error("", "", "CHMOD", "Permission is Undivined")
	}

	cmd := fmt.Sprintf(CHMOD, permission, path)
	_, e := s.GetOutputCommandSsh(cmd)

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

func UploadFile(s SshSetting, content io.Reader, size int64, filename string, destination string) error {
	/*if destination == "" {
		return errorlib.Error("", "", "UPLOAD_FILE", "Destination is Undivined")
	}

	e := s.SshCopyByFile(content, size, os.FileMode.ModeAppend, filename, destination)

	if e != nil {
		e = errorlib.Error("", "", "UPLOAD_FILE", e.Error())
	}

	return e*/
	return nil
}
