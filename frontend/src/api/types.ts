export type Language = 'ru' | 'en';
export type Theme = 'dark' | 'light';
export type UserRole = 'seeker' | 'adept' | 'master' | 'keeper';

export interface TarotCardDto {
  id: string;
  arcana: 'major' | 'minor';
  roman: string;
  name: string;
  illustrationLabel: string;
  subtitle: string;
  imageUrl: string;
  backImageUrl: string;
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

export interface CardOfDayInterpretationDto {
  general: string;
  work: string;
  family: string;
  advice: string;
}

export interface CardOfDayInterpretationResponse {
  success: true;
  source: 'cache' | 'generated' | 'reading';
  interpretation: CardOfDayInterpretationDto;
}

export interface RegisterPayload {
  email: string;
  username?: string;
  password: string;
  timezone: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface UserDto {
  id: string;
  username: string;
  email: string;
  role: UserRole;
  status: 'active' | 'disabled';
  timezone: string;
  language: Language;
  theme: Theme;
  selected_deck_id?: string;
  selected_back_id?: string;
  created_at?: string;
  last_login_at?: string;
}

export interface AuthResponse {
  success: true;
  user: UserDto;
  token: string;
}

export interface CurrentUserResponse {
  success: true;
  user: UserDto;
}

export interface ApiErrorDetails {
  field?: string;
  reason?: string;
}

export interface BackendReadingCardDto {
  key: string;
  arcana: 'major' | 'minor';
  type: string;
  type_ru: string;
  name: string;
  name_ru: string;
  interpretation: string;
  interpretation_ru: string;
  tags: string;
  tags_ru: string;
  image_url: string;
  back_image_url?: string;
  position: number;
  reversed: boolean;
  viewed: boolean;
}

export interface BackendReadingDto {
  id?: string;
  type: 'card_of_day' | 'three_card';
  accounting_date: string;
  generated_at: string;
  cards: BackendReadingCardDto[];
  deck_id?: string;
  back_id?: string;
  next_available_at?: string;
}

export interface GuestCardOfDayResponse {
  success: true;
  source: 'guest';
  reading: BackendReadingDto;
  claim_token: string;
}

export interface ResolveCardOfDayResponse {
  success: true;
  source: 'existing' | 'guest_claimed' | 'generated';
  reading: BackendReadingDto;
}

export interface DailyCardApi {
  generateGuest: (payload: { guest_id: string; timezone: string }) => Promise<GuestCardOfDayResponse>;
  resolve: (payload: { timezone?: string; claim_token?: string }) => Promise<ResolveCardOfDayResponse>;
  getGuestInterpretation: (payload: { claim_token: string; language: Language }) => Promise<CardOfDayInterpretationResponse>;
  getUserInterpretation: (payload: { reading_id: string; language: Language }) => Promise<CardOfDayInterpretationResponse>;
}

export interface NebulaTarotApi {
  readings: { getDaily: (language: Language) => Promise<DailyReadingDto> };
  oracle: { ask: (question: string, language: Language) => Promise<OracleReplyDto> };
  auth: {
    register: (payload: RegisterPayload) => Promise<AuthResponse>;
    login: (payload: LoginPayload) => Promise<AuthResponse>;
    getCurrentUser: () => Promise<CurrentUserResponse>;
  };
}
