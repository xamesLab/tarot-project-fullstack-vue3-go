import type { ReactNode } from 'react';

export function Chip({ solid, children }: { solid?: boolean; children: ReactNode }) { return <span className={`chip ${solid ? 'chip-solid' : ''}`}>{children}</span>; }
export function Badge({ variant = 'gold', dot, children }: { variant?: 'gold' | 'soft' | 'dim'; dot?: boolean; children: ReactNode }) { return <span className={`badge badge-${variant}`}>{dot && <i className="badge-dot"/>}{children}</span>; }
