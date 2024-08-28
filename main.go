package main

import (
	"fmt"
	"net/netip"
	"os"

	flag "github.com/spf13/pflag"
)

var version = "STATIC"

func main() {
	extra := flag.IntP("extra", "x", 0, "Allow for N number of extra addresses in subnet")
	debug := flag.BoolP("debug", "d", false, "Enable debug")

	allowBase := flag.BoolP("base", "a", false, "Allow base as a valid address")
	allowBroadcast := flag.BoolP("broadcast", "z", false, "Allow broadcast as a valid address")

	options := flag.BoolP("options", "o", false, "Show additional subnet options")
	help := flag.BoolP("help", "h", false, "Show this usage")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Make CIDR notation from a list of IPs, https://github.com/pschou/mkcidr version: %s\n\nmkcidr [flag] [IPs...]\n", version)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		os.Exit(0)
	}

	list := make([]netip.Addr, flag.NArg())
	for i, ip := range flag.Args() {
		t, err := netip.ParseAddr(ip)
		if err != nil {
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid IP: %q  %s\n", ip, err)
			os.Exit(1)
		}
		if i > 0 && t.Is4() != list[0].Is4() {
			fmt.Fprintf(os.Stderr, "Mismatching IP space: %q\n", ip)
			os.Exit(1)
		}
		list[i] = t
	}

	bits := len(list[0].AsSlice()) * 8
	c := bits - 1

nextCIDR:
	for ; c > 0; c-- {
		mask, _ := list[0].Prefix(c)
		bcast := BroadcastAddr(mask)
		if *debug {
			fmt.Println("Testing mask", mask)
		}
		passed := true
		for _, ip := range list {
			passed = passed && mask.Contains(ip)
			if *debug {
				fmt.Printf("  ip: %s  contained: %v\n", ip, mask.Contains(ip))
			}
			if ip.Compare(mask.Addr()) == 0 && !*allowBase {
				if *debug {
					fmt.Printf("  ip: %s  is base for subnet, skipping\n", ip)
				}
				passed = false
			}
			if ip.Compare(bcast) == 0 && !*allowBroadcast {
				if *debug {
					fmt.Printf("  ip: %s  is broadcast for subnet, skipping\n", ip)
				}
				passed = false
			}
			if !passed {
				continue nextCIDR
			}
		}
		used := len(list)
		var size int
		if bits-c <= 42 {
			size = 1 << (bits - c)
			if !*allowBroadcast {
				size--
			}
			if !*allowBase {
				size--
			}
			if *debug {
				fmt.Printf("  used: %d  size: %d  remaining: %d\n", used, size, size-used-2)
			}
			if size-used > 0 && size-used < *extra {
				passed = false
			}
		} else {
			size = -1
		}

		a := mask.Addr()
		if !*allowBase {
			a = a.Next()
		}
		z := bcast
		if !*allowBroadcast {
			z = z.Prev()
		}

		if passed {
			if size >= 0 {
				fmt.Printf("%s  %s-%s  %s  %d\n", mask, a, z, bcast, size)
			} else {
				fmt.Printf("%s  %s-%s  %s\n", mask, a, z, bcast)
			}

			if !*options {
				break
			}
		}
	}
}

// BroadcastAddr returns the last address in the given network, or the broadcast address.
func BroadcastAddr(n netip.Prefix) netip.Addr {
	slice := n.Addr().AsSlice()
	//fmt.Println("bits", slice)
	for b := len(slice)*8 - 1; b >= n.Bits(); b-- {
		//fmt.Println("b", b, (b-1)/8, slice[(b-1)/8], 1<<(7-b%8), b%8)
		slice[b/8] |= 1 << (7 - b%8)
		//fmt.Println(" ==", slice)
	}
	a, _ := netip.AddrFromSlice(slice)
	return a
}
