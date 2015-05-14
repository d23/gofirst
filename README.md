Gofirst - PID 1 process for Docker containers.
==============================================

Wrapper to take the pid 1 responsibility from programs running as entrypoint or
command in your docker containers.

Takes care of:
* Reaping of orphaned processes
* Terminating any orphaned processes, when main process dies.
* Relay signals to running processes (so ^C will work!)

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
SIGTERM, SIGHUP, SIGINT and SIGQUIT are relayed to any running process willing
to listen. As soon as the main process exits, SIGTERM is broadcast, and after 10
seconds SIGKILL is sent to any remaining processes. Gofirst will terminate as
soon as all children have terminated.

This fits nicely with a one app, one container philosophy, and makes signals
behave like expected. If that's not what you're doing, you might consider a
proper init system, like runit or others.

Issues
------
Don't run it outside of a container and send it SIGTERM ;)

[1] https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/

