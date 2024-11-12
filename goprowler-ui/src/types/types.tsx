export interface NetworkInterface {
    Name: string;
    Description: string;
}

export interface PacketInfo {
    Timestamp: string;
    EthernetInfo: string;
    IPv4Info?: string;
    TCPInfo?: string;
    UDPInfo?: string;
}

export interface AppState {
    interfaces: NetworkInterface[];
    isCapturing: boolean;
    packets: PacketInfo[];
    selectedInterface: string;
}