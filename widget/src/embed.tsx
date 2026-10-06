import { render } from 'preact';
import { SupportWidget } from './components/SupportWidget';
import type { SupportWidgetOptions, WidgetPosition } from './types';
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

let rootContainer: HTMLElement | null = null;

export function init(options: SupportWidgetOptions) {
  if (!rootContainer) {
    rootContainer = document.createElement('div');
    rootContainer.id = 'support-chat-root';
    document.body.appendChild(rootContainer);
  }

  render(<SupportWidget {...options} />, rootContainer);
}

window.SupportChat = {
  init,
};

// Auto-initialize if loaded via <script data-api="..."> or if window.SupportChatOptions is defined
if (typeof document !== 'undefined') {
  const autoInit = () => {
    if (window.SupportChatOptions) {
      init(window.SupportChatOptions);
      return;
    }

    const currentScript = (document.currentScript as HTMLScriptElement) ||
      document.querySelector('script[data-api]');

    if (currentScript) {
      const apiUrl = currentScript.getAttribute('data-api') || '';
      if (!apiUrl) return;

      const title = currentScript.getAttribute('data-title') || undefined;
      const subtitle = currentScript.getAttribute('data-subtitle') || undefined;
      const primaryColor = currentScript.getAttribute('data-color') || undefined;
      const position = (currentScript.getAttribute('data-position') as WidgetPosition) || undefined;
      const lang = (currentScript.getAttribute('data-lang') as 'ru' | 'en' | 'es' | 'auto') || undefined;
      const autoOpen = currentScript.hasAttribute('data-auto-open');

      init({
        apiUrl,
        title,
        subtitle,
        primaryColor,
        position,
        lang,
        autoOpen,
      });
    }
  };

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', autoInit);
  } else {
    autoInit();
  }
}
