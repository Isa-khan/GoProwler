import React from 'react';
import { 
    FormControl, 
    InputLabel, 
    Select, 
    MenuItem, 
    Button, 
    Box 
} from '@mui/material';
import { NetworkInterface } from '../types/types';

interface Props {
    interfaces: NetworkInterface[];
    selectedInterface: string;
    isCapturing: boolean;
    onInterfaceChange: (value: string) => void;
    onStartCapture: () => void;
    onStopCapture: () => void;
}

const InterfaceSelector: React.FC<Props> = ({
    interfaces,
    selectedInterface,
    isCapturing,
    onInterfaceChange,
    onStartCapture,
    onStopCapture,
}) => {
    return (
        <Box sx={{ display: 'flex', gap: 2, alignItems: 'center', p: 2 }}>
            <FormControl sx={{ minWidth: 200 }}>
                <InputLabel>Network Interface</InputLabel>
                <Select
                    value={selectedInterface}
                    label="Network Interface"
                    onChange={(e) => onInterfaceChange(e.target.value)}
                    disabled={isCapturing}
                >
                    {interfaces.map((iface) => (
                        <MenuItem key={iface.Name} value={iface.Name}>
                            {iface.Name} - {iface.Description}
                        </MenuItem>
                    ))}
                </Select>
            </FormControl>
            <Button
                variant="contained"
                color={isCapturing ? "error" : "primary"}
                onClick={isCapturing ? onStopCapture : onStartCapture}
                disabled={!selectedInterface}
            >
                {isCapturing ? "Stop Capture" : "Start Capture"}
            </Button>
        </Box>
    );
};

export default InterfaceSelector;