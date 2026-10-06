import { render } from 'preact';
import { SupportWidget } from './components/SupportWidget';
import type { SupportWidgetOptions } from './types';
import './styles.css';

export { SupportWidget } from './components/SupportWidget';
export { SupportClient } from './api';
export * from './types';
export * from './i18n';

/**
 * Mount support chat widget to any DOM element
 */
export function mountSupportChat(options: SupportWidgetOptions, container?: HTMLElement): () => void {
  const target = container || (() => {
    let el = document.getElementById('support-chat-root');
    if (!el) {
      el = document.createElement('div');
      el.id = 'support-chat-root';
      document.body.appendChild(el);
    }
    return el;
  })();

  render(<SupportWidget {...options} />, target);

  return () => {
    render(null, target);
  };
}
