package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	extra := flag.Int("x", 0, "Allow for N number of extra addresses in subnet")
	debug := flag.Bool("d", false, "Enable debug")
	omitBase := flag.Bool("omit-base", true, "Disallow base")
	omitBroadcast := flag.Bool("omit-broadcast", true, "Disallow broadcast")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Make CIDR from IPs, https://github.com/pschou/mkcidr\n\nmkcidr [flag] [IPs...]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	list := make([]net.IP, len(flag.Args()))
	for i, ip := range flag.Args() {
		t := net.ParseIP(ip)
		if t == nil {
			fmt.Fprintf(os.Stderr, "Invalid IP: %q\n", ip)
			os.Exit(1)
		}
		if i > 0 && len(t) != len(list[0]) {
			fmt.Fprintf(os.Stderr, "Mismatching IP space: %q\n", ip)
			os.Exit(1)
		}
		list[i] = t
	}

	var c, bits int

	if list[0].To4() == nil {
		c = 127
		bits = 128
	} else {
		c = 31
		bits = 32
	}
nextCIDR:
	for ; c > 0; c-- {
		_, mask, err := net.ParseCIDR(fmt.Sprintf("%s/%d", list[0], c))
		if err != nil {
			log.Fatal("err", err)
		}
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
			if ip.Equal(mask.IP) && *omitBase {
				if *debug {
					fmt.Printf("  ip: %s  is base for subnet, skipping\n", ip)
				}
				passed = false
			}
			if ip.Equal(bcast) && *omitBroadcast {
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
		size := 1 << (bits - c)
		if *omitBroadcast {
			size--
		}
		if *omitBase {
			size--
		}
		if *debug {
			fmt.Printf("  used: %d  size: %d  remaining: %d\n", used, size, size-used-2)
		}
		if size-used > 0 && size-used < *extra {
			passed = false
		}
		if passed {
			fmt.Printf("%s  %s-%s\n", mask, mask.IP, bcast)
			break
		}
	}
}

// BroadcastAddr returns the last address in the given network, or the broadcast address.
func BroadcastAddr(n *net.IPNet) net.IP {
	// The golang net package doesn't make it easy to calculate the broadcast address. :(
	broadcast := net.IP(make([]byte, len(n.IP)))
	for i := 0; i < len(n.IP); i++ {
		broadcast[i] = n.IP[i] | ^n.Mask[i]
	}
	return broadcast
}
