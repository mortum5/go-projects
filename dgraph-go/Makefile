.PHONY: dgraph
dgraph:
	docker run -d -p 5080:5080 -p 6080:6080 -p 8080:8080 \
  	-p 9080:9080 -v ./data:/dgraph --name dgraph \
  	dgraph/standalone:latest

stop:
	docker rm -f dgraph
	
clean:
	rm -rf ./data