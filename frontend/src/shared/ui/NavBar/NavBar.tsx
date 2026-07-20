export interface NavItem { id: string; label: string; href?: string; }

export function NavBar({ items, active, onSelect, label = 'Navigation' }: { items: NavItem[]; active: string; onSelect: (item: NavItem, trigger: HTMLElement) => void; label?: string }) {
  return (
    <nav className="site-nav" aria-label={label}>
      {items.map((item) => {
        const className = `nav-item ${item.id === active ? 'active' : ''}`;
        return item.href ? (
          <a
            key={item.id}
            className={className}
            href={item.href}
            aria-current={item.id === active ? 'page' : undefined}
            onClick={(event) => {
              if (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
              event.preventDefault();
              onSelect(item, event.currentTarget);
            }}
          >
            {item.label}
          </a>
        ) : (
          <button type="button" key={item.id} className={className} onClick={(event) => onSelect(item, event.currentTarget)}>
            {item.label}
          </button>
        );
      })}
    </nav>
  );
}
