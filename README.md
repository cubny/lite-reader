Lite-Reader
===========
Read your feeds on your own machine with a simple and lite application.


Screenshot
----------
![ScreenShot](https://raw.github.com/cubny/lite-reader/master/public/images/screenshot.png)

Requirements to run
-------------------
- None, just download the binary from the releases page and run it.

Migration from legacy Lite Reader
---------------------------------
If you are using the legacy Lite Reader, you can migrate your data to the new Lite Reader.
1. Download the latest release of Lite Reader
2. Copy the data folder consisting of the agg.db file from the legacy Lite Reader to the new Lite Reader folder
3. Run the new Lite Reader

Requirements to build
---------------------
- SQLite3
- Golang 

INSTALL
--------
1. git clone or download
2. run `go get` to install all dependencies
3. run `go build` to build the application
4. run `./lite-reader` to start the application

that's it, enjoy a very lite and minimal feed aggregator: the lite-reader


Want to Contribute?
-------------------
- report bugs (https://github.com/cubny/lite-reader/issues/new)
- share your ideas (https://github.com/cubny/lite-reader/issues/new)
- view our roadmap (https://trello.com/b/ekJbxyCL/lite-reader)
- fork and make changes


Legacy Lite Reader
------------------
The legacy version of Lite Reader is available at [legacy-lite-reader](https://github.com/cubny/legacy-lite-reader) repository. It is written in PHP and uses MySQL as the database. It is no longer maintained and is not recommended for use.
