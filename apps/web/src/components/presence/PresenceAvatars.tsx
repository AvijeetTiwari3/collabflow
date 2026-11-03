'use client';

import { PresenceUser } from '@/hooks/usePresence';

interface PresenceAvatarsProps {
  users: PresenceUser[];
}

export function PresenceAvatars({ users }: PresenceAvatarsProps) {
  if (!users || users.length === 0) return null;

  return (
    <div className="flex items-center gap-2">
      <span className="text-xs font-medium text-slate-400 hidden sm:inline-block">Live:</span>
      <div className="flex items-center -space-x-2">
        {users.map((user, idx) => (
          <div
            key={user.id || idx}
            title={`${user.name} (Active)`}
            className="relative group cursor-pointer"
          >
            <div
              style={{ borderColor: user.color || '#3B82F6' }}
              className="w-7 h-7 rounded-full bg-slate-800 border-2 flex items-center justify-center text-[10px] font-bold text-white shadow-md transition-transform group-hover:scale-110 group-hover:z-10 uppercase"
            >
              {user.name ? user.name.substring(0, 2) : 'U'}
            </div>

            {/* Hover tooltip */}
            <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-1.5 hidden group-hover:block bg-[#1E2230] text-slate-200 text-xs px-2 py-0.5 rounded shadow-lg whitespace-nowrap z-20">
              {user.name}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
