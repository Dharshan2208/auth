import { useState } from 'react';
import { Button } from '../common/Button';
import { Input } from '../common/Input';
import { Card } from '../common/Card';
import { ErrorDisplay } from '../common/ErrorDisplay';
import { SuccessDisplay } from '../common/SuccessDisplay';

interface PasswordChangeFormProps {
  onChangePassword: (
    oldPassword: string,
    confirmOldPassword: string,
    newPassword: string,
  ) => Promise<void>;
}

export function PasswordChangeForm({
  onChangePassword,
}: PasswordChangeFormProps) {
  const [oldPassword, setOldPassword] = useState('');
  const [confirmOldPassword, setConfirmOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (oldPassword !== confirmOldPassword) {
      setError('Old password confirmation does not match');
      return;
    }

    if (oldPassword === newPassword) {
      setError('New password must be different from old password');
      return;
    }

    setLoading(true);
    try {
      await onChangePassword(oldPassword, confirmOldPassword, newPassword);
      setSuccess('Password updated! All sessions have been revoked.');
      setOldPassword('');
      setConfirmOldPassword('');
      setNewPassword('');
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Password change failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="p-6 sm:p-7">
      <div className="mb-4">
        <h3 className="text-xl font-bold text-black">Change Password</h3>
        <p className="text-sm text-gray-500 mt-1">
          Changing your password will revoke all active sessions.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Current Password"
          type="password"
          placeholder="Enter current password"
          value={oldPassword}
          onChange={(e) => setOldPassword(e.target.value)}
          required
        />
        <Input
          label="Confirm Current Password"
          type="password"
          placeholder="Re-enter current password"
          value={confirmOldPassword}
          onChange={(e) => setConfirmOldPassword(e.target.value)}
          required
        />
        <Input
          label="New Password"
          type="password"
          placeholder="12-72 chars: upper + lower + digit + symbol"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          required
        />

        {error && <ErrorDisplay message={error} />}
        {success && <SuccessDisplay message={success} />}

        <Button
          variant="danger"
          type="submit"
          loading={loading}
          disabled={!oldPassword || !confirmOldPassword || !newPassword}
        >
          Update Password
        </Button>
      </form>

      <div className="mt-4 p-3 bg-red-50 border-2 border-red-300">
        <p className="text-xs font-semibold text-red-800">
          ⚠️ This will revoke all active sessions. You will need to log in
          again on all devices.
        </p>
      </div>
    </Card>
  );
}
