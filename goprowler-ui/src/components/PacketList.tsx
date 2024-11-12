import React from 'react';
import { 
    Paper, 
    List, 
    ListItem, 
    ListItemText, 
    Typography, 
    Box 
} from '@mui/material';
import { PacketInfo } from '../types/types';

interface Props {
    packets: PacketInfo[];
}

const PacketList: React.FC<Props> = ({ packets }) => {
    return (
        <Paper sx={{ maxHeight: 'calc(100vh - 200px)', overflow: 'auto' }}>
            <List>
                {packets.map((packet, index) => (
                    <ListItem key={index} divider>
                        <Box sx={{ width: '100%' }}>
                            <Typography variant="caption" color="textSecondary">
                                {packet.Timestamp}
                            </Typography>
                            <ListItemText
                                primary={packet.EthernetInfo}
                                secondary={
                                    <React.Fragment>
                                        {packet.IPv4Info && (
                                            <Typography component="span" variant="body2" display="block">
                                                {packet.IPv4Info}
                                            </Typography>
                                        )}
                                        {packet.TCPInfo && (
                                            <Typography component="span" variant="body2" display="block">
                                                {packet.TCPInfo}
                                            </Typography>
                                        )}
                                        {packet.UDPInfo && (
                                            <Typography component="span" variant="body2" display="block">
                                                {packet.UDPInfo}
                                            </Typography>
                                        )}
                                    </React.Fragment>
                                }
                            />
                        </Box>
                    </ListItem>
                ))}
            </List>
        </Paper>
    );
};

export default PacketList;