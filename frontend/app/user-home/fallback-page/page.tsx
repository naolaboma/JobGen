'use client';

import { useDispatch, useSelector } from 'react-redux';
import { useFetchJobsQuery } from '@/lib/redux/api/JobApi';
import JobCard from '@/app/components/JobCard';
import SearchBar from '@/app/components/SearchBar';
import Filters from '@/app/components/Filters';
import Pagination from '@/app/components/Pagination';
import Link from 'next/link';
import { setFilters, setPage, setSort } from '@/lib/redux/slices/jobSlice';
import { RootState } from '@/store/store';

export default function FallbackPage() {
  const dispatch = useDispatch();
  const filters = useSelector((state: RootState) => state.job);

  const { data, isLoading, error } = useFetchJobsQuery(filters);
  const totalPages = data ? Math.ceil(data.length / (filters.limit ?? 10)) : 1;

  const handleSearch = (query: string) => {
    dispatch(setFilters({ query }));
  };

  const handleFilterChange = (newFilters: { skills: string; location: string; sponsorship?: boolean }) => {
    dispatch(setFilters(newFilters));
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const { name, value } = e.target;
    const sortOrder = name === 'sort_order' ? (value as 'asc' | 'desc') : filters.sort_order ?? 'desc';
    const sortBy = name === 'sort_by' ? value : filters.sort_by ?? 'posted_at';
    dispatch(setSort({ sort_by: sortBy, sort_order: sortOrder }));
  };

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
      <h1 className="text-2xl font-bold mb-6 text-gray-900">Job Listings</h1>
      <p className="mb-6 text-gray-700">Set up your profile for personalized recommendations! <Link href="/profile-setup" className="text-[#7BBFB3] hover:underline">Set Up Profile</Link></p>

      <SearchBar onSearch={handleSearch} initialQuery={filters.query ?? ''} />
      <Filters
        onFilterChange={handleFilterChange}
        initialSkills={filters.skills ?? ''}
        initialLocation={filters.location ?? ''}
        initialSponsorship={filters.sponsorship}
      />
      <div className="mb-6 bg-gray-50 p-4 rounded-lg">
        <select name="sort_by" value={filters.sort_by ?? 'posted_at'} onChange={handleSortChange} className="p-2 border rounded-md mr-2">
          <option value="posted_at">Posted Date</option>
          <option value="title">Title</option>
        </select>
        <select name="sort_order" value={filters.sort_order ?? 'desc'} onChange={handleSortChange} className="p-2 border rounded-md">
          <option value="desc">Descending</option>
          <option value="asc">Ascending</option>
        </select>
      </div>

      <div className="grid gap-6">
        {data?.map((job) => (
          <JobCard
            key={job.id}
            {...job}
            percentage={Math.floor(Math.random() * 100)} // Placeholder
          />
        ))}
      </div>
      <Pagination currentPage={filters.page ?? 1} totalPages={totalPages} onPageChange={handlePageChange} />
    </div>
  );
}