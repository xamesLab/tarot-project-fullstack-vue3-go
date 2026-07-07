import type { ReactNode } from 'react';

export function UiKitSection({
  number,
  title,
  description,
  children,
  id,
}: {
  number: string;
  title: string;
  description?: string;
  children: ReactNode;
  id: string;
}) {
  return (
    <section className="kit-section" id={id}>
      <header className="kit-section-head">
        <span>{number}</span>
        <h2>{title}</h2>
      </header>
      {description && <p className="kit-description">{description}</p>}
      <div className="kit-demo">{children}</div>
    </section>
  );
}

export function DemoRow({ children }: { children: ReactNode }) {
  return <div className="demo-row">{children}</div>;
}
export function CodeCaption({ children }: { children: ReactNode }) {
  return <div className="kit-caption">{children}</div>;
}
