import { type Dispatch, forwardRef, type SetStateAction, useImperativeHandle, useState } from 'react';
import {
    Button,
    DialogContent,
    Divider,
    FormControl,
    FormControlLabel,
    InputAdornment,
    Link,
    Radio,
    RadioGroup,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import { Cookie, FolderOpen } from '@mui/icons-material';
import type { ChildSave } from './index';
import type { TabPanelProps } from './TabPanelProps';
import { OpenCookiesPath } from '../../../wailsjs/go/main/App';
import { BrowserOpenURL } from '../../../wailsjs/runtime';
import { useSettingsStore } from '../../stores/settings';

export const CookiesSettings = forwardRef<ChildSave, TabPanelProps>(({ value, index }, ref) => {
    const cookiesType = useSettingsStore((state) => state.cookiesType);
    const cookiesPath = useSettingsStore((state) => state.cookiesPath);
    const setCookiesType = useSettingsStore((state) => state.setCookiesType);
    const setCookiesPath = useSettingsStore((state) => state.setCookiesPath);

    const [localCookiesType, setLocalCookiesType] = useState(cookiesType);
    const [localCookiesPath, setLocalCookiesPath] = useState(cookiesPath);

    useImperativeHandle(ref, () => ({
        save() {
            setCookiesType(localCookiesType);

            if (localCookiesType === 'manual') {
                setCookiesPath(localCookiesPath);
            } else {
                setCookiesPath('');
            }
        },
    }));

    return (
        <DialogContent sx={{ px: 2, py: 1.5 }} hidden={value !== index}>
            <FormControl>
                <RadioGroup
                    aria-labelledby="cookies-import"
                    defaultValue={cookiesType}
                    name="radio-buttons-group"
                    onChange={(e) => setLocalCookiesType(e.target.value)}
                >
                    <FormControlLabel value="disabled" control={<Radio />} label="Disabled" />

                    <Divider sx={{ my: 1 }} />

                    <FormControlLabel value="automatic" control={<Radio />} label="Automatic" />
                    <AutomaticCookies disabled={localCookiesType !== 'automatic'} />

                    <Divider sx={{ mt: 2.5, mb: 1.5 }} />

                    <FormControlLabel value="manual" control={<Radio />} label="Manual" />
                    <ManualCookies
                        cookiesPath={localCookiesPath}
                        setCookiesPath={setLocalCookiesPath}
                        disabled={localCookiesType !== 'manual'}
                    />
                </RadioGroup>
            </FormControl>
        </DialogContent>
    );
});

const AutomaticCookies = ({ disabled }: { disabled: boolean }) => {
    return (
        <Stack
            direction="row"
            spacing="2em"
            alignItems="center"
            sx={{
                opacity: disabled ? 0.3 : 1,
                pointerEvents: disabled ? 'none' : 'auto',
            }}
        >
            <Typography color="textPrimary" variant="body2">
                This option will try to automatically import the cookies saved in your browser, so make sure you're
                logged-in to the site that you want to download media from before selecting it.
            </Typography>
        </Stack>
    );
};

type ManualCookiesProps = {
    cookiesPath: string;
    setCookiesPath: Dispatch<SetStateAction<string>>;
    disabled: boolean;
};

const ManualCookies = ({ cookiesPath, setCookiesPath, disabled }: ManualCookiesProps) => {
    const handleCookiesClick = async () => {
        const newPath = await OpenCookiesPath(cookiesPath);
        setCookiesPath(newPath);
    };

    return (
        <Stack
            direction="column"
            spacing="1em"
            sx={{
                opacity: disabled ? 0.3 : 1,
                pointerEvents: disabled ? 'none' : 'auto',
            }}
        >
            <Typography color="textPrimary" variant="body2">
                Download your cookies to a{' '}
                <Link href="#" onClick={() => BrowserOpenURL('https://www.google.com/search?q=cookies.txt')}>
                    cookies.txt
                </Link>{' '}
                file and select its path below:
            </Typography>

            <Stack direction="row" spacing="1em" alignItems="center">
                <TextField
                    fullWidth
                    id="cookies-file"
                    label="Cookies file path"
                    value={cookiesPath}
                    size="small"
                    slotProps={{
                        input: {
                            startAdornment: (
                                <InputAdornment position="start">
                                    <Cookie />
                                </InputAdornment>
                            ),
                        },
                    }}
                    sx={{ flex: 0.85 }}
                />

                <Button
                    id="select"
                    variant="contained"
                    startIcon={<FolderOpen />}
                    onClick={handleCookiesClick}
                    sx={{ flex: 0.15 }}
                >
                    Select
                </Button>
            </Stack>
        </Stack>
    );
};
