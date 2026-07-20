import { useQuery } from '@tanstack/react-query';
import { api } from '@/api';
import type { UserDto } from '@/api/types';
import { useAuthStore } from '@/features/auth/store';

export function useProfile() {
  const sessionUser = useAuthStore((state) => state.user);
  const initialProfile = sessionUser?.created_at ? sessionUser : undefined;

  return useQuery<UserDto>({
    queryKey: ['current-user-profile', sessionUser?.id],
    queryFn: async () => (await api.auth.getCurrentUser()).user,
    enabled: Boolean(sessionUser?.id),
    initialData: initialProfile,
    staleTime: 5 * 60 * 1000,
  });
}
