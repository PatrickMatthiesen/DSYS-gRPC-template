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

	gRPC "github.com/DarkLordOfDeadstiny/DSYS-gRPC-template/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Same principle as in client. Flags allows for user specific arguments/values
var clientsName = flag.String("name", "default", "Senders name")
var serverPort = flag.String("server", "5400", "Tcp server")

var ctx context.Context                 //Client context
var server gRPC.TemplateServiceClient   //the server
var ServerConn *grpc.ClientConn 		//the server connection

func main() {
	//parse flag/arguments
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	//log to file instead of console
	//setLog()

	fmt.Println("--- join Server ---")
	joinServer()
	defer ServerConn.Close() //when main method exits, close the connection to the server.

	//start the biding
	parseInput()
}

func joinServer() {
	//connect to server

	//dial options
	//the server is not using TLS, so we use insecure credentials
	//(should be fine for local testing but not in the real world)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())) 

	//use context for timeout on the connection
	timeContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() //cancel the connection when we are done

	log.Printf("client %s: Attempts to dial on port %s\n", *clientsName, *serverPort)
	conn, err := grpc.DialContext(timeContext, fmt.Sprintf(":%s", *serverPort), opts...) //dials the port with the given timeout
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		return
	}
	server = gRPC.NewTemplateServiceClient(conn) //create a new gRPC client
	ServerConn = conn                  			 // saves the MessageServiceClient's to connection
	fmt.Println(conn.GetState().String()) //prints connected if it's connected (i think (☞ﾟヮﾟ)☞)

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
	//Convert string to int64, return error if the int is larger than 32bit
	val, err := strconv.ParseInt(in, 10, 32) 
	if err != nil {
		log.Fatal(err)
	}

	//create amount type
	amount := &gRPC.Amount{
		ClientName: *clientsName,
		Value:      int32(val), //cast from int64 to int32
	}

	if conReady(server) { //If the connection to the server is ready
		ack, err := server.Increment(ctx, amount) //Make gRPC call to server with amount, and recieve acknowlegdement back.
		if err != nil {
			log.Printf("Client %s: no response from the server, attempting to reconnect", *clientsName)
			log.Println(err)
		}
		if ack.NewValue >= val {
			fmt.Printf("Success, the new value is now %d\n", ack.NewValue)
		} else {

			fmt.Println("Oh no something went wrong :(") //Hopefully this will never be reached
		}
	} else {
		log.Printf("Client %s: something was wrong with the connection to the server :(", *clientsName)
	}
}

// Function which returns a true boolean if the connection to the server is ready, and false if it's not.
func conReady(s gRPC.TemplateServiceClient) bool {
	return ServerConn.GetState().String() == "READY"
}

// sets the logger to use a log.txt file instead of the console
func setLog() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

}
