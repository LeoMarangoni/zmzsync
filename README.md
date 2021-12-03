# zmzsync

ZmzSync was made to facillitate account migration between zimbra servers, inspired in [imapsync](https://github.com/imapsync/imapsync/).

You can migrate accounts

------

Building:
-----

> - Install a golang dev environment
> - Clone this repository
> - go install
> - go build <- build for linux
> - GOOS=windows GOARCH=amd64 go build <- build for windows

The default workpath is linux "/tmp" in windows you are obliged to set the dir that will be used to store temp files with -workdir option


If you want just to use zmzsync for migration purposes, download from it [releases](https://github.com/LeoMarangoni/zmzsync/releases)

### **Help Menu:**
```help
Usage of zmzsync:
  -authuser1 string
    	Admin user to authenticate in host1. (Optional)
  -authuser2 string
    	Admin user to authenticate in host2. (Optional)
  -dateend string
    	Migrate from that date back. Date format YYYY-MM-DD
  -datestart string
    	Migrate from that date forward. Date format YYYY-MM-DD
  -host1 string
    	Source zimbra server.
  -host2 string
    	Destination zimbra server.
  -password1 string
    	Password to authenticate in host1.
  -password2 string
    	Password to authenticate in host2.
  -port1 string
    	 zimbra Admin Port on host1. (default "7071")
  -port2 string
    	 zimbra Admin Port on host2. (default "7071")
  -types string
    	Which data type to migrate. default is to migrate ALL
    		Valid Types are:
    			message
    			conversation
    			contact
    			appointment
    			task
    		Example, migrating only calendar and contacts: -types contact,appointment
    	
  -user1 string
    	User to migrate from host1.
  -user2 string
    	User to migrate to host2.
  -workdir string
    	Directory to store exported tgz files. (default "/tmp")
```

