import { createApi,fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const JobApi = createApi({
    reducerPath: 'jobApi',
    baseQuery: fetchBaseQuery({ baseUrl: 'http://localhost:8080/api/v1' }),
    endpoints: (builder) => ({

        // fetch all jobs
        fetchJobs: builder.query<JobProps[], JobQueryParams>({
            query: ({ page = 1, limit = 10, query, skills, location, sponsorship, source, sort_by = 'posted_at', sort_order = 'desc' } = {}) => ({
                url: '/jobs',
                method: 'GET',
                params: { page, limit, query, skills, location, sponsorship, source, sort_by, sort_order },
            }),
        }),
    }),
});

export const { useFetchJobsQuery } = JobApi;

export default JobApi;