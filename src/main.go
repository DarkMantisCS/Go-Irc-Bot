package main

import (
    "fmt"
    "log"
    "crypto/tls"
    "bufio"
    s "strings"
)

type IrcData struct {
    commandChar string
    currentNick string
}

func authCheck(auth string) bool {
    authStr := s.Split(auth, " ")
    if len(authStr) > 1 && authStr[0] == ":Mantis!uid3257@Cybershade.org" {
        return true
    }

    return false
}

func sendPong(ping string, sockfd *tls.Conn) {
    pong := s.Replace(ping, "I", "O", 1)
    fmt.Println(pong)
    fmt.Fprintf(sockfd, pong)
}

func getData(str string, option string) string {
    split := s.Split(str, " ")

    if len(split) >= 4{
        switch(option){
            case "user":
                return s.Trim(s.Split(split[0], "!")[0], ":")

            case "channel":
                return s.Trim(s.Split(split[2], "!")[0], ":")

        }
    }
    return ""
}


func executeCommands(status string, sockfd *tls.Conn) {
    sender  := getData(status, "user")
    channel := getData(status, "channel")

    if s.Contains(status, ">changeNick") {

        newNick := s.Split(status, " ")

        if authCheck(status) && len(newNick) >= 5 {
            if newNick[4] != ""{
                fmt.Fprintf(sockfd, "NICK " + newNick[4])
            }
        } else {
            fmt.Fprintf(sockfd, "PRIVMSG " + channel + " :Sorry " + sender + " you do not have enough Kudos to do this\r\n")
        }
    }

    // Send pong
    if s.Contains(status, "PING"){
        sendPong(status, sockfd);
    }
}


func main() {
    var config tls.Config
    var IrcData IrcData


    IrcData.currentNick = "Goo"
    IrcData.commandChar = ">"

    sockfd, err := tls.Dial("tcp", "irc.darkscience.net:6697", &config)

    if err != nil {
        log.Fatal(err)
    }

    i := 0

    for {

        status, err := bufio.NewReader(sockfd).ReadString('\n')

        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(status)

        switch(i) {
            case 0:
                fmt.Fprintf(sockfd, "USER guest 0 * :DarkMantisBOT\r\n")
                fmt.Fprintf(sockfd, "NICK Goo\r\n")
                break;

            case 5:
                fmt.Fprintf(sockfd, "JOIN #bots\r\n")
                fmt.Fprintf(sockfd, "JOIN #golang\r\n")
                break;
        }

        executeCommands(status, sockfd)

        i++
    }
}
