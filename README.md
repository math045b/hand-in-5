## How to run

The following should be done in seperate terminal

Start the leader node on 5050
```sh
go run .\server_node.go
```
starts replication node 
```sh
go run .\server_node.go 5051
```

To setup clients run this in different terminal
```sh
go run .\client\client.go <name>
```
You can now run the Bid and Result commands in the clients and terminate either of the server processes and the program
will continue with expected funcitonaity. (The auction last 30 seconds)
