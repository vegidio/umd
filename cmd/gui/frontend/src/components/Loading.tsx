import { Backdrop, Box, Button, CircularProgress, Stack, Typography } from '@mui/material';
import { Cancel } from '@mui/icons-material';
import { StopQuery } from '../../wailsjs/go/main/App';
import { useAppStore } from '../stores/app';
import './Loading.css';

export const Loading = () => {
    const amountQuery = useAppStore((state) => state.amountQuery);
    const extractorName = useAppStore((state) => state.extractorName);
    const isQuerying = useAppStore((state) => state.isQuerying);

    const handleCancelClick = () => {
        StopQuery();
    };

    return (
        <Backdrop open={isQuerying}>
            <Stack id="loading-box" spacing="1em">
                <CircularProgressWithLabel value={amountQuery} />

                <Typography color="textPrimary" variant="body2">
                    Querying media from <strong>{extractorName}</strong>...
                </Typography>

                <Button variant="contained" startIcon={<Cancel />} onClick={handleCancelClick}>
                    Cancel
                </Button>
            </Stack>
        </Backdrop>
    );
};

type CircularProgressProps = {
    value: number;
};

const CircularProgressWithLabel = ({ value }: CircularProgressProps) => {
    return (
        <Box id="circular-progress">
            <CircularProgress color="primary" size="5em" />
            <Typography position="absolute">{value}</Typography>
        </Box>
    );
};
