import { type ChangeEvent, useEffect, useState } from 'react';
import { Button, Checkbox, FormControlLabel, InputAdornment, Stack, TextField } from '@mui/material';
import {
    Checklist,
    CloudDownload,
    FolderOpen,
    Image,
    ImageOutlined,
    SmartDisplay,
    SmartDisplayOutlined,
} from '@mui/icons-material';
import { OpenDirectory } from '../../wailsjs/go/main/App';
import { useAppStore } from '../stores/app';
import { DialogDownload } from './DialogDownload';

export const DirectoryDownload = () => {
    const directory = useAppStore((state) => state.directory);
    const media = useAppStore((state) => state.media);
    const selectedMedia = useAppStore((state) => state.selectedMedia);
    const setDirectory = useAppStore((state) => state.setDirectory);
    const setSelectedMedia = useAppStore((state) => state.setSelectedMedia);

    const [filter, setFilter] = useState('');
    const [checkboxImage, setCheckboxImage] = useState(false);
    const [checkboxVideo, setCheckboxVideo] = useState(false);
    const [startDownload, setStartDownload] = useState(false);

    const handleFilterChange = (e: ChangeEvent<HTMLInputElement>) => {
        setFilter(e.target.value);
    };

    const handleDirectoryClick = async () => {
        const newDir = await OpenDirectory(directory);
        setDirectory(newDir);
    };

    const handleDownloadClick = () => {
        setStartDownload(true);
    };

    useEffect(() => {
        const selected = media.filter((m) => {
            return (
                (filter.length > 0 && m.Url.toLowerCase().includes(filter.toLowerCase())) ||
                (checkboxImage && m.Type === 0) ||
                (checkboxVideo && m.Type === 1)
            );
        });

        setSelectedMedia(selected);
    }, [checkboxImage, checkboxVideo, filter, setSelectedMedia, media]);

    return (
        <>
            <Stack spacing="0.5em">
                <Stack direction="row" spacing="1em">
                    <TextField
                        fullWidth
                        id="filter"
                        label="Filter by URL"
                        value={filter}
                        autoComplete="off"
                        autoCapitalize="off"
                        size="small"
                        slotProps={{
                            input: {
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <Checklist />
                                    </InputAdornment>
                                ),
                            },
                        }}
                        onChange={handleFilterChange}
                        sx={{ flex: 0.85 }}
                    />

                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={checkboxImage}
                                onClick={() => setCheckboxImage(!checkboxImage)}
                                icon={<ImageOutlined />}
                                checkedIcon={<Image />}
                            />
                        }
                        label="Images"
                        sx={{ flex: 0.075 }}
                    />

                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={checkboxVideo}
                                onClick={() => setCheckboxVideo(!checkboxVideo)}
                                icon={<SmartDisplayOutlined />}
                                checkedIcon={<SmartDisplay />}
                            />
                        }
                        label="Videos"
                        sx={{ flex: 0.075 }}
                    />
                </Stack>

                <Stack direction="row" spacing="1em" style={{ marginTop: '0.75em' }}>
                    <TextField
                        fullWidth
                        id="directory"
                        label="Save directory"
                        value={directory}
                        size="small"
                        slotProps={{
                            input: {
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <FolderOpen />
                                    </InputAdornment>
                                ),
                            },
                        }}
                        onClick={handleDirectoryClick}
                        sx={{ flex: 0.85 }}
                    />

                    <Button
                        id="query"
                        variant="contained"
                        startIcon={<CloudDownload />}
                        onClick={handleDownloadClick}
                        disabled={selectedMedia.length === 0}
                        sx={{ flex: 0.15 }}
                    >
                        Download
                    </Button>
                </Stack>
            </Stack>

            {startDownload && <DialogDownload open={true} onClose={() => setStartDownload(false)} />}
        </>
    );
};
