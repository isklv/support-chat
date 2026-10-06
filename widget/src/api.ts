import type { SupportView } from './types';

const TOKEN_KEY = 'support_chat_token';
const EMAIL_KEY = 'support_chat_email';
const NAME_KEY = 'support_chat_name';

export function getStoredToken(): string | null {
  try {
    return localStorage.getItem(TOKEN_KEY);
  } catch {
    return null;
  }
}

export function saveStoredToken(token: string): void {
  try {
    localStorage.setItem(TOKEN_KEY, token);
  } catch {
    // Ignore storage errors in restricted contexts
  }
}

export function getStoredContact(): { email: string; name: string } {
  try {
    return {
      email: localStorage.getItem(EMAIL_KEY) || '',
      name: localStorage.getItem(NAME_KEY) || '',
    };
  } catch {
    return { email: '', name: '' };
  }
}

export function saveStoredContact(email: string, name: string): void {
  try {
    if (email) localStorage.setItem(EMAIL_KEY, email);
    if (name) localStorage.setItem(NAME_KEY, name);
  } catch {
    // Ignore storage errors
  }
}

export class SupportClient {
  private baseUrl: string;

  constructor(apiUrl: string) {
    this.baseUrl = apiUrl.replace(/\/+$/, '');
  }

  private headers(token?: string | null, userId?: string): Record<string, string> {
    const h: Record<string, string> = {};
    if (token) h['X-Support-Token'] = token;
    if (userId) h['X-User-Id'] = userId;
    return h;
  }

  async fetchView(token?: string | null, userId?: string): Promise<SupportView> {
    const res = await fetch(`${this.baseUrl}/api/v1/support`, {
      method: 'GET',
      headers: this.headers(token, userId),
    });
    if (!res.ok) {
      throw new Error(`fetch view failed: ${res.status}`);
    }
    return res.json();
  }

  async sendMessage(form: FormData, token?: string | null, userId?: string): Promise<SupportView> {
    const res = await fetch(`${this.baseUrl}/api/v1/support/messages`, {
      method: 'POST',
      body: form,
      headers: this.headers(token, userId),
    });
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      const err = new Error(data.error || `send message failed: ${res.status}`);
      (err as unknown as { code: string }).code = data.error || '';
      throw err;
    }
    return res.json();
  }

  async markSeen(token?: string | null, userId?: string): Promise<void> {
    await fetch(`${this.baseUrl}/api/v1/support/seen`, {
      method: 'POST',
      headers: this.headers(token, userId),
    }).catch(() => {});
  }

  async checkUnread(token?: string | null, userId?: string): Promise<number> {
    const res = await fetch(`${this.baseUrl}/api/v1/support/unread`, {
      method: 'GET',
      headers: this.headers(token, userId),
    });
    if (!res.ok) return 0;
    const data = await res.json();
    return typeof data.unread === 'number' ? data.unread : 0;
  }

  attachmentUrl(messageId: string, token?: string | null): string {
    const url = `${this.baseUrl}/api/v1/support/attachments/${messageId}`;
    return token ? `${url}?t=${encodeURIComponent(token)}` : url;
  }
}
