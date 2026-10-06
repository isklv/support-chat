export interface SupportThread {
  id: string;
  number: number;
  email: string;
  name: string;
  locale: string;
  status: 'open' | 'closed';
  context?: string;
  user_seen_at?: string;
  last_message_at: string;
  created_at: string;
  closed_at?: string;
}

export interface SupportMessage {
  id: string;
  thread_id: string;
  author: 'user' | 'operator' | 'system';
  operator_name?: string;
  text: string;
  has_attachment: boolean;
  created_at: string;
}

export interface SupportView {
  thread: SupportThread | null;
  messages: SupportMessage[];
  token?: string;
}

export type WidgetPosition = 'bottom-right' | 'bottom-left';
export type WidgetTheme = 'auto' | 'light' | 'dark';

export interface SupportWidgetOptions {
  /** Base URL of the support backend (e.g. 'http://localhost:8080' or 'https://chat.example.com') */
  apiUrl: string;
  /** Widget title (default: 'Поддержка' / 'Support') */
  title?: string;
  /** Subtitle (default: 'Мы всегда на связи') */
  subtitle?: string;
  /** Primary accent color (hex or CSS color, default: '#2563eb') */
  primaryColor?: string;
  /** Floating button position on screen */
  position?: WidgetPosition;
  /** UI Language: 'ru' | 'en' | 'es' | 'auto' (default: 'auto') */
  lang?: 'ru' | 'en' | 'es' | 'auto';
  /** Pre-filled email for registered users */
  userEmail?: string;
  /** Pre-filled name for registered users */
  userName?: string;
  /** User identifier (e.g. user ID from your system) */
  userId?: string;
  /** Custom metadata attached to the conversation card (e.g. plan, tariff, etc.) */
  metadata?: Record<string, unknown>;
  /** Auto open widget upon initialization */
  autoOpen?: boolean;
}
