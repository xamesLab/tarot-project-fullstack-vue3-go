import { CelestialStrip } from '@/widgets/CelestialStrip/CelestialStrip';

/** @deprecated Use CelestialStrip directly. Kept as a compatibility wrapper. */
export function DateStrip({ date }: { date: Date; label?: string }) {
  return <CelestialStrip now={date} />;
}
