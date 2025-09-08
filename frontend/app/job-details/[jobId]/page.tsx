// app/job-details/[jobId]/page.tsx
import JobDetails from "@/app/components/JobDetails";
import Navbar from "@/app/components/header";

export default function JobDetailsPage({
  params,
}: {
  params: { jobId: string };
}) {
  const { jobId } = params;
  
  return (<>
 <Navbar/>
  <JobDetails jobId={jobId} />

  </>);
}
