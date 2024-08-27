# mkcidr

Make a CIDR mask given a set of IPs

## Usage
```
$ ./mkcidr -h
Make CIDR from IPs, https://github.com/pschou/mkcidr

mkcidr [flag] [IPs...]
  -d    Enable debug
  -omit-base
        Disallow base (default true)
  -omit-broadcast
        Disallow broadcast (default true)
  -x int
        Allow for N number of extra addresses in subnet
```


## Example
```
mkcidr  1.2.3.4 1.2.3.67
1.2.3.0/25  1.2.3.0-1.2.3.127
```
