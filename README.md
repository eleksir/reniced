# Reniced

Helper utility intended to run as a service or as background process (say in tmux) that can make system more respositive
under certain type of load.

## Why it was made

In general case os can handle load well, fairly sharing resources among processes. But in some cases it makes it poorly.
As an example it can be ill-minded pseudo-security stuff like antiviruses or mediaanalysers that searches forbidden
content on your mac (due to gov.polices that apple have to accept). Or it can be backup that crawls fs in background.

To help resolve such situation this utility was made. Some malicious processes can be reniced to lower nicelevels (say,
19) and for some good processes nicelevel can be rised to higher values(-19 if reniced run as superuser or 0 if as
normal).

In some critical cases process can be killed (if you want to know more about it see "man kill" in terminal).

## How to build, configure, run it

You will need Go lang 1.21 or newer, gnu make, terminal and possibly internets to build this utility.

To build reniced just type in src dir

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

if you can sudo. Thats it.

## Special "thanks"

To security assholes that disrupt workflows and make this world worse.
