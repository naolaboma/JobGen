'use client';

import { ChangeEvent } from 'react';

interface SearchBarProps {
  onSearch: (query: string) => void;
  initialQuery: string;
}

export default function SearchBar({ onSearch, initialQuery }: SearchBarProps) {
  const handleSubmit = (e: ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSearch(e.target.query.value);
  };

  return (
    <form onSubmit={handleSubmit} className="mb-4">
      <input
        type="text"
        name="query"
        defaultValue={initialQuery}
        placeholder="Search titles/companies..."
        className="w-full p-2 border rounded-md"
      />
      <button type="submit" className="ml-2 px-4 py-2 bg-[#7BBFB3] text-white rounded-md hover:bg-[#6aa69e]">
        Search
      </button>
    </form>
  );
}