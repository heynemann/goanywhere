setup:
	@cat requirements | xargs go get

run:
	@go install && echo "GoAnywhere updated." && echo "GoAnywhere running in http://localhost:12345/" && ../../bin/goanywhere

kill_mongo:
	ps aux | awk '(/mongod/ && $$0 !~ /awk/){ system("kill -9 "$$2) }'

mongo: kill_mongo
	rm -rf /tmp/goanywhere/mongodata && mkdir -p /tmp/goanywhere/mongodata
	mongod --dbpath /tmp/goanywhere/mongodata --logpath /tmp/goanywhere/mongolog --port 8888 --quiet &
