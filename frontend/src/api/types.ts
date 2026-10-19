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
  message: string;
  relationships: string;
  finances: string;
  emotional_state: string;
  career_business: string;
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

export interface UpdateProfilePayload {
  username: string;
  timezone: string;
}

export interface UpdateSettingsPayload {
  language?: Language;
  theme?: Theme;
  selected_deck_id?: string;
  selected_back_id?: string;
}

export interface LocalizedTextDto {
  en: string;
  ru: string;
}

export interface DeckDto {
  id: string;
  name: LocalizedTextDto;
  description: LocalizedTextDto;
  author: string;
  year: number;
  tags: string[];
  status: string;
  card_count: number;
  back_image_path: string;
  card_image_base_path?: string;
  created_at?: string;
  updated_at?: string;
}

export interface CardBackDto {
  id: string;
  name: LocalizedTextDto;
  description: LocalizedTextDto;
  tags: string[];
  status: string;
  image_path: string;
  created_at?: string;
  updated_at?: string;
}

export interface PaginationDto {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface AdminDashboardMetricsDto {
  users: {
    total: number;
    by_role: Record<UserRole, number>;
    by_status: { active: number; disabled: number };
  };
  readings: {
    today: { total: number; card_of_day: number; three_card: number };
    this_week: { total: number; card_of_day: number; three_card: number };
  };
  decks: {
    total: number;
    by_status: { done: number; processing: number; disabled: number; archived: number };
  };
}

export interface AdminDashboardResponse {
  success: true;
  metrics: AdminDashboardMetricsDto;
}

export interface AdminUserDto extends UserDto {
  updated_at?: string;
}

export interface AdminUsersParams {
  search?: string;
  role?: UserRole | 'all';
  status?: 'active' | 'disabled' | 'all';
  page: number;
  limit: number;
}

export interface AdminUsersResponse {
  success: true;
  users: AdminUserDto[];
  pagination: PaginationDto;
}

export interface AdminRolePolicyDto {
  id: string;
  role: UserRole;
  display_name: LocalizedTextDto;
  three_card_limit: number;
  ai_enabled: boolean;
  ai_daily_limit: number;
  history_retention_days: number;
  created_at?: string;
  updated_at?: string;
}

export type UpdateAdminRolePolicyPayload = Pick<
  AdminRolePolicyDto,
  'three_card_limit' | 'ai_enabled' | 'ai_daily_limit' | 'history_retention_days'
>;

export interface AdminRolePoliciesResponse {
  success: true;
  policies: AdminRolePolicyDto[];
}

export interface AdminSettingsDto {
  id: string;
  maintenance_mode: boolean;
  maintenance_message: LocalizedTextDto;
  ai_enabled: boolean;
  default_role: Exclude<UserRole, 'keeper'>;
  default_deck_id: string;
  default_back_id: string;
  registration_enabled: boolean;
  created_at?: string;
  updated_at?: string;
}

export type UpdateAdminSettingsPayload = Omit<AdminSettingsDto, 'id' | 'created_at' | 'updated_at'>;

export interface AdminSettingsResponse {
  success: true;
  settings: AdminSettingsDto;
}

export interface AdminAssetsParams {
  scope: 'active' | 'archived';
  page: number;
  limit: number;
}

export interface AdminDecksResponse {
  success: true;
  decks: DeckDto[];
  pagination: PaginationDto;
}

export interface AdminCardBacksResponse {
  success: true;
  card_backs: CardBackDto[];
  pagination: PaginationDto;
}

export interface AdminApi {
  getDashboard: () => Promise<AdminDashboardResponse>;
  getUsers: (params: AdminUsersParams) => Promise<AdminUsersResponse>;
  updateUserRole: (userId: string, role: UserRole) => Promise<{ success: true; user: AdminUserDto }>;
  updateUserStatus: (userId: string, status: 'active' | 'disabled') => Promise<{ success: true; user: AdminUserDto }>;
  getRolePolicies: () => Promise<AdminRolePoliciesResponse>;
  updateRolePolicy: (role: UserRole, payload: UpdateAdminRolePolicyPayload) => Promise<{ success: true; policy: AdminRolePolicyDto }>;
  getSettings: () => Promise<AdminSettingsResponse>;
  updateSettings: (payload: UpdateAdminSettingsPayload) => Promise<AdminSettingsResponse>;
  getDecks: (params: AdminAssetsParams) => Promise<AdminDecksResponse>;
  createDeck: (payload: FormData) => Promise<{ success: true; deck: DeckDto; message?: string }>;
  archiveDeck: (deckId: string) => Promise<{ success: true; deck: DeckDto }>;
  restoreDeck: (deckId: string) => Promise<{ success: true; deck: DeckDto }>;
  deleteDeck: (deckId: string) => Promise<{ success: true }>;
  getCardBacks: (params: AdminAssetsParams) => Promise<AdminCardBacksResponse>;
  createCardBack: (payload: FormData) => Promise<{ success: true; card_back: CardBackDto; message?: string }>;
  archiveCardBack: (cardBackId: string) => Promise<{ success: true; card_back: CardBackDto }>;
  restoreCardBack: (cardBackId: string) => Promise<{ success: true; card_back: CardBackDto }>;
  deleteCardBack: (cardBackId: string) => Promise<{ success: true }>;
}

export interface AvailableDecksResponse {
  success: true;
  decks: DeckDto[];
  default_deck_id: string;
}

export interface AvailableCardBacksResponse {
  success: true;
  card_backs: CardBackDto[];
  default_back_id: string;
}

export interface ApiErrorDetails {
  field?: string;
  reason?: string;
  current_limit?: RemainingTodayDto;
  next_reset_at?: string;
}

export interface RemainingTodayDto {
  count: number;
  limit: number;
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

export interface BackendAIInterpretationDto {
  request_id?: string;
  provider?: string;
  model?: string;
  prompt_version?: string;
  language?: Language;
  status?: string;
  content?: string;
  error_code?: string;
  error_message?: string;
  completed_at?: string;
  last_requested_at?: string;
  updated_at?: string;
}

export interface BackendReadingDto {
  id?: string;
  type: 'card_of_day' | 'three_card';
  accounting_date: string;
  generated_at: string;
  cards: BackendReadingCardDto[];
  deck_id?: string;
  back_id?: string;
  ai?: BackendAIInterpretationDto;
  next_available_at?: string;
  remaining_today?: RemainingTodayDto;
}

export interface ReadingHistoryItemDto extends BackendReadingDto {
  id: string;
}

export type ReadingHistoryFilter = 'all' | 'card_of_day' | 'three_card';

export interface ReadingHistoryParams {
  type: ReadingHistoryFilter;
  page: number;
  limit: number;
}

export interface ReadingHistoryResponse {
  success: true;
  readings: ReadingHistoryItemDto[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
  retention_info: {
    role: UserRole;
    retention_days: number;
    oldest_date: string;
  };
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

export interface AskRequestDto {
  id: string;
  reading_id: string;
  provider: string;
  model: string;
  prompt_version: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  response_excerpt?: string;
  error_code?: string;
  error_message?: string;
  requested_at: string;
  completed_at?: string;
}

export interface SpreadsApi {
  generateThreeCard: (payload: { timezone?: string }) => Promise<{ success: true; reading: BackendReadingDto }>;
  revealCard: (readingId: string, position: number) => Promise<{ success: true; reading: BackendReadingDto }>;
  requestInterpretation: (readingId: string) => Promise<{ success: true; request: AskRequestDto }>;
}

export interface NebulaTarotApi {
  readings: { getDaily: (language: Language) => Promise<DailyReadingDto> };
  oracle: { ask: (question: string, language: Language) => Promise<OracleReplyDto> };
  auth: {
    register: (payload: RegisterPayload) => Promise<AuthResponse>;
    login: (payload: LoginPayload) => Promise<AuthResponse>;
    getCurrentUser: () => Promise<CurrentUserResponse>;
  };
  profile: {
    updateProfile: (payload: UpdateProfilePayload) => Promise<CurrentUserResponse>;
    updateSettings: (payload: UpdateSettingsPayload) => Promise<CurrentUserResponse>;
  };
  catalog: {
    getAvailableDecks: () => Promise<AvailableDecksResponse>;
    getAvailableCardBacks: () => Promise<AvailableCardBacksResponse>;
  };
  history: {
    getHistory: (params: ReadingHistoryParams) => Promise<ReadingHistoryResponse>;
  };
  spreads: SpreadsApi;
  admin: AdminApi;
}
