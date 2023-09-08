# Command Line Arguments

There are 2 most common ways to take command line arguments, the `Os.Args` and the `flags`` package.

## Os.Args

Consider the command:

```sh
go run .\client\ alice
```

Here we give "alice" as an argument to the program in the folder "client". 

To get access to the argument using `os.Args` you can use:

```go
os.Args[1]
```

Position 0 is the full path to the program, and 1 and onwards is the arguments given in the terminal.

> An online example: [gobyexample.com/command-line-arguments](https://gobyexample.com/command-line-arguments)

## flags

flags allow us to name arguments and makes them optional, but forces us to write longer commands.

```sh
go run .\client\ -name alice
```

Here we set the argument "name" by writing a dash (-) followed by the name of the argument.

We can use it by writing the following code:

```go
import "flag"

var clientsName = flag.String("name", "default", "Senders name")

func main() {
    //parse flag/arguments
    flag.Parse()
}
```

Here we define a flag with the name "name", the default value "default" and a description "Senders name". To use the flag we need to parse it, which we do in the beginning of our main method.

The value of the flag is parsed into `clientName` after the line in main has run. From the command above we set the value to "alice". If we on the other hand didn't set the value of "name" when we run the program, then the value would have been "default".
