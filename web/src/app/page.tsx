"use client";
import Image from "next/image";
import { BellAlertIcon, ShieldCheckIcon, AdjustmentsHorizontalIcon, ArrowRightIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-[#0a1a0f] to-[#1a3a1a] flex flex-col">
      {/* Header */}
      <header className="w-full flex items-center justify-between px-8 py-6 bg-white/10 backdrop-blur-md shadow-lg border-b border-white/10 z-10">
        <div className="flex items-center gap-2">
          <span className="inline-block animate-float">
            <Image src="/GitNotify.png" alt="GitNotify Logo" width={44} height={44} className="rounded-lg shadow-lg" style={{boxShadow: '0 4px 24px 0 rgba(0,255,120,0.18)'}} />
          </span>
          <span className="text-2xl font-extrabold text-white tracking-tight drop-shadow-lg">GitNotify</span>
        </div>
        <Link href="/config" className="bg-gradient-to-r from-green-500 via-emerald-600 to-black text-white px-6 py-2 rounded-xl font-semibold shadow-xl hover:from-green-400 hover:to-black transition flex items-center gap-2 border border-white/10 backdrop-blur-md animate-glow">
          Dashboard <ArrowRightIcon className="h-5 w-5" />
        </Link>
      </header>
      {/* Hero Section */}
      <main className="flex-1 flex flex-col items-center justify-center text-center px-4 relative">
        <div className="absolute inset-0 z-0 bg-gradient-to-br from-white/5 via-white/0 to-transparent pointer-events-none" />
        <div className="relative z-10 max-w-2xl w-full flex flex-col items-center">
          <div className="mb-8 flex flex-col items-center">
            <div className="p-2 rounded-2xl bg-white/10 backdrop-blur-lg shadow-2xl mb-4 border border-white/10 animate-float">
              <Image src="/GitNotify.png" alt="GitNotify Logo" width={96} height={96} className="rounded-xl shadow-xl" style={{boxShadow: '0 8px 32px 0 rgba(0,255,120,0.18)'}} />
            </div>
            <h1 className="text-5xl md:text-6xl font-extrabold text-white mb-4 drop-shadow-xl">
              <span className="bg-gradient-to-r from-green-400 via-emerald-400 to-white bg-clip-text text-transparent">Real-time, Low-Noise</span>
              <br />
              <span className="text-white/90">GitHub Notifications</span>
            </h1>
            <p className="text-xl text-white/80 mb-8 font-medium drop-shadow-lg">
              GitNotify is a powerful, configurable webhook listener that keeps your team in sync with issues and pull requests—without the noise. Secure, filterable, and easy to manage.
            </p>
            <Link href="/config" className="inline-flex items-center gap-2 bg-gradient-to-r from-green-500 via-emerald-600 to-black text-white px-10 py-4 rounded-2xl font-bold shadow-2xl hover:from-green-400 hover:to-black transition text-xl border border-white/10 backdrop-blur-md animate-glow">
              Get Started <ArrowRightIcon className="h-6 w-6" />
            </Link>
          </div>
        </div>
        {/* Features Section */}
        <section className="relative z-10 mt-16 grid grid-cols-1 md:grid-cols-2 gap-10 max-w-4xl w-full">
          <FeatureCard
            icon={<BellAlertIcon className="h-10 w-10 text-green-400" />}
            title="Real-Time Webhook Listener"
            desc="Instantly receive and process GitHub events for issues and pull requests."
          />
          <FeatureCard
            icon={<AdjustmentsHorizontalIcon className="h-10 w-10 text-emerald-400" />}
            title="Configurable Filtering"
            desc="Filter by event type, action, or repository to reduce notification noise."
          />
          <FeatureCard
            icon={<ShieldCheckIcon className="h-10 w-10 text-green-300" />}
            title="Secure by Default"
            desc="Validates webhook signatures and keeps your secrets safe."
          />
          <FeatureCard
            icon={<Image src="/GitNotify.png" alt="GitNotify Logo" width={36} height={36} className="rounded shadow-lg animate-float" style={{boxShadow: '0 2px 12px 0 rgba(0,255,120,0.12)'}} />}
            title="Easy Configuration UI"
            desc="Modern web dashboard for managing your notification rules."
          />
        </section>
      </main>
      {/* Footer */}
      <footer className="w-full py-8 text-center text-white/70 text-base flex flex-col items-center gap-2 mt-16 bg-white/5 backdrop-blur-lg border-t border-white/10">
        <span>&copy; {new Date().getFullYear()} <span className="font-bold text-white/90">GitNotify</span>. All rights reserved.</span>
        <span>Built with Next.js, TailwindCSS, and Go</span>
      </footer>
      {/* Animations */}
      <style jsx global>{`
        @keyframes float {
          0%, 100% { transform: translateY(0); }
          50% { transform: translateY(-10px); }
        }
        .animate-float {
          animation: float 3.5s ease-in-out infinite;
        }
        @keyframes glow {
          0%, 100% { box-shadow: 0 0 16px 2px #22c55e44, 0 0 0 0 #0000; }
          50% { box-shadow: 0 0 32px 8px #22c55e99, 0 0 0 0 #0000; }
        }
        .animate-glow {
          animation: glow 2.5s ease-in-out infinite;
        }
      `}</style>
    </div>
  );
}

function FeatureCard({ icon, title, desc }: { icon: React.ReactNode; title: string; desc: string }) {
  return (
    <div className="bg-white/10 backdrop-blur-xl rounded-2xl shadow-xl p-8 flex flex-col items-center text-center border border-white/10 hover:scale-105 hover:shadow-2xl transition-transform duration-300 group relative overflow-hidden">
      <div className="mb-4 group-hover:animate-float">{icon}</div>
      <h3 className="text-xl font-bold text-white mb-2 drop-shadow-lg">{title}</h3>
      <p className="text-white/80 text-base font-medium drop-shadow-sm">{desc}</p>
      <span className="absolute inset-0 pointer-events-none group-hover:bg-gradient-to-br group-hover:from-green-400/10 group-hover:to-transparent transition-all duration-300" />
    </div>
  );
}
