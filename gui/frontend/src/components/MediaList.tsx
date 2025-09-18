import { ErrorTwoTone, ImageTwoTone, ListAlt, SaveOutlined, SmartDisplayTwoTone } from '@mui/icons-material';
import {
    Box,
    IconButton,
    Link,
    Popover,
    Stack,
    Table,
    TableBody,
    TableCell,
    TableRow,
    Tooltip,
    Typography,
} from '@mui/material';
import {
    DataGridPremium,
    type GridColDef,
    type GridRowSelectionModel,
    type GridRowsProp,
} from '@mui/x-data-grid-premium';
import { type ReactElement, useState } from 'react';
import './MediaList.css';
import { types } from '../../wailsjs/go/models';
import { BrowserOpenURL } from '../../wailsjs/runtime';
import { useAppStore } from '../stores/app';
import Media = types.Media;

export const MediaList = () => {
    const store = useAppStore();

    const columns: GridColDef[] = [
        { field: 'url', headerName: 'URL', flex: 0.65, renderCell: (params) => <LinkCell url={params.value} /> },
        { field: 'extension', headerName: 'Extension', flex: 0.1 },
        { field: 'type', headerName: 'Type', flex: 0.15, renderCell: (params) => <TypeCell type={params.value} /> },
        {
            field: 'metadata',
            headerName: 'Metadata',
            sortable: false,
            flex: 0.1,
            renderCell: (params) => <MetadataCell metadata={params.value} />,
        },
    ];

    const rows: GridRowsProp = store.media.map((m, id) => {
        return { id, url: m.Url, extension: m.Extension, type: m.Type, metadata: m.Metadata };
    });

    const handleSelectionChange = (selected: GridRowSelectionModel) => {
        const selectedMedia = selected.map((s) => store.media[s.valueOf() as number]);
        store.setSelectedMedia(selectedMedia);
    };

    const handleSelectionModel = (selectedMedia: Media[]): number[] => {
        return store.selectedMedia.map((media) => store.media.indexOf(media));
    };

    return (
        <DataGridPremium
            rows={rows}
            columns={columns}
            checkboxSelection={true}
            disableRowSelectionOnClick
            density="compact"
            localeText={{
                footerRowSelected: (count) => `${count} media selected`,
                footerTotalRows: <TotalMedia isCached={store.isCached} />,
                noRowsLabel: '',
            }}
            rowSelectionModel={handleSelectionModel(store.selectedMedia)}
            onRowSelectionModelChange={handleSelectionChange}
        />
    );
};

const TotalMedia = ({ isCached }: { isCached: boolean }) => {
    return (
        <>
            {isCached && (
                <Stack alignItems="center" direction="row" gap={1} sx={{ color: 'orange' }}>
                    <SaveOutlined />
                    <Typography variant="body2" sx={{ fontWeight: 'bold' }}>
                        Cached results&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                    </Typography>
                </Stack>
            )}

            <Typography variant="body2">Total media found:&nbsp;</Typography>
        </>
    );
};

const LinkCell = ({ url }: { url: string }) => {
    const handleClick = () => BrowserOpenURL(url);

    return (
        <Link href="#" onClick={handleClick}>
            {url}
        </Link>
    );
};

export const TypeCell = ({ type }: { type: number }) => {
    let icon: ReactElement;
    let label: ReactElement;

    switch (type) {
        case 0:
            icon = <ImageTwoTone sx={{ color: '#e3d026' }} />;
            label = <Typography variant="body2">Image</Typography>;
            break;
        case 1:
            icon = <SmartDisplayTwoTone color="secondary" />;
            label = <Typography variant="body2">Video</Typography>;
            break;
        default:
            icon = <ErrorTwoTone color="error" />;
            label = <Typography variant="body2">Unknown</Typography>;
    }

    return (
        <Stack direction="row" alignItems="center">
            {icon}
            &nbsp;{label}
        </Stack>
    );
};

const MetadataCell = ({ metadata }: { metadata: [string: string] }) => {
    const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);

    const handlePopoverOpen = (e: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(e.currentTarget);
    };

    const handlePopoverClose = () => {
        setAnchorEl(null);
    };

    const open = Boolean(anchorEl);

    return (
        <Box id="metadata">
            <Tooltip title="View metadata" arrow>
                <IconButton color="info" onClick={handlePopoverOpen}>
                    <ListAlt />
                </IconButton>
            </Tooltip>

            <Popover open={open} anchorEl={anchorEl} onClose={handlePopoverClose}>
                <Table size="small">
                    <TableBody>
                        {Object.entries(metadata).map(([key, value]) => (
                            <TableRow key={key}>
                                <TableCell>
                                    <strong>{key}</strong>
                                </TableCell>
                                <TableCell>{value}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </Popover>
        </Box>
    );
};
