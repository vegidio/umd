import { immer } from 'zustand/middleware/immer';
import { create } from 'zustand/react';
import { fetch, types } from '../../wailsjs/go/models';
import Response = fetch.Response;
import Media = types.Media;

type AppStore = {
    isCached: boolean;
    isQuerying: boolean;
    extractorName: string;
    extractorType: string;
    extractorTypeName: string;
    amountQuery: number;
    directory: string;
    progress: number;
    media: Media[];
    selectedMedia: Media[];
    downloadedMedia: number;
    currentDownloads: Response[];

    clear: () => void;
    setIsCached: (cached: boolean) => void;
    setIsQuerying: (querying: boolean) => void;
    setAmountQuery: (amount: number) => void;
    setExtractorName: (name: string) => void;
    setExtractorType: (eType: string, name: string) => void;
    setDirectory: (directory: string) => void;
    setProgress: (percent: number) => void;
    setMedia: (media: Media[]) => void;
    setSelectedMedia: (media: Media[]) => void;
    clearDownloads: () => void;
    setDownloadedMedia: (amount: number) => void;
    setCurrentDownloads: (responses: Response[]) => void;
};

export const useAppStore = create(
    immer<AppStore>((set, get) => ({
        isCached: false,
        isQuerying: false,
        extractorName: '',
        extractorType: '',
        extractorTypeName: '',
        amountQuery: 0,
        directory: '.',
        progress: 0,
        media: [],
        selectedMedia: [],
        downloadedMedia: 0,
        currentDownloads: [],

        clear: () =>
            set((state) => {
                state.isCached = false;
                state.extractorName = '';
                state.extractorType = '';
                state.extractorTypeName = '';
                state.amountQuery = 0;
                state.progress = 0;
                state.media = [];
                state.selectedMedia = [];
                state.downloadedMedia = 0;
                state.currentDownloads = [];
            }),

        setIsCached: (cached: boolean) =>
            set((state) => {
                state.isCached = cached;
            }),

        setIsQuerying: (querying: boolean) =>
            set((state) => {
                state.isQuerying = querying;
            }),

        setAmountQuery: (amount: number) =>
            set((state) => {
                state.amountQuery = amount;
            }),

        setExtractorName: (name: string) =>
            set((state) => {
                state.extractorName = name;
            }),

        setExtractorType: (eType: string, name: string) =>
            set((state) => {
                state.extractorType = eType;
                state.extractorTypeName = name;
            }),

        setDirectory: (directory: string) => {
            localStorage.setItem('lastDirectory', directory);

            set((state) => {
                state.directory = directory;
            });
        },

        setProgress: (percent: number) =>
            set((state) => {
                state.progress = percent;
            }),

        setMedia: (media: Media[]) =>
            set((state) => {
                state.media = media;
            }),

        setSelectedMedia: (media: Media[]) =>
            set((state) => {
                state.selectedMedia = media;
            }),

        clearDownloads: () =>
            set((state) => {
                state.downloadedMedia = 0;
                state.currentDownloads = [];
            }),

        setDownloadedMedia: (amount: number) =>
            set((state) => {
                state.downloadedMedia = amount;
            }),

        setCurrentDownloads: (responses: Response[]) =>
            set((state) => {
                state.currentDownloads = responses;
            }),
    })),
);
