# GoProwler 🕵️‍♂️🔍

This is the third time I am making a packet analyzer, first was in Python, then C, now Go. Main ideas are the same, just different implementations.

GoProwler is a network packet analyzer written in Go that captures and analyzes network traffic in real-time. It provides detailed information about Ethernet frames, IPv4 packets, and TCP/UDP segments.

## Features

- 📦 Real-time packet capture
- 🔍 Detailed packet analysis including:
  - Ethernet frame information
  - IPv4 packet details
  - TCP segment analysis
  - UDP segment analysis
- 🖥️ Network interface selection
- 📊 Human-readable output format

## Prerequisites

Before running GoProwler, ensure you have the following installed:

1. Go (1.16 or later)
2. libpcap
   ```bash
   # On macOS
   brew install libpcap

   # On Ubuntu/Debian
   sudo apt-get install libpcap-dev

   # On CentOS/RHEL
   sudo yum install libpcap-devel
   ```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/GoProwler.git
   cd GoProwler
   ```

2. Initialize Go module:
   ```bash
   go mod init goprowler
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

Run GoProwler with sudo privileges (required for packet capture):

Use default network interface
```bash
sudo go run packet_analyzer.go
```

Specify a network interface
```bash
sudo go run packet_analyzer.go en0
```








