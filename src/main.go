package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/Southern/logger"
	"os/exec"
	"regexp"
	"strconv"
	s "strings"
)

var (
	Log            = logger.New().Log
	MatchUser      = regexp.MustCompile(`^:([A-Za-z0-9\-_]+)!`)
	URLMatch       = regexp.MustCompile(`\b(https?://)?(([0-9a-zA-Z_!~*'().&=+$%-]+:)?[0-9a-zA-Z_!~*'().&=+$%-]+\@)?(([0-9]{1,3}\.){3}[0-9]{1,3}|([0-9a-zA-Z_!~*'()-]+\.)*([0-9a-zA-Z][0-9a-zA-Z-]{0,61})?[0-9a-zA-Z]\.[a-zA-Z]{2,6})(:[0-9]{1,4})?((/[0-9a-zA-Z_!~*'().;?:\@&=+$,%#-]+)*/?)`)
	Config         tls.Config
	CommandChar    = '>'
	BotOwner       = "Mantis!Mantis@Cybershade.org"
	SQLMapLocation = "/home/rclifford/Downloads/sqlmap/"
)

func sendMessage(message string, location string, sockfd *tls.Conn) {
	fmt.Fprintf(sockfd, "PRIVMSG %s :%s\r\n", location, message)
}

func sendRaw(message string, sockfd *tls.Conn) {
	fmt.Fprintf(sockfd, "%s\r\n", message)
}

func getData(str string, option string) string {
	split := s.SplitN(str, " ", 4)

	if len(split) >= 4 {
		switch option {
		case "user":
			return s.Trim(s.Split(split[0], "!")[0], ":")

		case "userIdent":
			return s.TrimLeft(split[0], ":")
		case "channel":
			return s.Trim(s.Split(split[2], "!")[0], ":")

		case "message":
			ret := s.TrimLeft(split[3], ":")
			return ret
		}
	}
	return ""
}

func authCheck(status string, sockfd *tls.Conn) bool {
	user := getData(status, "userIdent")

	if user == BotOwner {

		Log("BotOwner Confirmed.")

		return true
	}

	Log("w", "Not valid bot owner.")

	sendMessage("Sorry, you do not have the correct privileges", getData(status, "channel"), sockfd)

	return false
}

func executeCommands(status string, sockfd *tls.Conn) {
	sender := getData(status, "user")
	channel := getData(status, "channel")
	message := getData(status, "message")

	if s.HasPrefix(message, string(CommandChar)) {

		if authCheck(status, sockfd) {

			args := s.Split(message, " ")

			if s.HasPrefix(message, fmt.Sprintf("%cchangeNick", CommandChar)) {
				if len(args) >= 2 {
					sendRaw(fmt.Sprintf("NICK %s", args[1]), sockfd)
				}
			}

			if s.HasPrefix(message, fmt.Sprintf("%cjoin", CommandChar)) {
				if len(args) >= 2 {
					sendRaw(fmt.Sprintf("JOIN %s", args[1]), sockfd)
				}
			}

			if s.HasPrefix(message, fmt.Sprintf("%cpart", CommandChar)) {
				if len(args) >= 2 {
					sendRaw(fmt.Sprintf("PART %s", args[1]), sockfd)
				}
			}


// Invite
//			info: :Mantis!Mantis@Cybershade.org INVITE Goo :#treehouse

			// if s.HasPrefix(message, fmt.Sprintf("%cping", CommandChar)) {
			// 	if len(args) >= 2 {

			// 		cmd, err := exec.Command("ping", args[1], "-c3").Output()

			// 		fmt.Println(cmd)

			// 		if err != nil {

			// 			Log("w", fmt.Sprintf("%s", err))

			// 			sendMessage(fmt.Sprintf("Sorry %s, something went wrong", sender), channel, sockfd)

			// 		} else {

			// 			fmt.Println(cmd)

			// 			sendMessage(fmt.Sprintf("The server is up and responding to pings, %s", sender), channel, sockfd)

			// 		}
			// 	}
			// }

			/**
			//
			// -- Need to make this method concurrent
			//
			*/
			if s.HasPrefix(message, fmt.Sprintf("%csqlMap", CommandChar)) {

				if len(args) >= 2 {

					out, err := exec.Command("/usr/bin/python", fmt.Sprintf("%ssqlmap.py", SQLMapLocation), fmt.Sprintf("-u \"%s\"", args[1]), "--random-agent", "--threads=2", "--batch").Output()

					fmt.Println(string(out))

					if err != nil {

						sendMessage(fmt.Sprintf("Sorry %v, this could not be executed.", sender), channel, sockfd)

					} else {

						sendMessage(strconv.QuoteToASCII(string(out)), sender, sockfd)

					}
				}
			}
		}
	}
}

func main() {

	nick := flag.String("nick", "Goo", "The Nickname for the bot")

	server := flag.String("server", "irc.darkscience.net", "Server for the bot to connect to")

	port := flag.Int("port", 6697, "The port of the server")

	flag.Parse()

	sockfd, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", *server, *port), &Config)

	if err != nil {
		Log("w", fmt.Sprintf("%s", err))
	}

	i := 0

	for {

		status, err := bufio.NewReader(sockfd).ReadString('\n')

		if err != nil {
			Log("w", fmt.Sprintf("%s", err))
		}

		Log(status)

		switch i {

		case 0:
			fmt.Fprintf(sockfd, "USER guest 0 * :DarkMantisBOT\r\n")
			fmt.Fprintf(sockfd, "NICK %s\r\n", *nick)
			break

		case 5:
			fmt.Fprintf(sockfd, "JOIN #bots\r\n")
			// fmt.Fprintf(sockfd, "JOIN #golang\r\n")
			// fmt.Fprintf(sockfd, "JOIN #treehouse\r\n")
			break
		}

		// Send pong
		if s.Contains(status, "PING") {
			pong := s.Replace(status, "I", "O", 1)

			Log(pong)

			fmt.Fprintf(sockfd, pong)
		}

		executeCommands(status, sockfd)

		i++
	}
}
