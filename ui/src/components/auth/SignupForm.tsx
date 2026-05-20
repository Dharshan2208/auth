import { useState } from 'react';
import { Button } from '../common/Button';
import { Input } from '../common/Input';
import { ErrorDisplay } from '../common/ErrorDisplay';
import { SuccessDisplay } from '../common/SuccessDisplay';

interface SignupFormProps {
  onSignup: (username: string, email: string, password: string) => Promise<void>;
  onSuccess?: () => void;
}

export function SignupForm({ onSignup, onSuccess }: SignupFormProps) {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    setLoading(true);
    try {
      await onSignup(username, email, password);
      setSuccess('Account created successfully! You can now sign in.');
      setUsername('');
      setEmail('');
      setPassword('');
      setConfirmPassword('');
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Signup failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <Input
        label="Username"
        placeholder="johndoe (3-31 chars, lowercase, alphanumeric)"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        required
      />
      <Input
        label="Email"
        type="email"
        placeholder="john@example.com"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        required
      />
      <Input
        label="Password"
        type="password"
        placeholder="12-72 chars: upper + lower + digit + symbol"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
      />
      <Input
        label="Confirm Password"
        type="password"
        placeholder="Repeat your password"
        value={confirmPassword}
        onChange={(e) => setConfirmPassword(e.target.value)}
        required
      />

      {error && <ErrorDisplay message={error} />}
      {success && <SuccessDisplay message={success} />}

      <Button
        variant="secondary"
        type="submit"
        loading={loading}
        disabled={!username || !email || !password || !confirmPassword}
        className="w-full"
      >
        Create Account
      </Button>

      <div className="p-3 bg-yellow-50 border-2 border-yellow-400">
        <p className="text-xs font-semibold text-yellow-800">
          ⚡ 3 requests / minute rate limit
        </p>
      </div>
    </form>
  );
}
