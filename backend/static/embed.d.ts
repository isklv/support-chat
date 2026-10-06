import type { SupportWidgetOptions } from './types';
import './styles.css';
declare global {
    interface Window {
        SupportChat?: {
            init: (options: SupportWidgetOptions) => void;
            open?: () => void;
            close?: () => void;
        };
        SupportChatOptions?: SupportWidgetOptions;
    }
}
export declare function init(options: SupportWidgetOptions): void;
