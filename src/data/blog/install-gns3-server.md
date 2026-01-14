---
author: James Pentz
pubDatetime: 2026-1-13T19:30:00Z
modDatetime: 2026-1-13T19:30:00Z
title: Installing a GNS3 Server on Ubuntu
slug: gns3-server
featured: true
draft: false
tags:
  - gns3
  - proxmox
  - neteng
description:
  Creating a GNS3 server that runs on an Ubuntu 22.04.5 LTS Proxmox VM.
---

Now that I finally have some free time, I am starting to build out my homelab.
One of the first tools that I am building is a GNS3 server that runs on an 
Ubuntu VM in Proxmox.


## Table of contents

## Version Info

- **Proxmox:** 9.1.4
- **GNS3 Server:** 2.2.55

This guide walks through setting up a **GNS3 Server** inside an **Ubuntu Server VM** running on **Proxmox**. This approach keeps GNS3 isolated, enables better performance with nested virtualization, and works well for home labs and remote access setups.

## 1. Create an Ubuntu Server VM on Proxmox

### Download Ubuntu Server ISO

Download an **Ubuntu Server install image** from the official Ubuntu website (for example, Ubuntu 22.04 LTS).

### Upload the ISO to Proxmox

You can upload the ISO either through the **Proxmox Web GUI** or via **secure copy (scp)**.

#### Option A: Proxmox Web GUI

1. Expand your Proxmox node in the left navigation pane
2. Click on the `local` storage
3. Open the **ISO Images** tab
4. Click **Upload** and select the Ubuntu Server ISO

#### Option B: Secure Copy (scp)

From your local machine, copy the ISO to the Proxmox ISO directory:

```bash
scp ubuntu-22.04.5-live-server-amd64.iso root@192.168.0.123:/var/lib/vz/template/iso/
```

### Create the Virtual Machine

In the Proxmox web interface:

1. Click **Create VM**
2. Select the uploaded Ubuntu ISO as the installation media
3. Allocate resources appropriate for GNS3

Recommended minimums:

* **Memory:** 8 GB RAM
* **Storage:** 32 GB

### Enable Nested Virtualization

During VM creation, when you reach the **CPU** configuration step:

* Enable **Nested Virtualization**
* Set **`nested-virt` = `on`** under **Extra CPU Flags**

This is required for running virtualized network appliances inside GNS3.

## 2. Install Ubuntu Server

1. Boot the VM
2. Open the console from the Proxmox UI
3. Follow the Ubuntu Server installer prompts

Once installation completes:

1. **Shut down** the VM
2. Start it again

This ensures the installation media is detached and the system boots from disk.

## 3. Install the GNS3 Server

After logging into the VM, install the GNS3 server using the official installation script.

### Run the Installation Script

As `root`, run:

```bash
cd /tmp
curl https://raw.githubusercontent.com/GNS3/gns3-server/master/scripts/remote-install.sh > gns3-remote-install.sh
bash gns3-remote-install.sh --with-iou --with-i386-repository
```

### Notes

* The `--with-iou` option enables IOU support
* The `--with-i386-repository` option is required for some legacy images
* The `--with-openvpn` option is **not** used here, since this setup assumes access from a trusted home network

## 4. Verify the GNS3 Server

Check that the GNS3 service is running:

```bash
systemctl status gns3
```

If running correctly:

* The server listens on **port 3080** by default
* It is reachable at the VM’s IP address

To confirm the listening port:

```bash
sudo ss -tulnp | grep gns3
```

## 5. Add the Remote Server to GNS3

On your local machine, open **GNS3** and configure the remote server:

1. Go to **Edit → Preferences → Server**
2. Uncheck **Enable local server**
3. Enter the VM’s IP address in the **Host** field
4. Enter the credentials for the GNS3 server

Default credentials:

* **Username:** `gns3`
* **Password:** `gns3`

Apply the changes and restart GNS3 if prompted.

## Future

* Integrating the server with Tailscale or VPN access
* Automating VM creation with Terraform or Ansible
