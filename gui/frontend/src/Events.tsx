import { useEffect } from 'react';
import { fetch } from '../wailsjs/go/models';
import { EventsOn } from '../wailsjs/runtime';
import { useAppStore } from './stores/app';
import Response = fetch.Response;

const Events = () => {
    const clear = useAppStore((state) => state.clear);
    const setAmountQuery = useAppStore((state) => state.setAmountQuery);
    const setCurrentDownloads = useAppStore((state) => state.setCurrentDownloads);
    const setDownloadedMedia = useAppStore((state) => state.setDownloadedMedia);
    const setExtractorName = useAppStore((state) => state.setExtractorName);
    const setExtractorType = useAppStore((state) => state.setExtractorType);
    const setIsCached = useAppStore((state) => state.setIsCached);

    // biome-ignore lint/correctness/useExhaustiveDependencies: this should run only once
    useEffect(() => {
        const unbindOnExtractorFound = EventsOn('OnExtractorFound', (name: string) => {
            clear();
            setExtractorName(name);
        });

        const unbindOnExtractorTypeFound = EventsOn('OnExtractorTypeFound', (eType: string, name: string) =>
            setExtractorType(eType, name),
        );

        const unbindOnMediaQueried = EventsOn('OnMediaQueried', (amount: number) => {
            setAmountQuery(amount);
        });

        const unbindOnQueryCompleted = EventsOn('OnQueryCompleted', (_: number, isCached: boolean) => {
            setIsCached(isCached);
        });

        const unbindOnMediaDownloaded = EventsOn('OnMediaDownloaded', (amount: number, responses: Response[]) => {
            setDownloadedMedia(amount);
            setCurrentDownloads(responses);
        });

        return () => {
            unbindOnExtractorFound();
            unbindOnExtractorTypeFound();
            unbindOnMediaQueried();
            unbindOnQueryCompleted();
            unbindOnMediaDownloaded();
        };
    }, []);

    return <></>;
};

export default Events;
