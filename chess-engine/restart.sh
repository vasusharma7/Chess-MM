pm2 delete engine
go build engine.go board.go pieces.go server.go
pm2 start engine

