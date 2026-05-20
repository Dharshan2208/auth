import { useState } from 'react';
import { Button } from '../common/Button';
import { Input } from '../common/Input';
import { ErrorDisplay } from '../common/ErrorDisplay';
import { LoadingSpinner } from '../common/LoadingSpinner';
import * as api from '../../api/client';
import { ApiError } from '../../api/client';

type DemoMode = 'health' | 'signup' | 'login' | 'swagger';

export function InteractiveDemo() {
  const [mode, setMode] = useState<DemoMode>('health');

  return (
    <section id="demo" className="py-20 scroll-mt-24">
      <div className="max-w-4xl mx-auto px-6">
        {/* Section Header */}
        <div className="mb-10">
          <h2 className="text-3xl sm:text-4xl font-black text-black">
            Interactive
            <br />
            <span className="bg-blue-500 text-white px-2 inline-block mt-1">
              API Demo
            </span>
          </h2>
          <p className="mt-4 text-gray-500 max-w-xl text-base">
            Test the backend endpoints directly from your browser. No
            authentication required for public endpoints.
          </p>
        </div>

        {/* Tabs */}
        <div className="flex border-2 border-black divide-x-2 divide-black mb-8 overflow-hidden">
          {[
            { id: 'health' as const, label: 'Health' },
            { id: 'signup' as const, label: 'Signup' },
            { id: 'login' as const, label: 'Login' },
            { id: 'swagger' as const, label: 'Swagger' },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setMode(tab.id)}
              className={`
                flex-1 py-3 px-4 text-sm font-bold transition-colors
                ${
                  mode === tab.id
                    ? 'bg-black text-white'
                    : 'bg-white text-black hover:bg-gray-100'
                }
              `}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Demo Content */}
        <div className="bg-white border-2 border-black shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] p-6 sm:p-8">
          {mode === 'health' && <HealthChecker />}
          {mode === 'signup' && <SignupDemo />}
          {mode === 'login' && <LoginDemo />}
          {mode === 'swagger' && <SwaggerLink />}
        </div>
      </div>
    </section>
  );
}

function HealthChecker() {
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const check = async () => {
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const res = await api.healthCheck();
      setResult(JSON.stringify(res, null, 2));
    } catch (e) {
      setError(e instanceof ApiError ? e.body : 'Connection failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-xl font-bold text-black">Health Check</h3>
          <p className="text-sm text-gray-500 mt-1">
            Verify the backend server is running and the database is
            reachable.
          </p>
        </div>
        <Button
          variant="primary"
          size="sm"
          onClick={check}
          loading={loading}
        >
          Ping Server
        </Button>
      </div>

      {error && <ErrorDisplay message={error} />}
      {loading && <LoadingSpinner size="sm" />}

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
    </div>
  );
}

function SignupDemo() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSignup = async () => {
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const res = await api.signup({ username, email, password });
      setResult(JSON.stringify(res, null, 2));
    } catch (e) {
      setError(e instanceof ApiError ? e.body : 'Signup failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="mb-6">
        <h3 className="text-xl font-bold text-black">Signup</h3>
        <p className="text-sm text-gray-500 mt-1">
          Register a new user account. Password must be 12-72 characters with
          uppercase, lowercase, digit, and symbol.
        </p>
      </div>

      <div className="space-y-4">
        <Input
          label="Username"
          placeholder="johndoe"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <Input
          label="Email"
          type="email"
          placeholder="john@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <Input
          label="Password"
          type="password"
          placeholder="Str0ng!Pass1"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <Button
          variant="secondary"
          onClick={handleSignup}
          loading={loading}
          disabled={!username || !email || !password}
        >
          Create Account
        </Button>
      </div>

      {error && (
        <div className="mt-4">
          <ErrorDisplay message={error} />
        </div>
      )}
      {loading && <LoadingSpinner size="sm" text="Registering user..." />}
      {result && (
        <div className="mt-4">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-2 h-2 bg-green-500 rounded-full" />
            <span className="text-xs font-bold text-green-700 uppercase">
              201 Created
            </span>
          </div>
          <pre className="bg-gray-50 border-2 border-black p-4 text-sm overflow-x-auto font-mono">
            {result}
          </pre>
        </div>
      )}

      <div className="mt-4 p-3 bg-yellow-50 border-2 border-yellow-400">
        <p className="text-xs font-semibold text-yellow-800">
          ⚡ Rate Limited: 3 requests per minute
        </p>
      </div>
    </div>
  );
}

function LoginDemo() {
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleLogin = async () => {
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const res = await api.login({ identifier, password });
      setResult(JSON.stringify(res, null, 2));
    } catch (e) {
      setError(e instanceof ApiError ? e.body : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="mb-6">
        <h3 className="text-xl font-bold text-black">Login</h3>
        <p className="text-sm text-gray-500 mt-1">
          Authenticate with your username or email. Returns an access token
          and refresh token pair.
        </p>
      </div>

      <div className="space-y-4">
        <Input
          label="Username or Email"
          placeholder="johndoe or john@example.com"
          value={identifier}
          onChange={(e) => setIdentifier(e.target.value)}
        />
        <Input
          label="Password"
          type="password"
          placeholder="Enter your password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <Button
          variant="primary"
          onClick={handleLogin}
          loading={loading}
          disabled={!identifier || !password}
        >
          Sign In
        </Button>
      </div>

      {error && (
        <div className="mt-4">
          <ErrorDisplay message={error} />
        </div>
      )}
      {loading && <LoadingSpinner size="sm" text="Authenticating..." />}
      {result && (
        <div className="mt-4">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-2 h-2 bg-green-500 rounded-full" />
            <span className="text-xs font-bold text-green-700 uppercase">
              200 OK · Token Pair
            </span>
          </div>
          <pre className="bg-gray-50 border-2 border-black p-4 text-sm overflow-x-auto font-mono">
            {result}
          </pre>
        </div>
      )}

      <div className="mt-4 p-3 bg-yellow-50 border-2 border-yellow-400">
        <p className="text-xs font-semibold text-yellow-800">
          ⚡ Rate Limited: 5 requests per minute
        </p>
      </div>
    </div>
  );
}

function SwaggerLink() {
  return (
    <div className="text-center py-8">
      <div className="w-16 h-16 bg-green-300 border-2 border-black flex items-center justify-center text-2xl font-bold mx-auto mb-6">
        📖
      </div>
      <h3 className="text-2xl font-bold text-black mb-3">
        OpenAPI Documentation
      </h3>
      <p className="text-gray-500 mb-6 max-w-md mx-auto">
        Browse the complete API specification with interactive request/response
        examples, schema definitions, and try-it-yourself functionality.
      </p>
      <a
        href="/swagger/"
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex items-center gap-2 px-8 py-3.5 bg-black text-white font-bold border-2 border-black hover:bg-gray-800 transition-colors"
      >
        Open Swagger UI →
      </a>
    </div>
  );
}
