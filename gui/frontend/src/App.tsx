import { Button, Stack } from '@mui/material';
import { SnackbarProvider, enqueueSnackbar } from 'notistack';
import React, { useEffect } from 'react';
import { GetHomeDirectory, IsOutdated } from '../wailsjs/go/main/App';
import { DirectoryDownload, InfoRow, Loading, MediaList, SearchBox } from './components';
import './App.css';
import DownloadIcon from '@mui/icons-material/Download';
import { BrowserOpenURL } from '../wailsjs/runtime';
import { useAppStore } from './stores/app';

const action = () => (
    <Button
        variant="outlined"
        color="inherit"
        size="small"
        endIcon={<DownloadIcon />}
        onClick={() => BrowserOpenURL('https://github.com/vegidio/umd/releases')}
    >
        Download
    </Button>
);

const App = () => {
    const setDirectory = useAppStore((s) => s.setDirectory);

    useEffect(() => {
        const getHomeDirectory = async () => {
            let dir = localStorage.getItem('lastDirectory');
            if (!dir) dir = await GetHomeDirectory();
            setDirectory(dir);
        };

        const checkVersion = async () => {
            const isOutdated = await IsOutdated();
            if (isOutdated)
                enqueueSnackbar('A new version of UMD is available. Please update', {
                    variant: 'warning',
                    autoHideDuration: 10_000,
                    preventDuplicate: true,
                    action,
                });
        };

        Promise.all([checkVersion(), getHomeDirectory()]);
    }, [setDirectory]);

    return (
        <SnackbarProvider maxSnack={3}>
            <Stack id="app" spacing="1em">
                <SearchBox />

                <InfoRow />

                <MediaList />

                <DirectoryDownload />
            </Stack>

            <Loading />
        </SnackbarProvider>
    );
};

export default App;
