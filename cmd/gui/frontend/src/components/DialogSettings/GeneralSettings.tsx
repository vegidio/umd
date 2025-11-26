import { type ChangeEvent, forwardRef, useImperativeHandle, useState } from 'react';
import { Checkbox, DialogContent, FormControlLabel, FormGroup } from '@mui/material';
import type { ChildSave } from './index';
import type { TabPanelProps } from './TabPanelProps';
import { useSettingsStore } from '../../stores/settings';

export const GeneralSettings = forwardRef<ChildSave, TabPanelProps>(({ value, index }, ref) => {
    const deepSearch = useSettingsStore((state) => state.deepSearch);
    const ignoreCache = useSettingsStore((state) => state.ignoreCache);
    const enableTelemetry = useSettingsStore((state) => state.enableTelemetry);
    const setDeepSearch = useSettingsStore((state) => state.setDeepSearch);
    const setIgnoreCache = useSettingsStore((state) => state.setIgnoreCache);
    const setEnableTelemetry = useSettingsStore((state) => state.setEnableTelemetry);

    const [localDeepSearch, setLocalDeepSearch] = useState(deepSearch);
    const [localIgnoreCache, setLocalIgnoreCache] = useState(ignoreCache);
    const [localEnableTelemetry, setLocalEnableTelemetry] = useState(enableTelemetry);

    useImperativeHandle(ref, () => ({
        save() {
            setDeepSearch(localDeepSearch);
            setIgnoreCache(localIgnoreCache);
            setEnableTelemetry(localEnableTelemetry);
        },
    }));

    return (
        <DialogContent sx={{ padding: 2 }} hidden={value !== index}>
            <FormGroup sx={{ rowGap: 2 }}>
                <FormControlLabel
                    sx={{ alignItems: 'flex-start' }}
                    control={
                        <Checkbox
                            checked={localDeepSearch}
                            onChange={(e: ChangeEvent<HTMLInputElement>) => setLocalDeepSearch(e.target.checked)}
                        />
                    }
                    label={
                        <>
                            <strong>Deep Search:</strong> expands the search of unknown URLs, attempting to discover
                            additional media files.
                        </>
                    }
                />
                <FormControlLabel
                    sx={{ alignItems: 'flex-start' }}
                    control={
                        <Checkbox
                            checked={localIgnoreCache}
                            onChange={(e: ChangeEvent<HTMLInputElement>) => setLocalIgnoreCache(e.target.checked)}
                        />
                    }
                    label={
                        <>
                            <strong>Ignore Cache:</strong> this option will bypass previously cached URLs and always
                            fetch new ones.
                        </>
                    }
                />
                <FormControlLabel
                    sx={{ alignItems: 'flex-start' }}
                    control={
                        <Checkbox
                            checked={localEnableTelemetry}
                            onChange={(e: ChangeEvent<HTMLInputElement>) => setLocalEnableTelemetry(e.target.checked)}
                        />
                    }
                    label={
                        <>
                            <strong>Enable Telemetry:</strong> this option will send anonymous data to help improve the
                            app stability.
                        </>
                    }
                />
            </FormGroup>
        </DialogContent>
    );
});
