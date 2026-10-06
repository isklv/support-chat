import { useEffect, useRef, useState } from 'preact/hooks';
import type { SupportMessage, SupportView, SupportWidgetOptions } from '../types';
import { resolveLanguage } from '../i18n';
import {
  getStoredContact,
  getStoredToken,
  saveStoredContact,
  saveStoredToken,
  SupportClient,
} from '../api';
import { CloseIcon, LifeBuoyIcon, MessageSquareIcon, PaperclipIcon, SendIcon } from './Icons';

export function SupportWidget(options: SupportWidgetOptions) {
  const {
    apiUrl,
    title,
    subtitle,
    primaryColor = '#2563eb',
    position = 'bottom-right',
    lang = 'auto',
    userEmail = '',
    userName = '',
    userId,
    metadata,
    autoOpen = false,
  } = options;

  const t = resolveLanguage(lang);
  const client = useRef(new SupportClient(apiUrl)).current;

  const [isOpen, setIsOpen] = useState(autoOpen);
  const [view, setView] = useState<SupportView | null>(null);
  const [unread, setUnread] = useState(0);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Form states
  const [token, setToken] = useState<string | null>(getStoredToken());
  const [text, setText] = useState('');
  const [image, setImage] = useState<File | null>(null);

  const initialContact = getStoredContact();
  const [email, setEmail] = useState(userEmail || initialContact.email);
  const [name, setName] = useState(userName || initialContact.name);
  const [lightboxUrl, setLightboxUrl] = useState<string | null>(null);

  const listRef = useRef<HTMLUListElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Handle URL deep-linking: ?support=<token> or ?feedback=<token>
  useEffect(() => {
    if (typeof window === 'undefined') return;
    const params = new URLSearchParams(window.location.search);
    const linkToken = params.get('support') || params.get('feedback');
    if (linkToken) {
      if (linkToken !== '1') {
        saveStoredToken(linkToken);
        setToken(linkToken);
      }
      setIsOpen(true);
      // Clean query params cleanly without reloading
      params.delete('support');
      params.delete('feedback');
      const q = params.toString();
      window.history.replaceState(null, '', window.location.pathname + (q ? `?${q}` : '') + window.location.hash);
    }
  }, []);

  // Sync token changes to storage
  useEffect(() => {
    if (token) saveStoredToken(token);
  }, [token]);

  // Load conversation data
  const loadChat = async () => {
    try {
      const res = await client.fetchView(token, userId);
      setView(res);
      if (isOpen && res.thread) {
        await client.markSeen(token, userId);
        setUnread(0);
      }
    } catch {
      // Keep previous view if offline
    }
  };

  // Poll for unread count
  const checkUnreadCount = async () => {
    try {
      const count = await client.checkUnread(token, userId);
      setUnread(isOpen ? 0 : count);
    } catch {
      // Ignore
    }
  };

  // Adaptive polling loop
  useEffect(() => {
    void (isOpen ? loadChat() : checkUnreadCount());

    const intervalMs = isOpen ? 6000 : 45000;
    const timer = setInterval(() => {
      void (isOpen ? loadChat() : checkUnreadCount());
    }, intervalMs);

    return () => clearInterval(timer);
  }, [isOpen, token, userId]);

  // Scroll to bottom when new messages arrive or when panel is opened
  useEffect(() => {
    if (isOpen) {
      setTimeout(() => {
        if (listRef.current) {
          listRef.current.scrollTop = listRef.current.scrollHeight;
        }
      }, 50);
    }
  }, [isOpen, view?.messages.length]);

  const handleOpen = () => {
    setIsOpen(true);
    setError(null);
    void loadChat();
  };

  const handleClose = () => {
    setIsOpen(false);
  };

  const isGuest = !userId;
  const isFirstMessage = !view?.thread;
  const needEmail = isGuest && isFirstMessage;
  const canSend = !sending && (text.trim() !== '' || image !== null) && (!needEmail || /\S+@\S+\.\S+/.test(email));

  const handleSend = async (e: Event) => {
    e.preventDefault();
    if (!canSend) return;

    setSending(true);
    setError(null);

    const formData = new FormData();
    formData.append('text', text.trim());
    if (image) formData.append('image', image);
    formData.append('page', typeof window !== 'undefined' ? window.location.pathname : '');
    formData.append('lang', lang === 'auto' ? (typeof navigator !== 'undefined' ? navigator.language : 'ru') : lang);

    if (email) formData.append('email', email.trim());
    if (name) formData.append('name', name.trim());
    if (metadata) formData.append('metadata', JSON.stringify(metadata));

    try {
      const result = await client.sendMessage(formData, token, userId);
      if (result.token) {
        saveStoredToken(result.token);
        setToken(result.token);
      }
      if (email || name) {
        saveStoredContact(email, name);
      }
      setView(result);
      setText('');
      setImage(null);
      setSending(false);
    } catch (err) {
      setSending(false);
      const code = (err as { code?: string })?.code || '';
      if (code === 'email_required') setError(t.errors.email);
      else if (code === 'invalid_image') setError(t.errors.image);
      else if (code === 'message_too_long') setError(t.errors.long);
      else setError(t.errors.failed);
    }
  };

  // Clipboard paste support (Ctrl+V / Cmd+V screenshot paste)
  const handlePaste = (e: ClipboardEvent) => {
    const items = e.clipboardData?.items;
    if (!items) return;
    for (let i = 0; i < items.length; i++) {
      if (items[i].type.startsWith('image/')) {
        const file = items[i].getAsFile();
        if (file) {
          setImage(file);
          e.preventDefault();
          break;
        }
      }
    }
  };

  const formatTime = (dateStr: string) => {
    try {
      return new Intl.DateTimeFormat(undefined, {
        hour: '2-digit',
        minute: '2-digit',
      }).format(new Date(dateStr));
    } catch {
      return '';
    }
  };

  const renderBubble = (m: SupportMessage) => {
    if (m.author === 'system') {
      return (
        <li key={m.id} className="sc-note">
          {m.text === 'closed' ? t.closedNote : t.reopenedNote} · {formatTime(m.created_at)}
        </li>
      );
    }

    const isMine = m.author === 'user';
    const imgSrc = m.has_attachment ? client.attachmentUrl(m.id, token) : null;

    return (
      <li key={m.id} className={`sc-bubble ${isMine ? 'sc-mine' : 'sc-theirs'}`}>
        {!isMine && <span className="sc-who">{m.operator_name || 'Поддержка'}</span>}
        {imgSrc && (
          <img
            src={imgSrc}
            alt=""
            className="sc-image"
            onClick={() => setLightboxUrl(imgSrc)}
          />
        )}
        {m.text && <p>{m.text}</p>}
        <span className="sc-time">{formatTime(m.created_at)}</span>
      </li>
    );
  };

  const posClass = position === 'bottom-left' ? 'sc-pos-left' : '';
  const widgetStyle = { '--sc-primary': primaryColor } as Record<string, string>;

  return (
    <div style={widgetStyle}>
      {/* Floating launcher button */}
      {!isOpen && (
        <button
          className={`sc-launcher ${posClass}`}
          onClick={handleOpen}
          aria-label={t.title}
          title={t.title}
        >
          <MessageSquareIcon />
          {unread > 0 && <span className="sc-badge">{unread}</span>}
        </button>
      )}

      {/* Main chat window */}
      {isOpen && (
        <section className={`sc-panel ${posClass}`} role="dialog" aria-label={t.title}>
          <header className="sc-header">
            <div className="sc-header-icon">
              <LifeBuoyIcon />
            </div>
            <div className="sc-header-text">
              <h3 className="sc-title">{title || t.title}</h3>
              <p className="sc-subtitle">{subtitle || (isGuest ? t.subtitleGuest : t.subtitle)}</p>
            </div>
            <button className="sc-close-btn" onClick={handleClose} title="Закрыть">
              <CloseIcon />
            </button>
          </header>

          <ul className="sc-messages" ref={listRef}>
            {(!view?.messages || view.messages.length === 0) && (
              <li className="sc-empty">{t.empty}</li>
            )}
            {view?.messages.map(renderBubble)}
            {view?.thread?.status === 'closed' && (
              <li className="sc-note">{t.closedHint}</li>
            )}
          </ul>

          <form className="sc-form" onSubmit={handleSend}>
            {needEmail && (
              <div className="sc-contacts">
                <input
                  className="sc-input"
                  type="email"
                  required
                  placeholder={t.email}
                  value={email}
                  onInput={(e) => setEmail((e.target as HTMLInputElement).value)}
                />
                <input
                  className="sc-input"
                  type="text"
                  placeholder={t.name}
                  value={name}
                  onInput={(e) => setName((e.target as HTMLInputElement).value)}
                />
                {/* Honeypot field for spam prevention */}
                <input className="sc-hp" name="website" tabIndex={-1} autoComplete="off" />
              </div>
            )}

            {image && (
              <div className="sc-attached-chip">
                <span>📷 {image.name}</span>
                <button
                  type="button"
                  className="sc-close-btn"
                  onClick={() => setImage(null)}
                  title={t.removeImage}
                >
                  <CloseIcon />
                </button>
              </div>
            )}

            {error && <p className="sc-error">{error}</p>}

            <div className="sc-compose">
              <button
                type="button"
                className="sc-action-btn"
                onClick={() => fileInputRef.current?.click()}
                title={t.attach}
              >
                <PaperclipIcon />
              </button>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/png,image/jpeg,image/webp"
                style={{ display: 'none' }}
                onChange={(e) => {
                  const files = (e.target as HTMLInputElement).files;
                  if (files && files[0]) setImage(files[0]);
                  (e.target as HTMLInputElement).value = '';
                }}
              />
              <textarea
                className="sc-textarea"
                rows={2}
                maxLength={4000}
                placeholder={t.placeholder}
                value={text}
                onInput={(e) => setText((e.target as HTMLTextAreaElement).value)}
                onPaste={handlePaste}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    void handleSend(e);
                  }
                }}
              />
              <button
                type="submit"
                className="sc-send-btn"
                disabled={!canSend}
                title={t.send}
              >
                <SendIcon />
              </button>
            </div>
            <p className="sc-foot">{t.foot}</p>
          </form>
        </section>
      )}

      {/* Lightbox image preview modal */}
      {lightboxUrl && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            background: 'rgba(0,0,0,0.85)',
            zIndex: 999999,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '24px',
          }}
          onClick={() => setLightboxUrl(null)}
        >
          <img
            src={lightboxUrl}
            alt="Preview"
            style={{
              maxWidth: '90vw',
              maxHeight: '90vh',
              borderRadius: '8px',
              objectFit: 'contain',
            }}
          />
        </div>
      )}
    </div>
  );
}
