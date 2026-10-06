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

export const translations: Record<string, Translations> = {
  ru: {
    title: 'Поддержка',
    subtitle: 'Оператор ответит вам здесь',
    subtitleGuest: 'Оставьте вопрос или пожелание',
    empty: 'Задайте любой вопрос или отправьте снимок экрана — мы на связи!',
    email: 'Ваш email (для ответа)',
    name: 'Как к вам обращаться',
    placeholder: 'Напишите сообщение... (Enter — отправить)',
    send: 'Отправить',
    attach: 'Прикрепить скриншот',
    removeImage: 'Удалить картинку',
    closedNote: 'Тема закрыта оператором',
    reopenedNote: 'Диалог возобновлен',
    closedHint: 'Вопрос решен. Напишите снова, если потребуется помощь.',
    foot: 'Ответы оператора появятся здесь и придут на email',
    errors: {
      email: 'Укажите корректный email для связи',
      image: 'Недопустимый формат или размер файла (макс. 4 МБ)',
      long: 'Слишком длинное сообщение (макс. 4000 знаков)',
      failed: 'Не удалось отправить сообщение. Попробуйте еще раз',
    },
  },
  en: {
    title: 'Support Chat',
    subtitle: 'Our team will answer here',
    subtitleGuest: 'Ask a question or leave feedback',
    empty: 'Ask anything or attach a screenshot — we are here to help!',
    email: 'Your email address',
    name: 'Your name',
    placeholder: 'Type a message... (Enter to send)',
    send: 'Send',
    attach: 'Attach screenshot',
    removeImage: 'Remove image',
    closedNote: 'Ticket closed by operator',
    reopenedNote: 'Ticket reopened',
    closedHint: 'This conversation is closed. Send a message to reopen.',
    foot: 'Replies will appear here and in your inbox',
    errors: {
      email: 'Please provide a valid email address',
      image: 'Unsupported image format or size (max 4 MB)',
      long: 'Message is too long (max 4,000 characters)',
      failed: 'Could not send message. Please try again',
    },
  },
  es: {
    title: 'Soporte',
    subtitle: 'Nuestro equipo responderá aquí',
    subtitleGuest: 'Haz una pregunta o deja comentarios',
    empty: 'Escribe cualquier duda o adjunta una captura de pantalla',
    email: 'Tu correo electrónico',
    name: 'Tu nombre',
    placeholder: 'Escribe un mensaje... (Enter para enviar)',
    send: 'Enviar',
    attach: 'Adjuntar imagen',
    removeImage: 'Eliminar imagen',
    closedNote: 'Conversación cerrada por el operador',
    reopenedNote: 'Conversación reabierta',
    closedHint: 'Esta conversación ha sido cerrada. Escribe para reabrirla.',
    foot: 'Las respuestas llegarán aquí y a tu correo',
    errors: {
      email: 'Introduce un correo electrónico válido',
      image: 'Formato o tamaño de imagen no admitido (máx. 4 MB)',
      long: 'Mensaje demasiado largo (máx. 4000 caracteres)',
      failed: 'No se pudo enviar el mensaje. Inténtalo de nuevo',
    },
  },
};

export function resolveLanguage(lang?: 'ru' | 'en' | 'es' | 'auto'): Translations {
  if (lang && lang !== 'auto' && translations[lang]) {
    return translations[lang];
  }
  const browserLang = (typeof navigator !== 'undefined' ? navigator.language || '' : '').toLowerCase();
  if (browserLang.startsWith('ru') || browserLang.startsWith('be') || browserLang.startsWith('uk') || browserLang.startsWith('kk')) {
    return translations.ru;
  }
  if (browserLang.startsWith('es')) {
    return translations.es;
  }
  return translations.en;
}
