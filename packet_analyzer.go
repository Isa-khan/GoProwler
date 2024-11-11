package main

import (
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
)

func main() {
    interfaces, err := pcap.FindAllDevs()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Available interfaces:")
    for i, iface := range interfaces {
        fmt.Printf("%d. %s (%s)\n", i+1, iface.Name, iface.Description)
    }

    deviceName := interfaces[0].Name
    if len(os.Args) > 1 {
        deviceName = os.Args[1]
    }

    handle, err := pcap.OpenLive(deviceName, 65535, true, pcap.BlockForever)
    if err != nil {
        log.Fatal(err)
    }
    defer handle.Close()

    fmt.Printf("Capturing packets on interface: %s\n", deviceName)

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        printPacketInfo(packet)
    }
}

func printPacketInfo(packet gopacket.Packet) {
    ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
    if ethernetLayer != nil {
        eth := ethernetLayer.(*layers.Ethernet)
        fmt.Printf("\nEthernet Frame:\n")
        fmt.Printf("Destination: %s\n", eth.DstMAC)
        fmt.Printf("Source: %s\n", eth.SrcMAC)
        fmt.Printf("Protocol: %v\n", eth.EthernetType)
    }

    ipLayer := packet.Layer(layers.LayerTypeIPv4)
    if ipLayer != nil {
        ip := ipLayer.(*layers.IPv4)
        fmt.Printf("\nIPv4 Packet:\n")
        fmt.Printf("Version: %d\n", ip.Version)
        fmt.Printf("TTL: %d\n", ip.TTL)
        fmt.Printf("Protocol: %d\n", ip.Protocol)
        fmt.Printf("Source IP: %s\n", ip.SrcIP)
        fmt.Printf("Destination IP: %s\n", ip.DstIP)
    }

    // TCP layer
    tcpLayer := packet.Layer(layers.LayerTypeTCP)
    if tcpLayer != nil {
        tcp := tcpLayer.(*layers.TCP)
        fmt.Printf("\nTCP Segment:\n")
        fmt.Printf("Source Port: %d\n", tcp.SrcPort)
        fmt.Printf("Destination Port: %d\n", tcp.DstPort)
        fmt.Printf("Sequence: %d\n", tcp.Seq)
        fmt.Printf("Acknowledgment: %d\n", tcp.Ack)
        fmt.Printf("Flags: URG=%v ACK=%v PSH=%v RST=%v SYN=%v FIN=%v\n",
            tcp.URG, tcp.ACK, tcp.PSH, tcp.RST, tcp.SYN, tcp.FIN)
    }

    udpLayer := packet.Layer(layers.LayerTypeUDP)
    if udpLayer != nil {
        udp := udpLayer.(*layers.UDP)
        fmt.Printf("\nUDP Segment:\n")
        fmt.Printf("Source Port: %d\n", udp.SrcPort)
        fmt.Printf("Destination Port: %d\n", udp.DstPort)
    }

    fmt.Println(strings.Repeat("-", 80))
}