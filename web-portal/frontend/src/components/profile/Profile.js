import React from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { UserIcon } from '@heroicons/react/24/outline';

const Profile = () => {
  const { user } = useAuth();

  return (
    <div className="space-y-6">
      <div className="bg-white shadow rounded-lg p-6">
        <div className="flex items-center space-x-4">
          <div className="w-16 h-16 bg-blue-500 rounded-full flex items-center justify-center">
            <UserIcon className="w-8 h-8 text-white" />
          </div>
          <div>
            <h2 className="text-xl font-bold text-gray-900">{user?.name}</h2>
            <p className="text-gray-600">{user?.email}</p>
            <p className="text-sm text-gray-500">{user?.role} • {user?.department}</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Profile;
