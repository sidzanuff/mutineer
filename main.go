package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/namsral/flag"
)

var conn net.Conn
var debug bool

func main() {
	QUESTION_PROMPT := string([]byte{58, 32})
	MENU_PROMPT := "Main Menu] > "
	FILES_PROMPT := string([]byte{47, 27, 91, 51, 55, 59, 49, 109, 62, 32, 27, 91, 51, 55, 109})
	EDIT_PROMPT := string([]byte{32, 109, 101, 33, 13, 10, 13, 10, 13, 10, 27, 91, 51, 50, 59, 49, 109})

	var user, pass, src, dst, desc string
	flag.StringVar(&user, "user", "", "")
	flag.StringVar(&pass, "pass", "", "")
	flag.StringVar(&src, "src", "", "")
	flag.StringVar(&dst, "dst", "upload.bas", "")
	flag.StringVar(&desc, "desc", "", "")
	flag.BoolVar(&debug, "debug", false, "")
	flag.Parse()

	conn, _ = net.Dial("tcp", "mutinybbs.com:2300")
	defer conn.Close()

	wait(QUESTION_PROMPT)
	write("\b")
	wait(QUESTION_PROMPT)
	write("n\n")
	wait(QUESTION_PROMPT)
	send(user)
	wait(QUESTION_PROMPT)
	send(pass)
	wait("Continue > ")
	write("\n")

	wait(MENU_PROMPT)
	write("t")
	wait(FILES_PROMPT)
	send("cd users/" + user)
	wait(FILES_PROMPT)
	send("edit " + dst)
	wait(EDIT_PROMPT, "loaded\r\n")
	sendwait("new")

	file, _ := os.Open(src)
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

	wait(MENU_PROMPT)
	write("o")
	wait(QUESTION_PROMPT)
	write("y")
}

func wait(prompts ...string) {
	if debug {
		fmt.Print("Waiting for: ")
		for _, prompt := range prompts {
			chars(prompt)
		}
	}
	buf := make([]byte, 256)
	for {
		n, _ := conn.Read(buf)
		if n == 0 {
			break
		}
		line := string(buf[:n])
		fmt.Println(line)
		if debug {
			fmt.Print("Got: ")
			chars(line)
		}
		for _, prompt := range prompts {
			if strings.HasSuffix(line, prompt) {
				return
			}
		}
	}
}

func write(text string) {
	conn.Write([]byte(text))
}

func send(text string) {
	write(text + "\r\n")
}

func sendwait(text string) {
	text = text + "\r\n"
	write(text)
	wait(text)
}

func chars(text string) {
	for _, c := range text {
		fmt.Printf("%d,", c)
	}
	fmt.Println()
}
