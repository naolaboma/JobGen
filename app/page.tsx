// app/page.tsx
import JobCard from "@/components/JobListing";

export default function Home() {
  return (
    <JobCard
      title="Frontend Developer"
      company="JobGen"
      location="Boston, USA"
      tags={["Full-time", "Mid-level", "Remote", "Visa Sponsorship"]}
      salary="$70,000 - $80,000/year"
      posted="3 days ago"
      deadline="May 30"
      applicants={12}
    />
  );
}
