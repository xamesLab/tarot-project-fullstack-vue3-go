import { LogoSigil } from '@/shared/icons';

export function Logo({ size = 'md', tagline }: { size?: 'sm' | 'md' | 'lg'; tagline?: string }) {
  const markSize = { sm: 28, md: 34, lg: 52 }[size];
  return <div className={`logo logo-${size}`} aria-label="NEBULA-TAROT"><LogoSigil className="logo-mark" width={markSize} height={markSize}/><div><div className="logo-word">NEBULA-TAROT</div>{tagline && <div className="logo-tagline">{tagline}</div>}</div></div>;
}
