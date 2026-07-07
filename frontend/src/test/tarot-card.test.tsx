import { fireEvent, render, screen, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { TarotCard } from '@/shared/ui/TarotCard';

const props = {
  roman: 'II',
  name: 'The High Priestess',
  faceImageUrl: 'https://assets.example/face.png',
  backImageUrl: 'https://assets.example/back.png',
};

describe('TarotCard', () => {
  it('starts closed, opens on click, and closes on the next click', () => {
    const { container } = render(<TarotCard {...props} flippable flipLabel="Flip card"/>);
    const card = screen.getByRole('button', { name: 'Flip card' });
    const back = container.querySelector<HTMLElement>('[data-card-side="back"]');

    expect(card).toHaveAttribute('aria-pressed', 'false');
    expect(card).not.toHaveClass('flipped');
    expect(card.firstElementChild).toBe(back);
    expect(within(back!).getByRole('presentation')).toHaveAttribute('src', props.backImageUrl);

    fireEvent.click(card);
    expect(card).toHaveAttribute('aria-pressed', 'true');
    expect(card).toHaveClass('flipped');

    fireEvent.click(card);
    expect(card).toHaveAttribute('aria-pressed', 'false');
    expect(card).not.toHaveClass('flipped');
  });

  it('supports Enter and Space keyboard flips', () => {
    const onFlipChange = vi.fn();
    render(<TarotCard {...props} flippable flipLabel="Flip card" onFlipChange={onFlipChange}/>);
    const card = screen.getByRole('button', { name: 'Flip card' });

    fireEvent.keyDown(card, { key: 'Enter' });
    expect(card).toHaveAttribute('aria-pressed', 'true');
    expect(onFlipChange).toHaveBeenLastCalledWith(true);
    fireEvent.keyDown(card, { key: ' ' });
    expect(card).toHaveAttribute('aria-pressed', 'false');
    expect(onFlipChange).toHaveBeenLastCalledWith(false);
  });

  it('shows the face immediately when it is not flippable', () => {
    const { container } = render(<TarotCard {...props}/>);

    expect(container.querySelector('[data-card-side="back"]')).not.toBeInTheDocument();
    expect(screen.getByRole('img', { name: props.name })).toHaveAttribute('src', props.faceImageUrl);
  });

  it('uses the SVG face when the image cannot be loaded', () => {
    const { container } = render(<TarotCard {...props}/>);
    fireEvent.error(screen.getByRole('img', { name: props.name }));

    const face = container.querySelector<HTMLElement>('[data-card-side="face"]');
    expect(within(face!).queryByRole('img', { name: props.name })).not.toBeInTheDocument();
    expect(face?.querySelector('svg.card-art')).toBeInTheDocument();
    expect(face?.querySelector('.tarot-frame')).toBeInTheDocument();
    expect(face).toHaveTextContent(props.name);
  });
});
