# gRPC 101

- [gRPC 101](#grpc-101)
  - [Todo](#todo)
  - [Setup of repository](#setup-of-repository)
  - [The Proto file](#the-proto-file)
    - [What is it?](#what-is-it)
    - [Required lines](#required-lines)
    - [Defining a service](#defining-a-service)
    - [Implementing the server methods](#implementing-the-server-methods)
  - [Prerequisites](#prerequisites)

## Todo

- do some groupings in the files for better understanding
- explain that the method interface in the server needs to match the one from the compiled proto
  - could be just told how to write it, as it is the same syntax all the time
- explain that the methods on the server and client are depending on the name of the proto file
- remember that all methods for implementing the API, needs to extend the server type



If you haven't installed google's protocol buffers, see the prerequisites part at the bottom.

## Setup of repository

1. Make ``go.mod`` file with:

    ``$ go mod init [link to repo without "https://"]``

    Your repo should be on the public github as it needs to be an accessable web page. 

2. Make a ``.proto`` file in a sub-directory, for example ``proto/template.proto`` and fill it with IDL.
    - Notice line 3 and 4.

        ```Go
        option go_package = "DSYS-gRPC-template/gRPC";
        package gRPC;
        ```

3. Run command:

    ``$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/template.proto``

    which should create the two ``pb.go`` files. Remember to change "proto/template" to your directory and file name.
4. run command:

    ``$ go mod tidy``

    to install dependencies and create the ``go.sum`` file.
5. create a ``client\client.go`` file with the ``client_template.txt`` as template.
    > **Tip!**
    >
    > When implementing your grpc methods, you should write the link without "https://" and with the package name at the end. If you used example.com, you should write ``"example.com/package"``. If you used a long name for your package, you can write a shorter name (alias) before the quotation marks, for example ``pckg "example.com/longpackagename"``.
6. create a ``server\server.go`` file with the ``server_template.txt`` as template.
7. switch out the "myPackage" with your actual package.
8. switch our the method names with actual method names.
9. add more methods to the ``client.go`` file, so that there's a method for each request in the ``.proto`` file.
10. when everything is compilable, open a terminal, change directory to the ``server`` folder, and run the command:

    ``$ go run .``

    this will start your server.
11. open a new terminal, change directory to the ``client`` folder and run the command:

    ``$ go run .``

    this will run the requests listed in the ``main`` method of the ``client`` file.

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

2. Add RPC methods to the service. The method will start with "rpc" and then a name of the method. We have not made them yet, but the next is the message types that we will send and then what will be returned from the call.

    ```proto
    service <add a name here>
    {
        rpc <Method name> (<Message name>) returns (<Response name>);
    }
    ```

    A method can also be streamed, so more than one message or response can be sent. You can add it to just one or both. Here is an example that does it to both:

     ```proto
        rpc <Method name> (stream <Message name>) returns (stream <Response name>);
    ```

3. To complete the method we need to make the message types.

    ```proto
    message <Message type> {
        <type> <variable> = 1;
        <type> <variable> = 2;
        ...
    }
    ```

- notes:
  - All names of services, methods and message types should be in CamelCase
  - an example can be found in [template.proto](proto\template.proto)
  
### Implementing the server methods

- TODO

## Prerequisites

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