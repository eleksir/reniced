# Reniced

Helper utility intended to run as a service or as background process (say, in tmux) that can make system more
respositive under certain type of load.

## Why it was made

In general case os can handle load well, fairly sharing resources among processes. But in some cases it makes it poorly.
As an example it can be ill-minded pseudo-security stuff like antiviruses or mediaanalysers that searches forbidden
content on your mac (due to gov.policies that apple have to accept). Or it can be backup that crawls fs in background.

To help resolve such situation this utility was made. Some malicious processes can be reniced to lower nicelevels (say,
19) and for some good processes nicelevel can be raisen to higher values (-19 if reniced run as superuser or 0 if as
normal).

In some critical cases process can be killed (if you want to know more about it see "man kill" in terminal).

## How to build, configure, run it

You will need Go lang 1.21 or newer, gnu make, terminal and possibly internets to build this utility.

To build reniced just type in project dir

```bash
make
```

this will produce binary named "reniced".

Copy to your home directory .reniced.yaml.example as .reniced.yaml and edit it to suit your needs. After you done, you
can run utility as:

```bash
./reniced
```

if you have no access to superuser or

```bash
sudo ./reniced
```

if you can sudo. That's it.

Note: negative niceness values can be set only by root.

Note2: MacOS has too tight (or intrusive) security measures, so running something as system daemon is very difficult
task. Reniced can be run via sudo as daemon or backbround process (say, in tmux). On linux systems usually there is no
such obstacle and if reniced run from user with id=0, its config should be placed at **/etc/reniced.yaml**

## Install as systemd service

Supplied systemd unit is ready to install in system as-is and should work ok. It expect to see reniced at
**/usr/sbin/reniced** and assume that reniced runs as foreground task - in config file "daemon" setting either absent or
set to false.

## Install as System V Init service

Unlike supplpied systemd unit, supplied System V Init script assumes reniced runs as daemon (in config file "daemon"
setting must be set to true). Script also expects to find pid-file at **/var/run/reniced.pid** and executable itself at
**/usr/sbin/reniced**.

## Note on ionice operation

Currently implemented only linux io prio. On all other reniced-supported systems it is noop. Please note that io
priority can be set only by digits. Refer to *man ionice* for specific io classes and priorities inplemented there.

Here is table with numbers for defined class and priorities.

| class name | class number | list of possible priorities |
|:-----------|-------------:|----------------------------:|
| NONE       |0             |0                            |
| RealTime   |1             |0, 1, 2, 3, 4, 5, 6, 7       |
| BestEffort |2             |0, 1, 2, 3, 4, 5, 6, 7       |
| Idle       |3             |0                            |

None is default io prio if process is not ioreniced. And it is the same value and behaviour as BestEffort with priority
of 4.

The lower priority value means higher priority and process will eat more io.

RealTime prio is not available for regular users but for root only.

## Note #2 on ionice operation on Mac OS / Darwin

The only working official method for Apple OS is via taskpolicy utility that comes with Mac OS itself. (Unlimited kernel
API calls requires that calling binary being signed by Apple Inc itself) In context of this program such method is
inefficient and beacuse of that mostly likely io task policy for mac os will not be implemented.

## *BSD porting help needed :)

Indeed. If you have some time and will to add this functionality PRs are welcome!

## Special "thanks"

To security assholes that disrupt workflows and make this world worse.
