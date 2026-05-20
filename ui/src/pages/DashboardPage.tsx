import { ProfileCard } from '../components/dashboard/ProfileCard';
import { TokenViewer } from '../components/dashboard/TokenViewer';
import { RefreshFlowTester } from '../components/dashboard/RefreshFlowTester';
import { PasswordChangeForm } from '../components/dashboard/PasswordChangeForm';
import { AdminPanel } from '../components/dashboard/AdminPanel';
import type { User } from '../types';

interface DashboardPageProps {
  user: User;
  accessToken: string;
  refreshToken: string;
  onRefresh: () => Promise<{ access_token: string; refresh_token: string } | undefined>;
  onChangePassword: (
    oldPassword: string,
    confirmOldPassword: string,
    newPassword: string,
  ) => Promise<void>;
}

export function DashboardPage({
  user,
  accessToken,
  refreshToken,
  onRefresh,
  onChangePassword,
}: DashboardPageProps) {
  return (
    <main className="py-10">
      <div className="max-w-6xl mx-auto px-6">
        {/* Page Header */}
        <div className="mb-10">
          <h1 className="text-3xl sm:text-4xl font-black text-black">
            Dashboard
          </h1>
          <p className="mt-2 text-gray-500">
            Account overview, token management, and security settings.
          </p>
        </div>

        {/* Grid Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Profile */}
          <ProfileCard user={user} />

          {/* Admin Panel */}
          <AdminPanel accessToken={accessToken} role={user.role} />

          {/* Token Viewer */}
          <div className="lg:col-span-2">
            <TokenViewer
              accessToken={accessToken}
              refreshToken={refreshToken}
            />
          </div>

          {/* Refresh Flow Tester */}
          <RefreshFlowTester onRefresh={onRefresh} />

          {/* Password Change */}
          <PasswordChangeForm onChangePassword={onChangePassword} />
        </div>

        {/* API Reference Banner */}
        <div className="mt-10 bg-black text-white p-6 sm:p-8 border-2 border-black shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div>
              <h3 className="text-xl font-bold">
                Explore the Full API
              </h3>
              <p className="text-gray-400 text-sm mt-1">
                View the complete OpenAPI specification with interactive
                documentation.
              </p>
            </div>
            <a
              href="/swagger/"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-6 py-3 bg-yellow-400 text-black font-bold border-2 border-yellow-400 hover:bg-yellow-300 transition-colors whitespace-nowrap"
            >
              Open Swagger UI →
            </a>
          </div>
        </div>
      </div>
    </main>
  );
}
