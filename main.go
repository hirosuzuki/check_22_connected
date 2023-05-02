package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func parse_addr_port(s string) (ip net.IP, port int, err error) {
	// Example "0100007F:0050" => 127.0.0.1:80
	ip = make(net.IP, net.IPv4len)
	err = fmt.Errorf("address:port Format error")
	if len(s) != 13 || s[8] != ':' {
		return
	}
	for i := 0; i < 4; i++ {
		if v, e := strconv.ParseInt(s[6-i*2:8-i*2], 16, 0); e != nil {
			return
		} else {
			ip[i] = byte(v)
		}
	}
	if v, e := strconv.ParseInt(s[9:13], 16, 0); e != nil {
		return
	} else {
		port = int(v)
	}
	err = nil
	return
}

func count_tcp_connection(port int) (int, error) {
	file, err := os.Open("/proc/net/tcp")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		_, local_port, err := parse_addr_port(fields[1])
		if err != nil {
			continue
		}
		peer_addr, _, err := parse_addr_port(fields[2])
		if err != nil {
			continue
		}
		status, err := strconv.ParseInt(fields[3], 16, 0)
		if err != nil {
			continue
		}
		if status == 1 && local_port == port && peer_addr.String() != "0.0.0.0" {
			count++
		}
	}
	return count, nil
}

// SSHポート(22番)で接続中のTCPコネクション数を標準出力に出力する
func main() {
	count, err := count_tcp_connection(22)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(count)
}
