// pages/profile.tsx
import React from 'react';
import ProfileSetup from '@/app/components/ProfileSetup';
import Navbar from '@/app/components/header';
const ProfilePage: React.FC = () => {
  return (
    <div>
        <Navbar/>
      <ProfileSetup />
    </div>
  );
};

export default ProfilePage;