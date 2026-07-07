export function Avatar({ initials, size = 'md', ring }: { initials: string; size?: 'sm' | 'md' | 'lg'; ring?: boolean }) { return <span className={`avatar avatar-${size} ${ring ? 'avatar-ring' : ''}`} aria-hidden="true">{initials}</span>; }
export function UserPill({ initials, name, sign, onClick }: { initials: string; name: string; sign: string; onClick?: () => void }) {
  const content = <><Avatar initials={initials} size="sm"/><span><span className="user-name">{name}</span><span className="user-sign">☽ {sign}</span></span></>;
  return onClick ? <button type="button" className="user-pill user-pill-button" onClick={onClick}>{content}</button> : <div className="user-pill">{content}</div>;
}
