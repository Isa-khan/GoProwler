import React, { useEffect, useState } from 'react';
import { Container, Typography, Box } from '@mui/material';
import axios from 'axios';
import InterfaceSelector from './components/InterfaceSelector';
import PacketList from './components/PacketList';
import { AppState, NetworkInterface, PacketInfo } from './types/types';

const API_BASE_URL = 'http://localhost:8080';

function App() {
    const [state, setState] = useState<AppState>({
        interfaces: [],
        isCapturing: false,
        packets: [],
        selectedInterface: '',
    });

    useEffect(() => {
        fetchInterfaces();
    }, []);

    useEffect(() => {
        let interval: NodeJS.Timeout;
        if (state.isCapturing) {
            interval = setInterval(fetchPackets, 1000);
        }
        return () => {
            if (interval) {
                clearInterval(interval);
            }
        };
    }, [state.isCapturing]);

    const fetchInterfaces = async () => {
        try {
            const response = await axios.get<NetworkInterface[]>(`${API_BASE_URL}/interfaces`);
            setState(prev => ({ ...prev, interfaces: response.data }));
        } catch (error) {
            console.error('Error fetching interfaces:', error);
        }
    };

    const fetchPackets = async () => {
        try {
            const response = await axios.get<PacketInfo[]>(`${API_BASE_URL}/packets`);
            setState(prev => ({ ...prev, packets: response.data }));
        } catch (error) {
            console.error('Error fetching packets:', error);
        }
    };

    const handleInterfaceChange = (value: string) => {
        setState(prev => ({ ...prev, selectedInterface: value }));
    };

    const handleStartCapture = async () => {
        try {
            await axios.post(`${API_BASE_URL}/capture`, {
                interface: state.selectedInterface
            });
            setState(prev => ({ ...prev, isCapturing: true }));
        } catch (error) {
            console.error('Error starting capture:', error);
        }
    };

    const handleStopCapture = async () => {
        try {
            await axios.post(`${API_BASE_URL}/stop`);
            setState(prev => ({ ...prev, isCapturing: false }));
        } catch (error) {
            console.error('Error stopping capture:', error);
        }
    };

    return (
        <Container maxWidth="lg">
            <Box sx={{ my: 4 }}>
                <Typography variant="h3" component="h1" gutterBottom align="center">
                    GoProwler 🕵️‍♂️
                </Typography>
                <InterfaceSelector
                    interfaces={state.interfaces}
                    selectedInterface={state.selectedInterface}
                    isCapturing={state.isCapturing}
                    onInterfaceChange={handleInterfaceChange}
                    onStartCapture={handleStartCapture}
                    onStopCapture={handleStopCapture}
                />
                <PacketList packets={state.packets} />
            </Box>
        </Container>
    );
}

export default App;