package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	gRPC "github.com/DarkLordOfDeadstiny/Exam-template/proto"

	"google.golang.org/grpc"
)

//Same principle as in client. Flags allows for user specific arguments/values
var clientsName = flag.String("name", "default", "Senders name")
var tcpServer = flag.String("server", "5400", "Tcp server")

var _ports [5]string = [5]string{*tcpServer, "5401", "5402", "5403", "5404"} //List of ports the client tries to connect to._ports

var ctx context.Context                                       //Client context
var servers []gRPC.MessageServiceClient                       //list of servers.
var ServerConn map[gRPC.MessageServiceClient]*grpc.ClientConn //Map of server connections

func main() {
	//parse flag/arguments
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	//connect to log file
	setLog()

	fmt.Println("--- join Server ---")
	ServerConn = make(map[gRPC.MessageServiceClient]*grpc.ClientConn) //make map of server connections
	joinServer()                                                      // Method call
	defer closeAll()                                                  //when main method exits, close all the connections to the servers.

	//start the biding
	parseInput()
}

func joinServer() {
	//connect to server

	//dial options
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	//use context for timeout on the connection
	timeContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() //cancel the connection when we are done

	for _, port := range _ports { //try to connect to all the ports in _ports with dial
		log.Printf("client %s: Attempts to dial on port %s\n", *clientsName, port)
		conn, err := grpc.DialContext(timeContext, fmt.Sprintf(":%s", port), opts...) //dials the port with the given timeout
		if err != nil {
			log.Printf("Fail to Dial : %v", err)
			continue
		}
		var s = gRPC.NewMessageServiceClient(conn)
		servers = append(servers, s)          //add the new MessageServiceClient
		ServerConn[s] = conn                  // maps the MessageServiceClient to its connection
		fmt.Println(conn.GetState().String()) //prints connected if it's connected (i think (☞ﾟヮﾟ)☞)
	}
	ctx = context.Background()
}

func parseInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type the amount you wish to increment with here. Type 0 to get the current value")
	fmt.Println("--------------------")

	//Infinite loop to listen for clients input.
	for {
		fmt.Print("-> ")

		in, err := reader.ReadString('\n') //Read input into var in
		if err != nil {
			log.Fatal(err)
		}
		in = strings.TrimSpace(in) //Trim input
		incrementVal(in)
	}
}

func incrementVal(in string) {

	val, err := strconv.ParseInt(in, 10, 32) //Convert string to int64, return error if the int is larger than 32bit
	if err != nil {
		log.Fatal(err)
	}

	//create amount type
	amount := &gRPC.Amount{
		ClientName: *clientsName,
		Value:      int32(val), //cast from int64 to int32
	}
	for _, s := range servers {
		if conReady(s) { //If the connection to the server is ready
			fmt.Println(s)
			ack, err := s.Increment(ctx, amount) //Make gRPC call to server with amount, and recieve acknowlegdement back.
			if err != nil {
				log.Printf("Client %s: no response from the server, attempting to reconnect", *clientsName)
				log.Println(err)
			}
			if ack.NewValue >= val {
				fmt.Printf("Success, the new value is now %d\n", ack.NewValue)
			} else {

				fmt.Println("Oh no something went wrong :(") //Hopefully this will never be reached
			}
		}
	}
}

//Function which returns a true boolean if the connection to the server is ready, and false if it's not.
func conReady(s gRPC.MessageServiceClient) bool {
	return ServerConn[s].GetState().String() == "READY"
}

func closeAll() {
	for _, c := range ServerConn {
		c.Close()
	}
}

func setLog() {

	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

}
