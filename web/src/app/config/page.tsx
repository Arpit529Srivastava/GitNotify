"use client";
import React, { useEffect, useState } from "react";
import { InformationCircleIcon, CheckCircleIcon, ExclamationCircleIcon, Cog6ToothIcon } from "@heroicons/react/24/outline";
import Image from "next/image";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/config";
const TOKEN = process.env.NEXT_PUBLIC_CONFIG_TOKEN || "devtoken";

// Types for config
export type Notification = {
  event_type: string;
  actions?: string[];
  repos?: string[];
};

export type GitHubApp = {
  app_id?: number;
  installation_id?: number;
  private_key_path?: string;
};

export type Config = {
  organization: string;
  port: number;
  webhook_secret: string;
  notifications: Notification[];
  github_app?: GitHubApp;
};

// Helper for floating label
function FloatingLabel({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="relative my-4">
      {children}
      <span className="absolute left-3 top-[-0.7rem] bg-black/80 px-1 text-xs text-green-400 font-medium pointer-events-none">
        {label}
      </span>
    </div>
  );
}

export default function ConfigPage() {
  const [config, setConfig] = useState<Config | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [notifText, setNotifText] = useState<string>("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetch(API_URL, {
      headers: { Authorization: `Bearer ${TOKEN}` },
    })
      .then(async (res) => {
        if (!res.ok) throw new Error(await res.text());
        return res.json();
      })
      .then((data) => {
        setConfig(data);
        setNotifText(JSON.stringify(data.notifications, null, 2));
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!config) return;
    setConfig({ ...config, [e.target.name]: e.target.value } as Config);
  };

  const handleNotifChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setNotifText(e.target.value);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setSaving(true);
    let notifications: Notification[];
    try {
      notifications = JSON.parse(notifText);
    } catch {
      setError("Notifications must be valid JSON");
      setSaving(false);
      return;
    }
    const updated = { ...config, notifications } as Config;
    const res = await fetch(API_URL, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${TOKEN}`,
      },
      body: JSON.stringify(updated),
    });
    if (res.ok) {
      setSuccess("Configuration updated successfully!");
    } else {
      setError(await res.text());
    }
    setSaving(false);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-[#0a1a0f] to-[#1a3a1a] flex">
      {/* Sidebar */}
      <aside className="w-64 bg-white/10 border-r border-white/10 flex flex-col items-center py-8 shadow-2xl backdrop-blur-xl z-10">
        <div className="flex flex-col items-center gap-2 mb-8">
          <span className="inline-block animate-float">
            <Image src="/GitNotify.png" alt="GitNotify Logo" width={48} height={48} className="rounded-lg shadow-lg" style={{boxShadow: '0 4px 24px 0 rgba(0,255,120,0.18)'}} />
          </span>
          <span className="text-xl font-extrabold text-white tracking-tight drop-shadow-lg">GitNotify</span>
        </div>
        <nav className="flex flex-col gap-4 w-full px-6">
          <a href="/config" className="flex items-center gap-2 text-green-400 font-semibold bg-green-900/20 rounded px-3 py-2">
            <Cog6ToothIcon className="h-5 w-5" />
            Configuration
          </a>
        </nav>
        <div className="mt-auto mb-4 text-green-900/60 text-xs">Powered by Next.js & TailwindCSS</div>
      </aside>
      {/* Main Content */}
      <main className="flex-1 flex flex-col items-center justify-start py-12 px-4">
        <div className="w-full max-w-2xl bg-white/10 rounded-2xl shadow-2xl p-8 border border-white/10 backdrop-blur-xl">
          <div className="flex items-center gap-3 mb-6">
            <Cog6ToothIcon className="h-7 w-7 text-green-400" />
            <h1 className="text-2xl font-bold text-white">Configuration Dashboard</h1>
          </div>
          <p className="text-green-200 mb-8 flex items-center gap-2">
            <InformationCircleIcon className="h-5 w-5 text-green-400" />
            Manage your GitNotify instance settings below. Changes are applied live.
          </p>
          {loading ? (
            <div className="flex items-center gap-2 text-green-400 animate-pulse">
              <Cog6ToothIcon className="h-5 w-5 animate-spin" /> Loading configuration...
            </div>
          ) : error ? (
            <div className="flex items-center gap-2 text-red-400">
              <ExclamationCircleIcon className="h-5 w-5" /> Error: {error}
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-6">
              <FloatingLabel label="Organization">
                <input
                  className="border-2 border-green-700 bg-black/60 text-green-100 rounded-lg px-3 py-3 w-full focus:outline-none focus:ring-2 focus:ring-green-400 transition"
                  name="organization"
                  value={config?.organization || ""}
                  onChange={handleChange}
                  required
                  autoComplete="off"
                />
              </FloatingLabel>
              <FloatingLabel label="Port">
                <input
                  className="border-2 border-green-700 bg-black/60 text-green-100 rounded-lg px-3 py-3 w-full focus:outline-none focus:ring-2 focus:ring-green-400 transition"
                  name="port"
                  type="number"
                  value={config?.port || 8080}
                  onChange={handleChange}
                  required
                  min={1}
                  max={65535}
                />
              </FloatingLabel>
              <FloatingLabel label="Webhook Secret">
                <input
                  className="border-2 border-green-700 bg-black/60 text-green-100 rounded-lg px-3 py-3 w-full focus:outline-none focus:ring-2 focus:ring-green-400 transition"
                  name="webhook_secret"
                  value={config?.webhook_secret || ""}
                  onChange={handleChange}
                  required
                  autoComplete="off"
                  type="password"
                />
              </FloatingLabel>
              <FloatingLabel label="Notifications (JSON)">
                <textarea
                  className="border-2 border-green-700 bg-black/60 text-green-100 rounded-lg px-3 py-3 w-full font-mono focus:outline-none focus:ring-2 focus:ring-green-400 transition"
                  rows={8}
                  value={notifText}
                  onChange={handleNotifChange}
                  required
                  spellCheck={false}
                />
                <span className="text-xs text-green-400 flex items-center gap-1 mt-1">
                  <InformationCircleIcon className="h-4 w-4" /> Paste or edit notification rules as JSON
                </span>
              </FloatingLabel>
              {/* Future: Add GitHub App section here */}
              <button
                type="submit"
                className="bg-gradient-to-r from-green-500 via-emerald-600 to-black text-white px-6 py-3 rounded-lg font-semibold shadow-xl hover:from-green-400 hover:to-black transition flex items-center gap-2 disabled:opacity-60 animate-glow"
                disabled={saving}
              >
                {saving ? (
                  <span className="flex items-center gap-2"><Cog6ToothIcon className="h-5 w-5 animate-spin" /> Saving...</span>
                ) : (
                  <span className="flex items-center gap-2"><CheckCircleIcon className="h-5 w-5" /> Save Changes</span>
                )}
              </button>
              {success && (
                <div className="flex items-center gap-2 text-green-400 mt-2">
                  <CheckCircleIcon className="h-5 w-5" /> {success}
                </div>
              )}
              {error && (
                <div className="flex items-center gap-2 text-red-400 mt-2">
                  <ExclamationCircleIcon className="h-5 w-5" /> {error}
                </div>
              )}
            </form>
          )}
        </div>
        <footer className="mt-8 text-green-900/60 text-xs flex items-center gap-2">
          <svg width="20" height="20" fill="currentColor" className="text-green-900/60"><path d="M10 .5a9.5 9.5 0 1 0 0 19 9.5 9.5 0 0 0 0-19Zm0 1.5a8 8 0 1 1 0 16 8 8 0 0 1 0-16Zm.25 2.25a.75.75 0 0 0-.75.75v2.5a.75.75 0 0 0 1.5 0v-2.5a.75.75 0 0 0-.75-.75Zm0 10a.75.75 0 0 0-.75.75v2.5a.75.75 0 0 0 1.5 0v-2.5a.75.75 0 0 0-.75-.75Zm-5-5a.75.75 0 0 0-.75.75v2.5a.75.75 0 0 0 1.5 0v-2.5a.75.75 0 0 0-.75-.75Zm10 0a.75.75 0 0 0-.75.75v2.5a.75.75 0 0 0 1.5 0v-2.5a.75.75 0 0 0-.75-.75ZM5.47 5.47a.75.75 0 0 0-1.06 1.06l1.77 1.77a.75.75 0 0 0 1.06-1.06L5.47 5.47Zm7.06 7.06a.75.75 0 0 0-1.06 1.06l1.77 1.77a.75.75 0 0 0 1.06-1.06l-1.77-1.77ZM2.75 10a.75.75 0 0 0-.75.75v2.5a.75.75 0 0 0 1.5 0v-2.5A.75.75 0 0 0 2.75 10Zm14.5 0a.75.75 0 0 0-.75.75v2.5a.75.75 0 0 0 1.5 0v-2.5a.75.75 0 0 0-.75-.75Z"/></svg>
          &copy; {new Date().getFullYear()} GitNotify. All rights reserved.
        </footer>
      </main>
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