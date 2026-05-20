import { Hero } from '../components/landing/Hero';
import { FeatureCards } from '../components/landing/FeatureCards';
import { InteractiveDemo } from '../components/landing/InteractiveDemo';

interface HomePageProps {
  onOpenAuth: () => void;
}

export function HomePage({ onOpenAuth }: HomePageProps) {
  const scrollToDemo = () => {
    const el = document.getElementById('demo');
    if (el) el.scrollIntoView({ behavior: 'smooth' });
  };

  return (
    <main>
      <Hero onOpenAuth={onOpenAuth} onScrollToDemo={scrollToDemo} />
      <FeatureCards />
      <InteractiveDemo />
    </main>
  );
}
