'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { clearStoredToken } from '@/lib/api';
import { Search, LogOut, CheckCircle2 } from 'lucide-react';

interface NavbarProps {
  userName?: string;
  isConnected?: boolean;
}

export function Navbar({ userName = 'Developer', isConnected = true }: NavbarProps) {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');

  const handleLogout = () => {
    clearStoredToken();
    router.push('/login');
  };

  return (
    <header className="h-16 border-b border-[#1E2230] bg-[#0D0F17] px-6 flex items-center justify-between sticky top-0 z-50">
      <div className="flex items-center gap-6">
        <Link href="/dashboard" className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-blue-600 flex items-center justify-center font-bold text-white shadow-lg shadow-blue-500/20">
            CF
          </div>
          <span className="font-semibold text-lg tracking-tight text-white">
            Collab<span className="text-blue-500">Flow</span>
          </span>
        </Link>

        {/* Real-Time WebSocket Status Indicator */}
        <div className="hidden sm:flex items-center gap-2 px-2.5 py-1 rounded-full text-xs font-medium border border-[#1E2230] bg-[#12141C]">
          <span
            className={`w-2 h-2 rounded-full ${
              isConnected ? 'bg-emerald-500 shadow-sm shadow-emerald-500/50' : 'bg-rose-500 animate-pulse'
            }`}
          />
          <span className="text-slate-400">
            {isConnected ? 'Real-Time Sync Active' : 'Connecting to Hub...'}
          </span>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {/* Search */}
        <div className="relative hidden md:block w-64">
          <Search className="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search tasks, docs..."
            className="w-full bg-[#12141C] border border-[#1E2230] text-sm text-slate-200 pl-9 pr-4 py-1.5 rounded-lg focus:outline-none focus:border-blue-500 transition-colors placeholder:text-slate-500"
          />
        </div>

        {/* User Pill */}
        <div className="flex items-center gap-3 pl-2 border-l border-[#1E2230]">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-full bg-slate-700 flex items-center justify-center text-xs font-semibold text-white uppercase border border-slate-600">
              {userName.substring(0, 2)}
            </div>
            <span className="text-sm font-medium text-slate-300 hidden lg:inline-block">
              {userName}
            </span>
          </div>

          <button
            onClick={handleLogout}
            title="Log out"
            className="p-1.5 text-slate-400 hover:text-rose-400 hover:bg-rose-500/10 rounded-lg transition-colors"
          >
            <LogOut className="w-4 h-4" />
          </button>
        </div>
      </div>
    </header>
  );
}
