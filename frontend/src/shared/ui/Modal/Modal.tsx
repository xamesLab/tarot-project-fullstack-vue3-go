import { type ReactNode, useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import { useTranslation } from 'react-i18next';
import { CloseIcon } from '@/shared/icons';

const focusables = 'button:not([disabled]), [href], input:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])';

export function Modal({ open, onClose, title, children, footer, initialFocusRef }: { open: boolean; onClose: () => void; title: string; children: ReactNode; footer?: ReactNode; initialFocusRef?: React.RefObject<HTMLElement | null> }) {
  const { t } = useTranslation();
  const panelRef = useRef<HTMLDivElement>(null);
  const restoreRef = useRef<HTMLElement | null>(null);

  useEffect(() => {
    if (!open) return;
    restoreRef.current = document.activeElement as HTMLElement;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    const timer = window.setTimeout(() => (initialFocusRef?.current ?? panelRef.current?.querySelector<HTMLElement>(focusables))?.focus(), 20);
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose();
      if (event.key !== 'Tab' || !panelRef.current) return;
      const nodes = Array.from(panelRef.current.querySelectorAll<HTMLElement>(focusables));
      if (!nodes.length) return;
      const first = nodes[0]; const last = nodes[nodes.length - 1];
      if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last.focus(); }
      if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first.focus(); }
    };
    document.addEventListener('keydown', onKeyDown);
    return () => { window.clearTimeout(timer); document.removeEventListener('keydown', onKeyDown); document.body.style.overflow = previousOverflow; restoreRef.current?.focus(); };
  }, [initialFocusRef, onClose, open]);

  if (!open) return null;
  return createPortal(<div className="modal-overlay open" onMouseDown={(event) => { if (event.target === event.currentTarget) onClose(); }}><div ref={panelRef} className="modal-panel" role="dialog" aria-modal="true" aria-labelledby="modal-title"><div className="modal-head"><div id="modal-title" className="modal-title">{title}</div><button className="modal-close" type="button" onClick={onClose} aria-label={t('close')}><CloseIcon width="14" height="14"/></button></div><div className="modal-body">{children}</div>{footer && <div className="modal-foot">{footer}</div>}</div></div>, document.body);
}
