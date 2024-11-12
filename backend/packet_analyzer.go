package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "sync"
    "time"
    "github.com/rs/cors"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
)

// Add this type for packet information
type PacketInfo struct {
    Timestamp    string `json:"timestamp"`
    EthernetInfo string `json:"ethernetInfo"`
    IPv4Info     string `json:"ipv4Info,omitempty"`
    TCPInfo      string `json:"tcpInfo,omitempty"`
    UDPInfo      string `json:"udpInfo,omitempty"`
}

// Add the handleInterfaces function
func handleInterfaces(w http.ResponseWriter, r *http.Request) {
    interfaces, err := pcap.FindAllDevs()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(interfaces)
}

func main() {
    // Initialize channels and state
    captureQuit := make(chan bool)
    var currentHandle *pcap.Handle
    var packets []PacketInfo
    var packetsLock sync.Mutex
    var isCapturing bool

    // Create router
    mux := http.NewServeMux()
    
    // Set up routes
    mux.HandleFunc("/interfaces", handleInterfaces)
    
    mux.HandleFunc("/capture", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Decode request body
        var req struct {
            Interface string `json:"interface"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Stop existing capture if any
        if isCapturing {
            captureQuit <- true
            currentHandle.Close()
        }

        // Start new capture
        handle, err := pcap.OpenLive(req.Interface, 65535, true, pcap.BlockForever)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        currentHandle = handle
        isCapturing = true

        // Start capture in goroutine
        go func() {
            packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
            for {
                select {
                case <-captureQuit:
                    return
                case packet := <-packetSource.Packets():
                    packetsLock.Lock()
                    packets = append([]PacketInfo{parsePacket(packet)}, packets...)
                    if len(packets) > 100 { // Keep only last 100 packets
                        packets = packets[:100]
                    }
                    packetsLock.Unlock()
                }
            }
        }()

        w.WriteHeader(http.StatusOK)
    })

    mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        if isCapturing {
            captureQuit <- true
            currentHandle.Close()
            isCapturing = false
        }
        w.WriteHeader(http.StatusOK)
    })

    mux.HandleFunc("/packets", func(w http.ResponseWriter, r *http.Request) {
        packetsLock.Lock()
        json.NewEncoder(w).Encode(packets)
        packetsLock.Unlock()
    })

    // Setup CORS
    handler := cors.Default().Handler(mux)

    // Start server
    fmt.Println("Server starting at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}

func parsePacket(packet gopacket.Packet) PacketInfo {
    info := PacketInfo{
        Timestamp: time.Now().Format("15:04:05.000"),
    }

    // Parse Ethernet
    if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
        eth := ethernetLayer.(*layers.Ethernet)
        info.EthernetInfo = fmt.Sprintf("Ethernet: %s → %s (%s)", 
            eth.SrcMAC, eth.DstMAC, eth.EthernetType)
    }

    // Parse IPv4
    if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
        ip := ipLayer.(*layers.IPv4)
        info.IPv4Info = fmt.Sprintf("IPv4: %s → %s", ip.SrcIP, ip.DstIP)
    }

    // Parse TCP
    if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
        tcp := tcpLayer.(*layers.TCP)
        info.TCPInfo = fmt.Sprintf("TCP: %d → %d [%s]", 
            tcp.SrcPort, tcp.DstPort, tcpFlags(tcp))
    }

    // Parse UDP
    if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
        udp := udpLayer.(*layers.UDP)
        info.UDPInfo = fmt.Sprintf("UDP: %d → %d", udp.SrcPort, udp.DstPort)
    }

    return info
}

func tcpFlags(tcp *layers.TCP) string {
    flags := ""
    if tcp.SYN { flags += "SYN " }
    if tcp.ACK { flags += "ACK " }
    if tcp.PSH { flags += "PSH " }
    if tcp.RST { flags += "RST " }
    if tcp.FIN { flags += "FIN " }
    return strings.TrimSpace(flags)
}