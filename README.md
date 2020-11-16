Dispatcher
================


### INSTALL

```shell script
# download the source code 
go get github.com/juxuny/dispatcher

# install dispatcher command line 
go install github.com/juxuny/dispatcher
```

### START SERVICE

```shell script

mkdir scripts # create a working directory
dispatcher -w scripts -l :8080 -b /bin/bash

```