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

export interface GatewayHealthResponse {
  status: "ok";
  uptime: number;
}

export interface GatewayHealthStatus {
  status: "success" | "error";
  title: string;
  message: string;
  uptime?: number;
}
