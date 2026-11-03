'use client';

import { useParams } from 'next/navigation';
import { Navbar } from '@/components/Navbar';
import { CollaborativeEditor } from '@/components/editor/CollaborativeEditor';

export default function DocumentEditorPage() {
  const params = useParams();
  const taskId = (params?.id as string) || '';

  return (
    <div className="min-h-screen bg-[#090A0F] text-slate-100 flex flex-col">
      <Navbar />
      <main className="flex-1 p-6 max-w-7xl mx-auto w-full">
        <CollaborativeEditor taskId={taskId} />
      </main>
    </div>
  );
}
