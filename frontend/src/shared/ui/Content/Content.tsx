import type { ReactNode } from 'react';
import { SphereIcon } from '@/shared/icons';

export function AspectBlock({ type, title, children }: { type: 'love' | 'career' | 'health' | 'spirit' | 'finance'; title: string; children: ReactNode }) { return <article className="aspect-block"><SphereIcon type={type} className="aspect-icon"/><h3>{title}</h3><div>{children}</div></article>; }
export function CornerPanel({ title, children }: { title: string; children: ReactNode }) { return <section className="corner-panel"><i className="corner tl"/><i className="corner tr"/><i className="corner bl"/><i className="corner br"/><h2>{title}</h2><div>{children}</div></section>; }
export function Tooltip({ label, children, below }: { label: ReactNode; children: ReactNode; below?: boolean }) { return <span className={`tooltip ${below ? 'below' : ''}`} tabIndex={0}>{children}<span className="tooltip-bubble" role="tooltip">{label}</span></span>; }
