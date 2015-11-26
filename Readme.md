# What is that for?

Web interface for the ad-blocking Pi-hole: **a black hole for Internet advertisements**.

## Developed in Golang, features built-in web server and single executable with all the assets included. No lighttpd or php-fpm required!

From this interface, you will be able to see stats on how well your Pi-hole is performing.  You will also be able to update the lists used to block ads.

![Pi-hole Web interface](http://i.imgur.com/x2iMfoc.png)
![Fully responsive](http://i.imgur.com/NyAIXm8.png)

# Installation

Stop the lighttpd daemon.

```
sudo service lighttpd stop
```

```
wget https://github.com/girishso/pi-hole-web/releases/download/v1.0/pi-holeweb.zip
unzip pi-holeweb.zip
sudo ./pi-holeweb
```

Open the admin console: http://raspberry-pi-ip/admin

# Building

Uses https://github.com/elazarl/go-bindata-assetfs for embedding the static assets in executable.

// for debug

go-bindata-assetfs -debug static/... templates


// for release

go-bindata-assetfs -nomemcopy static/... templates

env GOOS=linux GOARCH=arm GOARM=6 go build


