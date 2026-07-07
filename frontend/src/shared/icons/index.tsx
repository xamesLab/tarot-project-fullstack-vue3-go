import type { SVGProps } from 'react';

type IconProps = SVGProps<SVGSVGElement>;

export function LogoSigil(props: IconProps) {
  return <svg viewBox="0 0 36 36" fill="none" aria-hidden="true" {...props}><circle cx="18" cy="18" r="17" stroke="currentColor" strokeWidth=".75" opacity=".6"/><circle cx="18" cy="18" r="11" stroke="currentColor" strokeWidth=".5" opacity=".4"/><ellipse cx="18" cy="18" rx="6" ry="4" stroke="currentColor" strokeWidth=".75"/><circle cx="18" cy="18" r="2" fill="currentColor" opacity=".8"/><path d="M18 3v4m0 22v4M3 18h4m22 0h4" stroke="currentColor" strokeWidth=".75" opacity=".5"/></svg>;
}

export function StarIcon(props: IconProps) {
  return <svg viewBox="0 0 20 20" fill="none" aria-hidden="true" {...props}><path d="m10 2 2 6h6l-5 4 2 6-5-4-5 4 2-6-5-4h6l2-6Z" stroke="currentColor" strokeWidth=".75" opacity=".7"/><circle cx="10" cy="10" r="2" fill="currentColor" opacity=".45"/></svg>;
}

export function CelestialCrossIcon(props: IconProps) {
  return <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}><circle cx="12" cy="12" r="3" stroke="currentColor" strokeWidth=".75" opacity=".6"/><path d="M12 2v4m0 12v4M2 12h4m12 0h4" stroke="currentColor" strokeWidth=".75" opacity=".4"/></svg>;
}

export function MoonPhases(props: IconProps) {
  return <svg viewBox="0 0 280 22" fill="none" stroke="currentColor" aria-hidden="true" {...props}><circle cx="20" cy="11" r="6" strokeWidth=".7" opacity=".4"/><circle cx="54" cy="11" r="6" strokeWidth=".7" opacity=".5"/><path d="M54 5a6 6 0 0 1 0 12 9 7 0 0 0 0-12Z" fill="currentColor" stroke="none" opacity=".4"/><circle cx="88" cy="11" r="6" strokeWidth=".7" opacity=".55"/><path d="M88 5a6 6 0 0 1 0 12Z" fill="currentColor" stroke="none" opacity=".5"/><circle cx="122" cy="11" r="6" strokeWidth=".7" opacity=".6"/><path d="M122 5a6 6 0 0 1 0 12 3.5 6 0 0 0 0-12Z" fill="currentColor" stroke="none" opacity=".6"/><circle cx="156" cy="11" r="7" fill="currentColor" stroke="none" opacity=".8"/><circle cx="190" cy="11" r="6" strokeWidth=".7" opacity=".6"/><path d="M190 5a6 6 0 0 0 0 12 3.5 6 0 0 1 0-12Z" fill="currentColor" stroke="none" opacity=".6"/><circle cx="224" cy="11" r="6" strokeWidth=".7" opacity=".5"/><path d="M224 5a6 6 0 0 0 0 12Z" fill="currentColor" stroke="none" opacity=".5"/><circle cx="258" cy="11" r="6" strokeWidth=".7" opacity=".4"/></svg>;
}

export function ZodiacWheel(props: IconProps) {
  return <svg viewBox="0 0 520 520" fill="none" stroke="currentColor" strokeWidth=".8" aria-hidden="true" {...props}><circle cx="260" cy="260" r="255"/><circle cx="260" cy="260" r="238" strokeWidth=".4"/><circle cx="260" cy="260" r="195" strokeWidth=".5"/><circle cx="260" cy="260" r="150" strokeWidth=".4"/><circle cx="260" cy="260" r="90"/><circle cx="260" cy="260" r="40" strokeWidth=".4"/><g strokeWidth=".4"><path d="M260 5v510M5 260h510"/><path d="M260 5v510" transform="rotate(30 260 260)"/><path d="M260 5v510" transform="rotate(60 260 260)"/><path d="M260 5v510" transform="rotate(120 260 260)"/><path d="M260 5v510" transform="rotate(150 260 260)"/></g><path d="m260 220 17 30-17 30-17-30 17-30Zm0 20 17-30-17-30-17 30 17 30Z" strokeWidth=".6"/></svg>;
}

export function Constellation({ variant = 1, ...props }: IconProps & { variant?: 1 | 2 }) {
  return variant === 1 ? <svg viewBox="0 0 180 140" fill="none" stroke="currentColor" aria-hidden="true" {...props}><polyline points="12,108 44,72 78,86 112,40 150,52 168,18" strokeWidth=".4" opacity=".4"/><circle cx="44" cy="72" r="2.2" fill="currentColor" stroke="none" opacity=".6"/><circle cx="112" cy="40" r="2.6" fill="currentColor" stroke="none" opacity=".7"/><circle cx="168" cy="18" r="1.8" fill="currentColor" stroke="none" opacity=".55"/></svg> : <svg viewBox="0 0 200 160" fill="none" stroke="currentColor" aria-hidden="true" {...props}><polyline points="18,30 56,58 94,44 128,84 110,124 158,140" strokeWidth=".4" opacity=".4"/><circle cx="94" cy="44" r="2.4" fill="currentColor" stroke="none" opacity=".7"/><circle cx="128" cy="84" r="1.5" fill="currentColor" stroke="none" opacity=".5"/><circle cx="158" cy="140" r="2.2" fill="currentColor" stroke="none" opacity=".6"/></svg>;
}

export function CloseIcon(props: IconProps) { return <svg viewBox="0 0 14 14" fill="none" aria-hidden="true" {...props}><path d="m1 1 12 12M13 1 1 13" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round"/></svg>; }
export function ChevronIcon(props: IconProps) { return <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}><path d="m6 9 6 6 6-6" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round"/></svg>; }
export function SendIcon(props: IconProps) { return <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}><path d="M4 12h16m-7-7 7 7-7 7" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round"/></svg>; }
export function LogoutIcon(props: IconProps) { return <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}><path d="M13 4h4a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-4M10 16l4-4-4-4m4 4H5" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round"/></svg>; }

export function SphereIcon({ type, ...props }: IconProps & { type: 'love' | 'career' | 'health' | 'spirit' }) {
  const art = {
    love: <><path d="M14 3C8 3 3 8 3 14s5 11 11 11 11-5 11-11"/><circle cx="14" cy="14" r="4"/></>,
    career: <><rect x="6" y="8" width="16" height="13" rx="2"/><path d="M10 8V6q4-3 8 0v2"/></>,
    health: <><ellipse cx="14" cy="14" rx="8" ry="10"/><path d="M6 14q4-4 8 0t8 0"/></>,
    spirit: <polygon points="14,3 17,11 26,11 19,16 22,24 14,19 6,24 9,16 2,11 11,11"/>,
  }[type];
  return <svg viewBox="0 0 28 28" fill="none" stroke="currentColor" strokeWidth=".8" aria-hidden="true" {...props}>{art}</svg>;
}
