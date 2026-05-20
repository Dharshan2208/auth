import { useState } from 'react';
import { Button } from '../common/Button';
import { Input } from '../common/Input';
import { ErrorDisplay } from '../common/ErrorDisplay';

interface LoginFormProps {
  onLogin: (identifier: string, password: string) => Promise<void>;
  onSuccess: () => void;
}

export function LoginForm({ onLogin, onSuccess }: LoginFormProps) {
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await onLogin(identifier, password);
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <Input
        label="Username or Email"
        placeholder="johndoe or john@example.com"
        value={identifier}
        onChange={(e) => setIdentifier(e.target.value)}
        required
      />
      <Input
        label="Password"
        type="password"
        placeholder="Enter your password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
      />

      {error && <ErrorDisplay message={error} />}

      <Button
        variant="primary"
        type="submit"
        loading={loading}
        disabled={!identifier || !password}
        className="w-full"
      >
        Sign In
      </Button>

      <div className="p-3 bg-yellow-50 border-2 border-yellow-400">
        <p className="text-xs font-semibold text-yellow-800">
          ⚡ 5 requests / minute rate limit
        </p>
      </div>
    </form>
  );
}
