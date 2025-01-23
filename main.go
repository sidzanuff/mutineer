package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/namsral/flag"
)

var c net.Conn
var d bool

func main() {
	QUESTION_PROMPT := string([]byte{58, 32})
	FILES_PROMPT := string([]byte{47, 27, 91, 51, 55, 59, 49, 109, 62, 32, 27, 91, 51, 55, 109})
	EDIT_PROMPT := string([]byte{32, 109, 101, 33, 13, 10, 13, 10, 13, 10, 27, 91, 51, 50, 59, 49, 109})

	var user, pass, path, name, desc string
	flag.StringVar(&user, "user", "", "")
	flag.StringVar(&pass, "pass", "", "")
	flag.StringVar(&path, "src", "", "")
	flag.StringVar(&name, "dst", "upload.bas", "")
	flag.StringVar(&desc, "desc", "", "")
	flag.BoolVar(&d, "debug", false, "")
	flag.Parse()

	conn, _ := net.Dial("tcp", "mutinybbs.com:2300")
	c = conn
	defer conn.Close()

	wait("press backspace/delete: ")
	write("\b")
	wait("oes your terminal support ANSI color? (Y)es, (N)o or don't know: ")
	write("n\n")
	wait("who are you?: ")
	send(user)
	wait("password: ")
	send(pass)
	wait("Continue > ")
	write("\n")

	wait("Main Menu] > ")
	write("t")
	wait(FILES_PROMPT)
	send("cd users/hexen")
	wait(FILES_PROMPT)
	send("edit " + name)
	wait(EDIT_PROMPT, "loaded\r\n")
	sendwait("new")

	file, _ := os.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sendwait(scanner.Text())
	}

	send("save")
	wait(QUESTION_PROMPT)
	send(desc)
	wait("program saved\r\n")
	send("quit")
	wait(FILES_PROMPT)
	send("quit")

	wait("Main Menu] > ")
	write("o")
	wait("Logoff? (Y)es, (W)ith Message, (N)o: ")
	write("y")
}

func wait(s ...string) {
	if d {
		fmt.Print(("Waiting for: "))
		for _, p := range s {
			chars(p)
		}
	}
	buf := make([]byte, 256)
	for {
		n, _ := c.Read(buf)
		if n == 0 {
			break
		}
		l := string(buf[:n])
		fmt.Println(l)
		if d {
			fmt.Print("Got: ")
			chars(l)
		}
		for _, p := range s {
			if strings.HasSuffix(l, p) {
				return
			}
		}
	}
}

func write(s string) {
	c.Write([]byte(s))
}

func send(s string) {
	write(s + "\r\n")
}

func sendwait(s string) {
	s = s + "\r\n"
	write(s)
	wait(s)
}

func chars(s string) {
	for _, c := range s {
		fmt.Printf("%d,", c)
	}
	fmt.Println()
}
