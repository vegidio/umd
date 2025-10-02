import { CssBaseline, createTheme, ThemeProvider, useMediaQuery } from '@mui/material';
import { LicenseInfo } from '@mui/x-license';
import { createRoot } from 'react-dom/client';
import App from './App';
import Events from './Events';

const container = document.getElementById('root');
// biome-ignore lint/style/noNonNullAssertion: <explanation>
const root = createRoot(container!);

const lightTheme = createTheme({
    palette: {
        mode: 'light',
    },
});

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});

const licenseKey = import.meta.env.VITE_MUI_LICENSE;
LicenseInfo.setLicenseKey(licenseKey);

const Main = () => {
    const isDarkMode = useMediaQuery('(prefers-color-scheme: dark)');

    return (
        <ThemeProvider theme={isDarkMode ? darkTheme : lightTheme}>
            <Events />
            <CssBaseline />
            <App />
        </ThemeProvider>
    );
};

root.render(<Main />);
