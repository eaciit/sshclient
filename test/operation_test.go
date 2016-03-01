package sshclient

import (
	. "github.com/frezadev/sshclient"
	// "os"
	// "path/filepath"
	"testing"
)

var SshClient SshSetting

func TestSshConnect(t *testing.T) {
	SshClient = SshSetting{
		SSHAuthType: SSHAuthType_Password,
		SSHHost:     "127.0.0.1:2222",
		SSHUser:     "root",
		SSHPassword: "Im4m.Bonj0l",
	}
}

func TestList(t *testing.T) {
	res, e := List(SshClient, "/root")

	if e != nil {
		t.Errorf("Error, %s \n", e)
	} else {
		t.Logf("RUN, %s \n", res)
	}
}

func TestMakeDir(t *testing.T) {
	e := MakeDir(SshClient, "/root/colony/test", "")

	if e != nil {
		t.Errorf("Error, %s \n", e)
	} else {
		t.Logf("RUN")
	}
}

func TestRenameDir(t *testing.T) {
	e := Rename(SshClient, "/root/colony/test", "/root/colony/testchange")

	if e != nil {
		t.Errorf("Error, %s \n", e)
	} else {
		t.Logf("RUN")
	}
}

func TestMakeFile(t *testing.T) {
	e := MakeFile(SshClient, "ini isinya ya", "/root/colony/testchange/file.txt", "")

	if e != nil {
		t.Errorf("Error, %s \n", e)
	} else {
		t.Logf("RUN")
	}
}

func TestRenameFile(t *testing.T) {
	e := Rename(SshClient, "/root/colony/testchange/file.txt", "/root/colony/testchange/file-new.txt")

	if e != nil {
		t.Errorf("Error, %s \n", e)
	} else {
		t.Logf("RUN")
	}
}
