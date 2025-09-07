'use client';

import { useDispatch, useSelector } from 'react-redux';
import { useFetchMatchedJobsQuery } from '@/lib/redux/api/JobApi';
import JobCard from '@/app/components/JobCard';
import Pagination from '@/app/components/Pagination';
import Link from 'next/link';
import { setPage } from '@/lib/redux/slices/jobSlice';
import { RootState } from '@/store/store';

export default function PersonalizedJobsPage() {
  const dispatch = useDispatch();
  const filters = useSelector((state: RootState) => state.job);

  const { data, isLoading, error } = useFetchMatchedJobsQuery({ page: filters.page, limit: filters.limit });
  const totalPages = data ? Math.ceil(data.length / (filters.limit ?? 10)) : 1;

  const handlePageChange = (newPage: number) => {
    dispatch(setPage(newPage));
  };

  if (isLoading) return <div className="text-center py-10">Loading...</div>;
  if (error) {
    const errorMessage = 'error' in error ? error.error : 'Unknown error occurred';
    return <div className="text-center py-10 text-red-600">Error: {errorMessage}</div>;
  }

  return (
    <div className="min-h-screen bg-white p-6">
      <h1 className="text-2xl font-bold mb-6 text-gray-900">Personalized Job Recommendations</h1>
      <p className="mb-6 text-gray-700">Jobs matched to your profile! <Link href="/profile-setup" className="text-[#7BBFB3] hover:underline">Update Profile</Link></p>

      <div className="grid gap-6">
        {data?.map((job) => (
          <JobCard
            key={job.id}
            {...job}
            percentage={Math.floor(Math.random() * 100)} // Placeholder for match percentage
          />
        ))}
      </div>
      <Pagination currentPage={filters.page ?? 1} totalPages={totalPages} onPageChange={handlePageChange} />
    </div>
  );
}