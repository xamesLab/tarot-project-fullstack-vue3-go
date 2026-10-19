import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { usePreferencesStore } from '@/features/preferences/store';
import { SendIcon } from '@/shared/icons';

interface Message {
  id: string;
  from: 'oracle' | 'me';
  text: string;
  seeded?: boolean;
}

const ORACLE_SYMBOL = '\u2736';

export function OracleChat() {
  const { t } = useTranslation('home');
  const language = usePreferencesStore((state) => state.language);
  const [draft, setDraft] = useState('');
  const [messages, setMessages] = useState<Message[]>([
    { id: 'greeting', from: 'oracle', text: t('oracleGreeting'), seeded: true },
  ]);
  const logRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setMessages((current) => current.map((message) => (
      message.seeded ? { ...message, text: t('oracleGreeting') } : message
    )));
  }, [language, t]);

  useEffect(() => {
    if (logRef.current) logRef.current.scrollTop = logRef.current.scrollHeight;
  }, [messages]);

  const send = () => {
    const text = draft.trim();
    if (!text) return;

    const questionId = crypto.randomUUID();
    setMessages((current) => [
      ...current,
      { id: questionId, from: 'me', text },
      { id: questionId + '-oracle', from: 'oracle', text: t('oracleSilenceReply') },
    ]);
    setDraft('');
  };

  return (
    <section className="chat-section">
      <div className="chat-heading"><span>{t('chatTitle')}</span></div>
      <div className="chat-box">
        <div className="chat-log" ref={logRef} aria-live="polite">
          {messages.map((message) => (
            <div key={message.id} className={'chat-row ' + message.from}>
              <div className={'mini-avatar ' + message.from}>
                {message.from === 'oracle' ? ORACLE_SYMBOL : language === 'ru' ? '\u041c\u0410' : 'MA'}
              </div>
              <div className="bubble">{message.text}</div>
            </div>
          ))}
        </div>
        <div className="chat-input">
          <input
            value={draft}
            onChange={(event) => setDraft(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === 'Enter') send();
            }}
            placeholder={t('chatPlaceholder')}
          />
          <button type="button" onClick={send} aria-label={t('send')}>
            <SendIcon />
          </button>
        </div>
      </div>
    </section>
  );
}
