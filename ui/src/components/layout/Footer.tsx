export function Footer() {
  return (
    <footer className="border-t-2 border-black bg-white mt-24">
      <div className="max-w-6xl mx-auto px-6 py-8 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <div className="w-5 h-5 bg-yellow-400 border border-black flex items-center justify-center font-bold text-[10px]">
            A
          </div>
          <span className="text-sm font-semibold text-gray-500">
            Auth System · Go + PostgreSQL + JWT
          </span>
        </div>

        <div className="flex items-center gap-6 text-sm text-gray-400">
          <span>Backend Engineering Demo</span>
          <span className="hidden sm:inline">·</span>
          <span className="hidden sm:inline">
            HMAC-SHA256 · bcrypt · Token Rotation
          </span>
        </div>
      </div>
    </footer>
  );
}
