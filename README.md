## How to run

The following should be done in seperate terminals and in this order
  - Run: go run .\server_node.go 5052 5053
  - Run: go run .\server_node.go 5051 5052
  - Run: go run .\client\client.go client_one 5051
  - Run: go run .\client\client.go client_two 5051

You can now run the Bid and Result commands in the clients and terminate either of the server processes and the program
will continue with expected funcitonaity. (The auction last 30 seconds)
