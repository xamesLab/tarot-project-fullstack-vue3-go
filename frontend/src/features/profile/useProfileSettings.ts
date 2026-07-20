import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/api';
import type { CurrentUserResponse, UpdateProfilePayload, UpdateSettingsPayload } from '@/api/types';
import { useAuthStore } from '@/features/auth/store';

function useSyncedProfileMutation<TPayload>(mutationFn: (payload: TPayload) => Promise<CurrentUserResponse>) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn,
    onSuccess: ({ user }) => {
      useAuthStore.getState().syncUser(user);
      queryClient.setQueryData(['current-user-profile', user.id], user);
    },
  });
}

export function useUpdateProfile() {
  return useSyncedProfileMutation<UpdateProfilePayload>(api.profile.updateProfile);
}

export function useUpdateSettings() {
  return useSyncedProfileMutation<UpdateSettingsPayload>(api.profile.updateSettings);
}

export function useProfileCatalog(enabled = true) {
  return useQuery({
    queryKey: ['profile-catalog'],
    queryFn: async () => {
      const [decks, cardBacks] = await Promise.all([
        api.catalog.getAvailableDecks(),
        api.catalog.getAvailableCardBacks(),
      ]);
      return { decks, cardBacks };
    },
    enabled,
    staleTime: 30 * 60 * 1000,
  });
}
