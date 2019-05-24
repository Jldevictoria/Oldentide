Oldentide Dedicated Server
==
The *Oldentide Dedicated Server* is an open source project, built in [Go][1], to implement the
backend game state for the multiplayer online role-playing game [*Oldentide*][2].  This
directory contains all of the code necessary to build and run the dedicated server.

The "from scratch" build process breaks down into five main steps, with an additional
step used during the development process:

1. Download dependencies.
2. Clone the repository into your $GOPATH/src/ folder.
3. Create the database and generate all of the necessary tables.
4. Populate the database template information. (Races, Professions, Items, and NPCs)
5. Compile the dedicated server.
6. Compile the test client. (If needed for debug.)

Requirements
------------
The *Oldentide Dedicated Server* currently supports a Linux development and runtime environment. (Although it should work on Windows with a little tinkering)

Operating Systems
------------
All development and testing is currently done on a 64-Bit linux environment.
[Ubuntu 18.04 LTS][3] is what the @jldevictoria is using.
Other distributions should work as long as you can run Go, SQLite and Git.

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
In linux (*assuming you have at least your $GOPATH variable exported*):

1. **Download necessary dependencies (see Dependencies above):**

    ```
    sudo apt install --y <...>
    go get <...>
    ```

2. **Clone the Oldentide repository into your $GOPATH/src directory:**

    ```
    cd $GOPATH/src/
    git clone https://github.com/Jldevictoria/Oldentide.git
    ```

3. **Change to the db directory and initialize the Oldentide DB (I call it "oldentide.db"):**

    ```
    cd Oldentide/db/
    sqlite3 oldentide.db < init_db.sql
    ```
    
4. **Populate the db with values from CSV files:**

    ```
    chmod +x init_npcs.sh init_item_templates.sh init_profession_templates.sh init_race_templates.sh
    ./init_npcs.sh
    ./init_item_templates.sh
    ./init_profession_templates.sh
    ./init_race_templates.sh
    ```

5. **Build the Oldentide dedicated server:**

    ```
    go install Oldentide/dedicated_server
    ```

6. **Build the Oldentide test client:**

    ```
    go install Oldentide/test_client
    ```

If everything built properly, the executable for the server should be found in your $GOBIN directory ($GOPATH/bin)

You can run the executable directly, as go compiles to machine code, but you will need to pass in several arguments:

* **Game port (-gport=1337)**
* **Web address (-webadd=http://imp.oldentide.com)**
* **Web port (-wport=80)**

If you want to use email authentication for new accounts, then you need to pass in a gmail username and password:

* **Gmail username (-email=oldentide@gmail.com)**
* **Gmail password (-epass=SuPeRsEcReTpAsSwOrD)**

If you placed your game database in a different location than $GOPATH/src/Oldentide/db/oldentide.db, you will need to specify that as well:

* **DB path (-dbpath=/home/coolguy/secret_db_folder/oldentide.db)**

If you run wil the -h flag should pull up the parameter list.

if you passed in all those arguments, the command line for running the server would be:

    $GOPATH/bin/dedicated_server -gport=1337 -webadd=http://imp.oldentide.com -wport=80 -email=oldentide@gmail.com -epass=SuPeRsEcReTpAsSwOrD

The test client is basically the same as the dedicated server, but you will need to specify the server address and a test case from within test_client.go as well such as:

    $GOPATH/bin/test_client -server=imp.oldentide.com -port=1337 -test=2

Good luck!

Let me know if this is too cryptic and I can update it.

[1]: http://golang.org/ "The Go Language"
[2]: http://www.oldentide.com/ "Oldentide, a game where you can be anyone!"
[3]: http://www.ubuntu.com/ "Ubuntu · The world's most popular free OS"
[4]: https://golang.org/cmd/go/ "Go Cmd Documentation"
[5]: https://git-scm.com/ "Git"
[6]: https://www.sqlite.org/ "SQLite 3"
[7]: https://github.com/msgpack/msgpack-go/ "msgpack-go"
[7]: https://github.com/mattn/go-sqlite3 "go-sqlite3"
