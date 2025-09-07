// app/components/JobDetails.tsx
"use client";

import { useEffect, useState } from "react";
import JobCard from "./JobList"; // now using JobCard directly

interface JobData {
  id: string;
  title: string;
  company_name?: string;
  company?: string;
  location?: string;
  description?: string;
  salary?: string;
  posted_at?: string;
  tags?: string[];
  applicants?: number;
  apply_url?: string;
  job_type?: string;
}

const JobDetails = ({ jobId }: { jobId: string }) => {
  const [job, setJob] = useState<JobData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchJob() {
      setLoading(true);
      try {
        const res = await fetch(`http://localhost:8080/api/v1/jobs/${jobId}`);
        const data = await res.json();
        const jobFromResponse =
          data?.data?.job ?? data?.data ?? data?.job ?? null;
        setJob(jobFromResponse);
      } catch (err) {
        console.error("Error fetching job", err);
        setJob(null);
      } finally {
        setLoading(false);
      }
    }
    if (jobId) fetchJob();
  }, [jobId]);

  if (loading) return <div className="text-center py-8">Loading job...</div>;
  if (!job)
    return <div className="text-center py-8 text-gray-600">Job not found.</div>;

  return (
    <div className="p-6 bg-gray-50 min-h-screen">
      <JobCard
        title={job.title}
        company={job.company_name ?? job.company ?? "Company"}
        location={job.location ?? "Remote"}
        tags={job.tags ?? []}
        salary={job.salary ?? "Competitive"}
        posted={
          job.posted_at ? new Date(job.posted_at).toLocaleDateString() : "N/A"
        }
        deadline={"N/A"}
        applicants={job.applicants ?? 0}
        description={job.description}
        applyUrl={job.apply_url}
      />
    </div>
  );
};

export default JobDetails;
