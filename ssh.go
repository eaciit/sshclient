package sshclient

//Handle all ssh connection and run command

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SSHAuthTypeEnum int

const (
	//Type authentication login in ssh, it can using password or using public-private key
	SSHAuthType_Password SSHAuthTypeEnum = iota
	SSHAuthType_Certificate
)

type SshSetting struct {
	//Setting information for ssh connection
	SSHHost        string
	SSHUser        string
	SSHPassword    string
	SSHKeyLocation string
	SSHAuthType    SSHAuthTypeEnum
}

/*
Parsing private key certicate using for connection over ssh
*/
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

/*
Build connection ssh client to ssh server
*/
func (S *SshSetting) Connect() (*ssh.Client, error) {
	var (
		cfg *ssh.ClientConfig
	)

	if S.SSHAuthType == SSHAuthType_Certificate {
		cfg = &ssh.ClientConfig{
			User: S.SSHUser,
			Auth: []ssh.AuthMethod{
				PublicKeyFile(S.SSHKeyLocation),
			},
		}
	} else {
		cfg = &ssh.ClientConfig{
			User: S.SSHUser,
			Auth: []ssh.AuthMethod{
				ssh.Password(S.SSHPassword),
			},
		}
	}

	client, e := ssh.Dial("tcp", S.SSHHost, cfg)
	return client, e

}

/*
Handle input and output into terminal using channel
*/
func TermInOut(w io.Writer, r io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 1)
	out := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()
	go func() {
		var (
			buf [1024 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				close(in)
				close(out)
				return
			}
			t += n
			if buf[t-2] == '$' {
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()
	return in, out
}

/*
Create new session
*/
func (s *SshSetting) NewSession() (*ssh.Client, *ssh.Session, error) {
	c, e := s.Connect()
	if e != nil {
		e = fmt.Errorf("Unable to connect: %s", e.Error())
		return c, nil, e
	}

	Ses, e := c.NewSession()
	if e != nil {
		e = fmt.Errorf("Unable to start new session: %s", e.Error())
		return c, Ses, e
	}

	return c, Ses, e
}

/*
Build connection and run ssh script, catch the output or give error message if any
*/
func (S *SshSetting) RunCommandSsh(cmds ...string) (string, error) {
	var (
		res string
		err error
	)

	c, Ses, e := S.NewSession()
	if e != nil {
		err = fmt.Errorf("Unable to connect: %s", e.Error())
		return res, err
	}
	defer c.Close()
	defer Ses.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if e = Ses.RequestPty("xterm", 80, 40, modes); e != nil {
		err = fmt.Errorf("Unable to start term: %s", e.Error())
		return res, err
	}

	w, _ := Ses.StdinPipe()
	r, _ := Ses.StdoutPipe()

	in, out := TermInOut(w, r)
	if e = Ses.Shell(); e != nil {
		err = fmt.Errorf("Unable to start shell: %s", e.Error())
		return res, err
	}
	<-out

	cmds = append(cmds, "exit")
	cmdtemp := ""

	for _, cmd := range cmds {
		in <- cmd
		outs := strings.Split(<-out, "\n")
		if len(outs) > 1 {
			outtemp := strings.Trim(strings.Join(outs[:len(outs)-1], "\n"), " ")
			res = res + "Output of " + cmdtemp + " : " + outtemp
		}
		cmdtemp = cmd
	}
	Ses.Wait()

	return res, err
}

type RunCommandResult struct {
	CMD    string
	Output string
}

/*
Build connection and run ssh script, catch the output or give error message if any
*/
func (s *SshSetting) RunCommandSshAsMap(cmds ...string) ([]RunCommandResult, error) {
	result := []RunCommandResult{}
	client, sess, err := s.NewSession()
	if err != nil {
		return result, fmt.Errorf("Unable to connect: %s", err.Error())
	}
	defer client.Close()
	defer sess.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err = sess.RequestPty("xterm", 80, 40, modes); err != nil {
		return result, fmt.Errorf("Unable to start term: %s", err.Error())
	}

	writer, _ := sess.StdinPipe()
	reader, _ := sess.StdoutPipe()

	in, out := TermInOut(writer, reader)
	if err = sess.Shell(); err != nil {
		return result, fmt.Errorf("Unable to start shell: %s", err.Error())
	}
	<-out

	cmds = append(cmds, "exit")

	for _, cmd := range cmds {
		in <- cmd
		if cmd != "exit" {
			res := <-out
			res = strings.Split(res, `]0;`)[0]
			res = strings.Split(res, string(0x1b))[0]
			res = strings.TrimSpace(res)
			result = append(result, RunCommandResult{cmd, res})
		}
	}
	sess.Wait()

	return result, nil
}

/*
Run single command, get the output
*/
func (s *SshSetting) GetOutputCommandSsh(cmd string) (string, error) {
	c, Ses, e := s.NewSession()
	if e != nil {
		e = fmt.Errorf("Unable to connect: %s", e.Error())
		return "", e
	}
	defer c.Close()
	defer Ses.Close()

	var out bytes.Buffer
	Ses.Stdout = &out
	var err bytes.Buffer
	Ses.Stderr = &err
	if e := Ses.Run(cmd); e != nil {
		return "", errors.New(fmt.Sprintf("%s. %s", e.Error(), err.String()))
	}

	return out.String(), nil
}

// Copy file adopted from https://github.com/tmc/scp/blob/master/scp.go
func (S *SshSetting) SshCopyByPath(filePath, destinationPath string) error {
	var (
		err error
	)

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return err
	}

	err = S.SshCopyByFile(f, s.Size(), s.Mode().Perm(), filepath.Base(f.Name()), destinationPath)

	return nil
}

func (S *SshSetting) SshCopyByFile(content io.Reader, size int64, perm os.FileMode, filename string, destination string) error {
	var (
		err error
	)

	c, Ses, e := S.NewSession()
	if e != nil {
		err = fmt.Errorf("Unable to connect: %s", e.Error())
		return err
	}
	defer c.Close()
	defer Ses.Close()

	go func() {
		w, _ := Ses.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C%#o %d %s\n", perm, size, filename)
		io.Copy(w, content)
		fmt.Fprint(w, "\x00")
	}()

	cmd := fmt.Sprintf("scp -t %s", destination)

	if err = Ses.Run(cmd); err != nil {
		return err
	}

	return nil
}

func (S *SshSetting) SshGetFile(path string) (res bytes.Buffer, e error) {
	c, Ses, e := S.NewSession()
	if e != nil {
		e = fmt.Errorf("Unable to connect: %s", e.Error())
		return
	}
	defer c.Close()
	defer Ses.Close()

	cmd := fmt.Sprintf(CAT, path)

	Ses.Stdout = &res
	if e = Ses.Run(cmd); e != nil {
		return
	}

	return
}
