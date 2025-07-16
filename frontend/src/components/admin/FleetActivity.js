import React from 'react';
import { useQuery } from 'react-query';
import { telemetryAPI } from '../../services/api';
import { ComputerDesktopIcon, UsersIcon } from '@heroicons/react/24/outline';

const FleetActivity = () => {
  const { data: fleetData, isLoading } = useQuery(
    'fleet-activity',
    telemetryAPI.getFleetActivity,
    {
      refetchInterval: 30000,
    }
  );

  const activities = fleetData?.data || [];

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
            Fleet Activity
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            Monitor activity across all devices in your organization
          </p>
        </div>
      </div>

      {/* Activity List */}
      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <ul className="divide-y divide-gray-200">
          {activities.map((activity, index) => (
            <li key={index} className="px-6 py-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <ComputerDesktopIcon className="h-6 w-6 text-gray-400" />
                  </div>
                  <div className="ml-4">
                    <div className="text-sm font-medium text-gray-900">
                      {activity.device_hostname}
                    </div>
                    <div className="text-sm text-gray-500">
                      {activity.activity_type}: {activity.description}
                    </div>
                  </div>
                </div>
                <div className="flex items-center text-sm text-gray-500">
                  <UsersIcon className="h-4 w-4 mr-1" />
                  {activity.user}
                  <span className="mx-2">•</span>
                  {new Date(activity.timestamp).toLocaleString()}
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default FleetActivity;
