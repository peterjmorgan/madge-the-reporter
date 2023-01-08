build:
	go generate
	go build -v -o madge

clean:
	go clean
	rm *.exe
	rm go_build*

showcase3:
	op run -- ./madge jira -j VULN -p dd937163-c655-4ee2-bb16-32fbc48a75f7 -c test_projects/showcase3/madge_config.yaml

showcase2:
	op run -- ./madge jira -j PROJ -p 2f2d5f17-c4e7-4a9a-b1ea-79bffb72d9c8 -c test_projects/showcase2/madge_config.yaml

fail:
	op run -- ./madge jira -d -j VULN -p dd937163-c655-4ee2-bb16-32fbc48a75f7 -c test_projects/fail/madge_config.yaml

