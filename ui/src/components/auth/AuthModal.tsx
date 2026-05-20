import { useState } from 'react';
import { LoginForm } from './LoginForm';
import { SignupForm } from './SignupForm';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
  onLogin: (identifier: string, password: string) => Promise<void>;
  onSignup: (username: string, email: string, password: string) => Promise<void>;
  onAuthSuccess: () => void;
}

type Tab = 'login' | 'signup';

export function AuthModal({
  isOpen,
  onClose,
  onLogin,
  onSignup,
  onAuthSuccess,
}: AuthModalProps) {
  const [tab, setTab] = useState<Tab>('login');

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/30 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative w-full max-w-md bg-white border-2 border-black shadow-[10px_10px_0px_0px_rgba(0,0,0,1)]">
        {/* Header */}
        <div className="flex border-b-2 border-black">
          <button
            onClick={() => setTab('login')}
            className={`
              flex-1 py-3.5 text-sm font-bold transition-colors
              ${
                tab === 'login'
                  ? 'bg-black text-white'
                  : 'bg-white text-black hover:bg-gray-100'
              }
            `}
          >
            Sign In
          </button>
          <button
            onClick={() => setTab('signup')}
            className={`
              flex-1 py-3.5 text-sm font-bold transition-colors
              ${
                tab === 'signup'
                  ? 'bg-black text-white'
                  : 'bg-white text-black hover:bg-gray-100'
              }
            `}
          >
            Create Account
          </button>
          <button
            onClick={onClose}
            className="px-4 py-3.5 text-sm font-bold bg-white text-black border-l-2 border-black hover:bg-gray-100 transition-colors"
          >
            ×
          </button>
        </div>

        {/* Body */}
        <div className="p-6">
          {tab === 'login' ? (
            <LoginForm onLogin={onLogin} onSuccess={onAuthSuccess} />
          ) : (
            <SignupForm onSignup={onSignup} />
          )}
        </div>
      </div>
    </div>
  );
}
