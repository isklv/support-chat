import React from 'react';
import { SupportWidget } from 'support-chat-widget';

export function App() {
  // Текущий пользователь (если авторизован на вашем сайте)
  const currentUser = {
    id: 'user_12345',
    email: 'alex@example.com',
    name: 'Александр',
  };

  return (
    <div style={{ padding: 40, fontFamily: 'sans-serif' }}>
      <h1>Мое веб-приложение на React</h1>
      <p>Виджет поддержки подключен как React-компонент.</p>

      <SupportWidget
        apiUrl="http://localhost:8080"
        title="Поддержка клиентов"
        subtitle="Ответим в течение 5 минут"
        primaryColor="#2563eb"
        position="bottom-right"
        lang="ru"
        userId={currentUser.id}
        userEmail={currentUser.email}
        userName={currentUser.name}
        metadata={{
          plan: 'Pro',
          role: 'Owner',
        }}
      />
    </div>
  );
}

export default App;
