export interface Translations {
    title: string;
    subtitle: string;
    subtitleGuest: string;
    empty: string;
    email: string;
    name: string;
    placeholder: string;
    send: string;
    attach: string;
    removeImage: string;
    closedNote: string;
    reopenedNote: string;
    closedHint: string;
    foot: string;
    errors: {
        email: string;
        image: string;
        long: string;
        failed: string;
    };
}
export declare const translations: Record<string, Translations>;
export declare function resolveLanguage(lang?: 'ru' | 'en' | 'es' | 'auto'): Translations;
