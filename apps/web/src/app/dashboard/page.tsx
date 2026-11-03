'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { Navbar } from '@/components/Navbar';
import { KanbanBoard } from '@/components/board/KanbanBoard';
import { PresenceAvatars } from '@/components/presence/PresenceAvatars';
import { useWebSocket, WSMessage } from '@/hooks/useWebSocket';
import { usePresence } from '@/hooks/usePresence';
import { Project, Task, apiRequest, getStoredToken } from '@/lib/api';
import { FolderKanban, Plus, Sparkles } from 'lucide-react';

export default function DashboardPage() {
  const router = useRouter();
  const [user, setUser] = useState<{ name: string; email: string } | null>(null);
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [newProjectName, setNewProjectName] = useState('');
  const [showNewProjectModal, setShowNewProjectModal] = useState(false);

  // Authenticate & load user
  useEffect(() => {
    const token = getStoredToken();
    if (!token) {
      router.push('/login');
      return;
    }

    apiRequest<{ name: string; email: string }>('/auth/me')
      .then(setUser)
      .catch(() => {
        router.push('/login');
      });
  }, [router]);

  // Load projects
  useEffect(() => {
    apiRequest<Project[]>('/projects')
      .then((data) => {
        setProjects(data);
        if (data.length > 0) {
          setSelectedProject(data[0]);
        }
      })
      .catch((err) => console.error('fetch projects error:', err));
  }, []);

  const currentProjectId = selectedProject?.id || '';

  // Presence hook
  const { users, handlePresenceEvent } = usePresence(currentProjectId);

  // Real-Time WebSocket event handler
  const handleWSMessage = useCallback(
    (msg: WSMessage) => {
      handlePresenceEvent(msg);

      if (msg.type === 'task_updated' && msg.payload) {
        const { task_id, status, position } = msg.payload;
        setTasks((prev) =>
          prev.map((t) => (t.id === task_id ? { ...t, status, position } : t))
        );
      } else if (msg.type === 'task_created' && msg.payload) {
        setTasks((prev) => {
          if (prev.some((t) => t.id === msg.payload.id)) return prev;
          return [...prev, msg.payload];
        });
      } else if (msg.type === 'task_deleted' && msg.payload) {
        setTasks((prev) => prev.filter((t) => t.id !== msg.payload.task_id));
      }
    },
    [handlePresenceEvent]
  );

  const { isConnected, sendMessage } = useWebSocket(currentProjectId, handleWSMessage);

  // Load tasks for selected project
  useEffect(() => {
    if (!currentProjectId) return;
    apiRequest<Task[]>(`/projects/${currentProjectId}/tasks`)
      .then(setTasks)
      .catch((err) => console.error('fetch tasks error:', err));
  }, [currentProjectId]);

  const handleCreateProject = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newProjectName.trim()) return;

    try {
      const created = await apiRequest<Project>('/projects', {
        method: 'POST',
        body: JSON.stringify({ name: newProjectName.trim() }),
      });
      setProjects([created, ...projects]);
      setSelectedProject(created);
      setNewProjectName('');
      setShowNewProjectModal(false);
    } catch (err) {
      console.error('create project error:', err);
    }
  };

  return (
    <div className="min-h-screen bg-[#090A0F] text-slate-100 flex flex-col">
      <Navbar userName={user?.name || 'Developer'} isConnected={isConnected} />

      <main className="flex-1 p-6 max-w-7xl mx-auto w-full space-y-6">
        {/* Project Header and Controls */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 border-b border-[#1E2230] pb-5">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-blue-600/10 border border-blue-500/20 flex items-center justify-center text-blue-400">
              <FolderKanban className="w-5 h-5" />
            </div>

            <div>
              <div className="flex items-center gap-2">
                <select
                  value={selectedProject?.id || ''}
                  onChange={(e) => {
                    const p = projects.find((proj) => proj.id === e.target.value);
                    if (p) setSelectedProject(p);
                  }}
                  className="bg-transparent font-bold text-lg md:text-xl text-white focus:outline-none cursor-pointer border-b border-transparent hover:border-slate-700 transition-colors"
                >
                  {projects.map((proj) => (
                    <option key={proj.id} value={proj.id} className="bg-[#12141C] text-white">
                      {proj.name}
                    </option>
                  ))}
                  {projects.length === 0 && (
                    <option value="" className="bg-[#12141C] text-white">
                      Sample Collaboration Board
                    </option>
                  )}
                </select>

                <button
                  onClick={() => setShowNewProjectModal(true)}
                  className="p-1 text-slate-400 hover:text-white hover:bg-slate-800 rounded-lg transition-colors"
                  title="Create new project"
                >
                  <Plus className="w-4 h-4" />
                </button>
              </div>
              <p className="text-xs text-slate-400 mt-0.5">
                Real-Time Kanban Board • Fractional Indexing O(1) Reordering
              </p>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <PresenceAvatars users={users} />
          </div>
        </div>

        {/* Real-Time Kanban Board */}
        <KanbanBoard
          projectId={currentProjectId || 'default-board'}
          tasks={tasks}
          onTasksChange={setTasks}
          broadcastUpdate={sendMessage}
        />
      </main>

      {/* New Project Modal */}
      {showNewProjectModal && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-[#12141C] border border-[#1E2230] rounded-2xl p-6 w-full max-w-md shadow-2xl space-y-4">
            <h3 className="text-lg font-bold text-white">Create New Project Board</h3>
            <form onSubmit={handleCreateProject} className="space-y-4">
              <input
                type="text"
                autoFocus
                required
                placeholder="e.g. Mobile App Redesign"
                value={newProjectName}
                onChange={(e) => setNewProjectName(e.target.value)}
                className="w-full bg-[#0D0F17] border border-[#1E2230] text-sm text-white px-4 py-2.5 rounded-xl focus:outline-none focus:border-blue-500"
              />
              <div className="flex items-center justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowNewProjectModal(false)}
                  className="text-xs text-slate-400 hover:text-white px-3 py-2"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="text-xs bg-blue-600 hover:bg-blue-500 text-white font-semibold px-4 py-2 rounded-lg"
                >
                  Create Board
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
