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
    const backImage = within(back!).getByRole('presentation');
    expect(backImage).toHaveAttribute('src', props.backImageUrl);
    expect(backImage).toHaveStyle({ objectFit: 'cover', objectPosition: 'center center' });

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

  it('supports server-controlled one-way reveals', () => {
    const onReveal = vi.fn();
    const { rerender } = render(<TarotCard {...props} flippable revealed={false} flipLabel="Reveal card" onReveal={onReveal}/>);
    const card = screen.getByRole('button', { name: 'Reveal card' });

    fireEvent.click(card);
    expect(onReveal).toHaveBeenCalledTimes(1);
    expect(card).not.toHaveClass('flipped');

    rerender(<TarotCard {...props} flippable revealed flipLabel="Revealed card" onReveal={onReveal}/>);
    expect(screen.getByRole('button', { name: 'Revealed card' })).toHaveClass('flipped');
    fireEvent.click(screen.getByRole('button', { name: 'Revealed card' }));
    expect(onReveal).toHaveBeenCalledTimes(1);
  });

  it('shows the face immediately when it is not flippable', () => {
    const { container } = render(<TarotCard {...props}/>);

    expect(container.querySelector('[data-card-side="back"]')).not.toBeInTheDocument();
    expect(screen.getByRole('img', { name: props.name })).toHaveAttribute('src', props.faceImageUrl);
  });

  it('matches the card container to the loaded face image proportions', () => {
    const { container } = render(<TarotCard {...props}/>);
    const image = screen.getByRole('img', { name: props.name });
    Object.defineProperties(image, {
      naturalWidth: { configurable: true, value: 340 },
      naturalHeight: { configurable: true, value: 562 },
    });

    fireEvent.load(image);

    expect(container.querySelector('.tarot')).toHaveStyle({ aspectRatio: '340 / 562' });
  });

  it('matches fallback art proportions when face and back images are unavailable', () => {
    const { container } = render(<TarotCard roman="II" name={props.name} flippable/>);
    const card = container.querySelector('.tarot');
    const artworks = container.querySelectorAll('svg.card-art');

    expect(card).toHaveStyle({ aspectRatio: '260 / 450' });
    expect(artworks).toHaveLength(2);
    artworks.forEach((artwork) => expect(artwork).toHaveAttribute('viewBox', '0 0 260 450'));
    expect(artworks[0]).toHaveAttribute('preserveAspectRatio', 'xMidYMid slice');
  });

  it('crops the SVG back fallback to fill a card sized by its face image', () => {
    const { container } = render(<TarotCard {...props} flippable/>);
    const faceImage = screen.getByRole('img', { name: props.name });
    Object.defineProperties(faceImage, {
      naturalWidth: { configurable: true, value: 340 },
      naturalHeight: { configurable: true, value: 562 },
    });
    fireEvent.load(faceImage);

    const backImage = container.querySelector<HTMLImageElement>('.card-image-back');
    fireEvent.error(backImage!);

    expect(container.querySelector('.tarot')).toHaveStyle({ aspectRatio: '340 / 562' });
    expect(container.querySelector('[data-card-side="back"] svg.card-art')).toHaveAttribute('preserveAspectRatio', 'xMidYMid slice');
  });

  it('uses the SVG face when the image cannot be loaded', () => {
    const { container } = render(<TarotCard {...props}/>);
    fireEvent.error(screen.getByRole('img', { name: props.name }));

    const face = container.querySelector<HTMLElement>('[data-card-side="face"]');
    expect(container.querySelector('.tarot')).toHaveStyle({ aspectRatio: '260 / 450' });
    expect(within(face!).queryByRole('img', { name: props.name })).not.toBeInTheDocument();
    expect(face?.querySelector('svg.card-art')).toBeInTheDocument();
    expect(face?.querySelector('svg.card-art')).toHaveAttribute('viewBox', '0 0 260 450');
    expect(face?.querySelector('.tarot-frame')).toBeInTheDocument();
    expect(face).toHaveTextContent(props.name);
  });
});
