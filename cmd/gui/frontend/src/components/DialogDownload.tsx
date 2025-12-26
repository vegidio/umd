import { useEffect, useState } from 'react'
import {
    Box,
    Button,
    CircularProgress,
    type CircularProgressProps,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    LinearProgress,
    Link,
    Paper,
    Stack,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography
} from '@mui/material'
import { enqueueSnackbar } from 'notistack'
import { CancelDownloads, StartDownload } from '../../wailsjs/go/main/App'
import { useAppStore } from '../stores/app'
import { TypeCell } from './MediaList'

type Props = {
    open: boolean;
    onClose: () => void;
};

export const DialogDownload = ({ open, onClose }: Props) => {
    const currentDownloads = useAppStore((state) => state.currentDownloads);
    const directory = useAppStore((state) => state.directory);
    const downloadedMedia = useAppStore((state) => state.downloadedMedia);
    const progress = useAppStore((state) => state.progress);
    const selectedMedia = useAppStore((state) => state.selectedMedia);
    const setProgress = useAppStore((state) => state.setProgress);

    const [isDownloading, setIsDownloading] = useState(true);

    // biome-ignore lint/correctness/useExhaustiveDependencies: should fire only once
    useEffect(() => {
        let mounted = true;
        (async () => {
            await StartDownload(selectedMedia, directory, 5);
            setIsDownloading(false);
            enqueueSnackbar('Download completed', { variant: 'success' });
        })();

        return () => {
            mounted = false;
        };
    }, []);

    useEffect(() => {
        const percentage = (downloadedMedia * 100) / selectedMedia.length;
        setProgress(Number.isNaN(percentage) ? 0 : percentage);
    }, [downloadedMedia, selectedMedia, setProgress]);

    const onCancelDownload = async () => {
        await CancelDownloads();
        setIsDownloading(false);
        enqueueSnackbar('Download cancelled', { variant: 'warning' });
    };

    return (
        <Dialog
            open={open}
            onClose={(_, reason) => {
                if (reason !== 'backdropClick') {
                    onClose();
                }
            }}
            maxWidth={false}
            disableEscapeKeyDown
            fullWidth
            slotProps={{
                paper: {
                    sx: {
                        width: 768,
                    },
                },
            }}
        >
            <DialogTitle sx={{ m: 0, p: 1.5 }} id="customized-dialog-title">
                Downloading
            </DialogTitle>

            <DialogContent dividers>
                <TableContainer component={Paper}>
                    <Table sx={{ tableLayout: 'fixed' }} size="small">
                        <TableHead>
                            <TableRow>
                                <TableCell sx={{ width: '12%' }}>Progress</TableCell>
                                <TableCell sx={{ width: '68%' }}>URL</TableCell>
                                <TableCell sx={{ width: '20%' }}>Type</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {currentDownloads.map((response) => {
                                return (
                                    <TableRow
                                        key={response.Request?.Url}
                                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                                    >
                                        <TableCell component="th" scope="row" align="center">
                                            <CircularProgressWithLabel value={response.Progress * 100} />
                                        </TableCell>
                                        <TableCell
                                            sx={{
                                                maxWidth: '200px',
                                                overflow: 'hidden',
                                                textOverflow: 'ellipsis',
                                                whiteSpace: 'nowrap',
                                            }}
                                        >
                                            <Link href="#">{response.Request?.Url}</Link>
                                        </TableCell>
                                        <TableCell>
                                            <TypeCell type={getMediaType(response.Request?.FilePath ?? '')} />
                                        </TableCell>
                                    </TableRow>
                                );
                            })}
                        </TableBody>
                    </Table>
                </TableContainer>

                <Stack direction="row" spacing="1em" alignItems="center" sx={{ marginTop: '2em' }}>
                    <Typography variant="body2" align="left" sx={{ flex: 0.1 }}>
                        {downloadedMedia} / {selectedMedia.length}
                    </Typography>

                    <LinearProgress variant="determinate" value={progress} sx={{ flex: 0.8 }} />

                    <Typography variant="body2" align="right" sx={{ flex: 0.1 }}>
                        {progress.toFixed(1)}%
                    </Typography>
                </Stack>
            </DialogContent>

            <DialogActions>
                <Button
                    color={isDownloading ? 'error' : 'primary'}
                    onClick={() => {
                        if (isDownloading) {
                            onCancelDownload();
                        } else {
                            onClose();
                        }
                    }}
                    sx={{ fontWeight: 'bold' }}
                >
                    {isDownloading ? 'Cancel' : 'Close'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

const CircularProgressWithLabel = (props: CircularProgressProps & { value: number }) => {
    return (
        <Box sx={{ position: 'relative', display: 'inline-flex' }}>
            <CircularProgress variant="determinate" {...props} />
            <Box
                sx={{
                    top: 0,
                    left: 0,
                    bottom: 0,
                    right: 0,
                    position: 'absolute',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                }}
            >
                <Typography
                    variant="caption"
                    component="div"
                    sx={{ color: 'text.secondary' }}
                >{`${Math.round(props.value)}%`}</Typography>
            </Box>
        </Box>
    );
};

const getExtension = (filePath: string) => {
    const idx = filePath.lastIndexOf('.');

    if (idx === -1 || idx === filePath.length - 1) {
        return idx === filePath.length - 1 ? '' : filePath;
    }

    return filePath.slice(idx + 1);
};

const getMediaType = (filePath: string) => {
    const ext = getExtension(filePath).toLowerCase();

    switch (ext) {
        case 'avif':
        case 'gif':
        case 'jpg':
        case 'jpeg':
        case 'png':
        case 'webp':
            return 0; // Image

        case 'gifv':
        case 'm4v':
        case 'mkv':
        case 'mov':
        case 'mp4':
        case 'webm':
            return 1; // Video

        default:
            return -1; // Unknown
    }
};
