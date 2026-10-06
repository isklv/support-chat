import type { SupportWidgetOptions } from './types';
import './styles.css';
export { SupportWidget } from './components/SupportWidget';
export { SupportClient } from './api';
export * from './types';
export * from './i18n';
/**
 * Mount support chat widget to any DOM element
 */
export declare function mountSupportChat(options: SupportWidgetOptions, container?: HTMLElement): () => void;
