![pipeline](https://gitlab.com/onurkeskin/go-http-fs-client/badges/master/pipeline.svg)
![coverage](https://gitlab.com/onurkeskin/go-http-fs-client/badges/master/coverage.svg)

# Client for standard - Go http file server

A minimal http file server with a tcp client. Has a lot of boilerplate http code for simpler development and more modular calling structure.
Only signifcant file is https://gitlab.com/onurkeskin/go-http-fs-client/-/blob/master/domain/core/file_service/file_service.go

# Local Setup

You can build, run and test if go is already installed on host. If go is not available, but you have a running docker environment, you can use the devcontainer I use in vscode. It will setup a nice development environment and ide with settings I prefer to use. Note that if you are using windows machines it will need linux containers.

# Container Setup

Can be run with 
``` make fs-server-run ``` or
``` make fs-client-run ```
to setup specific containers. However these two commands will require further network sharing between containers to communicate efficiently. For simpler debugging I recommend running one on host pc.

# Usage

There is a single path based endpoint which takes a single lower case letter. Searches the files in the file-server for earliest position for the specified character.
After getting the client application running on the localhost, to start a file transfer simply calling
``` localhost:8080/z ```
would start downloading files that contains the earliest z.

However if this is being run on a container, to see the files being transferred simply mounting a host volume to the container would suffice.

# Code Coverage
Coverage Reports can be accessed at 
**https://gitlab.com/onurkeskin/go-http-fs-client/**
