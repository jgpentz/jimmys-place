# Cloudflare Zero Trust

This documents explains how to configure a web server with Cloudflare Zero
Trust.  I have a cloud server hosted with Hetzner, and I will be using NGINX
to serve my web content over a Cloudflare tunnel.

## Setting Up NGINX

### Installing a Prebuilt Debian Package from Official NGINX Repo

1) Install prerequisites

```bash
sudo apt install curl gnupg2 ca-certificates lsb-release debian-archive-keyring
```

2) Import an official nginx signing key so apt could verify the packages 
authenticity. Fetch the key:

```bash
curl https://nginx.org/keys/nginx_signing.key | gpg --dearmor \
    | sudo tee /usr/share/keyrings/nginx-archive-keyring.gpg >/dev/null
```

3) Verify that the downloaded file contains the proper key:

```bash
gpg --dry-run --quiet --no-keyring --import --import-options import-show /usr/share/keyrings/nginx-archive-keyring.gpg
```

4) To set up the apt repository for stable nginx packages, run the following
command:

```bash
echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] \
http://nginx.org/packages/debian `lsb_release -cs` nginx" \
    | sudo tee /etc/apt/sources.list.d/nginx.list
```

5) Set up repository pinning to prefer our packages over distribution-provided ones:

```bash
echo -e "Package: *\nPin: origin nginx.org\nPin: release o=nginx\nPin-Priority: 900\n" \
    | sudo tee /etc/apt/preferences.d/99nginx
```

6) Install the NGINX package:

```bash
sudo apt update
sudo apt install nginx
```

7) Check if NGINX is running with `sudo systemctl status nginx`.  If it is not
running, then you start it with `sudo systemctl start nginx` and then enable
it to start at boot time with `sudo systemctl enable nginx`

8) Verify that NGINX Open Source is up and running

```bash
curl -I 127.0.0.1
HTTP/1.1 200 OK
Server: nginx/1.27.0
```

### Setting up a Virtualhost

On my system, I have a user named 'bob' and I would like bob to be able to 
create/modify all files for my website.  My approach to doing this so that
bob doesn't have to use sudo everytime, is to create a group `webgroup-mysite`, 
add bob and nginx to webgroup-mysite, and then enable the setGID bit on the 
website in `/var/www/` directory:

#### Create the Website Directory

1) Create a new web group:

```bash
sudo groupadd webgroup-mysite
```

2) Add users to the web group:

```bash
sudo usermod -aG webgroup-mysite bob
sudo usermod -aG webgroup-mysite nginx
```

3) Change directory ownership

```bash
sudo chown -R root:webgroup-mysite /var/www/mysite
```

4) Set permissions (setGID bit, all for owner, all for group, rx for everyone else)

```bash
sudo chmod -R 2775 /var/www/mysite
```

5) Verify changes

> Note: If you are logged in as Bob you may need to log out and log back in 
before bob can create a file in the mysite directory

```bash
ls -ld /var/www/mysite
groups bob
groups nginx
```

6) Add a dummy `index.html` file to your site directory with these contents:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>HTML 5 Boilerplate</title>
  </head>
  <body>
          <h1>Success! mysite.com is working!</h1>
  </body>
</html>
```

#### Create the Virtualhost

1) Copy the default config as a starting point:

```bash
sudo cp /etc/nginx/sites-available/default /etc/nginx/sites-available/mysite.com
```

2) Delete the default server config and uncomment the example virtualhost config.
We need to modify the config so that it looks like the following:

```bash
server {
        listen 80;
        listen [::]:80;

        server_name mysite.com www.mysite.com;

        root /var/www/mysite;
        index index.html index.htm;

        location / {
                try_files $uri $uri/ =404;
        }
}
```

3) Create a symbolic link from sites-enabled to sites-avaialable

```bash
sudo ln -s /etc/nginx/sites-available/mysite.com /etc/nginx/sites-enabled/
```

4) Check the `/etc/nginx/nginx.conf` file for the following line.  If it does
not include this line, then add it in:

```bash
include /etc/nginx/sites-enabled/*;
```

## Creating the Tunnel

On the Cloudflare dashboard, go to zero trust -> networks -> tunnels, create
a tunnel and follow the instructions for naming, installing, and routing.
