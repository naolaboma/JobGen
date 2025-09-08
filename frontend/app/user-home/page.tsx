import React from 'react'
import Navbar from '../components/NavBar';
import Footer from '../components/footer';
import JobCard from '../components/JobCard';

const page = () => {
  return (
    <div>
        <Navbar />
        <div className='ml-60 mt-8'>
            <h1 className="text-4xl font-bold">Find Your Next Job</h1>
            <h3 className="mb-6">A curated list of job openings for you</h3>
        </div>
        <Footer />  
    </div>
  )
}

export default page