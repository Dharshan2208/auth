import { useState } from 'react';
import { useAuth } from './hooks/useAuth';
import { Navbar } from './components/layout/Navbar';
import { Footer } from './components/layout/Footer';
import { AuthModal } from './components/auth/AuthModal';
import { HomePage } from './pages/HomePage';
import { DashboardPage } from './pages/DashboardPage';
import { LoadingSpinner } from './components/common/LoadingSpinner';
import type { Page } from './types';

export default function App() {
  const { auth, login, signup, logout, refreshTokens } = useAuth();
  const [page, setPage] = useState<Page>('home');
  const [showAuthModal, setShowAuthModal] = useState(false);

  const handleLogin = async (identifier: string, password: string) => {
    await login(identifier, password);
  };

  const handleSignup = async (
    username: string,
    email: string,
    password: string,
  ) => {
    await signup(username, email, password);
  };

  const handleAuthSuccess = () => {
    setShowAuthModal(false);
    setPage('dashboard');
  };

  const handleLogout = async () => {
    await logout();
    setPage('home');
  };

  const handleChangePassword = async (
    oldPassword: string,
    confirmOldPassword: string,
    newPassword: string,
  ) => {
    if (auth.status !== 'authenticated') return;

    const res = await fetch('/api/v1/password/change', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${auth.accessToken}`,
      },
      body: JSON.stringify({
        old_password: oldPassword,
        confirm_old_password: confirmOldPassword,
        new_password: newPassword,
      }),
    });

    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(err.error || 'Password change failed');
    }

    // Password change revokes all sessions → log out
    await logout();
    setPage('home');
    throw new Error(
      'Password updated. All sessions revoked. Please log in again.',
    );
  };

  // Loading state
  if (auth.status === 'loading') {
    return (
      <div className="min-h-screen flex items-center justify-center bg-white">
        <LoadingSpinner size="lg" text="Loading..." />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white bg-dots font-sans flex flex-col">
      <Navbar
        auth={auth}
        onOpenAuth={() => setShowAuthModal(true)}
        onLogout={handleLogout}
        onGoHome={() => setPage('home')}
        onGoDashboard={() => setPage('dashboard')}
        page={page}
      />

      <div className="flex-1">
        {page === 'home' && (
          <HomePage onOpenAuth={() => setShowAuthModal(true)} />
        )}

        {page === 'dashboard' && auth.status === 'authenticated' && (
          <DashboardPage
            user={auth.user}
            accessToken={auth.accessToken}
            refreshToken={auth.refreshToken}
            onRefresh={refreshTokens}
            onChangePassword={handleChangePassword}
          />
        )}
      </div>

      <Footer />

      {/* Auth Modal */}
      <AuthModal
        isOpen={showAuthModal}
        onClose={() => setShowAuthModal(false)}
        onLogin={handleLogin}
        onSignup={handleSignup}
        onAuthSuccess={handleAuthSuccess}
      />
    </div>
  );
}
