all:
	go build worker.go
	go build main.go
try:
	curl -G "127.0.0.1:8888/add" --data-urlencode "a=7" --data-urlencode "b=5"
