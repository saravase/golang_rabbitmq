#01-simple
send-1:
	go run 01-simple/sender/main.go
receive-1:
	go run 01-simple/receiver/main.go

#02-work-queue
send-2:
	go run 02-work-queue/sender/main.go
receive-2:
	go run 02-work-queue/receiver/main.go

#03-exchange
send-3:
	go run 03-exchange/sender/main.go
receive-3:
	go run 03-exchange/receiver/main.go

#04-routing
send-4:
	go run 04-routing/sender/main.go
receive-4:
	go run 04-routing/receiver/main.go