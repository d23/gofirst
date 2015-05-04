Gofirst - PID 1 process for Docker containers.
==============================================

Wrapper to take the pid 1 responsibility from arbitrary programs running as
entrypoint or command in your docker containers.

Takes care of:
* Reaping of orphaned children
* Terminating any orphaned processes when main process dies or SIGTERM is recieved.
* Relay SIGINT to main process.

For more information see the phusion blog[1] and their my_init script. This is similar
but simpler and without dependencies.

Usage
-----
Add the binary to the (path of the) container and set it as entrypoint in your Dockerfile

ENTRYPOINT ["gofirst"]

or as first element of your CMD array

CMD ["gofirst", "command", "arg1", "arg2"]


Issues
------
Don't run it outside of a container and send it SIGTERM ;)

[1] https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/

