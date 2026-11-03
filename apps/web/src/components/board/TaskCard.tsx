'use client';

import { Task } from '@/lib/api';
import Link from 'next/link';
import { FileText, MoreVertical, Trash2 } from 'lucide-react';

interface TaskCardProps {
  task: Task;
  onDragStart: (e: React.DragEvent, task: Task) => void;
  onDelete?: (id: string) => void;
}

const priorityColors: Record<string, string> = {
  urgent: 'text-rose-400 bg-rose-500/10 border-rose-500/20',
  high: 'text-amber-400 bg-amber-500/10 border-amber-500/20',
  medium: 'text-blue-400 bg-blue-500/10 border-blue-500/20',
  low: 'text-slate-400 bg-slate-500/10 border-slate-500/20',
};

export function TaskCard({ task, onDragStart, onDelete }: TaskCardProps) {
  return (
    <div
      draggable
      onDragStart={(e) => onDragStart(e, task)}
      className="p-3.5 bg-[#12141C] border border-[#1E2230] hover:border-slate-700 rounded-xl shadow-sm hover:shadow-md transition-all cursor-grab active:cursor-grabbing group select-none"
    >
      <div className="flex items-start justify-between gap-2 mb-2">
        <span
          className={`text-[10px] font-semibold uppercase px-2 py-0.5 rounded-full border ${
            priorityColors[task.priority] || priorityColors.medium
          }`}
        >
          {task.priority}
        </span>

        {onDelete && (
          <button
            onClick={() => onDelete(task.id)}
            className="opacity-0 group-hover:opacity-100 p-1 text-slate-500 hover:text-rose-400 transition-opacity"
            title="Delete task"
          >
            <Trash2 className="w-3.5 h-3.5" />
          </button>
        )}
      </div>

      <h4 className="text-sm font-medium text-slate-200 line-clamp-2 mb-3">
        {task.title}
      </h4>

      <div className="flex items-center justify-between text-xs text-slate-400 pt-2 border-t border-[#1E2230]/60">
        <Link
          href={`/dashboard/editor/${task.id}`}
          className="flex items-center gap-1.5 text-blue-400 hover:text-blue-300 transition-colors"
          title="Open real-time CRDT document"
        >
          <FileText className="w-3.5 h-3.5" />
          <span>Live Doc</span>
        </Link>

        <span className="font-mono text-[10px] text-slate-500">
          pos: {task.position.toFixed(1)}
        </span>
      </div>
    </div>
  );
}
