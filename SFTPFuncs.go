package main

import "golang.org/x/crypto/ssh"
import "github.com/pkg/sftp"
import "log"
import "os"
import "strings"
import "path/filepath"
import "bufio"

//import "golang.org/x/crypto/ssh/terminal"
//import "golang.org/x/sys/windows"
import "io/ioutil"

//import "crypto"

func GetPrivateKey() ssh.Signer {
	key, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	return signer
}

func HostKeyCheck(HostName string) ssh.PublicKey {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], HostName) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey for %s", HostName)
	}
	return hostKey
}

func GetDataFile(FiletoGet string, SaveAs string, Server string, Username string) bool {
	config := &ssh.ClientConfig{
		User:            Username,
		HostKeyCallback: ssh.FixedHostKey(HostKeyCheck(Server)),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(GetPrivateKey()),
		},
	}

	config.SetDefaults()

	client, err := ssh.Dial("tcp", Server+":22", config)
	if err != nil {
		log.Println("Fail 1")
		log.Print(err)
		return false
	}

	sftp, err := sftp.NewClient(client)
	if err != nil {
		log.Println("Fail 2")
		log.Print(err)
		return false
	}

	srcFile := FiletoGet
	dstFile := SaveAs

	srcGet, err := sftp.Open(srcFile)
	if err != nil {
		log.Println("Fail 3")
		log.Print(err)
		return false
	}

	dstPut, err := os.Create(dstFile)
	if err != nil {
		log.Println("Fail 4")
		log.Print(err)
		return false
	}

	srcGet.WriteTo(dstPut)

	srcGet.Close()
	dstPut.Close()
	sftp.Close()

	return true
}
