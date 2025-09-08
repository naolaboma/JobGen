interface FiltersProps {
  onFilterChange: (filters: { skills: string; location: string; sponsorship?: boolean }) => void;
  initialSkills: string;
  initialLocation: string;
  initialSponsorship: boolean | undefined;
}