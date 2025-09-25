import { Divider, Stack, Typography, useTheme } from '@mui/material';
import './InfoRow.css';
import { useAppStore } from '../stores/app';

export const InfoRow = () => {
    const theme = useTheme();
    const extractorName = useAppStore((state) => state.extractorName);
    const extractorType = useAppStore((state) => state.extractorType);
    const extractorTypeName = useAppStore((state) => state.extractorTypeName);

    return (
        <Stack
            id="info-bar"
            textAlign="center"
            direction="row"
            divider={<Divider orientation="vertical" />}
            sx={{ borderColor: theme.palette.divider, borderWidth: 1, borderStyle: 'solid' }}
        >
            <Typography variant="body2" sx={{ flex: 0.5 }}>
                Site name: <strong>{extractorName || '-'}</strong>
            </Typography>

            <Typography variant="body2" sx={{ flex: 0.5 }}>
                Source type:{' '}
                <strong>
                    {extractorType || '-'}
                    {extractorTypeName && ` (${extractorTypeName})`}
                </strong>
            </Typography>
        </Stack>
    );
};
