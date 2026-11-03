'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { getStoredToken } from '@/lib/api';

export interface WSMessage {
  type: string;
  room_id?: string;
  user_id?: string;
  user_name?: string;
  avatar_url?: string;
  payload?: any;
  timestamp?: string;
}

export function useWebSocket(roomId: string, onMessage?: (msg: WSMessage) => void) {
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>();
  const reconnectAttempts = useRef(0);

  const connect = useCallback(() => {
    if (!roomId) return;

    const token = getStoredToken() || '';
    const wsBase = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';
    const wsUrl = `${wsBase}?room=${encodeURIComponent(roomId)}&token=${encodeURIComponent(token)}`;

    try {
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        setIsConnected(true);
        reconnectAttempts.current = 0;
      };

      ws.onmessage = (event) => {
        try {
          const lines = event.data.split('\n');
          for (const line of lines) {
            if (!line.trim()) continue;
            const data: WSMessage = JSON.parse(line);
            onMessage?.(data);
          }
        } catch (err) {
          console.error('[ws] parsing error:', err);
        }
      };

      ws.onclose = () => {
        setIsConnected(false);
        // Exponential backoff auto-reconnect (max 10s)
        const delay = Math.min(1000 * Math.pow(1.5, reconnectAttempts.current), 10000);
        reconnectAttempts.current += 1;
        reconnectTimeoutRef.current = setTimeout(connect, delay);
      };

      ws.onerror = (err) => {
        console.warn('[ws] connection error:', err);
        ws.close();
      };
    } catch (e) {
      console.error('[ws] failed to initiate:', e);
    }
  }, [roomId, onMessage]);

  useEffect(() => {
    connect();

    return () => {
      if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);
      if (wsRef.current) {
        wsRef.current.onclose = null; // Prevent reconnect loop on intentional unmount
        wsRef.current.close();
      }
    };
  }, [connect]);

  const sendMessage = useCallback((msg: Omit<WSMessage, 'timestamp'>) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    }
  }, []);

  return { isConnected, sendMessage };
}
