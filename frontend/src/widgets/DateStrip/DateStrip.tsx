import { MoonPhases } from '@/shared/icons';

export function DateStrip({ label }: { date: Date; label: string }) { return <div className="date-strip"><div className="moon-row"><MoonPhases width="290" height="22"/></div><div className="date-line">{label}</div></div>; }
