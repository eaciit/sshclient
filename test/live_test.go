package sshclient

import (
	. "github.com/eaciit/sshclient"
	"testing"
)

func TestSshKey(t *testing.T) {
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
