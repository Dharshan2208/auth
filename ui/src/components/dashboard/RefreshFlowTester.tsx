import { useState } from 'react';
import { Button } from '../common/Button';
import { Card } from '../common/Card';
import { ErrorDisplay } from '../common/ErrorDisplay';

interface RefreshFlowTesterProps {
  onRefresh: () => Promise<{ access_token: string; refresh_token: string } | undefined>;
}

export function RefreshFlowTester({ onRefresh }: RefreshFlowTesterProps) {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [count, setCount] = useState(0);

  const handleRefresh = async () => {
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const tokens = await onRefresh();
      if (tokens) {
        setResult(JSON.stringify(tokens, null, 2));
        setCount((c) => c + 1);
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Refresh failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="p-6 sm:p-7">
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="text-xl font-bold text-black">Refresh Flow</h3>
          <p className="text-sm text-gray-500 mt-1">
            Test token rotation. Each refresh invalidates the old refresh
            token and issues a new pair.
          </p>
        </div>
        {count > 0 && (
          <span className="text-xs font-bold bg-blue-100 text-blue-800 px-2 py-1 border border-blue-300">
            Rotated: {count}
          </span>
        )}
      </div>

      <Button
        variant="primary"
        size="md"
        onClick={handleRefresh}
        loading={loading}
      >
        Rotate Tokens
      </Button>

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
              200 OK · New Token Pair
            </span>
          </div>
          <pre className="bg-gray-50 border-2 border-black p-4 text-sm overflow-x-auto font-mono">
            {result}
          </pre>
        </div>
      )}

      <div className="mt-4 p-3 bg-yellow-50 border-2 border-yellow-400">
        <p className="text-xs font-semibold text-yellow-800">
          ⚡ Rotating twice with the same refresh token triggers automatic
          revocation (reuse detection)
        </p>
      </div>
    </Card>
  );
}
