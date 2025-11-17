"use client";

import React, { useCallback, useEffect, useRef, useState } from "react";
import Link from "next/link";

type Suggestion = {
  id: string;
  type: string;
  content: string;
  applied?: boolean;
};

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
  const [isDragging, setIsDragging] = useState(false);

  const onFileChange = (f: File | null) => {
    if (!f) return;
    if (f.size > 5 * 1024 * 1024) {
      setError("File too large (max 5MB)");
      return;
    }
    setError(null);
    setFile(f);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0] || null;
    onFileChange(f);
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
        if (!r.ok)
          throw new Error(j?.error || j?.message || "Failed to fetch status");
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

      const id = j.jobId || j.id;
      if (id) {
        setJobId(id);
        pollStatus(id);
      } else {
        setResult(j);
      }
    } catch (e: any) {
      setError(e.message || "Unexpected error");
    } finally {
      setUploading(false);
    }
  };

  const handleDrop: React.DragEventHandler<HTMLLabelElement> = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      onFileChange(e.dataTransfer.files[0]);
      e.dataTransfer.clearData();
    }
  };

  const score = result?.score ?? null;
  const suggestions = (result?.suggestions as Suggestion[] | undefined) || [];

  const circleR = 44;
  const circleC = 2 * Math.PI * circleR;
  const dashOffset =
    score != null ? circleC - (score / 100) * circleC : circleC;

  return (
    <div className="min-h-screen bg-gradient-to-br from-teal-50 via-white to-cyan-50">
      {/* Hero */}
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 pt-10 pb-4">
        <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-black">
          CV Analyzer
        </h1>
        <p className="mt-2 text-gray-600">
          Upload your CV (PDF/DOC/DOCX) to get an instant score and actionable
          suggestions.
        </p>
      </div>

      {/* Content */}
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 pb-16 grid gap-8 lg:grid-cols-3">
        {/* Upload card */}
        <div className="lg:col-span-2">
          <div className="rounded-2xl border border-gray-200 bg-white shadow-sm p-6">
            <label
              htmlFor="cv-file"
              onDragOver={(e) => {
                e.preventDefault();
                setIsDragging(true);
              }}
              onDragLeave={() => setIsDragging(false)}
              onDrop={handleDrop}
              className={`block cursor-pointer rounded-xl border-2 border-dashed p-8 text-center transition-colors ${
                isDragging
                  ? "border-[#44C3BB] bg-[#44C3BB]/5"
                  : "border-gray-300 hover:border-gray-400"
              }`}
            >
              <div className="mx-auto w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center mb-4">
                📄
              </div>
              <p className="text-black font-semibold">
                Drag & drop your CV here
              </p>
              <p className="text-gray-500 text-sm mt-1">
                or click to select a file (PDF, DOC, DOCX, up to 5MB)
              </p>
              <input
                id="cv-file"
                type="file"
                accept=".pdf,.doc,.docx"
                className="sr-only"
                onChange={handleInputChange}
              />
            </label>

            {file && (
              <div className="mt-4 flex items-center justify-between rounded-lg border border-gray-200 p-3">
                <div className="text-sm text-gray-700 truncate mr-4">
                  <span className="font-medium text-black">Selected:</span>{" "}
                  {file.name}
                </div>
                <button
                  onClick={upload}
                  disabled={uploading}
                  className="px-4 py-2 rounded-xl bg-[#44C3BB] text-white font-semibold shadow-sm hover:bg-[#3bb3ac] disabled:opacity-60"
                >
                  {uploading ? "Uploading…" : "Analyze CV"}
                </button>
              </div>
            )}

            {error && <p className="mt-4 text-sm text-red-600">{error}</p>}

            {jobId &&
              (result?.status === "Pending" ||
                result?.status === "Processing" ||
                !result) && (
                <div className="mt-4 flex items-center gap-2 text-sm">
                  <span className="inline-flex items-center gap-2 rounded-full bg-amber-50 text-amber-700 px-3 py-1 border border-amber-200">
                    Job {jobId}
                  </span>
                  <span className="text-gray-600">
                    {result?.status || "Queued"}…
                  </span>
                </div>
              )}
          </div>
        </div>

        {/* Tips / CTA */}
        <div>
          <div className="rounded-2xl border border-gray-200 bg-white shadow-sm p-6">
            <h3 className="text-lg font-bold text-black">
              Tips for a higher score
            </h3>
            <ul className="mt-3 space-y-2 text-sm text-gray-700 list-disc list-inside">
              <li>Quantify impact (e.g., "Increased conversion by 12%")</li>
              <li>Align skills with job description keywords</li>
              <li>Keep formatting simple and machine-readable</li>
            </ul>
            <div className="mt-6 grid gap-3">
              <Link
                href="/jobs"
                className="w-full text-center rounded-xl bg-gray-900 text-white py-2.5 font-semibold hover:bg-black"
              >
                Browse Jobs
              </Link>
              <Link
                href="/user-home/personalized-jobs"
                className="w-full text-center rounded-xl bg-[#44C3BB] text-white py-2.5 font-semibold hover:bg-[#3bb3ac]"
              >
                View My Matches
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Results */}
      {result && (
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 pb-20">
          <div className="rounded-2xl border border-gray-200 bg-white shadow-sm p-6">
            <div className="flex flex-col lg:flex-row lg:items-center gap-8">
              {/* Score ring */}
              <div className="relative w-28 h-28">
                <svg className="w-28 h-28 transform -rotate-90">
                  <circle
                    cx="56"
                    cy="56"
                    r={circleR}
                    stroke="#e5e7eb"
                    strokeWidth="8"
                    fill="none"
                  />
                  <circle
                    cx="56"
                    cy="56"
                    r={circleR}
                    stroke="#44C3BB"
                    strokeWidth="8"
                    fill="none"
                    strokeDasharray={circleC}
                    strokeDashoffset={dashOffset}
                    strokeLinecap="round"
                  />
                </svg>
                <span className="absolute inset-0 flex items-center justify-center text-2xl font-extrabold text-black">
                  {result.score ?? 0}%
                </span>
              </div>

              <div className="flex-1 min-w-0">
                <h2 className="text-xl font-bold text-black">
                  Overall CV Score
                </h2>
                <p className="text-sm text-gray-600 mt-1">
                  {jobId ? `CV ID: ${jobId}` : ""}
                </p>
                {result.status !== "Completed" && (
                  <p className="mt-2 text-amber-700 bg-amber-50 inline-block px-2.5 py-1 rounded-full text-xs border border-amber-200">
                    Status: {result.status}
                  </p>
                )}
              </div>
            </div>

            {!!suggestions.length && (
              <div className="mt-8">
                <h3 className="font-semibold text-black mb-3">Suggestions</h3>
                <ul className="grid sm:grid-cols-2 gap-3">
                  {suggestions.map((s) => (
                    <li
                      key={s.id}
                      className="p-4 rounded-xl border border-gray-200 bg-gray-50"
                    >
                      <p className="text-xs text-gray-500 uppercase tracking-wide">
                        {s.type}
                      </p>
                      <p className="text-sm text-gray-800 mt-1">{s.content}</p>
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
