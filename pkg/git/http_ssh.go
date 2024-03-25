package git

import (
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/gvc/conf"
)

/*
Use http proxy for ssh.
*/
const (
	bufsize = 4096
)

// Returns the command that executes the current program
func GetParentCommand() string {
	ppid := os.Getppid()
	fileName := "/proc/" + strconv.Itoa(ppid) + "/cmdline"
	dat, err := os.ReadFile(fileName)
	if err != nil {
		gprint.PrintError("Can't read file '%s'\n", fileName)
		return ""
	}
	output := ""
	for i, elem := range dat {
		if elem == byte(0) {
			output += " "
		} else {
			output += string(dat[i])
		}
	}
	return output
}

// Parse COMMAND and returns in format 'ssh://'
func GetRepositoryFromCommand(command string) (string, string) {
	var (
		destination     string
		port            string = "22"
		repository      string
		gitShellCommand string
		arg             []string = strings.Split(command, " ")
	)
	/// "usr/bin/ssh git@github.com git-upload-pack 'bitexploder/timmy.git'"
	// "usr/bin/ssh -p 7999 git@globaldevtools.bbva.com git-receive-pack '/uqnwi/bitbucket_lifecycle.git'"

	for i, s := range arg {
		if strings.HasPrefix(s, "-p") && len(s) == 2 {
			port = arg[i+1]
		} else if strings.HasPrefix(s, "-p") && len(s) > 2 {
			port = s[2:]
		}
		if strings.HasPrefix(s, "git@") {
			destination = s
		}
		if strings.HasPrefix(s, "'") && strings.HasSuffix(s, ".git'") {
			repository = s[1 : len(s)-1]
		} else if strings.HasPrefix(s, "/") && strings.HasSuffix(s, ".git") {
			repository = s
		}
		if strings.HasPrefix(s, "git-") {
			gitShellCommand = s
		}
	}
	if !(strings.HasPrefix(repository, "/")) {
		repository = "/" + repository
	}
	if len(destination) > 0 && len(port) > 0 && len(repository) > 0 {
		return "ssh://" + destination + ":" + port + repository, gitShellCommand
	}
	return "", ""
}

// returns URI for proxy connection when there's no authentication
func GetURINoAuth(desthost, destport string) string {
	var (
		parentComand    string
		repository      string
		gitShellCommand string
	)
	parentComand = GetParentCommand()
	// log.Print("Parent command:", parentComand)
	repository, gitShellCommand = GetRepositoryFromCommand(parentComand)
	// log.Print("Repository:", repository)
	// log.Print("Git Shell Command: ", gitShellCommand)
	return "CONNECT " +
		desthost + ":" + destport + " HTTP/1.0\nHost: " +
		desthost + ":" + destport + "\nRepository: " +
		repository + "\nGitShellCommand: " +
		gitShellCommand + "\nUserAuth: NULL" + "\r\n\r\n"
}

// returns socket connection with proxyaddr (string with format `host:port`)
func CreateNetSocket(host string, port string, timeout int) net.Conn {
	proxyaddr := host + ":" + port
	conn, err := net.DialTimeout("tcp", proxyaddr, time.Duration(timeout)*time.Second)
	if err != nil {
		gprint.PrintError("No connection to Proxy '%s:%s'\n", host, port)
		os.Exit(1)
	}
	return conn
}

/*
use http proxy for ssh.
*/
func GrokscrewHttpSSH(destHost, destPort string, proxyTimeout ...int) {
	cfg := conf.NewGVConfig()
	proxyURI := cfg.GetLocalProxy()

	authURI := GetURINoAuth(destHost, destPort)
	var (
		buffer = make([]byte, bufsize)
		read   int
		write  int
		setup  int = 0
	)
	if u, err := url.Parse(proxyURI); err == nil {
		proxyhost := u.Hostname()
		proxyport := u.Port()
		pTimeout := 5
		if len(proxyTimeout) > 0 && proxyTimeout[0] > 0 {
			pTimeout = proxyTimeout[0]
		}
		conn := CreateNetSocket(proxyhost, proxyport, pTimeout)
		defer conn.Close()
		for {
			if setup == 0 {
				write, _ = conn.Write([]byte(authURI))
				if write <= 0 {
					break
				}
				read, _ = conn.Read(buffer)
				if read <= 0 {
					break
				}
				statusStr := strings.Split(string(buffer[:]), " ")[1]
				statusCode, _ := strconv.Atoi(statusStr)
				if statusCode >= 200 && statusCode < 300 {
					gprint.PrintInfo("Connection stablished. STATUS CODE: %d\n", statusCode)
					setup = 1
				} else if statusCode >= 407 {
					gprint.PrintError("Proxy could not open connection. STATUS CODE: %d\n", statusCode)
					os.Exit(1)
				}
			} else {
				FeelTheMagic(conn)
			}
		}
	} else {
		gprint.PrintError("Invalid proxy URI: %s", proxyURI)
	}
}

type Progress struct {
	bytes uint64
}

// this function redirects data between socket, stdin and stdout
func FeelTheMagic(con net.Conn) {
	c := make(chan Progress)

	// Read from Reader and write to Writer until EOF
	copy := func(r io.ReadCloser, w io.WriteCloser) {
		defer func() {
			r.Close()
			w.Close()
		}()
		n, _ := io.Copy(w, r)
		c <- Progress{bytes: uint64(n)}
	}

	go copy(con, os.Stdout)
	go copy(os.Stdin, con)

	p := <-c
	gprint.PrintInfo("[%s]: Connection has been closed by remote peer, %d bytes has been received\n", con.RemoteAddr(), p.bytes)
	p = <-c
	gprint.PrintInfo("[%s]: Local peer has been stopped, %d bytes has been sent\n", con.RemoteAddr(), p.bytes)
}
