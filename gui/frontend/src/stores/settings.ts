import { persist } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import { create } from 'zustand/react';

type SettingsStore = {
    // General
    deepSearch: boolean;
    ignoreCache: boolean;
    enableTelemetry: boolean;
    setDeepSearch: (deepSearch: boolean) => void;
    setIgnoreCache: (ignoreCache: boolean) => void;
    setEnableTelemetry: (enableTelemetry: boolean) => void;

    // Cookies
    cookiesType: string;
    cookiesPath: string;
    setCookiesType: (cookiesType: string) => void;
    setCookiesPath: (cookiesPath: string) => void;
};

export const useSettingsStore = create(
    persist(
        immer<SettingsStore>((set, get) => ({
            // General
            deepSearch: true,
            ignoreCache: false,
            enableTelemetry: true,

            setDeepSearch: (deepSearch: boolean) => {
                set((state) => {
                    state.deepSearch = deepSearch;
                });
            },

            setIgnoreCache: (ignoreCache: boolean) => {
                set((state) => {
                    state.ignoreCache = ignoreCache;
                });
            },

            setEnableTelemetry: (enableTelemetry: boolean) => {
                set((state) => {
                    state.enableTelemetry = enableTelemetry;
                });
            },

            // Cookies
            cookiesType: 'disabled',
            cookiesPath: '',

            setCookiesType: (cookiesType: string) => {
                set((state) => {
                    state.cookiesType = cookiesType;
                });
            },
            setCookiesPath: (cookiesPath: string) => {
                set((state) => {
                    state.cookiesPath = cookiesPath;
                });
            },
        })),
        {
            name: 'settings-storage',
        },
    ),
);
