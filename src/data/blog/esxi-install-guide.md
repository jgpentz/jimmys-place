---
author: James Pentz
pubDatetime: 2026-2-6T11:30:00Z
modDatetime: 2026-2-6T11:30:00Z
title: ESXI Installation Guide (Plus OPNSense)
slug: esxi-installation
featured: false
draft: false
tags:
  - esxi
  - hypervisor
  - virtualization
  - opnsense
  - neteng
description:
  A step-by-step walkthrough for setting up VMware ESXi on a NUC and deploying OPNsense as your virtualized firewall solution.
---


# Installing ESXi and OPNsense: A Complete Guide

A step-by-step walkthrough for setting up VMware ESXi on a NUC and deploying OPNsense as your virtualized firewall solution.

---

## Table of contents

---

## Part 1: Downloading and Preparing ESXi

### Downloading ESXi from Broadcom

1. Navigate to the [Broadcom website](https://www.broadcom.com) and create an account (or log in if you already have one)
2. Click on the **My Downloads** tab
3. Select the **Free Downloads** button located in the middle of the interface
4. Click on **VMware vSphere Hypervisor**
5. Choose the latest version and download the ISO

### Creating a Bootable USB Drive (macOS)

#### Find Your USB Drive

```bash
diskutil list
```

Identify your USB drive from the list. In this example, we'll use `disk4`.

#### Format the Drive

```bash
diskutil eraseDisk MS-DOS "ESXI" MBR disk4
```

#### Unmount the Drive

```bash
diskutil unmountDisk /dev/disk4
```

#### Configure the Boot Partition

Enter fdisk interactive mode:

```bash
sudo fdisk -e /dev/disk4
```

Make the first partition active and exit:

```
f 1
quit
```

#### Copy ESXi Files

1. Mount the ESXi ISO file (double-click in Finder)
2. Copy all contents from the mounted ISO to your USB drive (you can use Finder for this)

#### Modify the Boot Configuration

Navigate to the USB drive:

```bash
cd /Volumes/ESXI
```

Edit the boot configuration file using your preferred text editor:

```bash
vi ISOLINUX.CFG
```

Find the line that reads:

```
APPEND -c boot.cfg
```

Modify it to:

```
APPEND -c boot.cfg -p 1
```

Rename the configuration file:

```bash
mv ISOLINUX.CFG SYSLINUX.CFG
```

Return to the root directory and unmount the drive:

```bash
cd /
diskutil unmountDisk /dev/disk4
```

---

## Part 2: Installing ESXi on Your NUC

### BIOS Configuration

1. Insert the USB drive into your NUC
2. Power on the NUC and enter the BIOS (typically by pressing F2 or DEL during startup)
3. Navigate to the boot order settings
4. Set the USB drive as the first boot priority
5. Save and reboot

### ESXi Installation

1. The NUC will boot from the USB drive and launch the ESXi installer
2. Follow the on-screen prompts to complete the installation
3. Once installation is complete, **remove the USB drive**
4. Reboot the system

Your ESXi host should now be up and running. Access the web interface using the IP address displayed on the console.

---

## Part 3: Installing OPNsense

### Prepare the Network Configuration

In the ESXi web interface, create port groups for your network segmentation:

- Create a port group for **WAN** (external network)
- Create a port group for **LAN** (internal network)

### Download OPNsense

1. Visit the [OPNsense website](https://opnsense.org/download/)
2. Download the **DVD ISO** file

### Upload to ESXi Datastore

1. In the ESXi web interface, navigate to **Storage**
2. Select your datastore
3. Upload the OPNsense ISO file

### Create the OPNsense Virtual Machine

Create a new VM with the following specifications:

**General Settings:**
- Guest OS Family: **Other**
- Guest OS Version: **FreeBSD 14 or later**

**Hardware:**
- CPU: **2 cores**
- RAM: **4 GB**
- Hard Disk: **20 GB**
- SCSI Controller: **LSI Logic Parallel**

**Network:**
- Add **two network adapters**:
  - Network Adapter 1: Connected to LAN port group
  - Network Adapter 2: Connected to WAN port group

**Installation Media:**
- Add the OPNsense ISO to the CD/DVD drive

### Start the Installation

1. Power on the VM
2. The OPNsense installer will boot from the ISO
3. Follow the installation wizard to complete the setup

---

## Next Steps

After installation, you'll want to:

- Configure OPNsense network interfaces (assign WAN and LAN)
- Set up firewall rules
- Configure DHCP and DNS services
- Enable additional features as needed

Your virtualized network infrastructure is now ready to go!

---

