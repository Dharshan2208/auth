import { useState } from 'react';
import { Card } from '../common/Card';

interface TokenViewerProps {
  accessToken: string;
  refreshToken: string;
}

function decodeJWT(token: string): Record<string, unknown> | null {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) return null;
    const payload = parts[1];
    const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'));
    return JSON.parse(decoded) as Record<string, unknown>;
  } catch {
    return null;
  }
}

function TokenDisplay({
  label,
  token,
  color,
}: {
  label: string;
  token: string;
  color: 'yellow' | 'blue';
}) {
  const [copied, setCopied] = useState(false);
  const claims = decodeJWT(token);
  const truncated = token.length > 60 ? token.slice(0, 30) + '...' + token.slice(-20) : token;

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(token);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Fallback
    }
  };

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-bold text-black">{label}</h4>
        <button
          onClick={handleCopy}
          className="text-xs font-semibold text-gray-500 hover:text-black transition-colors"
        >
          {copied ? 'Copied!' : 'Copy'}
        </button>
      </div>
      <div
        className={`font-mono text-xs break-all p-3 border-2 border-black ${
          color === 'yellow' ? 'bg-yellow-50' : 'bg-blue-50'
        }`}
      >
        {truncated}
      </div>

      {claims && (
        <div className="bg-gray-50 border-2 border-black p-3">
          <p className="text-xs font-bold text-gray-500 mb-2 uppercase tracking-wider">
            Decoded Claims
          </p>
          <pre className="text-xs font-mono text-black whitespace-pre-wrap">
            {JSON.stringify(claims, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}

export function TokenViewer({ accessToken, refreshToken }: TokenViewerProps) {
  return (
    <Card className="p-6 sm:p-7">
      <h3 className="text-xl font-bold text-black mb-6">Token Viewer</h3>

      <div className="space-y-6">
        <TokenDisplay
          label="Access Token"
          token={accessToken}
          color="yellow"
        />
        <div className="border-t-2 border-gray-100" />
        <TokenDisplay
          label="Refresh Token"
          token={refreshToken}
          color="blue"
        />
      </div>

      <div className="mt-6 p-3 bg-gray-50 border-2 border-gray-300">
        <p className="text-xs text-gray-500">
          <strong>Access Token TTL:</strong> 15 minutes ·{' '}
          <strong>Refresh Token TTL:</strong> 7 days
        </p>
        <p className="text-xs text-gray-500 mt-1">
          Tokens are JWTs signed with HMAC-SHA256. Refresh tokens are stored as
          SHA-256 hashes in the database.
        </p>
      </div>
    </Card>
  );
}
