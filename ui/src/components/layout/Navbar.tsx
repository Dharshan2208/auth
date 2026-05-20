import type { AuthState } from '../../types';
import { Button } from '../common/Button';
import { Badge } from '../common/Badge';

interface NavbarProps {
  auth: AuthState;
  onOpenAuth: () => void;
  onLogout: () => void;
  onGoHome: () => void;
  onGoDashboard: () => void;
  page: 'home' | 'dashboard';
}

export function Navbar({
  auth,
  onOpenAuth,
  onLogout,
  onGoHome,
  onGoDashboard,
  page,
}: NavbarProps) {
  return (
    <nav className="sticky top-0 z-50 bg-white border-b-2 border-black">
      <div className="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
        {/* Logo */}
        <button
          onClick={onGoHome}
          className="flex items-center gap-3 hover:opacity-70 transition-opacity"
        >
          <div className="w-8 h-8 bg-yellow-400 border-2 border-black flex items-center justify-center font-bold text-sm">
            A
          </div>
          <span className="font-bold text-lg hidden sm:inline">
            Auth System
          </span>
        </button>

        {/* Nav Links */}
        <div className="flex items-center gap-6">
          <button
            onClick={onGoHome}
            className={`text-sm font-semibold transition-colors hover:text-yellow-600 ${
              page === 'home' ? 'text-black' : 'text-gray-500'
            }`}
          >
            Home
          </button>

          <a
            href="/swagger/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm font-semibold text-gray-500 hover:text-black transition-colors"
          >
            API Docs
          </a>

          {auth.status === 'authenticated' && (
            <button
              onClick={onGoDashboard}
              className={`text-sm font-semibold transition-colors hover:text-yellow-600 ${
                page === 'dashboard' ? 'text-black' : 'text-gray-500'
              }`}
            >
              Dashboard
            </button>
          )}

          {/* Auth Action */}
          {auth.status === 'loading' ? null : auth.status ===
            'authenticated' ? (
            <div className="flex items-center gap-3">
              <Badge variant="green">{auth.user.username}</Badge>
              <Button variant="outline" size="sm" onClick={onLogout}>
                Logout
              </Button>
            </div>
          ) : (
            <Button variant="primary" size="sm" onClick={onOpenAuth}>
              Sign In
            </Button>
          )}
        </div>
      </div>
    </nav>
  );
}
