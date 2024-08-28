# mkcidr

Make a CIDR mask given a set of IPs

## Usage
```
$ mkcidr -h
Make CIDR notation from a list of IPs, https://github.com/pschou/mkcidr version: 0.1.20240828.1318

mkcidr [flag] [IPs...]
  -a, --base        Allow base as a valid address
  -z, --broadcast   Allow broadcast as a valid address
  -d, --debug       Enable debug
  -x, --extra int   Allow for N number of extra addresses in subnet
  -h, --help        Show this usage
  -o, --options     Show additional subnet options
```


## Example

Note the output format is: CIDR  Range  Mask  Size
```
$ mkcidr 1.2.3.224 1.2.3.222 1.2.4.5
1.2.0.0/21  1.2.0.1-1.2.7.254  0.0.7.255/255.255.248.0  2046
```
