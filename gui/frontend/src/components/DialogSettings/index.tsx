import { type RefObject, type SyntheticEvent, useRef, useState } from 'react';
import { Box, Button, Dialog, DialogActions, DialogTitle, Divider, IconButton, Tab, Tabs } from '@mui/material';
import { Close } from '@mui/icons-material';
import { CookiesSettings } from './CookiesSettings';
import { GeneralSettings } from './GeneralSettings';

type Props = {
    open: boolean;
    onClose: () => void;
};

export type ChildSave = { save: () => void };

export const DialogSettings = ({ open, onClose }: Props) => {
    const generalRef = useRef<ChildSave>(null);
    const cookiesRef = useRef<ChildSave>(null);

    const onSave = () => {
        generalRef.current?.save();
        cookiesRef.current?.save();
        onClose();
    };

    return (
        <Dialog
            open={open}
            onClose={(_, reason) => {
                if (reason !== 'backdropClick') {
                    onClose();
                }
            }}
            disableEscapeKeyDown
            maxWidth="md"
        >
            <DialogTitle sx={{ m: 0, p: 1.5 }} id="customized-dialog-title">
                Settings
            </DialogTitle>

            <IconButton
                aria-label="close"
                onClick={onClose}
                sx={(theme) => ({
                    position: 'absolute',
                    right: 8,
                    top: 8,
                    color: theme.palette.grey[500],
                })}
            >
                <Close />
            </IconButton>

            <TabContent generalRef={generalRef} cookiesRef={cookiesRef} />
            <Divider />

            <DialogActions>
                <Button color="error" onClick={onClose} sx={{ fontWeight: 'bold' }}>
                    Cancel
                </Button>
                <Button autoFocus onClick={onSave} sx={{ fontWeight: 'bold' }}>
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    );
};

type TabPanelProps = {
    generalRef: RefObject<ChildSave>;
    cookiesRef: RefObject<ChildSave>;
};

const TabContent = ({ generalRef, cookiesRef }: TabPanelProps) => {
    const [value, setValue] = useState(0);

    const handleChange = (_: SyntheticEvent, newValue: number) => {
        setValue(newValue);
    };

    return (
        <Box sx={{ width: 640, height: 374 }}>
            <Tabs value={value} onChange={handleChange} variant="scrollable" scrollButtons="auto">
                <Tab label="General" />
                <Tab label="Cookies" />
            </Tabs>

            <Divider />

            <GeneralSettings ref={generalRef} value={value} index={0} />
            <CookiesSettings ref={cookiesRef} value={value} index={1} />
        </Box>
    );
};
