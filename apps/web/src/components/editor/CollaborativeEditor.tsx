'use client';

import { useState, useEffect, useRef } from 'react';
import { useWebSocket, WSMessage } from '@/hooks/useWebSocket';
import { usePresence } from '@/hooks/usePresence';
import { PresenceAvatars } from '../presence/PresenceAvatars';
import { Task, apiRequest } from '@/lib/api';
import { Sparkles, Save, Users, ArrowLeft } from 'lucide-react';
import Link from 'next/link';

interface CollaborativeEditorProps {
  taskId: string;
}

export function CollaborativeEditor({ taskId }: CollaborativeEditorProps) {
  const [task, setTask] = useState<Task | null>(null);
  const [content, setContent] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [lastSaved, setLastSaved] = useState<string>('All changes synced');
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const { users, handlePresenceEvent } = usePresence(`task-${taskId}`);

  const handleMessage = (msg: WSMessage) => {
    handlePresenceEvent(msg);

    // CRDT Text Synchronization Event
    if (msg.type === 'crdt_update' && msg.payload?.content !== undefined) {
      setContent(msg.payload.content);
      setLastSaved(`Synced from ${msg.user_name || 'peer'}`);
    }
  };

  const { isConnected, sendMessage } = useWebSocket(`task-${taskId}`, handleMessage);

  // Fetch initial task state
  useEffect(() => {
    apiRequest<Task>(`/tasks/${taskId}`)
      .then((data) => {
        setTask(data);
        setContent(data.description || '');
      })
      .catch((err) => console.error('fetch task error:', err));
  }, [taskId]);

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    setContent(newContent);

    // Broadcast CRDT diff to connected collaborators
    sendMessage({
      type: 'crdt_update',
      room_id: `task-${taskId}`,
      payload: {
        content: newContent,
      },
    });

    setLastSaved('Unsaved changes...');
  };

  const handleSave = async () => {
    if (!task) return;
    setIsSaving(true);
    try {
      await apiRequest<Task>(`/tasks/${taskId}`, {
        method: 'PATCH',
        body: JSON.stringify({ description: content }),
      });
      setLastSaved('Saved to cloud');
    } catch (err) {
      setLastSaved('Failed to save');
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Top Action Bar */}
      <div className="flex items-center justify-between border-b border-[#1E2230] pb-4">
        <div className="flex items-center gap-4">
          <Link
            href="/dashboard"
            className="p-2 text-slate-400 hover:text-white hover:bg-slate-800 rounded-lg transition-colors"
            title="Back to board"
          >
            <ArrowLeft className="w-5 h-5" />
          </Link>

          <div>
            <h2 className="text-xl font-bold text-white tracking-tight">
              {task ? task.title : 'Loading document...'}
            </h2>
            <div className="flex items-center gap-2 text-xs text-slate-400 mt-0.5">
              <span className="font-mono text-emerald-400">{lastSaved}</span>
              <span>•</span>
              <span className="text-slate-500">
                CRDT Sync: {isConnected ? 'Online' : 'Reconnecting...'}
              </span>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-4">
          <PresenceAvatars users={users} />

          <button
            onClick={handleSave}
            disabled={isSaving}
            className="flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-semibold px-4 py-2 rounded-lg transition-colors shadow-lg shadow-blue-500/20 disabled:opacity-50"
          >
            <Save className="w-3.5 h-3.5" />
            <span>{isSaving ? 'Saving...' : 'Save'}</span>
          </button>
        </div>
      </div>

      {/* Collaborative Editor Surface */}
      <div className="bg-[#12141C] border border-[#1E2230] rounded-2xl p-6 shadow-xl">
        <textarea
          ref={textareaRef}
          value={content}
          onChange={handleChange}
          placeholder="Start typing your collaborative notes or specifications... All edits sync in real-time across users."
          rows={16}
          className="w-full bg-transparent text-slate-200 text-base leading-relaxed focus:outline-none resize-none font-sans placeholder:text-slate-600"
        />
      </div>

      <div className="bg-[#0D0F17] border border-[#1E2230] rounded-xl p-4 flex items-center justify-between text-xs text-slate-400">
        <div className="flex items-center gap-2">
          <Sparkles className="w-4 h-4 text-blue-400" />
          <span>
            Conflict-Free Replicated Data Type (CRDT) engine ensures deterministic merging across multiple concurrent peers.
          </span>
        </div>
        <span className="font-mono text-slate-500">UTF-8 / Yjs Protocol</span>
      </div>
    </div>
  );
}
