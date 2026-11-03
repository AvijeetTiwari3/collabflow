import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'CollabFlow — Real-Time Collaborative Workspace',
  description: 'High-concurrency project management and collaborative document suite powered by CRDTs, WebSockets, and Go.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className="bg-[#090A0F] text-slate-100 min-h-screen antialiased font-sans selection:bg-blue-600 selection:text-white">
        {children}
      </body>
    </html>
  );
}
