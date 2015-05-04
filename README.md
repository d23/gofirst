Gofirst - PID 1 process for Docker containers.
==============================================

Wrapper to take the pid 1 responsibility from arbitrary programs running as
entrypoint or command in your docker containers.

Takes care of:
* Reaping of orphaned processes
* Terminating any orphaned processes when main process dies or SIGTERM is recieved.
* Relay SIGINT to main process.

For more information see the phusion blog[1] and their my_init script. This is
similar but simpler and without dependencies.

Usage
-----
Add the binary to the (path of the) container and set it as first argument to
the entrypoint in your Dockerfile:

    ENTRYPOINT ["gofirst"]

or

    ENTRYPOINT ["gofirst", "some-command"]


Gofirst will execute the command with the arguments supplied.
As soon as that process dies, it will send SIGTERM
to any remaining children after 10 seconds. Any recieved SIGTERM will be
broadcasted immediately. In either case, it will send SIGKILL 30 seconds after
that.

Issues
------
Don't run it outside of a container and send it SIGTERM ;)

[1] https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/

