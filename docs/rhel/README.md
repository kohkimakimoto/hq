# HQ yum repository

The yum repository for hosting HQ rpm packages.

## el7-x86_64

### Installation

Install Repository

```
$ sudo sh -c 'echo "[hq]
name=hq
baseurl=https://kohkimakimoto.github.io/hq/rhel/7/x86_64
gpgcheck=0
enabled=1
" > /etc/yum.repos.d/hq.repo'
```

Install HQ

```bash
yum install hq
```
