package sshclient

import (
	. "github.com/eaciit/sshclient"
	"os"
	"path/filepath"
	"testing"
)

func TestSshKey(t *testing.T) {
	t.Skip("Skip : Comment this line to do test")
	var SshClient SshSetting

	SshClient.SSHAuthType = SSHAuthType_Certificate
	SshClient.SSHHost = "192.168.56.101:22"
	SshClient.SSHUser = "alip"
	SshClient.SSHKeyLocation = "C:\\Users\\User\\.ssh\\id_rsa"

	ps := []string{"sudo service mysql status"}
	res, e := SshClient.RunCommandSsh(ps...)

	if e != nil {
		t.Errorf("Error, %s \n", e)
	} else {
		t.Logf("RUN, %s \n", res)
	}
}

func TestSshUsername(t *testing.T) {
	// t.Skip("Skip : Comment this line to do test")
	var SshClient SshSetting

	SshClient.SSHAuthType = SSHAuthType_Password
	SshClient.SSHHost = "192.168.56.101:22"
	SshClient.SSHUser = "alip"
	SshClient.SSHPassword = "Bismillah"

	ps := []string{"sudo service mysql status"}
	res, e := SshClient.RunCommandSsh(ps...)

	if e != nil {
		t.Errorf("Error, %s \n", e)
	} else {
		t.Logf("RUN, %s \n", res)
	}
}

func TestSshCopyFilePath(t *testing.T) {
	t.Skip("Skip : Comment this line to do test")
	var SshClient SshSetting

	SshClient.SSHAuthType = SSHAuthType_Password
	SshClient.SSHHost = "192.168.56.101:22"
	SshClient.SSHUser = "alip"
	SshClient.SSHPassword = "Bismillah"

	filepath := "E:\\goproject\\src\\github.com\\eaciit\\sshclient\\test\\live_test.go"
	destination := "/home/alip"

	e := SshClient.SshCopyByPath(filepath, destination)
	if e != nil {
		t.Errorf("Error, %s \n", e)
	} else {
		t.Logf("Copy File Success")
	}
}

func TestSshCopyFileDirect(t *testing.T) {
	// t.Skip("Skip : Comment this line to do test")
	var SshClient SshSetting

	SshClient.SSHAuthType = SSHAuthType_Password
	SshClient.SSHHost = "192.168.56.101:22"
	SshClient.SSHUser = "alip"
	SshClient.SSHPassword = "Bismillah"

	filePath := "E:\\goproject\\src\\github.com\\eaciit\\sshclient\\test\\live_test.go"
	//Prepare File=============
	f, err := os.Open(filePath)
	if err != nil {
		t.Errorf(err.Error())
	}

	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		t.Errorf(err.Error())
	}
	//========================
	destination := "/home/alip"

	e := SshClient.SshCopyByFile(f, s.Size(), s.Mode().Perm(), filepath.Base(f.Name()), destination)
	if e != nil {
		t.Errorf("Error, %s \n", e)
	} else {
		t.Logf("Copy File Success")
	}
}
