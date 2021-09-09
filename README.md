# FTPUP

A tiny server serving files from FTP over HTTP. Written as a workaround for a
(hopefully temporary) networking issue.

## Usage

```
$ ftpup -h
Usage of ftpup:
  -T duration
        ftp timeout (default 10s)
  -l string
        hostport to listen on (default "localhost:15201")
  -p string
        ftp host to proxy to (default "ftp.ncbi.nlm.nih.gov:21")
  -u value
        username and password (default anonymous:anonymous)
```

After starting the server, you can request files by their path:

```
$ curl http://localhost:15201/pub/pmc/readme.txt
On March 18, 2019, PMC will no longer provide bulk packages of Open Access (OA)
Subset text and XML at the top level directory of the FTP Service. These files
were superseded in August 2016 by the Commercial Use and Non-Commercial Use
bulk packages located in the oa_bulk subdirectory. One set comprises articles
that may be used for commercial purposes (the Commercial Use Collection); the
other contains articles that can be used only for non-commercial purposes.
Anyone planning to use OA subset content for non-commercial purposes will need
to download both ?non_comm_use.*.tar.gz? and ?comm_use.*.tar.gz? to access the
complete collection. See the Open Access Subset page
(https://www.ncbi.nlm.nih.gov/pmc/tools/openftlist/ for additional details.
Questions should be directed to pubmedcentral@ncbi.nlm.nih.gov.

See http://www.ncbi.nlm.nih.gov/pmc/tools/ftp/
```
