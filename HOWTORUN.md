# How our program works

## How to open a client

To run a client, simply run the following command

```shl
go run main.go
```

The new client will automatically connect to the other clients

## How to close a client

To close a client, it is important that you type `-q` in the console and hit return and *DO NOT* use ctrl + c, since this will mess with the port-file system.

### Accidental ctrl + c

If a client is accidentally closed with ctrl + c, the program can be fixed by:

1. Close all other running clients

2. Delete the `ports.txt` file in the `ressources` directory.
