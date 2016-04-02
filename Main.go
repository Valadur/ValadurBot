package main

import (
	"net"
	"bufio"
	"log"
	"time"
	"strings"
	"github.com/valadur/ValadurBot/reference"
)

func main() {
	serverAndPort := "irc.twitch.tv:6667"

	conn, err := net.Dial("tcp", serverAndPort)
	defer conn.Close()
	if (err != nil) {
		log.Fatal("Couldn't establish connection to server and port " + serverAndPort)
	}
	log.Println("Established connection to " + serverAndPort)

	connectBot(conn, reference.Psw, reference.Nickname, reference.Username)
	log.Printf("Starting to listen in 15 seconds..")
	time.Sleep(15 * time.Second)
	joinChannel(conn, "wingsofdeath")
	go listenToConn(conn)
	for true{
		time.Sleep(5 * time.Minute)
	}
}

func connectBot(conn net.Conn, psw string, nickname string, username string){
	sendMessage(conn, "PASS " + psw)
	sendMessage(conn, "NICK " + nickname)
}

func sendMessage(conn net.Conn, msg string)  {
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(msg + "\r\n")
	writer.Flush()
	if err != nil {
		log.Fatal("An error occurred during writing to connection")
	} else {
		log.Printf("out: " + msg)
	}
}

func listenToConn(conn net.Conn)  {
	reader := bufio.NewReader(conn)
	for(true) {
		time.Sleep(50 * time.Millisecond)
		msg, err := reader.ReadString('\n')
		if (err != nil) {
			log.Fatal("An error ocurred during reading from connection")
		}
		msgArray := strings.Split(msg, " ")
		if (msgArray[0] == "PING") {
			sendMessage(conn, "PONG " + msgArray[1])
		}
		log.Printf(getNickFromRawMessage(msg) + ": " + getMessageFromRawMessage(msg))
	}
}
func joinChannel(conn net.Conn, channel string)  {
	sendMessage(conn, "JOIN #" + channel)
}

func sendMessageToChannel(conn net.Conn, sendTo string, msg string)  {
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString("PRIVMSG " + sendTo + " :" + msg + "\r\n")
	writer.Flush()
	if err != nil {
		log.Fatal("An error occurred during writing to connection")
	} else {
		log.Printf("out: " + msg)
	}
}

func getNickFromRawMessage(msg string) (string) {
	returnMsg := strings.TrimPrefix(msg, ":")
	returnMsgArray := strings.Split(returnMsg, "!")
	if(len(returnMsgArray) > 0){
		returnMsg = returnMsgArray[0]
	}
	return returnMsg
}

func getMessageFromRawMessage(msg string) (string) {
	returnMsg := strings.TrimPrefix(msg, ":")
	returnMsgArray := strings.Split(returnMsg, ":")
	if(len(returnMsgArray) > 1){
		return returnMsgArray[1]
	} else {
		return returnMsgArray[0]
	}
}