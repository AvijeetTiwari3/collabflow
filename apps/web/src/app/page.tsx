import Link from 'next/link';
import { Sparkles, Zap, Users, ArrowRight, Layers, Lock, GitPullRequest } from 'lucide-react';

export default function HomePage() {
  return (
    <div className="min-h-screen bg-[#090A0F] text-slate-100 flex flex-col justify-between">
      {/* Navigation */}
      <nav className="h-20 border-b border-[#1E2230] px-8 flex items-center justify-between max-w-7xl mx-auto w-full">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-xl bg-blue-600 flex items-center justify-center font-bold text-white shadow-lg shadow-blue-500/30">
            CF
          </div>
          <span className="font-bold text-xl tracking-tight">
            Collab<span className="text-blue-500">Flow</span>
          </span>
        </div>

        <div className="flex items-center gap-4">
          <Link
            href="/login"
            className="text-sm font-medium text-slate-300 hover:text-white px-4 py-2 transition-colors"
          >
            Sign In
          </Link>
          <Link
            href="/register"
            className="text-sm font-semibold bg-blue-600 hover:bg-blue-500 text-white px-5 py-2.5 rounded-xl shadow-lg shadow-blue-500/25 transition-all hover:scale-105"
          >
            Get Started Free
          </Link>
        </div>
      </nav>

      {/* Hero Section */}
      <main className="max-w-6xl mx-auto px-8 py-20 text-center space-y-8">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-blue-500/30 bg-blue-500/10 text-xs font-semibold text-blue-400">
          <Sparkles className="w-3.5 h-3.5" />
          <span>Google SDE-Grade Distributed Architecture</span>
        </div>

        <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight text-white max-w-4xl mx-auto leading-tight">
          Real-time collaborative workspaces at <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-indigo-500">hyperscale</span>.
        </h1>

        <p className="text-lg md:text-xl text-slate-400 max-w-2xl mx-auto leading-relaxed">
          Engineered with Next.js 14, Go, WebSockets, Redis Pub/Sub, and CRDTs. 
          Experience conflict-free live document editing and O(1) fractional-indexed Kanban boards.
        </p>

        <div className="flex flex-col sm:flex-row items-center justify-center gap-4 pt-4">
          <Link
            href="/dashboard"
            className="w-full sm:w-auto flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-500 text-white font-semibold px-8 py-4 rounded-xl shadow-xl shadow-blue-500/20 transition-all hover:scale-105 text-base"
          >
            <span>Launch Live Workspace</span>
            <ArrowRight className="w-5 h-5" />
          </Link>
          <Link
            href="/register"
            className="w-full sm:w-auto flex items-center justify-center gap-2 bg-[#12141C] hover:bg-[#1E2230] text-slate-300 font-semibold px-8 py-4 rounded-xl border border-[#1E2230] transition-colors text-base"
          >
            Create Free Account
          </Link>
        </div>

        {/* Feature Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 pt-16 text-left">
          <div className="p-6 bg-[#0D0F17] border border-[#1E2230] rounded-2xl space-y-3">
            <div className="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center text-blue-400">
              <Zap className="w-5 h-5" />
            </div>
            <h3 className="text-lg font-bold text-white">CRDT-Powered Real-Time</h3>
            <p className="text-sm text-slate-400 leading-relaxed">
              Conflict-Free Replicated Data Types (Yjs) enable simultaneous multi-user document editing with mathematically proven convergence.
            </p>
          </div>

          <div className="p-6 bg-[#0D0F17] border border-[#1E2230] rounded-2xl space-y-3">
            <div className="w-10 h-10 rounded-xl bg-indigo-500/10 flex items-center justify-center text-indigo-400">
              <Layers className="w-5 h-5" />
            </div>
            <h3 className="text-lg font-bold text-white">O(1) Fractional Indexing</h3>
            <p className="text-sm text-slate-400 leading-relaxed">
              Figma-inspired fractional index drag-and-drop ordering. Reorders single tasks in O(1) without expensive O(N) database bulk updates.
            </p>
          </div>

          <div className="p-6 bg-[#0D0F17] border border-[#1E2230] rounded-2xl space-y-3">
            <div className="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center text-emerald-400">
              <Users className="w-5 h-5" />
            </div>
            <h3 className="text-lg font-bold text-white">Multi-Pod WebSocket Fan-Out</h3>
            <p className="text-sm text-slate-400 leading-relaxed">
              Horizontally scalable Go WebSocket hubs synchronized via Redis Pub/Sub, delivering sub-30ms presence and cursor tracking.
            </p>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-[#1E2230] py-8 text-center text-xs text-slate-500">
        <p>CollabFlow Engine • Built with Go 1.22, Next.js 14, WebSockets & PostgreSQL</p>
      </footer>
    </div>
  );
}
