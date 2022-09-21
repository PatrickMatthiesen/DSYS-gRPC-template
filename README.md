# gRPC 101

if you haven't installed google's protocol buffers, see the prerequisites part at the bottom.

## Setup of repository

1. make ``go.mod`` file with:

    ``$ go mod init [link to repo without "https://"]``

    your repo should be on the public github. i couldn't get it to work on the itu instance.
    > **Tip!**
    > 
    > You can use ``example.com``, if you don't have a repo.
2. make a ``.proto`` file in a sub-directory, for example ``proto/template.proto`` and fill it with IDL.
    - notice line 3 and 4.

        ```Go
        option go_package = "GRPC-template/gRPC";
        package gRPC;
        ```
    
3. run command:

    ``$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/template.proto``

    which should create the two ``pb.go`` files. Remember to change "proto/template" to your directory and file name.
4. run command:

    ``$ go mod tidy``

    to install dependencies and create the ``go.sum`` file.
5. create a ``client\client.go`` file with the ``client_template.txt`` as template.
    > **Tip!**
    >
    > When implementing your grpc methods, you should write the link without "https://" and with the packagename at the end. If you used example.com, you should write ``"example.com/package"``. If you used a long name for your package, you can write a shorter name (alias) before the quotation marks, for example ``pckg "example.com/longpackagename"``.
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


## Prerequisites

1. before starting, install google's protocol buffers:
    - go to this link: <https://developers.google.com/protocol-buffers/docs/downloads>
    - click on the "release page" link.
    - find the version you need and download it.
    - as per october 2021, if your on windows, it's the third from the bottom, ``protoc-3.18.1-win64.zip``.
2. unzip the downloaded file somewhere "safe".
    - on my windows machine, i placed it in ``C:\Program Files\Protoc``
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