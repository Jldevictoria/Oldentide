![Oldentide Logo](concepts/project/Logo_1.png?raw=true "Oldentide")

# Oldentide

**A Game Where You Can Be Anyone.**

=================================

## Project Information:

**Project name**: Oldentide<br>
**Initial Starting Date**: 2/9/2013<br>
**Founding Author**: Joseph DeVictoria<br>
**Client Platform**: [G3N - Go 3D Game Engine](https://github.com/g3n/engine)<br>
**Languages**: Go, SQLite3<br>
**Project Website**: [www.oldentide.com](http://www.oldentide.com)<br>
**Wiki**: [Github Wiki](https://github.com/Oldentide/Oldentide/wiki)<br>
**Development Server**: [imp.oldentide.com](imp.oldentide.com)<br>
**Development Game Port**: Port 1337<br>
**Contact**: jldevictoria@gmail.com<br>

=================================

## Directories:

**assets**:            A place where all game assets (Maps, Models, Textures, Scripts etc...) go.<br>
**concepts**:          A place where any concept art, drawings, sketches or purely conceptual code goes.<br>
**db**:                A place where the database and its helper scripts are stored.<br>
**lib**:               A place where external open source dlls are stored for G3N engine functionality. Please copy these to bin/<br>
**server**:            A place where the dedicated server source code goes.<br>
**shared**:            A place where shared client-server code is stored.<br>
**test_client**:       A place where the source for a command line test client is stored.  This will be phased out.<br>


**A readme for server and website operation is found within the server directory.**<br>

**Please contact me immediately if you see any bugs or want to contribute.  I need help developing the game!**<br>

[Reading Material](https://www.copyright.gov/fls/fl108.pdf)

=================================

## Client Build Instructions (During DEVELOPMENT):

In order to play the Oldentide client you will need to make sure you follow these steps:

1. Install [Git.](https://git-scm.com/)
2. Install [Go.](https://golang.org/cmd/go/)
3. [Set your $GOPATH.](https://github.com/golang/go/wiki/SettingGOPATH)
4. **cd $GOPATH/src/**
5. **git clone https://github.com/Jldevictoria/Oldentide.git**
6. **go get -u github.com/g3n/engine/...**
7. **go get -u github.com/vmihailenco/msgpack**
7. **go install Oldentide**
8. The Oldentide executable can be run from your $GOBIN directory ($GOPATH/bin).
