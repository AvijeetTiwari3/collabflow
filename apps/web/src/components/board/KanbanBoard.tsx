'use client';

import { useState } from 'react';
import { Task, apiRequest } from '@/lib/api';
import { calculatePosition } from '@/lib/fractionalIndex';
import { TaskCard } from './TaskCard';
import { Plus } from 'lucide-react';

interface KanbanBoardProps {
  projectId: string;
  tasks: Task[];
  onTasksChange: (tasks: Task[]) => void;
  broadcastUpdate: (msg: any) => void;
}

const COLUMNS: { id: Task['status']; title: string; countColor: string }[] = [
  { id: 'todo', title: 'To Do', countColor: 'bg-slate-800 text-slate-300' },
  { id: 'in_progress', title: 'In Progress', countColor: 'bg-blue-900/50 text-blue-300' },
  { id: 'in_review', title: 'In Review', countColor: 'bg-amber-900/50 text-amber-300' },
  { id: 'done', title: 'Done', countColor: 'bg-emerald-900/50 text-emerald-300' },
];

export function KanbanBoard({ projectId, tasks, onTasksChange, broadcastUpdate }: KanbanBoardProps) {
  const [draggedTask, setDraggedTask] = useState<Task | null>(null);
  const [newTaskColumn, setNewTaskColumn] = useState<string | null>(null);
  const [newTitle, setNewTitle] = useState('');

  const handleDragStart = (e: React.DragEvent, task: Task) => {
    setDraggedTask(task);
    e.dataTransfer.setData('text/plain', task.id);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  const handleDropOnColumn = async (targetStatus: Task['status']) => {
    if (!draggedTask) return;

    // Get all tasks in target column
    const columnTasks = tasks
      .filter((t) => t.status === targetStatus && t.id !== draggedTask.id)
      .sort((a, b) => a.position - b.position);

    const prevPos = columnTasks.length > 0 ? columnTasks[columnTasks.length - 1].position : null;
    const newPos = calculatePosition(prevPos, null);

    applyReorder(draggedTask.id, targetStatus, newPos, prevPos, null);
  };

  const applyReorder = async (
    taskId: string,
    newStatus: Task['status'],
    newPos: number,
    prevPos: number | null,
    nextPos: number | null
  ) => {
    // 1. Optimistic local state update (Instant UI feedback)
    const updatedTasks = tasks.map((t) =>
      t.id === taskId ? { ...t, status: newStatus, position: newPos } : t
    );
    onTasksChange(updatedTasks);
    setDraggedTask(null);

    // 2. Broadcast via WebSocket to all other connected peers
    broadcastUpdate({
      type: 'task_updated',
      room_id: projectId,
      payload: {
        task_id: taskId,
        status: newStatus,
        position: newPos,
      },
    });

    // 3. Persist to Backend with Fractional Indexing
    try {
      await apiRequest<Task>(`/tasks/${taskId}/reorder`, {
        method: 'POST',
        body: JSON.stringify({
          status: newStatus,
          prev_pos: prevPos,
          next_pos: nextPos,
        }),
      });
    } catch (err) {
      console.error('[board] failed to persist reorder:', err);
    }
  };

  const handleCreateTask = async (status: Task['status']) => {
    if (!newTitle.trim()) return;

    try {
      const created = await apiRequest<Task>(`/projects/${projectId}/tasks`, {
        method: 'POST',
        body: JSON.stringify({
          title: newTitle.trim(),
          status,
          priority: 'medium',
        }),
      });

      onTasksChange([...tasks, created]);
      broadcastUpdate({
        type: 'task_created',
        room_id: projectId,
        payload: created,
      });

      setNewTitle('');
      setNewTaskColumn(null);
    } catch (err) {
      console.error('[board] create task error:', err);
    }
  };

  const handleDeleteTask = async (id: string) => {
    try {
      await apiRequest(`/tasks/${id}`, { method: 'DELETE' });
      onTasksChange(tasks.filter((t) => t.id !== id));
      broadcastUpdate({
        type: 'task_deleted',
        room_id: projectId,
        payload: { task_id: id },
      });
    } catch (err) {
      console.error('[board] delete task error:', err);
    }
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5 items-start">
      {COLUMNS.map((col) => {
        const colTasks = tasks
          .filter((t) => t.status === col.id)
          .sort((a, b) => a.position - b.position);

        return (
          <div
            key={col.id}
            onDragOver={handleDragOver}
            onDrop={() => handleDropOnColumn(col.id)}
            className="bg-[#0D0F17]/90 border border-[#1E2230] rounded-2xl p-4 min-h-[500px] flex flex-col gap-3 transition-colors"
          >
            {/* Column Header */}
            <div className="flex items-center justify-between pb-2 border-b border-[#1E2230]">
              <div className="flex items-center gap-2">
                <h3 className="font-semibold text-sm text-slate-200">{col.title}</h3>
                <span className={`text-xs px-2 py-0.5 rounded-full font-mono ${col.countColor}`}>
                  {colTasks.length}
                </span>
              </div>

              <button
                onClick={() => setNewTaskColumn(col.id)}
                className="p-1 text-slate-400 hover:text-white hover:bg-slate-800 rounded-lg transition-colors"
                title="Add task"
              >
                <Plus className="w-4 h-4" />
              </button>
            </div>

            {/* In-line Task Creation Form */}
            {newTaskColumn === col.id && (
              <div className="p-3 bg-[#12141C] border border-blue-500/40 rounded-xl space-y-2">
                <input
                  type="text"
                  autoFocus
                  placeholder="Task title..."
                  value={newTitle}
                  onChange={(e) => setNewTitle(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') handleCreateTask(col.id);
                    if (e.key === 'Escape') setNewTaskColumn(null);
                  }}
                  className="w-full bg-transparent text-sm text-white focus:outline-none placeholder:text-slate-500"
                />
                <div className="flex items-center justify-end gap-2 pt-1">
                  <button
                    onClick={() => setNewTaskColumn(null)}
                    className="text-xs text-slate-400 hover:text-white px-2 py-1"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={() => handleCreateTask(col.id)}
                    className="text-xs bg-blue-600 hover:bg-blue-500 text-white font-medium px-3 py-1 rounded-md"
                  >
                    Add Card
                  </button>
                </div>
              </div>
            )}

            {/* Task Cards */}
            <div className="space-y-3 flex-1">
              {colTasks.map((task) => (
                <TaskCard
                  key={task.id}
                  task={task}
                  onDragStart={handleDragStart}
                  onDelete={handleDeleteTask}
                />
              ))}

              {colTasks.length === 0 && !newTaskColumn && (
                <div className="h-32 border-2 border-dashed border-[#1E2230] rounded-xl flex items-center justify-center text-xs text-slate-500">
                  Drop tasks here
                </div>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
