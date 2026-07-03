export type Language = 'ru' | 'en';

export interface TarotCardDto {
  id: string;
  arcana: 'major' | 'minor';
  roman: string;
  name: string;
  illustrationLabel: string;
  subtitle: string;
}

export interface ReadingSphereDto {
  id: 'love' | 'career' | 'health' | 'spirit';
  title: string;
  text: string;
}

export interface DailyReadingDto {
  date: string;
  card: TarotCardDto;
  keywords: string[];
  brief: string;
  spheres: ReadingSphereDto[];
  message: { paragraph: string; closing: string; emphasis: string };
}

export interface OracleReplyDto {
  id: string;
  text: string;
}

export interface RegisterPayload {
  name: string;
  email: string;
  password: string;
  birthDate?: string;
}

export interface UserDto {
  id: string;
  name: string;
  email: string;
  birthDate?: string;
}

export interface NebulaTarotApi {
  readings: { getDaily: (language: Language) => Promise<DailyReadingDto> };
  oracle: { ask: (question: string, language: Language) => Promise<OracleReplyDto> };
  auth: { register: (payload: RegisterPayload) => Promise<UserDto> };
}
