import { FeatureCard } from './FeatureCard';

const features = [
  {
    icon: '🔐',
    title: 'JWT Authentication',
    description:
      'Access and refresh token pairs signed with HMAC-SHA256. Stateless access tokens with 15-minute TTL, refresh tokens with 7-day rotation.',
    color: 'yellow' as const,
  },
  {
    icon: '🔄',
    title: 'Token Rotation',
    description:
      'Automatic refresh token rotation with reuse detection. Compromised tokens trigger automatic revocation of all device sessions.',
    color: 'blue' as const,
  },
  {
    icon: '💾',
    title: 'PostgreSQL Storage',
    description:
      'User accounts and session data stored in PostgreSQL with connection pooling, prepared statements, and efficient indexing.',
    color: 'green' as const,
  },
  {
    icon: '🛡️',
    title: 'Password Policies',
    description:
      'Bcrypt hashing with 12-72 character password rules requiring uppercase, lowercase, digits, and symbols. No plaintext storage.',
    color: 'pink' as const,
  },
  {
    icon: '🧩',
    title: 'Session Management',
    description:
      'Server-side session tracking with IP and user agent logging. Atomic token rotation, revocation, and multi-device support.',
    color: 'yellow' as const,
  },
  {
    icon: '👑',
    title: 'Role-Based Access',
    description:
      'JWT-embedded role claims for fine-grained authorization. Admin endpoints with middleware-level role verification.',
    color: 'blue' as const,
  },
  {
    icon: '📋',
    title: 'Security Headers',
    description:
      'Comprehensive HTTP security headers including CSP, HSTS, X-Frame-Options, and Permissions-Policy with Swagger-aware exceptions.',
    color: 'green' as const,
  },
  {
    icon: '⏱️',
    title: 'Rate Limiting',
    description:
      'Per-IP sliding window rate limiter with configurable limits per endpoint. Cleanup goroutine for automatic entry expiration.',
    color: 'pink' as const,
  },
];

export function FeatureCards() {
  return (
    <section className="py-20">
      <div className="max-w-6xl mx-auto px-6">
        {/* Section Header */}
        <div className="mb-14">
          <h2 className="text-3xl sm:text-4xl font-black text-black">
            Engineered for
            <br />
            <span className="bg-yellow-400 px-2 inline-block mt-1">
              Security & Scale
            </span>
          </h2>
          <p className="mt-4 text-gray-500 max-w-xl text-base">
            Every component is designed with security-first principles and
            production readiness in mind.
          </p>
        </div>

        {/* Cards Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {features.map((f) => (
            <FeatureCard key={f.title} {...f} />
          ))}
        </div>
      </div>
    </section>
  );
}
