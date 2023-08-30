# gRPC guide

> **Disclaimers**
>
> This is NOT to say how you NEED to do it, but a guide that can show you how YOU might WANT to do it.
>
> A **simpler** and older version of this guide can be found at:
> <https://github.com/NaddiNadja/grpc101>
>
> This guide is not a step by step guide and it might be a bit overcomplicated for some people. I hoped that it would be better to give more than too little.
> The guide may hold text like `<add a name here>`, this should be understood as take it all and replace it with a name so `<add a name here>` would become something like `Name`.
>
> Example:
> `rpc <Method name> (<Message name>) returns (<Response name>);`
> Becomes `rpc SayHi (Greeting) returns (Farewell);`

- [Intro](#intro)
- [Setup of a new repository](#setup-of-a-new-repository)
- [The Proto file](#the-proto-file)
  - [What is it?](#what-is-it)
  - [Required lines](#required-lines)
  - [Defining a service](#defining-a-service)
- [Implementation](#implementation)
  - [Implementing the server methods](#implementing-the-server-methods)
  - [Calling endpoints from the client](#calling-endpoints-from-the-client)
- [Prerequisites](#prerequisites)
  - [Download protoc on Windows](#download-protoc-on-windows)
  - [Download protoc on Mac OS](#download-protoc-on-mac-os)
  - [MY recommended VSCode extensions](#my-recommended-vscode-extensions)

If you haven't installed google's protocol buffers, see the prerequisites part at the bottom.

## Intro

>I will complete this at a later time (TODO)

This guide assumes that you at least have downloaded go.

You can start by following the [Setup of new repository](#setup-of-a-new-repository), which should show you how to start making your own repository structure and how to run the different things. While setting up your repository, make sure to go to the appropriate sections explaining the basics of what to do in each of the files. If the sections don't explain it well enough, then try to compare it to the working code example in this repository. I would very much appreciate feedback if you have any 🙂.

## Setup of a new repository

1. Make ``go.mod`` file with:

    ```sh
    go mod init [link to repo without "https://"]
    ```

    Your repo should be on the public github as it needs to be an accessable web page.
    If you want to keep it simple, but differ from the standard, then you can just call it something short like "Incrementer", as that would fit the name of this example.

2. Make a ``.proto`` file in a sub-directory, for example ``proto/template.proto`` and make your service.
    > See [The Proto file](#the-proto-file) for info on what to add in the ``.proto`` file
3. Run command:

    ```sh
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/template.proto
    ```

    The command should create the **two** ``.pb.go`` files. Remember to change "proto/template" to your directory and file name.
4. Run command:

    ```sh
    go mod tidy
    ```

    to install dependencies and create the ``go.sum`` file.
5. Implement your client and server. Refer to [Implementation](#implementation) for instructions.
6. Open a terminal for each client and server and run them with:

    The Client: `$ go run .\client\client.go`

    The Server: `$ go run .\server\server.go`

## The Proto file

### What is it?

The proto file is a file that is used to compile methods for the actual gRPC. So if you make changes in the proto file, then you will have to recompile it, before you can see any changes in your code.

### Required lines

1. To make sure to use the newer proto syntax, you will need to add the syntax as the first line in the file.

   ```protobuf
    syntax = "proto3";
   ```

2. A path for the compiler to know what the package name of the pb.go files needs to have, then an option go_package will need to be set. It is gonna use the last part of the path as the name of the package, but it matters little, it just needs to be unique. The example below will make the package name "proto"

    ```protobuf
    option go_package = "github.com/PatrickMatthiesen/DSYS-gRPC-template/proto";
    ```

3. Add a package name that is the same as the last folder from above ("proto" in this example)

    ```protobuf
    package proto;
    ```

### Defining a service

1. Make a service with a telling name. Here the name of the service is `Template`, It's a bad name but it is a name.

    ```protobuf
    service Template
    {
    }
    ```

2. Add RPC endpoints/methods to the service. The endpoint will start with `rpc`, followed by the name of the endpoint, a message type to send, and a message type to receive. We have not made the message types yet, but we will do that next.

   Here we make an endpoint called `SayHi`, which sends a message type called `Greeting` and expects to receive a response of a type called `Farewell`.

    ```protobuf
    service Template
    {
        rpc SayHi (Greeting) returns (Farewell);
    }
    ```

    An endpoint can also be streamed, so more than one message or response can be sent. You can add either the message or response, or both if you need. Here is an example that does it to both:

     ```protobuf
        rpc SayHi (stream Greeting) returns (stream Farewell);
    ```

3. To complete the endpoint we need to make the message types.

    ```protobuf
    message Greeting {
        string clientName = 1;
        string message = 2;
        ...
    }

    message Farewell {
        string message = 1;
    }
    ```

   - notes:
     - All names of services, endpoints and message types should be in CamelCase
     - an example can be found in [template.proto](proto/template.proto)

4. Now that we have defined our service we want to compile the go files. We do so by running the following command:

   ```sh
   protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/template.proto
   ```

   > remember to change the path `proto/template.proto` to the path of your proto file.

## Implementation

For this section we go over how to implement the server and client parts of the service you created in [Defining a service](#defining-a-service). This guide shows how to make a client-server architecture, but this guide will not show how to do a client to client architecture.

### Implementing the server methods

1. Make a "server" struct
   - It does not need to be called server, but its just what makes the most sense in this case.

    ```go
    type Server struct {
        // an interface that the server needs to have
        gRPC.UnimplementedTemplateServer
        
        // here you can impliment other fields that you want
    }
    ```

2. Make a method that matches the endpoint in your proto file and add `(s *Server)`. This attaches the method to the server struct so it can run when someone calls our endpoint.
   - For an endpoint that does no streaming, then we need to give the method a context and the input type. For the return we need to return a pair of your return type and an error.

    ```go
    func (s *Server) <endpoint name>(ctx context.Context, <name> *<input type>) (*<the return type>, error) {
        //some code here
        ...
        ack :=  // make an instance of your return type
        return (ack, nil)
    }
    ```

    - For an endpoint that streams messages, we need to give the method a stream and return an error. In this case, you get the input from the stream and send the return type back over the stream too.

    ```go
    func (s *Server) <endpoint name>(msgStream gRPC.<service name>_<endpoint name>Server) error {
        for {
            // get the next message from the stream
            msg, err := msgStream.Recv()
            if err == io.EOF {
                break
            }
        }

        ack := // make an instance of your return type
        msgStream.SendAndClose(ack)

        return nil
    }
    ```

3. Serving the endpoint

   - Make a main method and add the following

   - Listen on a port. If you add localhost then it is only running locally (on your computer), if you remove it, then it will open the port on your computer to others. If the port is opened then you should get a firewall prompt that you should approve.

    ```go
    list, _ := net.Listen("tcp", "localhost:5400")
    ```

   - Make a grpc server

    ```go
    grpcServer := grpc.NewServer(opts...)
    ```

   - Make an instance of your struct

    ```go
    server := &Server{
        // set fields here 
    }
    ```

    - Register your server struct. Remember that you need to change "Template" to what your service is called.

    ```go
    gRPC.RegisterTemplateServer(grpcServer, server)
    ```

    - Start serving

    ```go
    grpcServer.Serve(list)
    ```

### Calling endpoints from the client

1. Make dialing options with insecure credentials, as we don't have something called a TLS certificate (which is used for encryption).

    ```go
    var opts []grpc.DialOption
    opts = append(
        opts, grpc.WithBlock(), 
        grpc.WithTransportCredentials(insecure.NewCredentials())
    )

    timeContext, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    ```

2. Dial the server. Here we just go the port locally, but if you want to connect to another device, then you would just add the IP of the other device. (use ``ipconfig`` in the terminal or call ``GetOutboundIP()`` which can be found in [server](/server/server.go))

    ```go
    conn, err := grpc.Dial(":5400", opts...)
    ```

3. Make a client that has the endpoint methods

    ```go
    server = gRPC.NewTemplateClient(conn)
    ```

4. Now you can just say `server.<endpoint name>` to call the endpoint

5. For streamed endpoints.
   - You call the method with a context, to get the stream.

    ```go
    stream, err := server.SayHi(context.Background())
    ```

   - Then you can call send on the stream with the message type of the endpoint stream. Could be done in a loop or however you want.

    ```go
    message := // make an instance of your input ty
    stream.Send(message)
    ```

    - When done sending messages, you need to close the stream and wait for a response from the server.

    ```go
    farewell, err := stream.CloseAndRecv()
    ```

## Prerequisites

### Download protoc on Windows

1. Before starting, install Google's protocol buffers:
    - Go to the latest release of the protobuf repository: <https://github.com/protocolbuffers/protobuf/releases/latest>
    - Find the version you need and download it. 
    - As of July 2023, if you are on Windows, it's the third from the bottom, `protoc-21.7-win64.zip`.
2. Unzip the file and place it in a folder you won't move or delete.
    - On my Windows machine, I placed it in `C:\Users\<username>\go`
    - I chose to rename the folder to `Protoc`, so the path of the folder is `C:\Users\<username>\go\Protoc`
3. Add the path of the `bin` folder to your system variables.
    - On Windows, press the Windows key and search for `edit the system environment variables`.
    - Click on `Environment Variables...` at the bottom.
    - In the bottom list select the variable called `Path` and click on `Edit ...`
    - In the pop-up window, click on `New`
    - Paste the full path of the `bin` folder into the text field.

        My path is `C:\Users\<username>\go\Protoc\bin`.

    - Click `OK` twice.
4. Open a terminal and run the following commands:

    `$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

    `$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

### Download protoc on Mac OS

`$ brew install go`

`$ brew install protoc-gen-go`

`$ brew install protobuf`

`$ brew install protoc-gen-go-grpc`

### MY recommended VSCode extensions

It can be nice with some colors so here are my favorites.

- Go language support
    > vscode-proto3: <https://marketplace.visualstudio.com/items?itemName=golang.Go>

- Proto file syntax highlighting
    > vscode-proto3: <https://marketplace.visualstudio.com/items?itemName=zxh404.vscode-proto3>
