interface JobQueryParams {
  page?: number;
  limit?: number;
  query?: string;
  skills?: string;
  location?: string;
  sponsorship?: boolean;
  source?: string;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}