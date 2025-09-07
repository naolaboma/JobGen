'use client';

import { ChangeEvent } from 'react';

export default function Filters({ onFilterChange, initialSkills, initialLocation, initialSponsorship }: FiltersProps) {
  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    onFilterChange({ skills: name === 'skills' ? value : initialSkills, location: name === 'location' ? value : initialLocation, sponsorship: name === 'sponsorship' ? checked : initialSponsorship });
  };

  return (
    <div className="mb-6 bg-gray-50 p-4 rounded-lg">
      <input
        name="skills"
        defaultValue={initialSkills}
        onChange={handleChange}
        placeholder="Skills (comma-separated)..."
        className="w-full p-2 mb-2 border rounded-md"
      />
      <input
        name="location"
        defaultValue={initialLocation}
        onChange={handleChange}
        placeholder="Location..."
        className="w-full p-2 mb-2 border rounded-md"
      />
      <label className="block mb-2">
        <input
          type="checkbox"
          name="sponsorship"
          defaultChecked={initialSponsorship || false}
          onChange={handleChange}
        /> Sponsorship Available
      </label>
    </div>
  );
}