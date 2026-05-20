import { useState } from 'react';
import { Button } from '../common/Button';
import { Card } from '../common/Card';
import { ErrorDisplay } from '../common/ErrorDisplay';
import * as api from '../../api/client';
import { ApiError } from '../../api/client';

interface AdminPanelProps {
  accessToken: string;
  role: string;
}

export function AdminPanel({ accessToken, role }: AdminPanelProps) {
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const isAdmin = role === 'admin';

  const handleAdminCheck = async () => {
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const res = await api.adminCheck(accessToken);
      setResult(JSON.stringify(res, null, 2));
    } catch (e) {
      setError(e instanceof ApiError ? e.body : 'Admin check failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="p-6 sm:p-7">
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="text-xl font-bold text-black">Admin Access</h3>
          <p className="text-sm text-gray-500 mt-1">
            Test role-based access control. Only admin users can access this
            endpoint.
          </p>
        </div>
        <span
          className={`text-xs font-bold px-2 py-1 border-2 ${
            isAdmin
              ? 'bg-green-100 text-green-800 border-green-500'
              : 'bg-gray-100 text-gray-500 border-gray-300'
          }`}
        >
          {isAdmin ? 'Admin' : 'User'}
        </span>
      </div>

      <Button
        variant={isAdmin ? 'secondary' : 'outline'}
        onClick={handleAdminCheck}
        loading={loading}
      >
        {isAdmin ? 'Test Admin Access' : 'Attempt Admin Access'}
      </Button>

      {!isAdmin && (
        <p className="mt-3 text-xs text-gray-400">
          Your account has the &ldquo;user&rdquo; role. This endpoint will
          return a 403 error.
        </p>
      )}

      {error && (
        <div className="mt-4">
          <ErrorDisplay message={error} />
        </div>
      )}

      {result && (
        <div className="mt-4">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-2 h-2 bg-green-500 rounded-full" />
            <span className="text-xs font-bold text-green-700 uppercase">
              200 OK
            </span>
          </div>
          <pre className="bg-gray-50 border-2 border-black p-4 text-sm overflow-x-auto font-mono">
            {result}
          </pre>
        </div>
      )}
    </Card>
  );
}
