import type { ReactNode } from 'react';

interface FeatureCardProps {
  icon: ReactNode;
  title: string;
  description: string;
  color: 'yellow' | 'blue' | 'green' | 'pink';
}

const colorMap: Record<string, string> = {
  yellow: 'bg-yellow-400',
  blue: 'bg-blue-500',
  green: 'bg-green-300',
  pink: 'bg-pink-300',
};

export function FeatureCard({ icon, title, description, color }: FeatureCardProps) {
  return (
    <div className="group bg-white border-2 border-black shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] transition-all duration-200 hover:-translate-y-1 hover:shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
      <div className="p-6 sm:p-7">
        {/* Icon Box */}
        <div
          className={`
            w-12 h-12 ${colorMap[color]}
            border-2 border-black flex items-center justify-center
            text-lg font-bold mb-5
          `}
        >
          {icon}
        </div>

        {/* Title */}
        <h3 className="text-xl font-bold text-black mb-2">{title}</h3>

        {/* Description */}
        <p className="text-sm text-gray-500 leading-relaxed">{description}</p>
      </div>
    </div>
  );
}
