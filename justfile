build:
	go build -v -o madge

clean:
	go clean
	rm *.exe
	rm go_build*

showcase3:
	pushd test_projects/showcase3; ./../../madge jira -j VULN -p 2f2d5f17-c4e7-4a9a-b1ea-79bffb72d9c8; popd

showcase2:
	pushd test_projects/showcase2; ./../../madge jira -j PROJ -p dd937163-c655-4ee2-bb16-32fbc48a75f7; popd
