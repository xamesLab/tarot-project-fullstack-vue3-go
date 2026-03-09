export type CardState = "open" | "closed";

export interface ReadingCard {
  id: string;
  name: string;
  state: CardState;
  imageUrl: string;
  backImageUrl: string;
  description: string;
  interpretation: string;
}

export interface TodayReading {
  date: string;
  timezoneUsed: string;
  cardOfDay: ReadingCard;
  spreadA: ReadingCard[];
  spreadB: ReadingCard[];
  spreadGenerationsUsed: number;
  spreadGenerationsLimit?: number;
  spreadGenerationsRemaining?: number;
  canGenerateMore: boolean;
}
