import React from 'react';
import { useQuery } from 'react-query';
import { dashboardAPI } from '../../services/api';
import {
  UsersIcon,
  ComputerDesktopIcon,
  ExclamationTriangleIcon,
  ChartBarIcon
} from '@heroicons/react/24/outline';

const AdminDashboard = () => {
  const { data: adminData, isLoading } = useQuery(
    'admin-overview',
    dashboardAPI.getAdminOverview,
    {
      refetchInterval: 30000,
    }
  );

  const stats = adminData?.data || {
    totalDevices: 0,
    activeUsers: 0,
    totalThreats: 0,
    avgHealthScore: 0,
  };

  const StatCard = ({ title, value, icon: Icon, color = 'blue' }) => (
    <div className="bg-white overflow-hidden shadow rounded-lg">
      <div className="p-5">
        <div className="flex items-center">
          <div className="flex-shrink-0">
            <Icon className={`h-6 w-6 text-${color}-600`} />
          </div>
          <div className="ml-5 w-0 flex-1">
            <dl>
              <dt className="text-sm font-medium text-gray-500 truncate">{title}</dt>
              <dd className="text-2xl font-semibold text-gray-900">{value}</dd>
            </dl>
          </div>
        </div>
      </div>
    </div>
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="md:flex md:items-center md:justify-between">
        <div className="flex-1 min-w-0">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Admin Dashboard
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Fleet-wide overview and management
          </p>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Total Devices"
          value={stats.totalDevices}
          icon={ComputerDesktopIcon}
        />
        <StatCard
          title="Active Users"
          value={stats.activeUsers}
          icon={UsersIcon}
          color="green"
        />
        <StatCard
          title="Total Threats"
          value={stats.totalThreats}
          icon={ExclamationTriangleIcon}
          color="red"
        />
        <StatCard
          title="Avg Health Score"
          value={`${stats.avgHealthScore}%`}
          icon={ChartBarIcon}
          color="yellow"
        />
      </div>

      {/* Quick Actions */}
      <div className="bg-white shadow rounded-lg p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Quick Actions</h3>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <button className="flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
            View Fleet Activity
          </button>
          <button className="flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
            Manage Users
          </button>
          <button className="flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
            Security Reports
          </button>
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;
