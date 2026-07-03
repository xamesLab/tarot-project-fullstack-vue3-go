import { useMutation } from '@tanstack/react-query';
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { api } from '@/api';
import type { Language } from '@/api/types';
import { usePreferencesStore } from '@/features/preferences/store';
import { SendIcon } from '@/shared/icons';

interface Message { id: string; from: 'oracle' | 'me'; text: string; seeded?: boolean; }

export function OracleChat() {
  const { t } = useTranslation('home');
  const language = usePreferencesStore((state) => state.language);
  const [draft, setDraft] = useState('');
  const [messages, setMessages] = useState<Message[]>([{ id: 'greeting', from: 'oracle', text: t('oracleGreeting'), seeded: true }]);
  const logRef = useRef<HTMLDivElement>(null);
  const mutation = useMutation({ mutationFn: ({ question, lang }: { question: string; lang: Language }) => api.oracle.ask(question, lang), onSuccess: (reply) => setMessages((current) => [...current, { id: reply.id, from: 'oracle', text: reply.text }]) });

  useEffect(() => { setMessages((current) => current.map((message) => message.seeded ? { ...message, text: t('oracleGreeting') } : message)); }, [language, t]);
  useEffect(() => { if (logRef.current) logRef.current.scrollTop = logRef.current.scrollHeight; }, [messages, mutation.isPending]);

  const send = () => { const text = draft.trim(); if (!text || mutation.isPending) return; setMessages((current) => [...current, { id: crypto.randomUUID(), from: 'me', text }]); setDraft(''); mutation.mutate({ question: text, lang: language }); };
  return <section className="chat-section"><div className="chat-heading"><span>{t('chatTitle')}</span></div><div className="chat-box"><div className="chat-log" ref={logRef} aria-live="polite">{messages.map((message) => <div key={message.id} className={`chat-row ${message.from}`}><div className={`mini-avatar ${message.from}`}>{message.from === 'oracle' ? '✶' : language === 'ru' ? 'МА' : 'MA'}</div><div className="bubble">{message.text}</div></div>)}{mutation.isPending && <div className="chat-row oracle"><div className="mini-avatar oracle">✶</div><div className="bubble typing">· · ·</div></div>}</div><div className="chat-input"><input value={draft} onChange={(event) => setDraft(event.target.value)} onKeyDown={(event) => { if (event.key === 'Enter') send(); }} placeholder={t('chatPlaceholder')}/><button type="button" onClick={send} disabled={mutation.isPending} aria-label={t('send')}><SendIcon/></button></div></div></section>;
}
