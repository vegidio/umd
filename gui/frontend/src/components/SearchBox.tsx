import { type ChangeEvent, useCallback, useState } from 'react';
import { Button, IconButton, InputAdornment, Stack, TextField } from '@mui/material';
import { Cookie, Public, Search, Settings as SettingsIcon } from '@mui/icons-material';
import { enqueueSnackbar } from 'notistack';
import { QueryMedia } from '../../wailsjs/go/main/App';
import { BrowserOpenURL } from '../../wailsjs/runtime';
import { useAppStore } from '../stores/app';
import { useSettingsStore } from '../stores/settings';
import { DialogSettings } from './DialogSettings';

export const SearchBox = () => {
    const directory = useAppStore((state) => state.directory);
    const setIsQuerying = useAppStore((state) => state.setIsQuerying);
    const setMedia = useAppStore((state) => state.setMedia);

    const deep = useSettingsStore((state) => state.deepSearch);
    const noCache = useSettingsStore((state) => state.ignoreCache);
    const enableTelemetry = useSettingsStore((state) => state.enableTelemetry);
    const cookiesType = useSettingsStore((state) => state.cookiesType);
    const cookiesPath = useSettingsStore((state) => state.cookiesPath);

    const [url, setUrl] = useState('');
    const [limit, setLimit] = useState(99_999);
    const [openSettings, setOpenSettings] = useState(-1);

    const action = useCallback(
        () => (
            <Button
                variant="outlined"
                color="inherit"
                size="small"
                endIcon={<Cookie />}
                onClick={() => setOpenSettings(1)}
            >
                Settings
            </Button>
        ),
        [],
    );

    const handleUrlChange = (e: ChangeEvent<HTMLInputElement>) => {
        setUrl(e.target.value);
    };

    const handleLimitChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (e.target.value === '') {
            setLimit(1);
            return;
        }

        const value = Number.parseInt(e.target.value, 10);
        setLimit(value < 1 ? 1 : value);
    };

    const handleQueryClick = async () => {
        setIsQuerying(true);

        try {
            const media = await QueryMedia(
                url,
                directory,
                limit,
                deep,
                noCache,
                cookiesType,
                cookiesPath,
                enableTelemetry,
            );

            setMedia(media);
        } catch (e: any) {
            if (e.endsWith('cookies')) {
                enqueueSnackbar(e.replace('the', 'The'), {
                    variant: 'error',
                    action,
                });
            } else {
                enqueueSnackbar('Error querying the media from this URL', { variant: 'error' });
            }
        } finally {
            setIsQuerying(false);
        }
    };

    return (
        <>
            <Stack id="search-box" direction="row" spacing="1em">
                <TextField
                    id="url"
                    label="Enter a URL"
                    value={url}
                    size="small"
                    autoComplete="off"
                    autoCapitalize="off"
                    slotProps={{
                        input: {
                            startAdornment: (
                                <InputAdornment position="start">
                                    <Public />
                                </InputAdornment>
                            ),
                        },
                    }}
                    onChange={handleUrlChange}
                    sx={{ flex: 0.72 }}
                />

                <TextField
                    id="limit"
                    label="Limit results"
                    type="number"
                    value={limit}
                    size="small"
                    onChange={handleLimitChange}
                    sx={{ flex: 0.14 }}
                />

                <Button
                    id="query"
                    variant="outlined"
                    startIcon={<Search />}
                    disabled={url.trim() === ''}
                    onClick={handleQueryClick}
                    sx={{ flex: 0.14 }}
                >
                    Query
                </Button>

                <IconButton onClick={() => setOpenSettings(0)}>
                    <SettingsIcon />
                </IconButton>
            </Stack>

            {openSettings !== -1 && (
                <DialogSettings open={true} onClose={() => setOpenSettings(-1)} tabIndex={openSettings} />
            )}
        </>
    );
};
