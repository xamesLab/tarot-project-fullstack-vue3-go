import { describe, expect, it } from 'vitest';
import { mockDailyCardApi } from '@mock/daily-card';

describe('mock daily card API', () => {
  it('generates a guest card and resolves the same card after claiming it', async () => {
    const guest = await mockDailyCardApi.generateGuest({
      guest_id: 'guest-test-id',
      timezone: 'UTC',
    });

    expect(guest.success).toBe(true);
    expect(guest.source).toBe('guest');
    expect(guest.reading.cards).toHaveLength(1);
    const faceImageUrl = guest.reading.cards[0].image_url;
    const backImageUrl = guest.reading.cards[0].back_image_url ?? '';
    expect(new URL(faceImageUrl).pathname).toContain('mock-card-face.png');
    expect(new URL(backImageUrl).pathname).toContain('mock-card-back.png');

    const claimed = await mockDailyCardApi.resolve({
      timezone: 'UTC',
      claim_token: guest.claim_token,
    });

    expect(claimed.source).toBe('guest_claimed');
    expect(claimed.reading.id).toBeTruthy();
    expect(claimed.reading.cards[0].key).toBe(guest.reading.cards[0].key);
  });

  it('returns localized interpretations for guest and authenticated readings', async () => {
    const guest = await mockDailyCardApi.generateGuest({
      guest_id: 'guest-localization-test',
      timezone: 'Europe/Moscow',
    });
    const guestInterpretation = await mockDailyCardApi.getGuestInterpretation({
      claim_token: guest.claim_token,
      language: 'ru',
    });

    expect(guestInterpretation.interpretation.general).toMatch(/\p{Script=Cyrillic}/u);
    expect(guestInterpretation.interpretation.general).not.toBe(guestInterpretation.interpretation.message);
    expectInterpretationWordLimits(guestInterpretation.interpretation);

    const account = await mockDailyCardApi.resolve({ timezone: 'UTC' });
    const accountInterpretation = await mockDailyCardApi.getUserInterpretation({
      reading_id: account.reading.id ?? '',
      language: 'en',
    });

    expect(accountInterpretation.interpretation.general).toMatch(/[A-Za-z]/);
    expect(accountInterpretation.interpretation.general).not.toBe(accountInterpretation.interpretation.message);
    expectInterpretationWordLimits(accountInterpretation.interpretation);
  });
});

function expectInterpretationWordLimits(interpretation: {
  general: string;
  message: string;
  relationships: string;
  finances: string;
  emotional_state: string;
  career_business: string;
}) {
  const generalWordCount = interpretation.general.trim().split(/\s+/).filter(Boolean).length;
  expect(generalWordCount, 'general should contain no more than 20 words').toBeLessThanOrEqual(20);
  const longBlocks = Object.entries(interpretation).filter(([block]) => block !== 'general');
  for (const [block, text] of longBlocks) {
    const wordCount = text.trim().split(/\s+/).filter(Boolean).length;
    expect(wordCount, `${block} should contain at least 50 words`).toBeGreaterThanOrEqual(50);
  }
}
