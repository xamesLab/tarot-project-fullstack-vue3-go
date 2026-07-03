import type { ReactNode } from 'react';
import { CelestialCrossIcon, StarIcon } from '@/shared/icons';
export function Ornament({ variant = 'star', className = '' }: { variant?: 'star' | 'cross'; className?: string }) { const Icon = variant === 'cross' ? CelestialCrossIcon : StarIcon; return <div className={`ornament ${className}`}><Icon width={variant === 'cross' ? 24 : 20} height={variant === 'cross' ? 24 : 20}/></div>; }
export function Divider({ children, dashed }: { children?: ReactNode; dashed?: boolean }) { return <div className={`divider ${dashed ? 'divider-dashed' : ''}`}>{children}</div>; }
export function SectionLabel({ children }: { children: ReactNode }) { return <div className="section-label">{children}</div>; }
