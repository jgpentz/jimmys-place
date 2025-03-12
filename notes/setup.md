# Setting up a debian server on hetzner

## Creating a Resource (Server)

Everything in brackets **[ example ]** are my selections

1) From the dashboard, click on **Create Resource** and
select the **Servers** option.

2) Select the server location **[Hillsboro, OR]**

3) Select the image **[Debian]**

4) Select the server type **[Share vCPU -> x86]**

5) Determine if you want to use both IPv4 and IPv6 **[IPv4 IPv6]**

6) Add an SSH key (you can always reset the root password later if needed)

- First generate the key:

```bash
ssh-keygen -t ed25519 -C "your@email.com"
```

- Then copy the generated public key `/home/you/.ssh/id_ed25519.pub` and add it
to the prompt.  I also made this key my default, and you have the option
of adding multiple keys

7) Add volumes **[None]**

8) Create daily backups **[None]**

9) Placement groups **[None]**

10) Add Labels **[None]**

11) Use cloud-init to execute init scripts **[None]**

12) Give your server a name **[my-cloud-box]**

## SSH Login as Root

Now that the server is created, you should be able to log in as **root**. Your
public IP should be available on your server dashboard, so use that with the
key you generated when creating the server.

```bash
ssh -i /home/you/.ssh/id_ed25519 root@1.2.3.4
```

### Create an SSH Alias

While that's fun, it is more convenient to just type in something simpler, like
`ssh my-cloud-box`. 

If it doesn't already exist, create a file in the .ssh
directory named config `/home/you/.ssh/config`, and then add the following
entry:

```bash
Host cloud-box
    HostName 1.2.3.4
    User root
    IdentityFile ~/.ssh/id_ed25519
```

Now try logging in with your new alias!

```bash
ssh cloud-box
```

## Creating a new user

1) Create a user group for the new user you are about to add:

```bash
groupadd -g 3000 cool_cat
```

2) Create a user, and add them to the user group you just created:

```bash
useradd -u 3000 -g 3000 -m -s /bin/bash -k /etc/skel cool_cat
```
**Options:**

- -u provides the **User ID**
- -g provides the **Group ID**
- -c provides a comment for **Gecos** (here it's the user's full name)
- -m to create a home directory (default of **'/home/cool_cat'**)
- -s to provide the default shell
- -k to provide the default 'skeleton' files to add to the user's home directory

The last argument provided is the **Name** (username) of the user.

A couple of general practices:

1) It has become standard when creating a user to create a corresponding primary
group (with a group name the same as the username), and to match the **User ID**
and **Group ID**. This way, when the user is creating new files and directories,
the default group of the new files and directories is that of the user.  If the
permissions are set correctly, other users won't get access to those new files 
by default. The default group will be overriden in directories that have the
set GID permission set.
2) In order to add the user to a group, the group must exist before the 
user is created. 
3) When the user's home directory is created, the owner and group of the home directory
are those specified in the **useradd** command. It's wise to change the permissions
of that home directory to **700** so that only the owner of the directory (
i.e. the user who created it) has access.

### Add Password for New User

Now that you have a new user, you can give them a password with the following:

```bash
passwd cool_cat
```

### Viewing User Accounts
The easiest way to view all user accounts on a system is to look at the contents
of the file `/etc/passwd`

Here is what that looks like on an example machine, (the actual output has been paired down)

```text
$ cat /etc/passwd

root:x:0:0:root:/root:/bin/bash
...
...
cool_cat:x:3015:3015:Cool Cat:/home/cool_cat:/bin/bash
```

The fields that it displays are:

`Name : Password : UserID : GroupID : Gecos : HomeDirectory : Shell`

- **Name:** user's login name
- **Password:** the value of 'x' indicating that the password for this exists
in /etc/shadow and shouldn't be modified. In legacy systems, passwords use
to exist in this field, but they've been replaced by encrypted passwords
that have entries in /etc/shadow which correspond to each user.
- **User ID:** user's unique numeric ID.
- **Group ID:** user's primary group ID. This must be the numeric ID of a group
in the user  database.
- **Gecos:** general information about the user that is not needed by the system.
This could be full name,phone number, email address, etc..
- **HomeDirectory:** the full path name of the user's home directory. The value
is a character string.
- **Shell:** the initial program or shell that is executed after a user invokes
the login command or su command. If the user does not have a defined shell ,
**/usr/bin/sh**, the system shell, is used. Value is a character string that
may contain arguments to pass to initial program.

### Viewing Groups

Viewing all groups on a system can be viewed from the `/etc/group` file:

```text
$ cat /etc/group

root:x:0:
bin:x:1:
daemon:x:2:
sys:x:3:
...
...
cool_cat:x:3000:
```

The fields that it displays are:

`Name : Password : ID : User1,User2,...,Usern`

- **Name:** the group name that is unique on the system.
- **Password:** not used. Group administrators are provided instead of
group passwords.
- **ID:** the group ID. This is a unique decimal integer string.
- **Users:** a list of zero or more users. Separate group member names with
commas.

### SUDO

In order to minimize the damage that you'll do, make sure that you add your
new user to the sudo group and login as the new user from now on.

```bash
usermod -aG sudo cool_cat

# Verify that your user now exists in the group
cat /etc/group | grep sudo
```

You can also check the permissions given to sudoers by viewing the sudoers file
`/etc/sudoers` and ensuring the following line is in there:

```bash
# Allow members of group sudo to execute any command
%sudo   ALL=(ALL:ALL) ALL
```

### SSH Login as New User

Now that you have a new user, generate a new ssh key and update your 
`~/.ssh/config` with your new user and identity file:

```bash
Host cloud-box-user
    HostName 1.2.3.4
    User cool_cat
    IdentityFile ~/.ssh/id_ed25519-cool_cat
```

Then you have to transfer the public key to the new user on the server.  The
easiest way to do this is with the `ssh-copy-id` command, which you should 
always use since modifying the `authorized_keys` file by hand can easily
mess up other existing keys, and you could introduce weird whitespace characters.

```bash
ssh-copy-id -i ~/.ssh/id_ed25519-cool_cat.pub cool_cat@1.2.3.4
```

### Disable Password Login

Edit the `/etc/ssh/sshd_config` by changing the line `PasswordAuthentication yes`
to `PasswordAuthentication no`

Then restart the SSH service by running `sudo service ssh restart`
