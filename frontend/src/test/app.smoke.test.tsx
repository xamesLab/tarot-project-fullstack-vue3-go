import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { App } from '@/app/App';
import { router } from '@/app/router';
import { usePreferencesStore } from '@/features/preferences/store';

describe('NEBULA-TAROT smoke', () => {
  beforeEach(() => {
    localStorage.clear();
    usePreferencesStore.setState({ theme: 'dark', language: 'ru' });
    window.history.replaceState({}, '', '/');
  });

  it('renders the branded home and opens registration from navigation', async () => {
    render(<App />);
    expect(screen.getByLabelText('NEBULA-TAROT')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Карта дня' }));
    expect(await screen.findByRole('dialog')).toHaveTextContent('Присоединиться к NEBULA-TAROT');
  });

  it('switches language and theme', () => {
    render(<App />);
    fireEvent.click(screen.getByRole('button', { name: 'EN' }));
    expect(screen.getByRole('button', { name: 'Card of the Day' })).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Switch to light theme' }));
    expect(usePreferencesStore.getState().theme).toBe('light');
  });

  it('registers the UI Kit route in development', async () => {
    await router.navigate('/ui-kit');
    render(<App />);
    expect(await screen.findByText('NEBULA-TAROT · DESIGN SYSTEM')).toBeInTheDocument();
    expect(screen.getByText('Foundations · Palette')).toBeInTheDocument();
  });
});
