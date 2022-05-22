package sshUtil

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"log"
)

func ConnectSSH(addr, username, password, command string) string {
	kexConfig := ssh.Config{}
	kexConfig.SetDefaults()
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         0,
	}
	kexConfig.KeyExchanges = append(kexConfig.KeyExchanges, "diffie-hellman-group1-sha1")
	kexConfig.Ciphers = append(kexConfig.Ciphers, "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc")
	config.Config = kexConfig

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	session, err1 := conn.NewSession()
	if err1 != nil {
		log.Fatal(err1)
	}
	defer session.Close()

	var buff bytes.Buffer
	session.Stdout = &buff
	if err2 := session.Run(command); err2 != nil {
		log.Fatal(err2)
	}
	return buff.String()
}
