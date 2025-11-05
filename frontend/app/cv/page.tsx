"use client";

import React, { useCallback, useEffect, useRef, useState } from "react";

type Suggestion = { id: string; type: string; content: string; applied?: boolean };

type CVResult = {
  id: string;
  status: "Pending" | "Processing" | "Completed" | "Failed";
  score?: number;
  suggestions?: Suggestion[];
  message?: string;
};

export default function CvPage() {
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [jobId, setJobId] = useState<string | null>(null);
  const [result, setResult] = useState<CVResult | null>(null);
  const pollTimer = useRef<NodeJS.Timeout | null>(null);

  const onFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0] || null;
    if (!f) return;
    if (f.size > 5 * 1024 * 1024) {
      setError("File too large (max 5MB)");
      return;
    }
    setError(null);
    setFile(f);
  };

  const stopPolling = useCallback(() => {
    if (pollTimer.current) {
      clearTimeout(pollTimer.current);
      pollTimer.current = null;
    }
  }, []);

  const pollStatus = useCallback(
    async (id: string) => {
      try {
        const r = await fetch(`/api/cv/${encodeURIComponent(id)}`);
        const j = await r.json();
        if (!r.ok) throw new Error(j?.error || j?.message || "Failed to fetch status");
        setResult(j);
        if (j.status === "Completed" || j.status === "Failed") {
          stopPolling();
        } else {
          pollTimer.current = setTimeout(() => pollStatus(id), 2000);
        }
      } catch (e: any) {
        setError(e.message || "Failed to fetch status");
        pollTimer.current = setTimeout(() => pollStatus(id), 4000);
      }
    },
    [stopPolling]
  );

  useEffect(() => {
    return () => stopPolling();
  }, [stopPolling]);

  const upload = async () => {
    if (!file) {
      setError("Please select a file first.");
      return;
    }
    setUploading(true);
    setError(null);
    setResult(null);
    setJobId(null);

    try {
      const form = new FormData();
      form.append("file", file);
      const r = await fetch("/api/cv/parse", { method: "POST", body: form });
      const j = await r.json();
      if (!r.ok) throw new Error(j?.error || j?.message || "Upload failed");

      // If backend returns just jobId, start polling.
      const id = j.jobId || j.id;
      if (id) {
        setJobId(id);
        pollStatus(id);
      } else {
        // If immediate result
        setResult(j);
      }
    } catch (e: any) {
      setError(e.message || "Unexpected error");
    } finally {
      setUploading(false);
    }
  };

  const score = result?.score ?? null;
  const suggestions = (result?.suggestions as Suggestion[] | undefined) || [];

  const circleR = 32;
  const circleC = 2 * Math.PI * circleR;
  const dashOffset = score != null ? circleC - (score / 100) * circleC : circleC;

  return (
    <div className="max-w-3xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">CV Scoring</h1>

      <div className="flex items-center gap-3 mb-4">
        <input type="file" accept=".pdf,.doc,.docx" onChange={onFileChange} />
        <button
          onClick={upload}
          disabled={!file || uploading}
          className="px-4 py-2 rounded bg-[#44C3BB] text-white disabled:opacity-50"
        >
          {uploading ? "Uploading..." : "Analyze CV"}
        </button>
      </div>

      {error && <p className="text-red-600 mb-3">{error}</p>}

      {jobId && (result?.status === "Pending" || result?.status === "Processing" || !result) && (
        <p className="text-sm text-gray-600 mb-4">Job {jobId}: {result?.status || "Queued"}...</p>
      )}

      {score !== null && (
        <div className="mb-6">
          <div className="flex items-center gap-4">
            <div className="relative w-20 h-20">
              <svg className="w-20 h-20 transform -rotate-90">
                <circle cx="40" cy="40" r={circleR} stroke="#e5e7eb" strokeWidth="6" fill="none" />
                <circle
                  cx="40"
                  cy="40"
                  r={circleR}
                  stroke="#44C3BB"
                  strokeWidth="6"
                  fill="none"
                  strokeDasharray={circleC}
                  strokeDashoffset={dashOffset}
                />
              </svg>
              <span className="absolute inset-0 flex items-center justify-center text-xl font-bold text-black">
                {score}%
              </span>
            </div>
            <div>
              <p className="text-gray-600">Overall CV Score</p>
              {jobId && <p className="text-xs text-gray-400">CV ID: {jobId}</p>}
            </div>
          </div>
        </div>
      )}

      {!!suggestions.length && (
        <div>
          <h2 className="font-semibold mb-2">Suggestions</h2>
          <ul className="space-y-2">
            {suggestions.map((s) => (
              <li key={s.id} className="p-3 rounded bg-gray-100">
                <p className="text-xs text-gray-500">{s.type}</p>
                <p className="text-sm">{s.content}</p>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
