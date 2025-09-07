// app/job-details/[jobId]/page.tsx
import JobDetails from "@/app/components/JobDetails";

export default function JobDetailsPage({
  params,
}: {
  params: { jobId: string };
}) {
  const { jobId } = params;
  // JobDetails is a client component (it does data fetching in the client)
  return <JobDetails jobId={jobId} />;
}
