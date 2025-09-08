import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface JobQueryParams {
  page?: number;
  limit?: number;
  query?: string;
  skills?: string;
  location?: string;
  sponsorship?: boolean;
  source?: string;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}

const initialState: JobQueryParams = {
  page: 1,
  limit: 2,
  query: "",
  skills: "",
  location: "",
  sponsorship: undefined,
  sort_by: "posted_at",
  sort_order: "desc",
};

const jobSlice = createSlice({
  name: "job",
  initialState,
  reducers: {
    setFilters: (state, action: PayloadAction<Partial<JobQueryParams>>) => {
      return { ...state, ...action.payload, page: 1 };
    },
    setPage: (state, action: PayloadAction<number>) => {
      state.page = action.payload;
    },
    setSort: (
      state,
      action: PayloadAction<{ sort_by: string; sort_order: "asc" | "desc" }>
    ) => {
      state.sort_by = action.payload.sort_by;
      state.sort_order = action.payload.sort_order;
    },
  },
});

export const { setFilters, setPage, setSort } = jobSlice.actions;
export default jobSlice.reducer;
