package main

//todo , kør den nye proto fil, så du kan få de nye types, ellers får du problemer med at metoderne ikke er de samme senere hen.

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	//this has to be the same as the go.mod module followed by the folder the proto file is in.
	gRPC "github.com/DarkLordOfDeadstiny/DSYS-gRPC-template/proto"

	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedTemplateServiceServer        //You need this line if you have a server
	name                                    string //Not required but useful if you want to name your server
	port                                    string //Not required but useful if your server needs to know what port it's listening to

	incrementValue int64      //value that clients can increment.
	mutex          sync.Mutex //used to lock the server to avoid race conditions.
}

// flags are used to get arguments from the terminal. Flags take a value, a default value and a description of the flag.
// to use a flag then just add it as an argument when running the program.
var serverName = flag.String("name", "default", "Senders name") // set with "-name <name>" in terminal
var port = flag.String("port", "5400", "Server port") // set with "-port <port>" in terminal

func main() {

	// setLog() //uncomment this line to log to a log.txt file instead of the console

	// This parses the flags and sets the correct/given corresponding values.
	flag.Parse()
	fmt.Println(".:server is starting:.")

	//starts a goroutine executing the launchServer method.
	go launchServer()

	//This makes sure that the main method is "kept alive"/keeps running
	for {
		time.Sleep(time.Second * 5)
	}
}

func launchServer() {
	log.Printf("Server %s: Attempts to create listener on port %s\n", *serverName, *port)

	// Create listener tcp on given port or default port 5400
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *port, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}

	//makes gRPC server using the options
	// you can add options here if you want or remove the options part entirely
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// makes a new server instance using the name and port from the flags.
	server := &Server{
		name:           *serverName,
		port:           *port,
		incrementValue: 0, // gives default value, but not sure if it is necessary
	}

	gRPC.RegisterTemplateServiceServer(grpcServer, server) //Registers the server to the gRPC server.

	log.Printf("Server %s: Listening on port %s\n", *serverName, *port)

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
	//code here is unreachable because serve occupies the current thread.
}

// creates a new instance of a server type
func newServer(serverPort *string) *Server {
	s := &Server{
		name:           *serverName, //* is used because of flags
		port:           *serverPort, // not sure why * isnt used here but it isn't
		incrementValue: 0,           // gives default value, but not sure if it is necessary
	}

	fmt.Println(s) //prints the server struct to console
	return s       //return server
}

// The method format can be found in the pb.go file. If the format is wrong, the server type will give an error.
func (s *Server) Increment(ctx context.Context, Amount *gRPC.Amount) (*gRPC.Ack, error) {
	s.mutex.Lock()         //locks the server ensuring no one else can increment the value
	defer s.mutex.Unlock() //unlocks the mutex when exiting the method

	s.incrementValue += int64(Amount.GetValue())      //add value from Amount.ctx
	return &gRPC.Ack{NewValue: s.incrementValue}, nil //create a new acknowlegdement with the new incremented value and returns it.
}

// sets the logger to use a log.txt file instead of the console
func setLog() {
	//Clears the log.txt file when a new server is started
	if err := os.Truncate("log.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	//This connects to the log file/changes the output of the log informaiton to the log.txt file.
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

}
