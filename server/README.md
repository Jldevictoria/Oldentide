Oldentide Dedicated Server
==
The *Oldentide Dedicated Server* is an open source project, built in [Go][1], to implement the
backend game state for the multiplayer online role-playing game [*Oldentide*][2].  This
directory contains all of the code necessary to build and run the dedicated server.

The "from scratch" build process breaks down into three main steps, with an additional
step used during the development process:

1. Create the database and generate all of the needed tables.
2. Populate the database with NPC information.
3. Compile the dedicated server.
4. Compile the test client. <IP>

Requirements
------------
The *Oldentide Dedicated Server* currently supports a Linux development and runtime environment. (Although it should work on Windows with some hacking)

Operating Systems
------------
All development and testing is currently done on a 64-Bit linux environment.
[Ubuntu 18.04 LTS][3] is recommended (being used by the primary developer).
Other distributions should work as long as you can run Go and Git.

This should also run on Windows if you have sqlite and administrator privileges, but I have always tested on Linux.

Compilers
------------
Building *Oldentide Dedicated Server* is consistent with the standard ["Go" build/install tools][4]

Dependencies
------------
* [git][5] - Needed for revision control, and for downloading and installing msgpack and go-sqlite3.
    * **sudo apt install -y git-all**
* [sqlite3][6] - The sqlite3 command-line tool, used to initialize and manage our sqlite databases.
    * **sudo apt install -y sqlite3 libsqlite3-dev**
    * **go get github.com/mattn/go-sqlite3**
* [msgpack-go][7] is used for data serialization for transmitting packets between server and client.
    * Msgpack recommends the Vmihailenco implementation of msgpack for go.
    * **go get github.com/vmihailenco/msgpack**

Server and Test Client Usage
------------
In linux:

Switch to the db directory and initialize the Oldentide DB:

    sqlite3 Oldentide.db < init_db.sql
    
Populate the db with values from CSV files:

    chmod +x init_npcs.sh init_item_templates.sh init_profession_templates.sh init_race_templates.sh
    ./init_npcs.sh
    ./init_item_templates.sh
    ./init_profession_templates.sh
    ./init_race_templates.sh

Make sure you have the Oldentide repository cloned into your $GOPATH/src folder.

Download necessary dependencies (see above)

    go get <...>

cd into the src/server folder and run

    go install
    or
    go install /path/to/Oldentide/server

If everything built properly, the executable for the server should be found in your $GOBIN directory ($GOPATH/bin)

You can run the executable, and you will need to pass in several arguments for the game port, web port, web address, if you are using email authentication, then you need to pass in the gmail username and password.

I believe the -h flag should pull up the parameter list.

[1]: http://golang.org/ "The Go Language"
[2]: http://www.oldentide.com/ "Oldentide, a game where you can be anyone!"
[3]: http://www.ubuntu.com/ "Ubuntu · The world's most popular free OS"
[4]: https://golang.org/cmd/go/ "Go Cmd Documentation"
[5]: https://git-scm.com/ "Git"
[6]: https://www.sqlite.org/ "SQLite 3"
[7]: https://github.com/msgpack/msgpack-go/ "msgpack-go"
[7]: https://github.com/mattn/go-sqlite3 "go-sqlite3"
