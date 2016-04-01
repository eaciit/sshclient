package sshclient

import (
	"fmt"
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
	t.Skip("Skip : Comment this line to do test")
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
	t.Skip("Skip : Comment this line to do test")
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

func TestSshReadRecHistory(t *testing.T) {
	// t.Skip("Skip : Comment this line to do test")
	var SshClient SshSetting

	SshClient.SSHHost = "192.168.56.101:22"
	SshClient.SSHAuthType = 0
	SshClient.SSHUser = "eaciit1"
	SshClient.SSHPassword = "12345"

	output, err := SshClient.GetOutputCommandSsh(`/home/eaciit1/src/github.com/eaciit/sedotan/sedotan.v2/sedotanread/sedotanread -readtype=rechistory -recfile=/home/eaciit1/src/github.com/eaciit/sedotan/sedotan.v2/sedotanread/irondcecomcn.Iron01-20160316022830.csv`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
}

func TestSshReadHistory(t *testing.T) {
	// t.Skip("Skip : Comment this line to do test")
	var SshClient SshSetting

	SshClient.SSHHost = "192.168.56.101:22"
	SshClient.SSHAuthType = 0
	SshClient.SSHUser = "eaciit1"
	SshClient.SSHPassword = "12345"

	output, err := SshClient.GetOutputCommandSsh(`/home/eaciit1/src/github.com/eaciit/sedotan/sedotan.v2/sedotanread/sedotanread -readtype=history -pathfile=/home/eaciit1/src/github.com/eaciit/sedotan/sedotan.v2/sedotanread/HIST-GRABDCE-20160316.csv`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
}

func TestSshReadSnapShot(t *testing.T) {
	// t.Skip("Skip : Comment this line to do test")
	var SshClient SshSetting

	SshClient.SSHHost = "192.168.56.101:22"
	SshClient.SSHAuthType = 0
	SshClient.SSHUser = "eaciit1"
	SshClient.SSHPassword = "12345"

	output, err := SshClient.GetOutputCommandSsh(`/home/eaciit1/src/github.com/eaciit/sedotan/sedotan.v2/sedotanread/sedotanread -readtype=snapshot -nameid=irondcecomcn -pathfile=/home/eaciit1/src/github.com/eaciit/sedotan/sedotan.v2/sedotanread/daemonsnapshot.csv`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)
}

func TestCommandAsMap(t *testing.T) {
	var SshClient SshSetting

	SshClient.SSHHost = "192.168.56.101:22"
	SshClient.SSHAuthType = 0
	SshClient.SSHUser = "eaciit1"
	SshClient.SSHPassword = "12345"

	res, err := SshClient.RunCommandSshAsMap("pwd", "echo 1", "pwd", "echo 1")
	fmt.Printf("--- %#v\n", err)
	fmt.Printf("--- %#v\n", res)
}
