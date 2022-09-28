# gRPC guide

> **Disclaimer**
>
> This is NOT to say how you NEED to do it, but a guide that can show you how YOU might WANT to do it.
>
> A simpler older version of this guide can be found at:
> https://github.com/NaddiNadja/grpc101

- [gRPC guide](#grpc-guide)
  - [Setup of new repository](#setup-of-new-repository)
  - [The Proto file](#the-proto-file)
    - [What is it?](#what-is-it)
    - [Required lines](#required-lines)
    - [Defining a service](#defining-a-service)
  - [Implementation](#implementation)
    - [Implementing the server methods](#implementing-the-server-methods)
    - [Calling the endpoints from client](#calling-the-endpoints-from-client)
  - [Prerequisites](#prerequisites)

If you haven't installed google's protocol buffers, see the prerequisites part at the bottom.

## Setup of new repository

1. Make ``go.mod`` file with:

    ``$ go mod init [link to repo without "https://"]``

    Your repo should be on the public github as it needs to be an accessable web page. 

2. Make a ``.proto`` file in a sub-directory, for example ``proto/template.proto`` and fill it with IDL.
    - Notice line 3 and 4.

        ```Go
        option go_package = "DSYS-gRPC-template/gRPC";
        package gRPC;
        ```

    > see [The Proto file](#the-proto-file) for info on what to add in the ``.proto`` file
3. Run command:

    ``$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/template.proto``

    which should create the two ``pb.go`` files. Remember to change "proto/template" to your directory and file name.
4. run command:

    ``$ go mod tidy``

    to install dependencies and create the ``go.sum`` file.
5. Implement your client and server. Refer to [Implementation](#implementation) for instructions.
6. open a terminal for each the client(s) and server(s) and run them with:

    The Client: `$ go run .\client\client.go`

    The Server: `$ go run .\server\server.go`

## The Proto file

### What is it?

The proto file is a file that is used to compile methods for the actual gRPC. So if you make changes in the proto file, then you will have to recompile it, before you can see any changes in your code.

### Required lines

1. To make sure to use the newer proto syntax that, then you will need to add the syntax as the first line in the file.

   ```proto
    syntax = "proto3";
   ```

2. A path for the compiler to know what the package name of the pb.go files needs to have, then an option go_package will need to be set. It is gonna use the last part of the path as the name of the package, but it matters little, it just needs to be unique. The example below will make the package name "proto"

    ```proto
    option go_package = "github.com/DarkLordOfDeadstiny/DSYS-gRPC-template/proto";
    ```

3. Add a package name that is the same as the last folder from above ("proto" in this example)

    ```proto
    package proto;
    ```

### Defining a service

1. Make a service with a telling name.

    ```proto
    service <add a name here>
    {
    }
    ```

2. Add RPC endpoints/methods to the service. The endpoint will start with "rpc" and then a name of the endpoint. We have not made them yet, but the next is the message types that we will send and then what will be returned from the call.

    ```proto
    service <add a name here>
    {
        rpc <Method name> (<Message name>) returns (<Response name>);
    }
    ```

    A endpoint can also be streamed, so more than one message or response can be sent. You can add it to just one or both. Here is an example that does it to both:

     ```proto
        rpc <endpoint name> (stream <Message name>) returns (stream <Response name>);
    ```

3. To complete the endpoint we need to make the message types.

    ```proto
    message <Message type> {
        <type> <variable> = 1;
        <type> <variable> = 2;
        ...
    }
    ```

- notes:
  - All names of services, endpoints and message types should be in CamelCase
  - an example can be found in [template.proto](proto/template.proto)

## Implementation

### Implementing the server methods

1. Make a "server" struct
   - It does not need to be called server, but its just what makes sense.

    ```go
    type Server struct {
        // an interface that the server needs to have
        gRPC.UnimplementedTemplateServer
        
        // here you can impliment other fields that you want
    }
    ```

2. Make a method that matches the endpoint in your proto file. To do this you need to give your struct the method. This is done by adding the "(s *Server)" part as in the examples below.
   - For an endpoint that does no streaming, then we need to give the method a context and the input type. For the return we need to return a pair of your return type and an error.

    ```go
    func (s *Server) <endpoint name>(ctx context.Context, <name> *<input type>) (*<the return type>, error) {
        //some code here
        ...
        ack :=  // make an instance of your return type
        return (ack, nil)
    }
    ```

    - For an endpoint that streams messages, then we need to give the method a stream and return an error.
    - In this case you get the input from the stream and send the return type back over the stream too.

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

    - Register your server struct

    ```go
    gRPC.RegisterTemplateServer(grpcServer, server)
    ```

    - Start serving

    ```go
    grpcServer.Serve(list)
    ```

### Calling the endpoints from client

1. Make dialing options 

    ```go
    var opts []grpc.DialOption
    opts = append(
        opts, grpc.WithBlock(), 
        grpc.WithTransportCredentials(insecure.NewCredentials())
    )

    timeContext, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    ```

2. Make a context with insecure credentials, as we don't have something called a TLS certificate

    ```go
    conn, err := grpc.DialContext(
        timeContext,
        fmt.Sprintf(":%s", *serverPort),
        opts...
    )
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

    - When done sending messages, then you need to close the stream and wait for a response from the server.

    ```go
    farewell, err := stream.CloseAndRecv()
    ```

## Prerequisites

> Feel free to ask for help if this doesn't work anymore

1. before starting, install google's protocol buffers:
    - go to this link: <https://developers.google.com/protocol-buffers/docs/downloads>
    - click on the "release page" link.
    - find the version you need and download it.
    - as per October 2021, if your on windows, it's the third from the bottom, ``protoc-3.18.1-win64.zip``.
2. unzip the downloaded file somewhere "safe".
    - on my windows machine, I placed it in ``C:\Program Files\Protoc``
3. add the path to the ``bin`` folder to your system variables.
    - on windows, click the windows key and search for "system", then there should come something up like "edit the system environment variables".
    - click the button "environment variables..." at the bottom.
    - in the bottom list select the variable called "path" and click "edit ..."
    - in the pop-up window, click "new..."
    - paste the total path to the ``bin`` folder into the text field.

        my path is ``C:\Program Files\Protoc\bin``.
    - click "ok".
4. open a terminal and run these commands:

    ``$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26``

    ``$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1``