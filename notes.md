# Garage Development Notes

* Garage should have a garage file explicitely passed in with -d <directory>
* otherwise, it will use default ~/.garage/ to install scripts
* it will also use the toumorokoshi/garage-packages as a default location for scripts

## configuration

each garage repository should have a .garagerc in it. It should contain configuration for:

* url(s) to a directory where garage files are kept

## Features

* remotely grabbing shell scripts from multiple locations
* garage butler, which will help you figure out what you want to do
