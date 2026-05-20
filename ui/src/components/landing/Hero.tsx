import { Badge } from '../common/Badge';
import { Button } from '../common/Button';

interface HeroProps {
  onOpenAuth: () => void;
  onScrollToDemo: () => void;
}

export function Hero({ onOpenAuth, onScrollToDemo }: HeroProps) {
  return (
    <section className="relative pt-20 pb-16 sm:pt-28 sm:pb-20">
      <div className="max-w-4xl mx-auto px-6 text-center">
        {/* Badge */}
        <div className="mb-6 flex justify-center">
          <Badge variant="yellow">
            <span className="text-[10px]">▲</span>
            v1.0.0 · Authentication System
          </Badge>
        </div>

        {/* Main Heading */}
        <h1 className="text-5xl sm:text-6xl md:text-7xl lg:text-8xl font-black leading-[0.9] tracking-tight text-black">
          Modern Auth
          <br />
          <span className="bg-yellow-400 px-3 inline-block mt-2">
            Infrastructure
          </span>
        </h1>

        {/* Tagline */}
        <p className="mt-8 text-lg sm:text-xl text-gray-500 max-w-2xl mx-auto leading-relaxed font-medium">
          A production-grade authentication system built with{' '}
          <strong className="text-black">Go</strong>,{' '}
          <strong className="text-black">PostgreSQL</strong>, and{' '}
          <strong className="text-black">JWT</strong>.
          <br />
          Secure session management, token rotation, and role-based access
          control.
        </p>

        {/* CTA Buttons */}
        <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button variant="primary" size="lg" onClick={onOpenAuth}>
            Get Started →
          </Button>
          <Button variant="outline" size="lg" onClick={onScrollToDemo}>
            View Demo ↓
          </Button>
        </div>

        {/* Tech Stack Badges */}
        <div className="mt-12 flex flex-wrap items-center justify-center gap-3">
          {['Go', 'PostgreSQL', 'JWT', 'bcrypt', 'HMAC-SHA256', 'REST'].map(
            (tech) => (
              <span
                key={tech}
                className="px-3 py-1.5 text-xs font-bold border-2 border-black bg-white hover:-translate-y-0.5 transition-transform"
              >
                {tech}
              </span>
            ),
          )}
        </div>
      </div>
    </section>
  );
}
