'use client';

import { useState, useEffect, useCallback } from 'react';
import { apiRequest } from '@/lib/api';
import { WSMessage } from './useWebSocket';

export interface PresenceUser {
  id: string;
  name: string;
  avatar_url: string;
  color?: string;
}

export function usePresence(roomId: string) {
  const [users, setUsers] = useState<PresenceUser[]>([]);

  // Fetch initial room presence via REST
  useEffect(() => {
    if (!roomId) return;
    apiRequest<PresenceUser[]>(`/rooms/${roomId}/presence`)
      .then((data) => setUsers(data || []))
      .catch(() => {});
  }, [roomId]);

  const handlePresenceEvent = useCallback((msg: WSMessage) => {
    if (msg.type === 'user_joined' && msg.user_id) {
      setUsers((prev) => {
        if (prev.some((u) => u.id === msg.user_id)) return prev;
        return [
          ...prev,
          {
            id: msg.user_id!,
            name: msg.user_name || 'Anonymous',
            avatar_url: msg.avatar_url || '',
          },
        ];
      });
    } else if (msg.type === 'user_left' && msg.user_id) {
      setUsers((prev) => prev.filter((u) => u.id !== msg.user_id));
    }
  }, []);

  return { users, handlePresenceEvent };
}
